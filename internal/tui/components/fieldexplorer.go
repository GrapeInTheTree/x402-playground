package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-demo/internal/tui"
)

// Field represents a data field with its value and description.
type Field struct {
	Name        string
	Value       string
	Description string
}

// FieldExplorer renders a list of fields; selected field shows its description.
type FieldExplorer struct {
	Fields []Field
	Cursor int
	Width  int
}

func NewFieldExplorer(fields []Field) FieldExplorer {
	return FieldExplorer{Fields: fields, Width: 60}
}

func (f *FieldExplorer) Up() {
	if f.Cursor > 0 {
		f.Cursor--
	}
}

func (f *FieldExplorer) Down() {
	if f.Cursor < len(f.Fields)-1 {
		f.Cursor++
	}
}

func (f FieldExplorer) View() string {
	var b strings.Builder

	nameStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary).Width(20)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))
	selectedName := lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true).Width(20)

	for i, field := range f.Fields {
		if i == f.Cursor {
			fmt.Fprintf(&b, "  > %s %s\n",
				selectedName.Render(field.Name+":"),
				lipgloss.NewStyle().Foreground(tui.ColorAccent).Render(field.Value))
		} else {
			fmt.Fprintf(&b, "    %s %s\n",
				nameStyle.Render(field.Name+":"),
				valueStyle.Render(field.Value))
		}
	}

	// Show description for selected field
	if f.Cursor >= 0 && f.Cursor < len(f.Fields) {
		desc := f.Fields[f.Cursor].Description
		if desc != "" {
			b.WriteString("\n")
			descStyle := lipgloss.NewStyle().
				Foreground(tui.ColorMuted).
				PaddingLeft(4).
				Width(f.Width - 8)
			b.WriteString(descStyle.Render("  " + desc))
			b.WriteString("\n")
		}
	}

	return b.String()
}
