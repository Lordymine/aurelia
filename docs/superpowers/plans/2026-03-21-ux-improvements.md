# UX Improvements — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Aurelia usable in real conversations: session continuity, image support, tool progress feedback, and smart agent routing.

**Architecture:** Session resume via SDK session IDs stored per Telegram chat. Images downloaded and passed as file paths to Bridge. Tool progress sent as editable Telegram messages. Agent routing via a lightweight Bridge classification call.

**Tech Stack:** Go, TypeScript, Claude Agent SDK, Telegram Bot API

---

## File Structure

### Files to modify

```
internal/bridge/protocol.go            — add Resume field (already exists, just wire it)
internal/bridge/events.go              — add ToolID field for tool_use tracking
internal/telegram/bot.go               — add sessions map for session resume
internal/telegram/input_pipeline.go    — wire session resume, progress messages, image handling
internal/telegram/input.go             — download photos, pass file paths to processInput
internal/agents/registry.go            — add ClassifyRoute method using Bridge
bridge/index.ts                        — no changes needed (resume already supported)
```

### Files to create

```
internal/telegram/sessions.go          — session store (chat_id → session_id, thread-safe)
internal/telegram/progress.go          — tool progress message formatting + editing
```

---

## Task 1: Session Resume (conversation continuity)

**Goal:** Each Telegram chat maintains a Claude Code session. Follow-up messages continue the same session instead of starting fresh.

**Files:**
- Create: `internal/telegram/sessions.go`
- Modify: `internal/telegram/input_pipeline.go`

- [ ] **Step 1: Create sessions.go — thread-safe session store**

```go
package telegram

import "sync"

// sessionStore maps Telegram chat IDs to Claude Code session IDs.
type sessionStore struct {
    mu       sync.RWMutex
    sessions map[int64]string // chat_id → session_id
}

func newSessionStore() *sessionStore {
    return &sessionStore{sessions: make(map[int64]string)}
}

func (s *sessionStore) Get(chatID int64) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.sessions[chatID]
}

func (s *sessionStore) Set(chatID int64, sessionID string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.sessions[chatID] = sessionID
}

func (s *sessionStore) Clear(chatID int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.sessions, chatID)
}
```

- [ ] **Step 2: Add sessions field to BotController**

In `internal/telegram/bot.go`, add `sessions *sessionStore` to the struct and initialize in constructor:

```go
sessions: newSessionStore(),
```

- [ ] **Step 3: Wire session resume in processInput**

In `internal/telegram/input_pipeline.go`, after building the bridge request:

```go
// Resume previous session if exists
if sessionID := bc.sessions.Get(c.Chat().ID); sessionID != "" {
    req.Options.Resume = sessionID
}
```

- [ ] **Step 4: Capture session_id from result event**

In `processBridgeEvents`, on the `"system"` event case, capture session_id. On `"result"` event, save session_id:

```go
case "system":
    if ev.SessionID != "" {
        bc.sessions.Set(c.Chat().ID, ev.SessionID)
    }
```

Also add `"system"` case to the switch (currently falls into `default`).

- [ ] **Step 5: Run tests**

```bash
go build ./...
go test ./internal/telegram/ -v
```

- [ ] **Step 6: Commit**

```bash
git commit -m "feat(telegram): add session resume for conversation continuity"
```

---

## Task 2: Tool Progress in Telegram

**Goal:** When the Bridge uses tools, show progress in Telegram (e.g., "🔧 Reading src/main.go...") as an editable message that updates.

**Files:**
- Create: `internal/telegram/progress.go`
- Modify: `internal/telegram/input_pipeline.go`

- [ ] **Step 1: Create progress.go — progress message manager**

