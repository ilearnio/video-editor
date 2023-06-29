package helpers

import (
	"strings"
)

func EscapeHTML(html string) string {
	html = strings.ReplaceAll(html, "&", "&amp;")
	html = strings.ReplaceAll(html, "<", "&lt;")
	html = strings.ReplaceAll(html, ">", "&gt;")
	html = strings.ReplaceAll(html, "\"", "&quot;")
	return html
}

func EscapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return s
}
