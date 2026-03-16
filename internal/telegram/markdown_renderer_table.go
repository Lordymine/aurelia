package telegram

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *telegramRenderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre>")
	} else {
		_, _ = w.WriteString("</pre>\n")
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_ = w.WriteByte('\n')
		columnCount := 0
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			columnCount++
		}
		for i := 0; i < columnCount; i++ {
			if i > 0 {
				_, _ = w.WriteString("-+-")
			}
			_, _ = w.WriteString("---")
		}
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_ = w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *telegramRenderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering && node.PreviousSibling() != nil {
		_, _ = w.WriteString(" | ")
	}
	return ast.WalkContinue, nil
}