```go
package telegram

import (
    "fmt"
    "log"
    "strings"
    "sync"

    "gopkg.in/telebot.v3"
)

// progressReporter sends and edits a single progress message in a Telegram chat.
type progressReporter struct {
    bot     *telebot.Bot
    chat    *telebot.Chat
    msg     *telebot.Message
    tools   []string
    mu      sync.Mutex
}

func newProgressReporter(bot *telebot.Bot, chat *telebot.Chat) *progressReporter {
    return &progressReporter{bot: bot, chat: chat}
}

// ReportTool adds a tool use to the progress message and sends/edits it.
func (p *progressReporter) ReportTool(toolName string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    label := toolDisplayName(toolName)
    p.tools = append(p.tools, label)

    // Keep last 5 tools to avoid message bloat
    display := p.tools
    if len(display) > 5 {
        display = display[len(display)-5:]
    }

    text := strings.Join(display, "\n")

    if p.msg == nil {
        sent, err := p.bot.Send(p.chat, text)
        if err != nil {
            log.Printf("Progress send error: %v", err)
            return
        }
        p.msg = sent
    } else {
        _, err := p.bot.Edit(p.msg, text)
        if err != nil {
            log.Printf("Progress edit error: %v", err)
        }
    }
}

// Delete removes the progress message from chat.
func (p *progressReporter) Delete() {
    p.mu.Lock()
    defer p.mu.Unlock()
    if p.msg != nil {
        _ = p.bot.Delete(p.msg)
        p.msg = nil
    }
}

func toolDisplayName(name string) string {
    switch name {
    case "Read":
        return "📖 Reading file..."
    case "Write":
        return "✍️ Writing file..."
    case "Edit":
        return "✏️ Editing file..."
    case "Bash":
        return "⚡ Running command..."
    case "Glob":
        return "🔍 Searching files..."
    case "Grep":
        return "🔎 Searching content..."
    case "WebSearch":
        return "🌐 Searching web..."
    case "WebFetch":
        return "🌐 Fetching page..."
    default:
        return fmt.Sprintf("🔧 %s...", name)
    }
}
```

- [ ] **Step 2: Wire progress in processBridgeEvents**

In `input_pipeline.go`, create a progress reporter and use it:

```go
func (bc *BotController) processBridgeEvents(c telebot.Context, ch <-chan bridge.Event, userText string) error {
    var assistantText strings.Builder
    progress := newProgressReporter(bc.bot, c.Chat())
    defer progress.Delete()

    for ev := range ch {
        switch ev.Type {
        case "system":
            if ev.SessionID != "" {
                bc.sessions.Set(c.Chat().ID, ev.SessionID)
            }

        case "tool_use":
            toolName := ev.Name
            if toolName == "" {
                toolName = "tool"
            }
            progress.ReportTool(toolName)

        // ... rest stays the same
```

Delete the progress message before sending the final result (the `defer progress.Delete()` handles it).

- [ ] **Step 3: Run tests**

```bash
go build ./...
go test ./internal/telegram/ -v
```

- [ ] **Step 4: Commit**

```bash
git commit -m "feat(telegram): show tool progress during bridge execution"
```

---

## Task 3: Image/Photo Support

**Goal:** When user sends photos, download them to a temp file and tell Claude Code to analyze them via the file path.

**Files:**
- Modify: `internal/telegram/input.go`

- [ ] **Step 1: Implement photo download in handlePhotoGroup**

Replace the TODO at line 70-71. Download each photo, save to temp file, and include the file paths in the prompt:

```go
func (bc *BotController) handlePhotoGroup(c telebot.Context, photos []albumPhoto, caption string) error {
    if len(photos) == 0 {
        return nil
    }

    stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.UploadingPhoto, 4*time.Second)
    defer stopTyping()

    text := strings.TrimSpace(caption)

    // Download photos to temp files
    var photoPaths []string
    for _, p := range photos {
        filePath, err := bc.downloadTelegramFile(&p.photo.File, fmt.Sprintf("photo_%d.jpg", p.messageID))
        if err != nil {
            log.Printf("Failed to download photo: %v", err)
            continue
        }
        photoPaths = append(photoPaths, filePath)
    }

    if len(photoPaths) == 0 {
        return bc.processInput(c, text, nil, false)
    }

    // Build prompt with file paths for Claude Code to read
    if text == "" {
        text = "Analise as imagens a seguir."
    }
    for _, path := range photoPaths {
        text += fmt.Sprintf("\n\nImagem: %s", path)
    }

    // Note: don't defer remove — Claude Code needs the files during execution.
    // They'll be cleaned up by the OS temp directory eventually.
    return bc.processInput(c, text, nil, false)
}
```

