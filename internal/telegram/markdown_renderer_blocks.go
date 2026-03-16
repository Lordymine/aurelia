package telegram

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *telegramRenderer) renderDocument(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.listDepth = 0
		r.orderedCount = 0
		r.inList = false
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderHeading(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<b>")
	} else {
		_, _ = w.WriteString("</b>\n")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderBlockquote(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<blockquote>")
	} else {
		_, _ = w.WriteString("</blockquote>")
		if n.NextSibling() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	list := node.(*ast.List)
	if entering {
		r.listDepth++
		r.inList = true
		if list.IsOrdered() {
			r.orderedCount = list.Start
		} else {
			r.orderedCount = 0
		}
	} else {
		r.listDepth--
		if r.listDepth == 0 {
			r.inList = false
		}
		if parent := node.Parent(); parent != nil {
			if parentList, ok := parent.(*ast.List); ok && parentList.IsOrdered() {
				r.orderedCount = 0
			}
		}
		if node.NextSibling() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		indent := strings.Repeat("  ", r.listDepth-1)
		_, _ = w.WriteString(indent)
		if r.orderedCount > 0 {
			_, _ = fmt.Fprintf(w, "%d. ", r.orderedCount)
			r.orderedCount++
		} else {
			_, _ = w.WriteString("\u2022 ")
		}
	} else {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, isListItem := n.Parent().(*ast.ListItem); isListItem {
			return ast.WalkContinue, nil
		}
		if n.NextSibling() != nil {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderTextBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering && n.NextSibling() != nil && n.FirstChild() != nil {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderThematicBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n")
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	lineCount := n.Lines().Len()
	for i := 0; i < lineCount; i++ {
		line := n.Lines().At(i)
		_, _ = w.Write(util.EscapeHTML(line.Value(source)))
	}
}

func (r *telegramRenderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre><code>")
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	block := node.(*ast.FencedCodeBlock)
	if entering {
		_, _ = w.WriteString("<pre><code")
		language := block.Language(source)
		if len(language) > 0 {
			_, _ = w.WriteString(` class="language-`)
			_, _ = w.Write(util.EscapeHTML(language))
			_, _ = w.WriteString(`"`)
		}
		_, _ = w.WriteString(">")
		r.writeLines(w, source, block)
	} else {
		_, _ = w.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	block := node.(*ast.HTMLBlock)
	lineCount := block.Lines().Len()
	for i := 0; i < lineCount; i++ {
		line := block.Lines().At(i)
		_, _ = w.Write(util.EscapeHTML(line.Value(source)))
	}
	if block.HasClosure() {
		_, _ = w.Write(util.EscapeHTML(block.ClosureLine.Value(source)))
	}
	return ast.WalkContinue, nil
}
