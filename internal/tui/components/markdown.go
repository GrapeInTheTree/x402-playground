package components

import (
	"github.com/charmbracelet/glamour"
)

// RenderMarkdown renders markdown text with glamour for terminal display.
func RenderMarkdown(content string, width int) string {
	if width < 20 {
		width = 80
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4),
	)
	if err != nil {
		return content
	}
	out, err := r.Render(content)
	if err != nil {
		return content
	}
	return out
}
