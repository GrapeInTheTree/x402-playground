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
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
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
	spinner  spinner.Model
	width    int
	height   int
}

// New creates a new dashboard model with the given dimensions and configuration.
func New(width, height int, cfg *config.ExplorerConfig) *Model {
	wallets := []demo.WalletInfo{}
	if cfg != nil && cfg.PayToAddress != "" {
		wallets = append(wallets, demo.WalletInfo{Name: "PAY_TO", Address: cfg.PayToAddress})
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tui.ColorAccent)

	return &Model{
		loading: true,
		cfg:     cfg,
		wallets: wallets,
		spinner: s,
		width:   width,
		height:  height,
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

// View renders the dashboard with network info and wallet balances.
func (m *Model) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		MarginLeft(2).
		Render("Dashboard — Wallet Balances")

	var content string

	if m.loading {
		content = lipgloss.NewStyle().MarginLeft(4).Render(m.spinner.View() + " Loading balances from chain...")
	} else if m.err != nil {
		content = tui.ErrorStyle.MarginLeft(4).Render(fmt.Sprintf("Error: %v", m.err))
	} else if len(m.balances) == 0 {
		content = tui.MutedStyle.MarginLeft(4).Render("No wallet data available.")
	} else {
		content = m.renderBalances()
	}

	network := ""
	if m.cfg != nil {
		network = m.cfg.Network
	}
	networkInfo := tui.MutedStyle.MarginLeft(4).Render(
		fmt.Sprintf("Network: %s  |  USDC: %s", network, m.usdcAddr()))

	hints := components.StatusBar{Width: m.width}.View("  r refresh  ? help  esc back")

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		networkInfo,
		"",
		content,
		"",
		hints,
	)
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

		b.WriteString(fmt.Sprintf("    %s %s\n",
			nameStyle.Render(bal.Wallet.Name),
			addrStyle.Render(addr)))

		ethVal := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(bal.ETH)
		usdcVal := lipgloss.NewStyle().Foreground(tui.ColorSuccess).Bold(true).Render(bal.USDC)

		b.WriteString(fmt.Sprintf("    %s%s  %s    %s%s\n\n",
			strings.Repeat(" ", nameWidth),
			ethLabel, ethVal,
			usdcLabel, usdcVal))
	}

	return b.String()
}
