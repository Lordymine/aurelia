package telegram

import (
	"fmt"
	"strings"
)

func EscapeHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return replacer.Replace(text)
}

func Bold(text string) string {
	return fmt.Sprintf("<b>%s</b>", EscapeHTML(text))
}

func TitleWithEmoji(emoji, title string) string {
	return emoji + " " + Bold(title)
}

func ErrorMessage(title, body string) string {
	return TitleWithEmoji("⚠️", title) + "\n\n" + EscapeHTML(body)
}
