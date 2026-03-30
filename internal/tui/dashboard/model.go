package dashboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/demo"
	"github.com/GrapeInTheTree/x402-playground/internal/quiz"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

type balancesMsg struct {
	balances []demo.WalletBalance
	err      error
}

// Model is the dashboard page TUI model showing wallet balances.
type Model struct {
	balances []demo.WalletBalance
	loading  bool
	err      error
	cfg      *config.ExplorerConfig
	wallets  []demo.WalletInfo
	progress *quiz.QuizProgress
	spinner  spinner.Model
	width    int
	height   int
}

// New creates a new dashboard model with the given dimensions and configuration.
func New(width, height int, cfg *config.ExplorerConfig, progress *quiz.QuizProgress) *Model {
	wallets := []demo.WalletInfo{}
	if cfg != nil && cfg.PayToAddress != "" {
		wallets = append(wallets, demo.WalletInfo{Name: "PAY_TO", Address: cfg.PayToAddress})
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tui.ColorAccent)

	return &Model{
		loading:  true,
		cfg:      cfg,
		wallets:  wallets,
		progress: progress,
		spinner:  s,
		width:    width,
		height:   height,
	}
}

// Init starts the initial balance fetch and spinner.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.fetchBalances(), m.spinner.Tick)
}

func (m *Model) fetchBalances() tea.Cmd {
	if m.cfg == nil || m.cfg.RPCURL == "" || len(m.wallets) == 0 {
		return func() tea.Msg {
			return balancesMsg{err: fmt.Errorf("configuration incomplete — set RPC_URL, PAY_TO_ADDRESS in .env")}
		}
	}

	cfg := m.cfg
	wallets := m.wallets
	return func() tea.Msg {
		client, err := ethclient.Dial(cfg.RPCURL)
		if err != nil {
			return balancesMsg{err: err}
		}
		defer client.Close()

		bals, err := demo.QueryBalances(context.Background(), client, cfg.USDCAddress, wallets)
		return balancesMsg{balances: bals, err: err}
	}
}

// Update handles balance results, key events, and spinner ticks.
func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case balancesMsg:
		m.loading = false
		m.balances = msg.balances
		m.err = msg.err

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return tui.BackMsg{} }
		case "r":
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.fetchBalances(), m.spinner.Tick)
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// SetSize updates the model dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the dashboard with side-by-side wallet and progress panels.
func (m *Model) View() string {
	// lipgloss Width = content+padding, border added outside.
	// RootModel padding = 4 chars. Each box border = 2 chars wide.
	// (leftW+2) + 1(gap) + (rightW+2) <= m.width - 4
	gap := 1
	innerTotal := m.width - 4 - gap - 4
	leftW := innerTotal / 2
	rightW := innerTotal - leftW

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorBorder).
		Padding(0, 1)

	// ── Left panel: Wallet Balances ──
	leftBox := boxStyle.Width(leftW).Render(m.renderWalletPanel())

	// ── Right panel: Quiz Progress ──
	rightContent := m.renderProgressPanel(rightW)
	rightBox := boxStyle.Width(rightW).Render(rightContent)

	// Side-by-side if wide enough, stacked otherwise
	if innerTotal >= 50 {
		return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
	}
	return lipgloss.JoinVertical(lipgloss.Left, leftBox, "", rightBox)
}

func (m *Model) renderWalletPanel() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
		Render("Wallet Balances")

	network := ""
	if m.cfg != nil {
		network = m.cfg.Network
	}
	netLine := tui.MutedStyle.Render(fmt.Sprintf("Network: %s", network))
	usdcLine := tui.MutedStyle.Render(fmt.Sprintf("USDC:    %s", m.usdcAddr()))

	var body string
	if m.loading {
		body = m.spinner.View() + " Loading balances..."
	} else if m.err != nil {
		body = tui.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	} else if len(m.balances) == 0 {
		body = tui.MutedStyle.Render("No wallet data.")
	} else {
		body = m.renderBalances()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title, "", netLine, usdcLine, "", body)
}

func (m *Model) usdcAddr() string {
	if m.cfg != nil {
		addr := m.cfg.USDCAddress
		if len(addr) > 14 {
			return addr[:10] + "..." + addr[len(addr)-4:]
		}
		return addr
	}
	return "unknown"
}

