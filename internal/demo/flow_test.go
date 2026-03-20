package demo

import (
	"testing"
)

func TestFlowExecutor_Init(t *testing.T) {
	f := &FlowExecutor{
		FacilitatorURL: "http://localhost:4022",
		ResourceURL:    "http://localhost:4021",
		EndpointPath:   "/weather",
		State:          NewFlowState("eip3009"),
	}

	if f.FacilitatorURL != "http://localhost:4022" {
		t.Errorf("unexpected FacilitatorURL: %s", f.FacilitatorURL)
	}
	if f.ResourceURL != "http://localhost:4021" {
		t.Errorf("unexpected ResourceURL: %s", f.ResourceURL)
	}
	if f.EndpointPath != "/weather" {
		t.Errorf("unexpected EndpointPath: %s", f.EndpointPath)
	}
	if f.State.TransferMethod != "eip3009" {
		t.Errorf("unexpected TransferMethod: %s", f.State.TransferMethod)
	}
}

func TestStepResult(t *testing.T) {
	result := &StepResult{
		Step:        2,
		Description: StepDescription(2),
		Data:        []byte(`{"test": true}`),
	}

	if result.Step != 2 {
		t.Errorf("expected step 2, got %d", result.Step)
	}
	if result.Description == "" {
		t.Error("expected non-empty description")
	}
	if result.Error != nil {
		t.Errorf("expected nil error, got %v", result.Error)
	}
}
