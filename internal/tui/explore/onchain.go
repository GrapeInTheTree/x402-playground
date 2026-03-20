package explore

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/GrapeInTheTree/x402-playground/internal/demo"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

type balancesMsg struct {
	balances []demo.WalletBalance
	err      error
}

type allowanceMsg struct {
	allowance string
	err       error
}

// OnChainModel shows on-chain state: balances, allowance, contract info.
type OnChainModel struct {
	balances  []demo.WalletBalance
	allowance string
	loading   bool
	err       error
	rpcURL    string
	usdcAddr  string
	wallets   []demo.WalletInfo
	width     int
	height    int
}

// NewOnChainModel creates a new on-chain state viewer.
func NewOnChainModel(width, height int) *OnChainModel {
	return &OnChainModel{
		width:  width,
		height: height,
	}
}

// Configure sets the chain connection parameters for fetching on-chain data.
func (m *OnChainModel) Configure(rpcURL, usdcAddr string, wallets []demo.WalletInfo) {
	m.rpcURL = rpcURL
	m.usdcAddr = usdcAddr
	m.wallets = wallets
}

// FetchBalances creates a tea.Cmd that fetches balances from chain.
func (m *OnChainModel) FetchBalances() tea.Cmd {
	if m.rpcURL == "" {
		return nil
	}
	rpcURL := m.rpcURL
	usdcAddr := m.usdcAddr
	wallets := m.wallets
	return func() tea.Msg {
		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			return balancesMsg{err: err}
		}
		defer client.Close()
		bals, err := demo.QueryBalances(context.Background(), client, usdcAddr, wallets)
		return balancesMsg{balances: bals, err: err}
	}
}

// Update handles balance and allowance messages.
func (m *OnChainModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case balancesMsg:
		m.loading = false
		m.balances = msg.balances
		m.err = msg.err
	case allowanceMsg:
		m.allowance = msg.allowance
		if msg.err != nil && m.err == nil {
			m.err = msg.err
		}
	}
	return nil
}

// SetSize updates the model dimensions.
func (m *OnChainModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the on-chain state including balances and allowances.
func (m *OnChainModel) View() string {
	title := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		MarginLeft(4).
		Render("On-Chain State")

	var b strings.Builder

	if m.rpcURL == "" {
		b.WriteString(tui.MutedStyle.Render("    RPC not configured. Set RPC_URL environment variable."))
		return lipgloss.JoinVertical(lipgloss.Left, title, "", b.String())
	}

	if m.loading {
		b.WriteString("    Loading on-chain data...")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", b.String())
	}

	if m.err != nil {
		b.WriteString(tui.ErrorStyle.Render(fmt.Sprintf("    Error: %v", m.err)))
		return lipgloss.JoinVertical(lipgloss.Left, title, "", b.String())
	}

	if len(m.balances) > 0 {
		nameStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary).Width(16)
		addrStyle := tui.MutedStyle
		valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))

		for _, bal := range m.balances {
			addr := bal.Wallet.Address
			if len(addr) > 14 {
				addr = addr[:10] + "..." + addr[len(addr)-4:]
			}
			b.WriteString(fmt.Sprintf("    %s %s\n",
				nameStyle.Render(bal.Wallet.Name+":"),
				addrStyle.Render(addr)))
			b.WriteString(fmt.Sprintf("      ETH: %s    USDC: %s\n\n",
				valStyle.Render(bal.ETH),
				valStyle.Render(bal.USDC)))
		}
	} else {
		b.WriteString(tui.MutedStyle.Render("    No balance data. Press 'r' to refresh."))
	}

	if m.allowance != "" {
		b.WriteString(fmt.Sprintf("    Permit2 Allowance: %s USDC\n",
			lipgloss.NewStyle().Foreground(tui.ColorAccent).Render(m.allowance)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, "", b.String())
}
