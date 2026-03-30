package practice

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/demo"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

type stepResultMsg struct {
	step int
	data string
}

type stepErrorMsg struct {
	step int
	err  error
}

// PaymentFlowModel manages a 10-step payment flow with live execution.
// Parameterized by flow type (eip3009 or permit2) and step descriptions.
type PaymentFlowModel struct {
	sm       stepManager
	executor *demo.LiveExecutor
	execErr  string // error creating executor
	running  bool
	spinner  spinner.Model
	results  [10]string
	errors   [10]string
	width    int
	height   int
	cfg      *config.ExplorerConfig
}

var eip3009Descriptions = []stepDesc{
	{"Check wallet addresses", "—", "—"},
	{"—", "Call GET /supported", "/supported response"},
	{"GET /weather (no payment)", "Returns 402", "—"},
	{"Decode PAYMENT-REQUIRED", "—", "—"},
	{"Create EIP-712 signature", "—", "—"},
	{"Send PAYMENT-SIGNATURE", "Parse header → /verify", "—"},
	{"—", "Forward /verify request", "Verify signature/balance/simulation"},
	{"Receive 200 + data", "Return data + /settle", "—"},
	{"—", "—", "Submit on-chain transaction"},
	{"Check final balances", "—", "—"},
}

var permit2Descriptions = []stepDesc{
	{"Check wallet addresses + Permit2 approve", "—", "—"},
	{"—", "Call GET /supported", "/supported response"},
	{"GET /weather (no payment)", "402 + assetTransferMethod:permit2", "—"},
	{"Decode PAYMENT-REQUIRED (Permit2)", "—", "—"},
	{"Create Permit2 EIP-712 signature", "—", "—"},
	{"Send PAYMENT-SIGNATURE", "Parse header → /verify", "—"},
	{"—", "Forward /verify request", "Permit2 signature + allowance verification"},
	{"Receive 200 + data", "Return data + /settle", "—"},
	{"—", "—", "Submit x402Permit2Proxy.settle()"},
	{"Check final balances", "—", "—"},
}

// NewEIP3009Flow creates a payment flow model for EIP-3009.
func NewEIP3009Flow(width, height int, cfg *config.ExplorerConfig) *PaymentFlowModel {
	return newPaymentFlowModel(width, height, cfg, "eip3009", eip3009Descriptions)
}

// NewPermit2Flow creates a payment flow model for Permit2.
func NewPermit2Flow(width, height int, cfg *config.ExplorerConfig) *PaymentFlowModel {
	return newPaymentFlowModel(width, height, cfg, "permit2", permit2Descriptions)
}

func newPaymentFlowModel(width, height int, cfg *config.ExplorerConfig, method string, descriptions []stepDesc) *PaymentFlowModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tui.ColorAccent)

	m := &PaymentFlowModel{
		sm:      newStepManager(demo.NewFlowState(method), descriptions),
		spinner: s,
		width:   width,
		height:  height,
		cfg:     cfg,
	}

	if cfg != nil {
		exec, err := demo.NewLiveExecutor(
			cfg.FacilitatorURL, cfg.ResourceURL, cfg.RPCURL,
			cfg.USDCAddress, cfg.PayToAddress, cfg.ClientPrivateKey, method,
		)
		if err != nil {
			m.execErr = err.Error()
		} else {
			m.executor = exec
		}
	} else {
		m.execErr = "Configuration missing — set CLIENT_PRIVATE_KEY, RESOURCE_URL, FACILITATOR_URL in .env"
	}

	return m
}

// Init starts the spinner tick.
func (m *PaymentFlowModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles key events, step results, and spinner ticks.
func (m *PaymentFlowModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.running {
			switch msg.String() {
			case "n":
				return m.executeCurrentStep()
			case "p":
				if m.sm.flow.CurrentStep > 0 {
					m.sm.prev()
				}
			}
		}
	case stepResultMsg:
		m.running = false
		m.results[msg.step] = msg.data
		m.sm.markStepDone()
		return nil
	case stepErrorMsg:
		m.running = false
		m.errors[msg.step] = msg.err.Error()
		m.sm.markStepError()
		return nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return cmd
	}
	return nil
}

func (m *PaymentFlowModel) executeCurrentStep() tea.Cmd {
	if m.executor == nil {
		return nil
	}
	step := m.sm.flow.CurrentStep
	if step >= m.sm.flow.TotalSteps {
		return nil
	}
	m.running = true
	m.sm.markStepRunning()
	executor := m.executor
	return func() tea.Msg {
		result, err := executor.RunStep(context.Background(), step)
		if err != nil {
			return stepErrorMsg{step: step, err: err}
		}
		return stepResultMsg{step: step, data: result}
	}
}

// View renders the flow panels and current step detail.
func (m *PaymentFlowModel) View() string {
	view := m.sm.view(m.width)

	step := m.sm.flow.CurrentStep
	if step >= m.sm.flow.TotalSteps {
		step = m.sm.flow.TotalSteps - 1
	}

	var detail string
	if m.running {
		detail = m.spinner.View() + " Executing..."
	} else if m.errors[step] != "" {
		detail = lipgloss.NewStyle().Foreground(tui.ColorError).Render("Error: " + m.errors[step])
	} else if m.results[step] != "" {
		detail = lipgloss.NewStyle().Width(m.width - 4).Render(m.results[step])
	} else if m.execErr != "" {
		detail = lipgloss.NewStyle().Foreground(tui.ColorError).Render(m.execErr)
	} else {
		detail = lipgloss.NewStyle().Foreground(tui.ColorMuted).Render("Press n to execute this step")
	}

	return lipgloss.JoinVertical(lipgloss.Left, view, "", detail)
}
