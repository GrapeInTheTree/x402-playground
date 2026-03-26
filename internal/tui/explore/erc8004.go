package explore

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

// ERC8004Model shows ERC-8004 registry structures with field exploration.
// Tab cycles through Identity, Reputation, and Validation registries.
type ERC8004Model struct {
	explorer components.FieldExplorer
	registry int // 0=Identity, 1=Reputation, 2=Validation
	width    int
	height   int
}

var registryNames = []string{"Identity Registry", "Reputation Registry", "Validation Registry"}

// NewERC8004Model creates a new ERC-8004 registry explorer.
func NewERC8004Model(width, height int) *ERC8004Model {
	return &ERC8004Model{
		explorer: components.NewFieldExplorer(identityFields()),
		width:    width,
		height:   height,
	}
}

func identityFields() []components.Field {
	return []components.Field{
		{Name: "agentId", Value: "uint256 (auto-increment)", Description: "ERC-721 token ID assigned on registration. Auto-incremented starting from 1. Each agent is a unique NFT owned by the registering address."},
		{Name: "owner", Value: "address", Description: "ERC-721 token owner. Has exclusive rights to set metadata, request validations, and transfer the agent. Standard ERC-721 ownership semantics apply."},
		{Name: "agentURI", Value: "ipfs://... or https://...", Description: "URI pointing to the agent's registration file (JSON). Contains name, description, services, capabilities, and x402 support configuration. Can use IPFS, HTTPS, or data: URIs."},
		{Name: "agentWallet", Value: "address (EIP-712 verified)", Description: "Verified wallet address for the agent. Set via EIP-712 signature proving wallet consent. Automatically cleared on NFT transfer to prevent stale wallet associations. Reserved metadata key — cannot be set via setMetadata()."},
		{Name: "metadata", Value: "mapping(string => string)", Description: "Arbitrary key-value metadata store per agent. Keys are strings, values are strings. The key 'agentWallet' is reserved and will revert if used via setMetadata(). Useful for storing capabilities, version info, etc."},
	}
}

func reputationFields() []components.Field {
	return []components.Field{
		{Name: "agentId", Value: "uint256", Description: "Target agent receiving the feedback. Must be a registered agent in the Identity Registry."},
		{Name: "provider", Value: "address (msg.sender)", Description: "Address giving the feedback. Automatically set to msg.sender. Self-feedback is blocked — the agent's owner or authorized wallet cannot give feedback to their own agent."},
		{Name: "value", Value: "int256 (WAD: 1e18 scale)", Description: "Feedback score using WAD fixed-point math (18 decimals). Positive values = good reputation, negative = bad. Example: 80 * 1e18 = very positive. Range bounded by MAX_ABS_VALUE = 1e38."},
		{Name: "tag", Value: "string", Description: "Category tag for filtering feedback. Examples: 'reliability', 'accuracy', 'speed', 'x402-payment'. Enables per-category reputation summaries."},
		{Name: "timestamp", Value: "uint64 (block.timestamp)", Description: "Block timestamp when feedback was recorded. Used for time-weighted reputation calculations and freshness filtering."},
		{Name: "revoked", Value: "bool", Description: "True if the original feedback provider revoked this feedback. Only the original provider can revoke. Revoked feedback is excluded from summary calculations."},
		{Name: "proofOfPayment", Value: "{txHash, nonce}", Description: "Optional x402 integration: links feedback to a verified on-chain payment. Contains transaction hash and nonce from an x402 settlement, creating cryptographic proof that the feedback is backed by a real payment interaction."},
	}
}

func validationFields() []components.Field {
	return []components.Field{
		{Name: "requestId", Value: "uint256 (auto-increment)", Description: "Unique validation request identifier. Auto-incremented per request. Used to track and respond to specific validation requests."},
		{Name: "agentId", Value: "uint256", Description: "Agent being validated. Only the agent's owner can create validation requests for their agent."},
		{Name: "requester", Value: "address", Description: "Agent owner who initiated the validation request. Must be ownerOf(agentId) to create a request."},
		{Name: "validator", Value: "address", Description: "Designated validator address. Only this specific address can respond to the request. Enables targeted validation from trusted parties (auditors, stakers, TEE oracles)."},
		{Name: "score", Value: "uint8 (0-100)", Description: "Validation score from the validator. 0 = complete failure, 100 = perfect. Can represent binary (0 or 100) or spectrum validation. Updated by the validator's response."},
		{Name: "reason", Value: "string", Description: "Free-text explanation from the validator justifying the score. Provides human-readable context for the validation result."},
		{Name: "status", Value: "Pending → Completed", Description: "Request lifecycle: created as Pending when owner requests, transitions to Completed when validator responds. Validators can update their response (progressive validation)."},
	}
}

// Update handles key events for field navigation and registry switching.
func (m *ERC8004Model) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			m.explorer.Up()
		case "down", "j":
			m.explorer.Down()
		case "tab":
			m.registry = (m.registry + 1) % 3
			switch m.registry {
			case 0:
				m.explorer = components.NewFieldExplorer(identityFields())
			case 1:
				m.explorer = components.NewFieldExplorer(reputationFields())
			case 2:
				m.explorer = components.NewFieldExplorer(validationFields())
			}
		}
	}
	return nil
}

// SetSize updates the model dimensions.
func (m *ERC8004Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.explorer.Width = width
}

// View renders the ERC-8004 registry field explorer with split-panel layout.
func (m *ERC8004Model) View() string {
	availW := m.width - 4
	gap := 1
	innerTotal := availW - gap - 4
	leftW := innerTotal * 2 / 5
	rightW := innerTotal - leftW

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorBorder).
		Padding(0, 1)

	// Left: field list
	leftTitle := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
		Render("ERC-8004 — " + registryNames[m.registry])
	tabHint := tui.MutedStyle.Render("Tab to switch registry")
	var fieldList strings.Builder
	nameW := max(leftW/2, 12)
	for i, f := range m.explorer.Fields {
		nameStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary).Width(nameW)
		valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))
		if i == m.explorer.Cursor {
			nameStyle = lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true).Width(nameW)
			valStyle = lipgloss.NewStyle().Foreground(tui.ColorAccent)
			fmt.Fprintf(&fieldList, " > %s %s\n", nameStyle.Render(f.Name+":"), valStyle.Render(f.Value))
		} else {
			fmt.Fprintf(&fieldList, "   %s %s\n", nameStyle.Render(f.Name+":"), valStyle.Render(f.Value))
		}
	}
	leftContent := lipgloss.JoinVertical(lipgloss.Left, leftTitle, tabHint, "", fieldList.String())
	leftBox := boxStyle.Width(leftW).Height(m.height - 4).Render(leftContent)

	// Right: description
	rightTitle := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render("Field Details")
	desc := ""
	if m.explorer.Cursor >= 0 && m.explorer.Cursor < len(m.explorer.Fields) {
		f := m.explorer.Fields[m.explorer.Cursor]
		name := lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).Render(f.Name)
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(f.Value)
		body := lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(rightW - 2).Render(f.Description)
		desc = lipgloss.JoinVertical(lipgloss.Left, name, val, "", body)
	}
	rightContent := lipgloss.JoinVertical(lipgloss.Left, rightTitle, "", desc)
	rightBox := boxStyle.Width(rightW).Height(m.height - 4).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
}
