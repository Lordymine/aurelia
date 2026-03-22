// index.ts
import { createInterface } from "node:readline";
import { query } from "@anthropic-ai/claude-agent-sdk";
function emit(obj) {
  process.stdout.write(JSON.stringify(obj) + "\n");
}
function log(msg) {
  process.stderr.write(`[bridge] ${msg}
`);
}
function buildSDKOptions(opts) {
  if (!opts) return {};
  const sdkOpts = {};
  if (opts.model) sdkOpts.model = opts.model;
  if (opts.cwd) sdkOpts.cwd = opts.cwd;
  if (opts.system_prompt) sdkOpts.systemPrompt = opts.system_prompt;
  if (opts.resume) sdkOpts.resume = opts.resume;
  if (opts.max_turns) sdkOpts.maxTurns = opts.max_turns;
  if (opts.permission_mode) {
    sdkOpts.permissionMode = opts.permission_mode;
    if (opts.permission_mode === "bypassPermissions") {
      sdkOpts.allowDangerouslySkipPermissions = true;
    }
  }
  if (opts.mcp_servers) sdkOpts.mcpServers = opts.mcp_servers;
  if (opts.allowed_tools) sdkOpts.allowedTools = opts.allowed_tools;
  if (opts.no_user_settings) {
    sdkOpts.settingSources = [];
  } else {
    sdkOpts.settingSources = ["user", "project", "local"];
  }
  if (opts.disabled_tools && opts.disabled_tools.length > 0) {
    sdkOpts.disallowedTools = opts.disabled_tools;
  }
  return sdkOpts;
}
function extractText(content) {
  if (!Array.isArray(content)) return "";
  return content.filter(
    (block) => typeof block === "object" && block !== null && "type" in block && block.type === "text" && "text" in block
  ).map((block) => block.text).join("");
}
async function handleQuery(req) {
  const reqId = req.request_id || "";
  const emitReq = (obj) => emit({ ...obj, request_id: reqId });
  const sdkOptions = buildSDKOptions(req.options);
  log(`query start \u2014 rid=${reqId} model=${sdkOptions.model ?? "default"} prompt="${req.prompt.slice(0, 80)}..."`);
  const timeoutMs = 10 * 60 * 1e3;
  const timeout = setTimeout(() => {
    log(`query timeout \u2014 rid=${reqId} no result after 10 minutes`);
    emitReq({ event: "error", message: "query timeout: no result after 10 minutes" });
  }, timeoutMs);
  try {
    const stream = query({
      prompt: req.prompt,
      options: sdkOptions
    });
    for await (const message of stream) {
      const msg = message;
      const msgType = msg.type;
      switch (msgType) {
        // ── System init ──────────────────────────────────────────────
        case "system": {
          emitReq({
            event: "system",
            session_id: msg.session_id,
            tools: msg.tools,
            model: msg.model
          });
          break;
        }
        // ── Assistant text + tool_use blocks ────────────────────────
        case "assistant": {
          const inner = msg.message;
          if (inner?.content && Array.isArray(inner.content)) {
            const text = extractText(inner.content);
            if (text) {
              emitReq({ event: "assistant", text });
            }
            for (const block of inner.content) {
              if (block.type === "tool_use") {
                emitReq({
                  event: "tool_use",
                  id: block.id,
                  name: block.name,
                  input: block.input
                });
              }
            }
          }
          break;
        }
        // ── Tool use summary ─────────────────────────────────────────
        case "tool_use_summary": {
          emitReq({
            event: "tool_result",
            content: msg.summary
          });
          break;
        }
        // ── Result (success or error) ────────────────────────────────
        case "result": {
          const subtype = msg.subtype;
          if (subtype === "success") {
            emitReq({
              event: "result",
              content: msg.result,
              cost_usd: msg.total_cost_usd,
              session_id: msg.session_id,
              duration_ms: msg.duration_ms,
              num_turns: msg.num_turns
            });
          } else {
            const errors = msg.errors;
            emitReq({
              event: "error",
              message: errors?.join("; ") ?? `result error: ${subtype}`,
              subtype: subtype ?? "unknown"
            });
          }
          break;
        }
        // ── All other message types (status, hooks, etc.) ───────────
        default: {
          break;
        }
      }
    }
  } catch (err) {
    const errMsg = err instanceof Error ? err.message : String(err);
    log(`query error: rid=${reqId} ${errMsg}`);
    emitReq({ event: "error", message: errMsg });
  } finally {
    clearTimeout(timeout);
  }
}
async function handleRequest(line) {
  let req;
  try {
    req = JSON.parse(line);
  } catch {
    emit({ event: "error", message: `invalid JSON: ${line.slice(0, 200)}` });
    return;
  }
  if (!req.command) {
    emit({ event: "error", request_id: req.request_id || "", message: "missing 'command' field" });
    return;
  }
  const reqId = req.request_id || "";
  switch (req.command) {
    case "query": {
      if (!req.prompt) {
        emit({ event: "error", request_id: reqId, message: "missing 'prompt' field for query command" });
        return;
      }
      await handleQuery(req);
      break;
    }
    case "ping": {
      emit({ event: "pong", request_id: reqId });
      break;
    }
    default: {
      emit({ event: "error", request_id: reqId, message: `unknown command: ${req.command}` });
    }
  }
}
function main() {
  log("bridge started \u2014 waiting for commands on stdin");
  const rl = createInterface({
    input: process.stdin,
    terminal: false
  });
  rl.on("line", (line) => {
    const trimmed = line.trim();
    if (!trimmed) return;
    handleRequest(trimmed).catch((err) => {
      const errMsg = err instanceof Error ? err.message : String(err);
      log(`unhandled error in request processing: ${errMsg}`);
      emit({ event: "error", message: `internal bridge error: ${errMsg}` });
    });
  });
  rl.on("close", () => {
    log("stdin closed \u2014 shutting down");
    process.exit(0);
  });
  process.on("unhandledRejection", (reason) => {
    const msg = reason instanceof Error ? reason.message : String(reason);
    log(`unhandled rejection: ${msg}`);
    emit({ event: "error", message: `unhandled rejection: ${msg}` });
  });
  process.on("uncaughtException", (err) => {
    log(`uncaught exception: ${err.message}`);
    emit({ event: "error", message: `uncaught exception: ${err.message}` });
    process.exit(1);
  });
}
main();
