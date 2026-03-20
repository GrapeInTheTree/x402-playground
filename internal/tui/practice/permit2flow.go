package practice

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/demo"
)

// Permit2FlowModel manages the Permit2 10-step payment flow.
type Permit2FlowModel struct {
	flow   *demo.FlowState
	steps  [10]stepState
	width  int
	height int
	cfg    *config.ExplorerConfig
}

func NewPermit2FlowModel(width, height int, cfg *config.ExplorerConfig) *Permit2FlowModel {
	m := &Permit2FlowModel{
		flow:   demo.NewFlowState("permit2"),
		width:  width,
		height: height,
		cfg:    cfg,
	}
	m.initSteps()
	return m
}

func (m *Permit2FlowModel) initSteps() {
	descriptions := []struct{ client, resource, facilitator string }{
		{"지갑 주소 + Permit2 approve 확인", "—", "—"},
		{"—", "GET /supported 호출", "/supported 응답"},
		{"GET /weather (결제 없음)", "402 + assetTransferMethod:permit2", "—"},
		{"PAYMENT-REQUIRED 디코딩 (Permit2)", "—", "—"},
		{"Permit2 EIP-712 서명 생성", "—", "—"},
		{"PAYMENT-SIGNATURE 전송", "헤더 파싱 → /verify", "—"},
		{"—", "/verify 요청 전달", "Permit2 서명 + allowance 검증"},
		{"200 + 데이터 수신", "데이터 반환 + /settle", "—"},
		{"—", "—", "x402Permit2Proxy.settle() 제출"},
		{"최종 잔액 확인", "—", "—"},
	}

	for i, d := range descriptions {
		m.steps[i] = stepState{
			client:      StepInfo{Actor: "Client", Action: d.client, Status: "pending"},
			resource:    StepInfo{Actor: "Resource", Action: d.resource, Status: "pending"},
			facilitator: StepInfo{Actor: "Facilitator", Action: d.facilitator, Status: "pending"},
		}
	}

	if len(descriptions) > 0 {
		m.steps[0].client.Status = "running"
	}
}

func (m *Permit2FlowModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "n":
			m.nextStep()
		case "p":
			m.prevStep()
		}
	}
	return nil
}

func (m *Permit2FlowModel) nextStep() {
	if m.flow.CurrentStep >= m.flow.TotalSteps-1 {
		return
	}

	step := &m.steps[m.flow.CurrentStep]
	step.client.Status = "done"
	step.resource.Status = "done"
	step.facilitator.Status = "done"

	m.flow.CurrentStep++

	newStep := &m.steps[m.flow.CurrentStep]
	if newStep.client.Action != "—" {
		newStep.client.Status = "running"
	}
	if newStep.resource.Action != "—" {
		newStep.resource.Status = "running"
	}
	if newStep.facilitator.Action != "—" {
		newStep.facilitator.Status = "running"
	}
}

func (m *Permit2FlowModel) prevStep() {
	if m.flow.CurrentStep <= 0 {
		return
	}

	step := &m.steps[m.flow.CurrentStep]
	step.client.Status = "pending"
	step.resource.Status = "pending"
	step.facilitator.Status = "pending"

	m.flow.CurrentStep--

	prevStep := &m.steps[m.flow.CurrentStep]
	if prevStep.client.Action != "—" {
		prevStep.client.Status = "running"
	} else {
		prevStep.client.Status = "done"
	}
	if prevStep.resource.Action != "—" {
		prevStep.resource.Status = "running"
	} else {
		prevStep.resource.Status = "done"
	}
	if prevStep.facilitator.Action != "—" {
		prevStep.facilitator.Status = "running"
	} else {
		prevStep.facilitator.Status = "done"
	}
}

func (m *Permit2FlowModel) View() string {
	clientSteps := make([]StepInfo, 10)
	resourceSteps := make([]StepInfo, 10)
	facilitatorSteps := make([]StepInfo, 10)

	for i, s := range m.steps {
		clientSteps[i] = s.client
		resourceSteps[i] = s.resource
		facilitatorSteps[i] = s.facilitator
	}

	return RenderFlowPanels(
		m.flow.CurrentStep, m.flow.TotalSteps,
		clientSteps, resourceSteps, facilitatorSteps,
		m.width,
	)
}
