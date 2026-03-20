package practice

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/GrapeInTheTree/x402-demo/internal/config"
	"github.com/GrapeInTheTree/x402-demo/internal/demo"
)

// EIP3009FlowModel manages the EIP-3009 10-step payment flow.
type EIP3009FlowModel struct {
	flow   *demo.FlowState
	steps  [10]stepState
	width  int
	height int
	cfg    *config.ExplorerConfig
}

type stepState struct {
	client      StepInfo
	resource    StepInfo
	facilitator StepInfo
}

func NewEIP3009FlowModel(width, height int, cfg *config.ExplorerConfig) *EIP3009FlowModel {
	m := &EIP3009FlowModel{
		flow:   demo.NewFlowState("eip3009"),
		width:  width,
		height: height,
		cfg:    cfg,
	}
	m.initSteps()
	return m
}

func (m *EIP3009FlowModel) initSteps() {
	descriptions := []struct{ client, resource, facilitator string }{
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

	for i, d := range descriptions {
		m.steps[i] = stepState{
			client:      StepInfo{Actor: "Client", Action: d.client, Status: "pending"},
			resource:    StepInfo{Actor: "Resource", Action: d.resource, Status: "pending"},
			facilitator: StepInfo{Actor: "Facilitator", Action: d.facilitator, Status: "pending"},
		}
	}

	// Mark first step as running
	if len(descriptions) > 0 {
		m.steps[0].client.Status = "running"
	}
}

func (m *EIP3009FlowModel) Update(msg tea.Msg) tea.Cmd {
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

func (m *EIP3009FlowModel) nextStep() {
	if m.flow.CurrentStep >= m.flow.TotalSteps-1 {
		return
	}

	// Mark current step as done
	step := &m.steps[m.flow.CurrentStep]
	step.client.Status = "done"
	step.resource.Status = "done"
	step.facilitator.Status = "done"

	m.flow.CurrentStep++

	// Mark new step as running
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

func (m *EIP3009FlowModel) prevStep() {
	if m.flow.CurrentStep <= 0 {
		return
	}

	// Mark current step as pending
	step := &m.steps[m.flow.CurrentStep]
	step.client.Status = "pending"
	step.resource.Status = "pending"
	step.facilitator.Status = "pending"

	m.flow.CurrentStep--

	// Mark prev step as running
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

func (m *EIP3009FlowModel) View() string {
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
