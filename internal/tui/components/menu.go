package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// MenuItem represents a selectable menu entry.
type MenuItem struct {
	Title       string
	Description string
	Icon        string
}

// Menu renders a vertical list of items with a cursor.
type Menu struct {
	Items    []MenuItem
	Cursor   int
	Width    int
}

// NewMenu creates a new menu with the given items.
func NewMenu(items []MenuItem) Menu {
	return Menu{Items: items, Width: 60}
}

// Up moves the cursor up one item.
func (m *Menu) Up() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

// Down moves the cursor down one item.
func (m *Menu) Down() {
	if m.Cursor < len(m.Items)-1 {
		m.Cursor++
	}
}

// Selected returns the index of the currently selected item.
func (m *Menu) Selected() int {
	return m.Cursor
}

// View renders the menu with cursor highlighting.
func (m Menu) View() string {
	var b strings.Builder

	for i, item := range m.Items {
		icon := item.Icon
		if icon == "" {
			icon = " "
		}

		title := item.Title
		desc := item.Description

		if i == m.Cursor {
			cursor := lipgloss.NewStyle().
				Foreground(tui.ColorPrimary).
				Bold(true).
				Render(">")
			title = lipgloss.NewStyle().
				Foreground(tui.ColorPrimary).
				Bold(true).
				Render(title)
			desc = lipgloss.NewStyle().
				Foreground(tui.ColorSecondary).
				Render(desc)
			fmt.Fprintf(&b, "%s %s %s\n", cursor, icon, title)
		} else {
			title = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D5DB")).
				Render(title)
			desc = tui.MutedStyle.Render(desc)
			fmt.Fprintf(&b, "  %s %s\n", icon, title)
		}

		if desc != "" {
			fmt.Fprintf(&b, "    %s\n", desc)
		}
		b.WriteString("\n")
	}

	return b.String()
}
