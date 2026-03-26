package learn

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/quiz"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

type page int

const (
	pageQuestionList page = iota
	pageQuestion
	pageAnimation
	pageResult
)

type editorFinishedMsg struct{ err error }

type testResultMsg struct{ result *quiz.Result }

type animTickMsg struct{}

const maxAnimTicks = 20

// Model is the quiz-based learning page model.
type Model struct {
	questions []quiz.Question
	cursor    int
	current   page
	goRunner  *quiz.Runner
	solRunner *quiz.Runner
	results       map[int]*quiz.Result
	score         quiz.Score
	progress      *quiz.QuizProgress
	errMsg        string // error message to display
	lastEditedIdx int    // last question index opened in editor (-1 = none)
	animTicks     int    // animation tick counter
	animResult    *quiz.Result // result being animated
	scrollOffset  int    // viewport scroll for question list
	width         int
	height        int
}

// New creates a new quiz learning model.
func New(width, height int, progress *quiz.QuizProgress) *Model {
	questions := quiz.AllQuestions()
	return &Model{
		questions:     questions,
		results:       make(map[int]*quiz.Result),
		score:         quiz.Score{Questions: len(questions)},
		progress:      progress,
		lastEditedIdx: -1,
		width:         width,
		height:        height,
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
		m.animResult = msg.result
		m.animTicks = 0
		m.current = pageAnimation
		m.syncProgress()
		return m, tickAnimation()

	case animTickMsg:
		if m.current != pageAnimation {
			return m, nil
		}
		m.animTicks++
		if m.animTicks >= maxAnimTicks {
			m.current = pageResult
			return m, nil
		}
		return m, tickAnimation()

	case tea.KeyMsg:
		switch m.current {
		case pageQuestionList:
			return m.updateList(msg)
		case pageQuestion:
			return m.updateQuestion(msg)
		case pageAnimation:
			// Any key skips to result
			m.current = pageResult
			return m, nil
		case pageResult:
			return m.updateResult(msg)
		}
	}
	return m, nil
}

func tickAnimation() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return animTickMsg{}
	})
}

func (m *Model) syncProgress() {
	if m.progress == nil {
		return
	}
	groups := buildLevelGroups(m.questions)
	modules := make([]quiz.ModuleProgress, len(groups))
	for i, g := range groups {
		mp := quiz.ModuleProgress{
			Name:  g.title,
			Total: len(g.indices),
		}
		for _, j := range g.indices {
			if r, ok := m.results[j]; ok {
				mp.Attempted++
				if r.Passed == r.Total && r.Total > 0 {
					mp.Passed++
				}
			}
		}
		modules[i] = mp
	}
	m.progress.Modules = modules
	m.progress.Score = m.score
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
			m.errMsg = "Failed to create runner: " + err.Error()
			return nil
		}
		m.setRunner(lang, runner)
	}

	// Write template only for a new question (preserve user's code on retry)
	if m.lastEditedIdx != m.cursor {
		if err := os.WriteFile(runner.TemplatePath(), []byte(q.Template), 0644); err != nil {
			m.errMsg = "Failed to write template: " + err.Error()
			return nil
		}
		m.lastEditedIdx = m.cursor
	}
	m.errMsg = ""

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
	case pageAnimation:
		return m.viewAnimation()
	case pageResult:
		return m.viewResult()
	default:
		return m.viewList()
	}
}

// levelGroup defines a section header for quiz questions.
type levelGroup struct {
	title   string
	indices []int // actual question indices in this group
}

