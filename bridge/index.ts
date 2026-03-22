import { createInterface } from "node:readline";
import { readFile } from "node:fs/promises";
import { homedir } from "node:os";
import { join } from "node:path";
import { query } from "@anthropic-ai/claude-agent-sdk";

// ── Types ────────────────────────────────────────────────────────────────────

interface MCPServerConfig {
  command?: string;
  args?: string[];
  env?: Record<string, string>;
  type?: string;
  url?: string;
  id?: string;
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
  continue?: boolean;
  agents?: Record<string, unknown>;
  no_user_settings?: boolean;
  disabled_tools?: string[];
}

// ── Cloud MCP (claude.ai) ───────────────────────────────────────────────────

interface CloudMCPServer {
  type: "claudeai-proxy";
  url: string;
  id: string;
}

interface CloudMCPCache {
  servers: Record<string, CloudMCPServer>;
  expiresAt: number;
}

const CLOUD_MCP_CACHE_TTL_MS = 5 * 60 * 1000; // 5 minutes
const CLOUD_MCP_API_TIMEOUT_MS = 5000;
const CLOUD_MCP_PROXY_BASE = "https://mcp-proxy.anthropic.com/v1/mcp";
const CLOUD_MCP_API_URL = "https://api.anthropic.com/v1/mcp_servers?limit=1000";
const CLOUD_MCP_BETA_HEADER = "mcp-servers-2025-12-04";

let cloudMcpCache: CloudMCPCache | null = null;

async function loadCloudMCPs(): Promise<Record<string, CloudMCPServer>> {
  // Return cached if still valid
  if (cloudMcpCache && Date.now() < cloudMcpCache.expiresAt) {
    return cloudMcpCache.servers;
  }

  try {
    const credsPath = join(homedir(), ".claude", ".credentials.json");
    const raw = await readFile(credsPath, "utf8");
    const creds = JSON.parse(raw);

    const oauth = creds?.claudeAiOauth;
    if (!oauth?.accessToken) {
      log("cloud-mcp: no OAuth token found");
      return {};
    }

    if (!oauth.scopes?.includes("user:mcp_servers")) {
      log("cloud-mcp: missing user:mcp_servers scope");
      return {};
    }

    // Check if token is expired
    if (oauth.expiresAt && Date.now() > oauth.expiresAt) {
      log("cloud-mcp: OAuth token expired");
      return {};
    }

    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), CLOUD_MCP_API_TIMEOUT_MS);

    const resp = await fetch(CLOUD_MCP_API_URL, {
      headers: {
        Authorization: `Bearer ${oauth.accessToken}`,
        "Content-Type": "application/json",
        "anthropic-beta": CLOUD_MCP_BETA_HEADER,
        "anthropic-version": "2023-06-01",
      },
      signal: controller.signal,
    });

    clearTimeout(timeout);

    if (!resp.ok) {
      log(`cloud-mcp: API returned ${resp.status}`);
      return {};
    }

    const body = (await resp.json()) as { data: Array<{ id: string; display_name: string }> };
    const servers: Record<string, CloudMCPServer> = {};

    for (const srv of body.data) {
      const name = `claude.ai ${srv.display_name}`;
      servers[name] = {
        type: "claudeai-proxy",
        url: `${CLOUD_MCP_PROXY_BASE}/${srv.id}`,
        id: srv.id,
      };
    }

    log(`cloud-mcp: fetched ${Object.keys(servers).length} servers`);
    cloudMcpCache = { servers, expiresAt: Date.now() + CLOUD_MCP_CACHE_TTL_MS };
    return servers;
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    log(`cloud-mcp: failed to load — ${msg}`);
    return {};
  }
}

interface Request {
  command: string;
  prompt: string;
  request_id?: string;
  options?: RequestOptions;
}

interface OutEvent {
  event: string;
  request_id?: string;
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

async function buildSDKOptions(opts: RequestOptions | undefined) {
  if (!opts) return {};

  const sdkOpts: Record<string, unknown> = {};

  if (opts.model) sdkOpts.model = opts.model;
  if (opts.cwd) sdkOpts.cwd = opts.cwd;
  if (opts.system_prompt) sdkOpts.systemPrompt = opts.system_prompt;
  if (opts.resume) sdkOpts.resume = opts.resume;
  if (opts.continue) sdkOpts.continue = opts.continue;
  if (opts.agents && Object.keys(opts.agents).length > 0) {
    sdkOpts.agents = opts.agents;
  }
  if (opts.max_turns) sdkOpts.maxTurns = opts.max_turns;
  if (opts.permission_mode) {
    sdkOpts.permissionMode = opts.permission_mode;
    if (opts.permission_mode === "bypassPermissions") {
      sdkOpts.allowDangerouslySkipPermissions = true;
    }
  }
  if (opts.allowed_tools) sdkOpts.allowedTools = opts.allowed_tools;

  // Load user settings unless explicitly disabled (e.g. cron jobs)
  if (opts.no_user_settings) {
    sdkOpts.settingSources = [];
  } else {
    sdkOpts.settingSources = ["user", "project", "local"];
  }

  // Merge agent MCP servers with cloud MCPs from claude.ai
  const cloudServers = opts.no_user_settings ? {} : await loadCloudMCPs();
  const agentServers = opts.mcp_servers ?? {};
  const merged = { ...cloudServers, ...agentServers }; // agent overrides cloud
  if (Object.keys(merged).length > 0) {
    sdkOpts.mcpServers = merged;
  }

  // Disable specific tools if requested
  if (opts.disabled_tools && opts.disabled_tools.length > 0) {
    sdkOpts.disallowedTools = opts.disabled_tools;
  }

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
  const reqId = req.request_id || "";
  const emitReq = (obj: OutEvent) => emit({ ...obj, request_id: reqId });

  const sdkOptions = await buildSDKOptions(req.options);

  log(`query start — rid=${reqId} model=${sdkOptions.model ?? "default"} prompt="${req.prompt.slice(0, 80)}..."`);

  const timeoutMs = 10 * 60 * 1000;
  const timeout = setTimeout(() => {
    log(`query timeout — rid=${reqId} no result after 10 minutes`);
    emitReq({ event: "error", message: "query timeout: no result after 10 minutes" });
    // Don't exit — just emit error and let the process continue.
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
          emitReq({
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
              emitReq({ event: "assistant", text });
            }
            for (const block of inner.content as Record<string, unknown>[]) {
              if (block.type === "tool_use") {
                emitReq({
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
          emitReq({
            event: "tool_result",
            content: msg.summary as string,
          });
          break;
        }

        // ── Result (success or error) ────────────────────────────────
        case "result": {
          const subtype = msg.subtype as string | undefined;
          if (subtype === "success") {
            emitReq({
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
            emitReq({
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
    log(`query error: rid=${reqId} ${errMsg}`);
    emitReq({ event: "error", message: errMsg });
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

// ── Main loop ────────────────────────────────────────────────────────────────

function main(): void {
  log("bridge started — waiting for commands on stdin");

  const rl = createInterface({
    input: process.stdin,
    terminal: false,
  });

  // Process requests concurrently — each query runs independently
  rl.on("line", (line: string) => {
    const trimmed = line.trim();
    if (!trimmed) return;

    // Fire and forget — each request runs in its own async context
    handleRequest(trimmed).catch((err: unknown) => {
      const errMsg = err instanceof Error ? err.message : String(err);
      log(`unhandled error in request processing: ${errMsg}`);
      emit({ event: "error", message: `internal bridge error: ${errMsg}` });
    });
  });

  rl.on("close", () => {
    log("stdin closed — shutting down");
    process.exit(0);
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
