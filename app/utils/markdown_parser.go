package utils

import (
	"bytes"
	"gopkg.in/russross/blackfriday.v2"
	"io"
)

type TextureRender struct {
	Flags blackfriday.HTMLFlags
}

func (r *TextureRender) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	buf := bytes.Buffer{}
	if node.Next != nil && node.Next.Literal != nil {
		buf.WriteString("\n")
	}
	if node.Literal != nil {
		buf.Write(node.Literal)
	}
	if buf.Len() > 0 {
		_, _ = w.Write(buf.Bytes())
	}
	return blackfriday.GoToNext
}

func (r *TextureRender) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

func (r *TextureRender) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}

func MarkdownToPlainText(md string, sliceSize int) []string {
	buf := blackfriday.Run([]byte(md), blackfriday.WithRenderer(&TextureRender{}))
	if len(buf) == 0 {
		return []string{}
	}
	bufStr := string(buf)
	length := len(bufStr)
	if length < sliceSize {
		return []string{bufStr}
	}
	var list []string
	for length > 0 {
		if length < sliceSize {
			sliceSize = length
		}
		list = append(list, bufStr[:sliceSize])
		bufStr = bufStr[sliceSize:]
		length = len(bufStr)
	}
	return list
}
