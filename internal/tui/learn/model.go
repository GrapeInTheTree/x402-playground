package learn

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/quiz"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

type page int

const (
	pageQuestionList page = iota
	pageQuestion
	pageResult
)

type editorFinishedMsg struct{ err error }

type testResultMsg struct{ result *quiz.Result }

// Model is the quiz-based learning page model.
type Model struct {
	questions []quiz.Question
	cursor    int
	current   page
	goRunner  *quiz.Runner
	solRunner *quiz.Runner
	results   map[int]*quiz.Result
	score     quiz.Score
	width     int
	height    int
}

// New creates a new quiz learning model.
func New(width, height int) *Model {
	questions := quiz.AllQuestions()
	return &Model{
		questions: questions,
		results:   make(map[int]*quiz.Result),
		score:     quiz.Score{Questions: len(questions)},
		width:     width,
		height:    height,
	}
}

// Init implements the SubModel interface.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles navigation, editor launch, and test results.
func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case editorFinishedMsg:
		if msg.err != nil {
			return m, nil
		}
		return m, m.runTests()

	case testResultMsg:
		m.results[m.cursor] = msg.result
		m.score.Answered++
		if msg.result.Passed == msg.result.Total && msg.result.Total > 0 {
			m.score.Correct++
		}
		m.current = pageResult
		return m, nil

	case tea.KeyMsg:
		switch m.current {
		case pageQuestionList:
			return m.updateList(msg)
		case pageQuestion:
			return m.updateQuestion(msg)
		case pageResult:
			return m.updateResult(msg)
		}
	}
	return m, nil
}

func (m *Model) updateList(msg tea.KeyMsg) (tui.SubModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.questions)-1 {
			m.cursor++
		}
	case "enter":
		m.current = pageQuestion
	case "esc", "q":
		return m, func() tea.Msg { return tui.BackMsg{} }
	}
	return m, nil
}

func (m *Model) updateQuestion(msg tea.KeyMsg) (tui.SubModel, tea.Cmd) {
	switch msg.String() {
	case "enter", "e":
		return m, m.openEditor()
	case "esc":
		m.current = pageQuestionList
	}
	return m, nil
}

func (m *Model) updateResult(msg tea.KeyMsg) (tui.SubModel, tea.Cmd) {
	switch msg.String() {
	case "enter", "r":
		return m, m.openEditor()
	case "n":
		if m.cursor < len(m.questions)-1 {
			m.cursor++
			m.current = pageQuestion
		} else {
			m.current = pageQuestionList
		}
	case "esc":
		m.current = pageQuestionList
	}
	return m, nil
}

func (m *Model) getRunner(lang quiz.Lang) *quiz.Runner {
	switch lang {
	case quiz.LangSolidity:
		return m.solRunner
	default:
		return m.goRunner
	}
}

func (m *Model) setRunner(lang quiz.Lang, r *quiz.Runner) {
	switch lang {
	case quiz.LangSolidity:
		m.solRunner = r
	default:
		m.goRunner = r
	}
}

func (m *Model) openEditor() tea.Cmd {
	q := m.questions[m.cursor]
	lang := q.Language
	if lang == "" {
		lang = quiz.LangGo
	}

	runner := m.getRunner(lang)
	if runner == nil {
		var err error
		switch lang {
		case quiz.LangSolidity:
			runner, err = quiz.NewSolidityRunner()
		default:
			runner, err = quiz.NewRunner()
		}
		if err != nil {
			return nil
		}
		m.setRunner(lang, runner)
	}

	if err := os.WriteFile(runner.TemplatePath(), []byte(q.Template), 0644); err != nil {
		return nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		for _, candidate := range []string{"nvim", "vim", "nano"} {
			if _, err := exec.LookPath(candidate); err == nil {
				editor = candidate
				break
			}
		}
		if editor == "" {
			editor = "vi"
		}
	}

	c := exec.Command(editor, runner.TemplatePath())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err}
	})
}

func (m *Model) runTests() tea.Cmd {
	q := m.questions[m.cursor]
	lang := q.Language
	if lang == "" {
		lang = quiz.LangGo
	}
	runner := m.getRunner(lang)
	if runner == nil {
		return nil
	}
	return func() tea.Msg {
		solution, err := os.ReadFile(runner.TemplatePath())
		if err != nil {
			return testResultMsg{result: &quiz.Result{Error: err.Error()}}
		}
		result := runner.Run(string(solution), q.TestCode)
		return testResultMsg{result: result}
	}
}

// SetSize updates the model dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the current quiz page.
func (m *Model) View() string {
	switch m.current {
	case pageQuestionList:
		return m.viewList()
	case pageQuestion:
		return m.viewQuestion()
	case pageResult:
		return m.viewResult()
	default:
		return m.viewList()
	}
}

// levelGroup defines a section header for quiz questions.
type levelGroup struct {
	title string
	start int // first question index (inclusive)
	end   int // last question index (exclusive)
}

// buildLevelGroups computes group boundaries from the question list.
// Groups are determined by the AllQuestions() ordering:
// level1Basics, level2Standards, level3Protocol, level4Advanced, SolidityQuestions.
func buildLevelGroups(questions []quiz.Question) []levelGroup {
	var groups []levelGroup
	type groupDef struct {
		label     string
		matchFunc func(q quiz.Question) bool
	}
	defs := []groupDef{
		{"LEVEL 1 \u2014 Basics", func(q quiz.Question) bool {
			return q.Difficulty == "easy" && q.Language != quiz.LangSolidity
		}},
		{"LEVEL 2 \u2014 Standards", func(q quiz.Question) bool {
			return q.Difficulty == "medium" && q.Language != quiz.LangSolidity && q.Category != "x402"
		}},
		{"LEVEL 3 \u2014 Protocol", func(q quiz.Question) bool {
			return q.Category == "x402" && q.Language != quiz.LangSolidity
		}},
		{"LEVEL 4 \u2014 Advanced", func(q quiz.Question) bool {
			return q.Difficulty == "hard" && q.Language != quiz.LangSolidity
		}},
		{"SOLIDITY", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity
		}},
	}

	assigned := make([]bool, len(questions))
	for _, def := range defs {
		start := -1
		end := -1
		for i, q := range questions {
			if !assigned[i] && def.matchFunc(q) {
				if start == -1 {
					start = i
				}
				end = i + 1
				assigned[i] = true
			}
		}
		if start != -1 {
			groups = append(groups, levelGroup{title: def.label, start: start, end: end})
		}
	}
	return groups
}

