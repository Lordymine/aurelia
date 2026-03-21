import { createInterface } from "node:readline";
import { query } from "@anthropic-ai/claude-agent-sdk";

// ── Types ────────────────────────────────────────────────────────────────────

interface MCPServerConfig {
  command: string;
  args?: string[];
  env?: Record<string, string>;
}

interface RequestOptions {
  model?: string;
  cwd?: string;
  system_prompt?: string;
  resume?: string;
  max_turns?: number;
  permission_mode?: string;
  mcp_servers?: Record<string, MCPServerConfig>;
  allowed_tools?: string[];
}

interface Request {
  command: string;
  prompt: string;
  options?: RequestOptions;
}

interface OutEvent {
  event: string;
  [key: string]: unknown;
}

// ── Helpers ──────────────────────────────────────────────────────────────────

function emit(obj: OutEvent): void {
  process.stdout.write(JSON.stringify(obj) + "\n");
}

function log(msg: string): void {
  process.stderr.write(`[bridge] ${msg}\n`);
}

// ── Map request options to SDK options ───────────────────────────────────────

function buildSDKOptions(opts: RequestOptions | undefined) {
  if (!opts) return {};

  const sdkOpts: Record<string, unknown> = {};

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

  return sdkOpts;
}

// ── Extract text from content blocks ─────────────────────────────────────────

function extractText(content: unknown): string {
  if (!Array.isArray(content)) return "";
  return content
    .filter(
      (block: unknown) =>
        typeof block === "object" &&
        block !== null &&
        "type" in block &&
        (block as Record<string, unknown>).type === "text" &&
        "text" in block,
    )
    .map((block: unknown) => (block as Record<string, string>).text)
    .join("");
}

// ── Handle a single query command ────────────────────────────────────────────

async function handleQuery(req: Request): Promise<void> {
  const sdkOptions = buildSDKOptions(req.options);

  log(`query start — model=${sdkOptions.model ?? "default"} prompt="${req.prompt.slice(0, 80)}..."`);

  const timeoutMs = 10 * 60 * 1000;
  const timeout = setTimeout(() => {
    log("query timeout — no result after 10 minutes");
    emit({ event: "error", message: "query timeout: no result after 10 minutes" });
    process.exit(1);
  }, timeoutMs);

  try {
    const stream = query({
      prompt: req.prompt,
      options: sdkOptions as Parameters<typeof query>[0]["options"],
    });

    for await (const message of stream) {
      const msg = message as Record<string, unknown>;
      const msgType = msg.type as string | undefined;

      switch (msgType) {
        // ── System init ──────────────────────────────────────────────
        case "system": {
          emit({
            event: "system",
            session_id: msg.session_id as string,
            tools: msg.tools as string[],
            model: msg.model as string,
          });
          break;
        }

        // ── Assistant text + tool_use blocks ────────────────────────
        case "assistant": {
          const inner = msg.message as Record<string, unknown> | undefined;
          if (inner?.content && Array.isArray(inner.content)) {
            const text = extractText(inner.content);
            if (text) {
              emit({ event: "assistant", text });
            }
            for (const block of inner.content as Record<string, unknown>[]) {
              if (block.type === "tool_use") {
                emit({
                  event: "tool_use",
                  id: block.id as string,
                  name: block.name as string,
                  input: block.input as Record<string, unknown>,
                });
              }
            }
          }
          break;
        }

        // ── Tool use summary ─────────────────────────────────────────
        case "tool_use_summary": {
          emit({
            event: "tool_result",
            content: msg.summary as string,
          });
          break;
        }

        // ── Result (success or error) ────────────────────────────────
        case "result": {
          const subtype = msg.subtype as string | undefined;
          if (subtype === "success") {
            emit({
              event: "result",
              content: msg.result as string,
              cost_usd: msg.total_cost_usd as number,
              session_id: msg.session_id as string,
              duration_ms: msg.duration_ms as number,
              num_turns: msg.num_turns as number,
            });
          } else {
            // error_max_turns, error_during_execution, etc.
            const errors = msg.errors as string[] | undefined;
            emit({
              event: "error",
              message: errors?.join("; ") ?? `result error: ${subtype}`,
              subtype: subtype ?? "unknown",
            });
          }
          break;
        }

        // ── All other message types (status, hooks, etc.) ───────────
        default: {
          // Not emitted — keep the protocol focused on what Go needs.
          break;
        }
      }
    }
  } catch (err: unknown) {
    const errMsg = err instanceof Error ? err.message : String(err);
    log(`query error: ${errMsg}`);
    emit({ event: "error", message: errMsg });
  } finally {
    clearTimeout(timeout);
  }
}

// ── Handle incoming request ──────────────────────────────────────────────────

async function handleRequest(line: string): Promise<void> {
  let req: Request;

  try {
    req = JSON.parse(line) as Request;
  } catch {
    emit({ event: "error", message: `invalid JSON: ${line.slice(0, 200)}` });
    return;
  }

  if (!req.command) {
    emit({ event: "error", message: "missing 'command' field" });
    return;
  }

  switch (req.command) {
    case "query": {
      if (!req.prompt) {
        emit({ event: "error", message: "missing 'prompt' field for query command" });
        return;
      }
      await handleQuery(req);
      break;
    }

    case "ping": {
      emit({ event: "pong" });
      break;
    }

    default: {
      emit({ event: "error", message: `unknown command: ${req.command}` });
    }
  }
}

// ── Main loop ────────────────────────────────────────────────────────────────

function main(): void {
  log("bridge started — waiting for commands on stdin");

  const rl = createInterface({
    input: process.stdin,
    terminal: false,
  });

  // Process one line at a time, sequentially
  let processing: Promise<void> = Promise.resolve();

  rl.on("line", (line: string) => {
    const trimmed = line.trim();
    if (!trimmed) return;

    // Chain requests sequentially — one query at a time
    processing = processing
      .then(() => handleRequest(trimmed))
      .catch((err: unknown) => {
        const errMsg = err instanceof Error ? err.message : String(err);
        log(`unhandled error in request processing: ${errMsg}`);
        emit({ event: "error", message: `internal bridge error: ${errMsg}` });
      });
  });

  rl.on("close", () => {
    log("stdin closed — shutting down");
    processing.then(() => process.exit(0));
  });

  // Catch unhandled rejections so the bridge never crashes silently
  process.on("unhandledRejection", (reason: unknown) => {
    const msg = reason instanceof Error ? reason.message : String(reason);
    log(`unhandled rejection: ${msg}`);
    emit({ event: "error", message: `unhandled rejection: ${msg}` });
  });

  process.on("uncaughtException", (err: Error) => {
    log(`uncaught exception: ${err.message}`);
    emit({ event: "error", message: `uncaught exception: ${err.message}` });
    process.exit(1);
  });
}

main();
