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

// EIP3009FlowModel manages the EIP-3009 10-step payment flow with live execution.
type EIP3009FlowModel struct {
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

// NewEIP3009FlowModel creates a new EIP-3009 flow model with live executor.
func NewEIP3009FlowModel(width, height int, cfg *config.ExplorerConfig) *EIP3009FlowModel {
	descriptions := []stepDesc{
		{"지갑 주소 확인", "—", "—"},
		{"—", "GET /supported 호출", "/supported 응답"},
		{"GET /weather (결제 없음)", "402 반환", "—"},
		{"PAYMENT-REQUIRED 디코딩", "—", "—"},
		{"EIP-712 서명 생성", "—", "—"},
		{"PAYMENT-SIGNATURE 전송", "헤더 파싱 → /verify", "—"},
		{"—", "/verify 요청 전달", "서명/잔액/시뮬레이션 검증"},
		{"200 + 데이터 수신", "데이터 반환 + /settle", "—"},
		{"—", "—", "온체인 트랜잭션 제출"},
		{"최종 잔액 확인", "—", "—"},
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(tui.ColorAccent)

	m := &EIP3009FlowModel{
		sm:      newStepManager(demo.NewFlowState("eip3009"), descriptions),
		spinner: s,
		width:   width,
		height:  height,
		cfg:     cfg,
	}

	// Try to create live executor
	if cfg != nil {
		exec, err := demo.NewLiveExecutor(
			cfg.FacilitatorURL, cfg.ResourceURL, cfg.RPCURL,
			cfg.USDCAddress, cfg.PayToAddress, cfg.ClientPrivateKey, "eip3009",
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
func (m *EIP3009FlowModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles key events, step results, and spinner ticks.
func (m *EIP3009FlowModel) Update(msg tea.Msg) tea.Cmd {
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

func (m *EIP3009FlowModel) executeCurrentStep() tea.Cmd {
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

// View renders the EIP-3009 flow panels and current step detail.
func (m *EIP3009FlowModel) View() string {
	view := m.sm.view(m.width)

	// Show current step result or spinner
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