func (m *Model) viewList() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render("Learn \u2014 x402 Protocol Quiz")
	scoreText := lipgloss.NewStyle().Foreground(tui.ColorMuted).
		Render(fmt.Sprintf("Score: %d/%d", m.score.Correct, m.score.Questions))
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("\u2500", min(m.width-8, 60)))

	groups := buildLevelGroups(m.questions)

	// Build a set of indices that start a new group, mapped to header text.
	groupHeaders := make(map[int]string)
	for _, g := range groups {
		groupHeaders[g.start] = g.title
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorSecondary).
		MarginTop(1)

	rowWidth := max(min(m.width-8, 64), 20)

	var items strings.Builder
	for i, q := range m.questions {
		// Insert group header if this index starts a new group.
		if header, ok := groupHeaders[i]; ok {
			if i > 0 {
				items.WriteString("\n")
			}
			items.WriteString(headerStyle.Render("  "+header) + "\n")
		}

		icon := "\u25cb"
		iconStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted)
		if r, ok := m.results[i]; ok {
			if r.Passed == r.Total && r.Total > 0 {
				icon = "\u2713"
				iconStyle = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
			} else {
				icon = "\u2717"
				iconStyle = lipgloss.NewStyle().Foreground(tui.ColorError)
			}
		}

		diffStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted)
		switch q.Difficulty {
		case "easy":
			diffStyle = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
		case "medium":
			diffStyle = lipgloss.NewStyle().Foreground(tui.ColorAccent)
		case "hard":
			diffStyle = lipgloss.NewStyle().Foreground(tui.ColorError)
		}

		langIcon := "Go"
		if q.Language == quiz.LangSolidity {
			langIcon = "Sol"
		}
		langStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted)
		catStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary)
		badge := langStyle.Render(langIcon) + " " + catStyle.Render(q.Category) + " " + diffStyle.Render("["+q.Difficulty+"]")

		if i == m.cursor {
			// Full-row highlight bar for selected item
			innerText := fmt.Sprintf(" \u25b8 %s %-28s %s", iconStyle.Render(icon), q.Title, badge)
			row := lipgloss.NewStyle().
				Background(tui.ColorHighlight).
				Foreground(lipgloss.Color("#A78BFA")).
				Bold(true).
				Width(rowWidth).
				Padding(0, 1).
				Render(innerText)
			items.WriteString(row + "\n")
		} else {
			name := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(q.Title)
			fmt.Fprintf(&items, "    %s %s  %s\n", iconStyle.Render(icon), name, badge)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, scoreText, divider, "", items.String())
}

func (m *Model) viewQuestion() string {
	q := m.questions[m.cursor]

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(fmt.Sprintf("Question %d/%d: %s", m.cursor+1, len(m.questions), q.Title))
	subtitle := lipgloss.NewStyle().Foreground(tui.ColorAccent).Render("[" + q.Difficulty + "]")
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", min(m.width-8, 60)))

	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).
		Width(min(m.width-10, 70)).Render(q.Description)

	prompt := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorAccent).
		Render("Press Enter to open $EDITOR...")

	var hintsView string
	if len(q.Hints) > 0 {
		var hb strings.Builder
		hb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).Render("Hints:") + "\n")
		for i, h := range q.Hints {
			hb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).
				Render(fmt.Sprintf("  %d. %s", i+1, h)) + "\n")
		}
		hintsView = hb.String()
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, divider, "", desc, "", prompt, "", hintsView)
}

func (m *Model) viewResult() string {
	q := m.questions[m.cursor]
	r := m.results[m.cursor]

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(fmt.Sprintf("Result: %s", q.Title))
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", min(m.width-8, 60)))

	var status string
	if r == nil {
		status = tui.ErrorStyle.Render("No result")
	} else if r.Error != "" {
		status = lipgloss.JoinVertical(lipgloss.Left,
			tui.ErrorStyle.Render("✗ "+r.Error),
			"",
			lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(min(m.width-10, 80)).Render(r.Output),
		)
	} else {
		var sb strings.Builder
		if r.Compiled {
			sb.WriteString(tui.SuccessStyle.Render("✓ Compilation: PASS") + "\n")
		} else {
			sb.WriteString(tui.ErrorStyle.Render("✗ Compilation: FAIL") + "\n")
		}
		if r.Passed == r.Total && r.Total > 0 {
			sb.WriteString(tui.SuccessStyle.Render(fmt.Sprintf("✓ Tests: %d/%d PASSED", r.Passed, r.Total)) + "\n")
		} else {
			sb.WriteString(tui.ErrorStyle.Render(fmt.Sprintf("✗ Tests: %d/%d passed", r.Passed, r.Total)) + "\n")
		}
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(min(m.width-10, 80)).Render(r.Output))
		status = sb.String()
	}

	scoreText := lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).
		Render(fmt.Sprintf("Score: %d/%d", m.score.Correct, m.score.Questions))

	return lipgloss.JoinVertical(lipgloss.Left, title, divider, "", status, "", scoreText)
}
