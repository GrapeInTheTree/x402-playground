package components

import "github.com/charmbracelet/lipgloss"

// TriPanel renders three panels side by side for the Practice view.
type TriPanel struct {
	Left   Panel
	Center Panel
	Right  Panel
	Width  int
}

const minTriPanelColWidth = 20

// View renders the three panels side by side.
func (t TriPanel) View() string {
	colWidth := max((t.Width-6)/3, minTriPanelColWidth)

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
