package telegram

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func (r *telegramRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = w.WriteString("<code>")
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		text, ok := child.(*ast.Text)
		if !ok {
			continue
		}
		value := text.Segment.Value(source)
		if bytes.HasSuffix(value, []byte("\n")) {
			_, _ = w.Write(util.EscapeHTML(value[:len(value)-1]))
			_, _ = w.WriteString(" ")
			continue
		}
		_, _ = w.Write(util.EscapeHTML(value))
	}
	_, _ = w.WriteString("</code>")
	return ast.WalkSkipChildren, nil
}

func (r *telegramRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	link := node.(*ast.AutoLink)
	_, _ = w.WriteString(`<a href="`)
	url := link.URL(source)
	label := link.Label(source)
	if link.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	_, _ = w.WriteString(`">`)
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	emphasis := node.(*ast.Emphasis)
	tag := "i"
	if emphasis.Level == 2 {
		tag = "b"
	}
	if entering {
		_, _ = w.WriteString("<" + tag + ">")
	} else {
		_, _ = w.WriteString("</" + tag + ">")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	image := node.(*ast.Image)
	altText := r.renderTextsToBytes(source, image)
	if len(altText) > 0 {
		_, _ = w.WriteString("[image: ")
		_, _ = w.Write(util.EscapeHTML(altText))
		_, _ = w.WriteString("]")
	} else {
		_, _ = w.WriteString("[image]")
	}
	return ast.WalkSkipChildren, nil
}

func (r *telegramRenderer) renderTextsToBytes(source []byte, n ast.Node) []byte {
	var buf bytes.Buffer
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if s, ok := child.(*ast.String); ok {
			buf.Write(util.EscapeHTML(s.Value))
			continue
		}
		if text, ok := child.(*ast.Text); ok {
			buf.Write(util.EscapeHTML(text.Segment.Value(source)))
			continue
		}
		buf.Write(r.renderTextsToBytes(source, child))
	}
	return buf.Bytes()
}

func (r *telegramRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	link := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString(`<a href="`)
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(link.Destination, true)))
		_, _ = w.WriteString(`">`)
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	raw := node.(*ast.RawHTML)
	for i := 0; i < raw.Segments.Len(); i++ {
		segment := raw.Segments.At(i)
		_, _ = w.Write(util.EscapeHTML(segment.Value(source)))
	}
	return ast.WalkSkipChildren, nil
}

func (r *telegramRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	text := node.(*ast.Text)
	value := text.Segment.Value(source)
	_, _ = w.Write(util.EscapeHTML(value))
	if text.HardLineBreak() || text.SoftLineBreak() {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	value := node.(*ast.String)
	if value.IsCode() {
		_, _ = w.Write(value.Value)
	} else {
		_, _ = w.Write(util.EscapeHTML(value.Value))
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<s>")
	} else {
		_, _ = w.WriteString("</s>")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	checkBox := node.(*extast.TaskCheckBox)
	if checkBox.IsChecked {
		_, _ = w.WriteString("\u2705 ")
	} else {
		_, _ = w.WriteString("\u2610 ")
	}
	return ast.WalkContinue, nil
}