// buildLevelGroups computes group boundaries from the question list.
// Go questions: 4 levels by difficulty/category (mutually exclusive matchers).
// Solidity questions: 6 modules by Category prefix "M1:"..."M6:".
func buildLevelGroups(questions []quiz.Question) []levelGroup {
	var groups []levelGroup
	type groupDef struct {
		label     string
		matchFunc func(q quiz.Question) bool
	}

	// Match functions are mutually exclusive to prevent cross-group leaking.
	// Go questions: exclude x402 and ERC-8004 categories from difficulty-based groups.
	defs := []groupDef{
		{"GO — LEVEL 1: Basics", func(q quiz.Question) bool {
			return q.Language != quiz.LangSolidity && q.Difficulty == "easy" && q.Category != "x402" && q.Category != "ERC-8004"
		}},
		{"GO — LEVEL 2: Standards", func(q quiz.Question) bool {
			return q.Language != quiz.LangSolidity && q.Difficulty == "medium" && q.Category != "x402" && q.Category != "ERC-8004"
		}},
		{"GO — LEVEL 3: Protocol", func(q quiz.Question) bool {
			return q.Language != quiz.LangSolidity && q.Category == "x402" && q.Difficulty != "hard"
		}},
		{"GO — LEVEL 4: Advanced", func(q quiz.Question) bool {
			return q.Language != quiz.LangSolidity && q.Difficulty == "hard" && q.Category != "ERC-8004"
		}},
		{"GO — LEVEL 5: Agents", func(q quiz.Question) bool {
			return q.Language != quiz.LangSolidity && q.Category == "ERC-8004"
		}},
		{"SOLIDITY — M1: Foundations", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M1:")
		}},
		{"SOLIDITY — M2: ERC-20", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M2:")
		}},
		{"SOLIDITY — M3: Signatures", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M3:")
		}},
		{"SOLIDITY — M4: Gasless", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M4:")
		}},
		{"SOLIDITY — M5: Advanced", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M5:")
		}},
		{"SOLIDITY — M6: x402", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M6:")
		}},
		{"SOLIDITY — M7: ERC-8004", func(q quiz.Question) bool {
			return q.Language == quiz.LangSolidity && strings.HasPrefix(q.Category, "M7:")
		}},
	}

	assigned := make([]bool, len(questions))
	for _, def := range defs {
		var indices []int
		for i, q := range questions {
			if !assigned[i] && def.matchFunc(q) {
				indices = append(indices, i)
				assigned[i] = true
			}
		}
		if len(indices) > 0 {
			groups = append(groups, levelGroup{title: def.label, indices: indices})
		}
	}
	return groups
}

func (m *Model) viewList() string {
	// lipgloss Width/Height = content+padding, border is ADDED outside.
	// RootModel wraps content with PaddingLeft(2)+PaddingRight(2).
	// So available outer space = (m.width - 4) wide, m.height tall.
	//
	// Each box: outer_w = Width + 2(border), outer_h = Height + 2(border).
	// Two boxes + 1 gap: (leftW+2) + 1 + (rightW+2) <= m.width - 4
	//   => leftW + rightW <= m.width - 9
	// Height: boxH + 2 <= m.height  =>  boxH <= m.height - 2

	gap := 1
	innerTotal := m.width - 4 - gap - 4 // subtract RootModel padding, gap, 2 borders
	leftW := innerTotal / 2
	rightW := innerTotal - leftW
	boxH := m.height - 2 // subtract top+bottom border

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorBorder).
		Padding(0, 1)

	// Padding(0,1) = 1 left + 1 right = 2 horizontal padding inside Width
	leftContentW := leftW - 2
	leftContent := m.buildQuestionList(leftContentW, boxH)
	leftBox := boxStyle.Width(leftW).Height(boxH).Render(leftContent)

	rightContentW := rightW - 2
	rightContent := m.buildQuestionPreview(rightContentW, boxH)
	rightBox := boxStyle.Width(rightW).Height(boxH).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
}