// shortModuleName maps full group titles to compact display names.
func shortModuleName(name string) string {
	switch {
	case strings.Contains(name, "LEVEL 1"):
		return "Go Basics"
	case strings.Contains(name, "LEVEL 2"):
		return "Go Standards"
	case strings.Contains(name, "LEVEL 3"):
		return "Go Protocol"
	case strings.Contains(name, "LEVEL 4"):
		return "Go Advanced"
	case strings.Contains(name, "LEVEL 5"):
		return "Go Agents"
	case strings.Contains(name, "M1:"):
		return "Sol Foundations"
	case strings.Contains(name, "M2:"):
		return "Sol ERC-20"
	case strings.Contains(name, "M3:"):
		return "Sol Signatures"
	case strings.Contains(name, "M4:"):
		return "Sol Gasless"
	case strings.Contains(name, "M5:"):
		return "Sol Advanced"
	case strings.Contains(name, "M6:"):
		return "Sol x402"
	case strings.Contains(name, "M7:"):
		return "Sol ERC-8004"
	default:
		return name
	}
}

// progressBar renders a uniform-width bar: filled portion + empty portion.
func progressBar(passed, total, width int) string {
	filled := 0
	if total > 0 {
		filled = passed * width / total
	}
	filledStr := lipgloss.NewStyle().Foreground(tui.ColorSuccess).
		Render(strings.Repeat("━", filled))
	emptyStr := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", width-filled))
	return filledStr + emptyStr
}

func (m *Model) renderProgressPanel(panelW int) string {
	if m.progress == nil {
		title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
			Render("Quiz Progress")
		hint := tui.MutedStyle.Render("Complete quizzes in Learn\nto see progress here.")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", hint)
	}

	modules := m.progress.GetModules()
	if len(modules) == 0 {
		title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
			Render("Quiz Progress")
		hint := tui.MutedStyle.Render("Complete quizzes in Learn\nto see progress here.")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", hint)
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
		Render("Quiz Progress")

	// Totals
	totalQ, totalPassed, totalAttempted := 0, 0, 0
	for _, mod := range modules {
		totalQ += mod.Total
		totalPassed += mod.Passed
		totalAttempted += mod.Attempted
	}

	pct := 0
	if totalQ > 0 {
		pct = totalPassed * 100 / totalQ
	}

	// Inner width = panelW - border(2) - padding(2)
	innerW := panelW - 4

	// ── Overall progress bar (full inner width) ──
	overallBarW := max(innerW-14, 8) // leave room for " 0% (0/36)"
	overallBar := progressBar(totalPassed, totalQ, overallBarW)
	overallLine := overallBar + tui.MutedStyle.Render(
		fmt.Sprintf(" %d%% (%d/%d)", pct, totalPassed, totalQ))

	// ── Per-module rows: [name] [bar] [count] [check] ──
	// Layout: nameCol + 1 + barW + 1 + countCol + 2 = innerW
	const nameCol = 16
	const countCol = 4
	barW := max(innerW-nameCol-countCol-4, 4)

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Width(nameCol)
	countStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(countCol).Align(lipgloss.Right)
	checkStyle := lipgloss.NewStyle().Foreground(tui.ColorSuccess)

	var rows strings.Builder
	for _, mod := range modules {
		bar := progressBar(mod.Passed, mod.Total, barW)
		check := "  "
		if mod.Passed == mod.Total && mod.Total > 0 {
			check = " " + checkStyle.Render("\u2713")
		}
		fmt.Fprintf(&rows, "%s %s %s%s\n",
			nameStyle.Render(shortModuleName(mod.Name)),
			bar,
			countStyle.Render(fmt.Sprintf("%d/%d", mod.Passed, mod.Total)),
			check)
	}

	// Summary
	failed := totalAttempted - totalPassed
	summary := tui.MutedStyle.Render(
		fmt.Sprintf("Attempted: %d  Passed: %d  Failed: %d", totalAttempted, totalPassed, failed))

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		overallLine,
		"",
		rows.String(),
		summary,
	)
}

func (m *Model) renderBalances() string {
	var b strings.Builder

	nameWidth := 16
	nameStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary).Bold(true).Width(nameWidth)
	addrStyle := tui.MutedStyle
	ethLabel := lipgloss.NewStyle().Foreground(tui.ColorMuted).Render("ETH: ")
	usdcLabel := lipgloss.NewStyle().Foreground(tui.ColorMuted).Render("USDC: ")

	for _, bal := range m.balances {
		addr := bal.Wallet.Address
		if len(addr) > 20 {
			addr = addr[:10] + "..." + addr[len(addr)-4:]
		}

		b.WriteString(fmt.Sprintf("%s %s\n",
			nameStyle.Render(bal.Wallet.Name),
			addrStyle.Render(addr)))

		ethVal := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(bal.ETH)
		usdcVal := lipgloss.NewStyle().Foreground(tui.ColorSuccess).Bold(true).Render(bal.USDC)

		b.WriteString(fmt.Sprintf("%s%s  %s    %s%s\n\n",
			strings.Repeat(" ", nameWidth),
			ethLabel, ethVal,
			usdcLabel, usdcVal))
	}

	return b.String()
}
