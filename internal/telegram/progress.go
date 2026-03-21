package telegram

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"gopkg.in/telebot.v3"
)

type progressReporter struct {
	bot   *telebot.Bot
	chat  *telebot.Chat
	msg   *telebot.Message
	tools []string
	mu    sync.Mutex
}

func newProgressReporter(bot *telebot.Bot, chat *telebot.Chat) *progressReporter {
	return &progressReporter{bot: bot, chat: chat}
}

func (p *progressReporter) ReportTool(toolName string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	label := toolDisplayName(toolName)
	p.tools = append(p.tools, label)

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