- [ ] **Step 2: Fix handleImageDocument similarly**

Replace the TODO at line 112-113:

```go
func (bc *BotController) handleImageDocument(c telebot.Context, doc *telebot.Document) error {
    stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.UploadingPhoto, 4*time.Second)
    defer stopTyping()

    text := strings.TrimSpace(c.Message().Caption)
    if text == "" {
        text = "Analise esta imagem."
    }

    filePath, err := bc.downloadTelegramFile(&doc.File, doc.FileID+"_"+doc.FileName)
    if err != nil {
        log.Printf("Failed to download image document: %v", err)
        return bc.processInput(c, text, nil, false)
    }

    text += fmt.Sprintf("\n\nImagem: %s", filePath)
    return bc.processInput(c, text, nil, false)
}
```

- [ ] **Step 3: Verify downloadTelegramFile exists and works**

Read the existing `downloadTelegramFile` method — it should already handle downloading files from Telegram to a temp directory.

- [ ] **Step 4: Run tests**

```bash
go build ./...
go test ./internal/telegram/ -v
```

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(telegram): download and pass photos to Bridge for analysis"
```

---

## Task 4: Smart Agent Routing

**Goal:** Instead of requiring `@agentname`, use a lightweight Bridge call to classify which agent should handle the message. Falls back to no agent (general) if no match.

**Files:**
- Modify: `internal/agents/registry.go`
- Modify: `internal/telegram/input_pipeline.go`

- [ ] **Step 1: Add ClassifyRoute to Registry**

In `internal/agents/registry.go`, add a method that builds a classification prompt:

```go
// ClassifyPrompt returns a prompt that asks an LLM to classify which agent
// should handle the given message. Returns empty string if no agents exist.
func (r *Registry) ClassifyPrompt(message string) string {
    if len(r.agents) == 0 {
        return ""
    }

    var sb strings.Builder
    sb.WriteString("Given these available agents:\n\n")
    for _, a := range r.Agents() {
        sb.WriteString(fmt.Sprintf("- %s: %s\n", a.Name, a.Description))
    }
    sb.WriteString(fmt.Sprintf("\nAnd this user message: %q\n\n", message))
    sb.WriteString("Reply with ONLY the agent name that best matches, or 'none' if no agent is a good match. Reply with a single word, no explanation.")
    return sb.String()
}
```

- [ ] **Step 2: Add routing logic to input_pipeline.go**

In `processInput`, replace the simple `Route(text)` with smart routing:

```go
// 1. Route to agent
agent := bc.agents.Route(text) // try @name first

// If no @name match and agents exist, try LLM classification
if agent == nil && bc.agents != nil && len(bc.agents.Agents()) > 0 {
    classifyPrompt := bc.agents.ClassifyPrompt(text)
    if classifyPrompt != "" {
        ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
        result, err := bc.bridge.ExecuteSync(ctx, bridge.Request{
            Command: "query",
            Prompt:  classifyPrompt,
            Options: bridge.RequestOptions{
                Model:          bc.config.DefaultModel,
                SystemPrompt:   "You are a message classifier. Reply with only the agent name or 'none'.",
                MaxTurns:       1,
                PermissionMode: "bypassPermissions",
            },
        })
        cancel()
        if err == nil && result.Type == "result" {
            name := strings.TrimSpace(strings.ToLower(result.Content))
            if name != "none" && name != "" {
                agent = bc.agents.Get(name)
            }
        }
    }
}
```

- [ ] **Step 3: Write test for ClassifyPrompt**

```go
func TestRegistry_ClassifyPrompt(t *testing.T) {
    reg := setupTestRegistry(t) // has "prospector" agent
    prompt := reg.ClassifyPrompt("busque leads em SP")
    if prompt == "" {
        t.Fatal("expected non-empty classify prompt")
    }
    if !strings.Contains(prompt, "prospector") {
        t.Fatal("expected prompt to contain agent name")
    }
}