// buildQuestionList renders the scrollable question list for the left panel.
func (m *Model) buildQuestionList(innerW, boxH int) string {
	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render("Questions")
	scoreText := lipgloss.NewStyle().Foreground(tui.ColorMuted).
		Render(fmt.Sprintf("Score: %d/%d", m.score.Correct, m.score.Questions))

	headerH := 3 // title + score + blank line

	groups := buildLevelGroups(m.questions)
	groupHeaders := make(map[int]string)
	for _, g := range groups {
		if len(g.indices) > 0 {
			groupHeaders[g.indices[0]] = g.title
		}
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary)
	rowWidth := max(innerW, 20)

	// Build all lines and track cursor position.
	var lines []string
	cursorLine := 0
	cursorGroupHeaderLine := 0
	for i, q := range m.questions {
		if hdr, ok := groupHeaders[i]; ok {
			if i > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, headerStyle.Render(" "+hdr))
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
		diffBadge := diffStyle.Render("[" + q.Difficulty + "]")

		if i == m.cursor {
			cursorLine = len(lines)
			if _, isGroupStart := groupHeaders[i]; isGroupStart {
				if i > 0 {
					cursorGroupHeaderLine = cursorLine - 2
				} else {
					cursorGroupHeaderLine = cursorLine - 1
				}
			} else {
				cursorGroupHeaderLine = cursorLine
			}
			row := lipgloss.NewStyle().
				Background(tui.ColorHighlight).
				Foreground(lipgloss.Color("#A78BFA")).
				Bold(true).
				Width(rowWidth).
				Render(fmt.Sprintf(" \u25b8 %s %s %s", iconStyle.Render(icon), q.Title, diffBadge))
			lines = append(lines, row)
		} else {
			name := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(q.Title)
			lines = append(lines, fmt.Sprintf("   %s %s %s", iconStyle.Render(icon), name, diffBadge))
		}
	}

	// Viewport scrolling: boxH minus border(2) minus header lines
	visibleH := max(boxH-2-headerH, 5)
	scrollTarget := max(cursorGroupHeaderLine, 0)
	if scrollTarget < m.scrollOffset {
		m.scrollOffset = scrollTarget
	}
	if cursorLine >= m.scrollOffset+visibleH {
		m.scrollOffset = cursorLine - visibleH + 1
	}
	m.scrollOffset = max(m.scrollOffset, 0)

	end := min(m.scrollOffset+visibleH, len(lines))
	visible := strings.Join(lines[m.scrollOffset:end], "\n")

	scrollHint := ""
	if m.scrollOffset > 0 || end < len(lines) {
		scrollHint = tui.MutedStyle.Render(
			fmt.Sprintf(" [%d-%d/%d]", m.scrollOffset+1, end, len(lines)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, scoreText+scrollHint, "", visible)
}

// buildQuestionPreview renders the right-panel preview for the selected question.
// innerH is the available inner height (box height minus border).
func (m *Model) buildQuestionPreview(innerW, innerH int) string {
	q := m.questions[m.cursor]
	w := max(innerW, 10)

	// Title
	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(q.Title)

	// Badges
	diffStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted)
	switch q.Difficulty {
	case "easy":
		diffStyle = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
	case "medium":
		diffStyle = lipgloss.NewStyle().Foreground(tui.ColorAccent)
	case "hard":
		diffStyle = lipgloss.NewStyle().Foreground(tui.ColorError)
	}
	langLabel := "Go"
	if q.Language == quiz.LangSolidity {
		langLabel = "Solidity"
	}
	badges := tui.MutedStyle.Render(langLabel) + "  " +
		lipgloss.NewStyle().Foreground(tui.ColorSecondary).Render(q.Category) + "  " +
		diffStyle.Render("["+q.Difficulty+"]")

	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", w))

	// Description (word-wrapped to panel width)
	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D1D5DB")).
		Width(w).
		Render(q.Description)

	// Result status if already attempted
	var resultView string
	if r, ok := m.results[m.cursor]; ok {
		if r.Passed == r.Total && r.Total > 0 {
			resultView = tui.SuccessStyle.Render(
				fmt.Sprintf("\u2713 PASSED %d/%d", r.Passed, r.Total))
		} else {
			resultView = tui.ErrorStyle.Render(
				fmt.Sprintf("\u2717 FAILED %d/%d", r.Passed, r.Total))
		}
	}

	prompt := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorAccent).
		Render("Enter \u2192 open editor")

	// Hints
	var hintsView string
	if len(q.Hints) > 0 {
		var hb strings.Builder
		hb.WriteString(tui.MutedStyle.Render("Hints:") + "\n")
		for i, h := range q.Hints {
			hb.WriteString(lipgloss.NewStyle().
				Foreground(tui.ColorMuted).
				Width(w).
				Render(fmt.Sprintf(" %d. %s", i+1, h)) + "\n")
		}
		hintsView = hb.String()
	}

	full := lipgloss.JoinVertical(lipgloss.Left,
		title, badges, divider, "", desc, "", resultView, "", prompt, "", hintsView)

	// Truncate to fit within the panel height
	allLines := strings.Split(full, "\n")
	if len(allLines) > innerH {
		allLines = allLines[:innerH]
	}
	return strings.Join(allLines, "\n")
}

