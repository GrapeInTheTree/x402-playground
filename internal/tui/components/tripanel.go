package components

import "github.com/charmbracelet/lipgloss"

// TriPanel renders three panels side by side for the Practice view.
type TriPanel struct {
	Left   Panel
	Center Panel
	Right  Panel
	Width  int
}

func (t TriPanel) View() string {
	colWidth := (t.Width - 6) / 3
	if colWidth < 20 {
		colWidth = 20
	}

	t.Left.Width = colWidth
	t.Center.Width = colWidth
	t.Right.Width = colWidth

	return lipgloss.JoinHorizontal(lipgloss.Top,
		t.Left.View(),
		" ",
		t.Center.View(),
		" ",
		t.Right.View(),
	)
}