func TestRegistry_ClassifyPrompt_Empty(t *testing.T) {
    reg, _ := Load(t.TempDir())
    prompt := reg.ClassifyPrompt("hello")
    if prompt != "" {
        t.Fatal("expected empty prompt for empty registry")
    }
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/agents/ -v
go test ./internal/telegram/ -v
go build ./...
```

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(agents): add LLM-based smart routing for agent classification"
```

---

## Task 5: Telegram /cwd Command

**Goal:** Let user set working directory per chat via `/cwd path` command, so they don't need an agent just to specify where to work.

**Files:**
- Modify: `internal/telegram/sessions.go` — add cwd storage
- Modify: `internal/telegram/bot_middleware.go` — register /cwd handler
- Modify: `internal/telegram/input_pipeline.go` — use stored cwd

- [ ] **Step 1: Add cwd storage to sessionStore**

```go
type sessionStore struct {
    mu       sync.RWMutex
    sessions map[int64]string
    cwds     map[int64]string // chat_id → working directory
}

func (s *sessionStore) GetCwd(chatID int64) string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.cwds[chatID]
}

func (s *sessionStore) SetCwd(chatID int64, cwd string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.cwds[chatID] = cwd
}
```

Initialize `cwds: make(map[int64]string)` in constructor.

- [ ] **Step 2: Add /cwd command handler**

In `bot_middleware.go`, add to `registerContentRoutes`:

```go
bc.bot.Handle("/cwd", bc.handleCwdCommand)
bc.bot.Handle("/reset", bc.handleResetCommand)
```

Implement handlers:

```go
func (bc *BotController) handleCwdCommand(c telebot.Context) error {
    args := strings.TrimSpace(c.Message().Payload)
    if args == "" {
        cwd := bc.sessions.GetCwd(c.Chat().ID)
        if cwd == "" {
            return SendText(bc.bot, c.Chat(), "Nenhum diretório configurado. Use: /cwd C:\\path\\to\\project")
        }
        return SendText(bc.bot, c.Chat(), fmt.Sprintf("Diretório atual: %s", cwd))
    }
    bc.sessions.SetCwd(c.Chat().ID, args)
    return SendText(bc.bot, c.Chat(), fmt.Sprintf("Diretório configurado: %s", args))
}

func (bc *BotController) handleResetCommand(c telebot.Context) error {
    bc.sessions.Clear(c.Chat().ID)
    return SendText(bc.bot, c.Chat(), "Sessão resetada. Próxima mensagem inicia conversa nova.")
}
```

- [ ] **Step 3: Wire cwd in processInput**

After building the request and applying agent overrides:

```go
// Apply chat-level cwd if no agent overrides it
if req.Options.Cwd == "" {
    if chatCwd := bc.sessions.GetCwd(c.Chat().ID); chatCwd != "" {
        req.Options.Cwd = chatCwd
    }
}
```

- [ ] **Step 4: Run tests**

```bash
go build ./...
go test ./internal/telegram/ -v
```

- [ ] **Step 5: Commit**

```bash
git commit -m "feat(telegram): add /cwd and /reset commands for session control"
```

---

## Task 6: Final Verification

- [ ] **Step 1: Full build and test**

```bash
go build ./...
go vet ./...
go test ./... -v -short
cd bridge && npx tsc --noEmit
```

- [ ] **Step 2: Verify no dead code introduced**

```bash
grep -rn "TODO" --include="*.go" internal/telegram/input.go
```

Expect: no remaining TODOs for photo handling.

- [ ] **Step 3: Commit any remaining changes**

```bash
git add -A && git status
```