func (m *Model) viewQuestion() string {
	q := m.questions[m.cursor]

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(fmt.Sprintf("Question %d/%d: %s", m.cursor+1, len(m.questions), q.Title))
	subtitle := lipgloss.NewStyle().Foreground(tui.ColorAccent).Render("[" + q.Difficulty + "]")
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", max(m.width-6, 20)))

	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).
		Width(max(m.width-6, 30)).Render(q.Description)

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

	var errView string
	if m.errMsg != "" {
		errView = tui.ErrorStyle.Render("Error: " + m.errMsg)
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, divider, "", desc, "", prompt, "", errView, hintsView)
}

func (m *Model) viewAnimation() string {
	r := m.animResult
	if r == nil {
		return ""
	}

	allPassed := r.Passed == r.Total && r.Total > 0

	var symbol, label, detail string
	var baseColor lipgloss.Color

	if allPassed {
		symbol = "\u2713"
		label = "ALL TESTS PASSED"
		detail = fmt.Sprintf("%d/%d tests", r.Passed, r.Total)
		baseColor = lipgloss.Color("#10B981")
	} else {
		symbol = "\u2717"
		label = "TESTS FAILED"
		detail = fmt.Sprintf("%d/%d passed", r.Passed, r.Total)
		baseColor = lipgloss.Color("#EF4444")
	}

	// Pulse effect: alternate brightness on even/odd ticks
	brightColor := baseColor
	dimColor := lipgloss.Color("#6B7280")
	color := brightColor
	if m.animTicks%2 == 1 {
		color = dimColor
	}

	symbolLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		Render(fmt.Sprintf("        %s %s", symbol, label))
	detailLine := lipgloss.NewStyle().
		Foreground(color).
		Render(fmt.Sprintf("          %s", detail))

	var skipHint string
	if !allPassed {
		skipHint = lipgloss.NewStyle().Foreground(tui.ColorMuted).
			Render("      Press any key for details...")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, "", "", symbolLine, detailLine, "", skipHint)

	// Center vertically
	contentH := lipgloss.Height(content)
	padTop := max((m.height-contentH)/2, 0)
	return strings.Repeat("\n", padTop) + content
}

func (m *Model) viewResult() string {
	q := m.questions[m.cursor]
	r := m.results[m.cursor]

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(fmt.Sprintf("Result: %s", q.Title))
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", max(m.width-6, 20)))

	var status string
	if r == nil {
		status = tui.ErrorStyle.Render("No result")
	} else if r.Error != "" {
		status = lipgloss.JoinVertical(lipgloss.Left,
			tui.ErrorStyle.Render("✗ "+r.Error),
			"",
			lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(max(m.width-4, 30)).Render(r.Output),
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
		sb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(max(m.width-4, 30)).Render(r.Output))
		status = sb.String()
	}

	scoreText := lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).
		Render(fmt.Sprintf("Score: %d/%d", m.score.Correct, m.score.Questions))

	return lipgloss.JoinVertical(lipgloss.Left, title, divider, "", status, "", scoreText)
}
