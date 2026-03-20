package demo

import (
	"testing"
)

func TestWalletInfo(t *testing.T) {
	w := WalletInfo{Name: "Client", Address: "0x1234567890abcdef"}
	if w.Name != "Client" {
		t.Errorf("expected name 'Client', got %q", w.Name)
	}
	if w.Address != "0x1234567890abcdef" {
		t.Errorf("expected address '0x1234567890abcdef', got %q", w.Address)
	}
}

func TestNewFlowState(t *testing.T) {
	fs := NewFlowState("eip3009")
	if fs.TransferMethod != "eip3009" {
		t.Errorf("expected transfer method 'eip3009', got %q", fs.TransferMethod)
	}
	if fs.TotalSteps != 10 {
		t.Errorf("expected 10 total steps, got %d", fs.TotalSteps)
	}
	if fs.CurrentStep != 0 {
		t.Errorf("expected current step 0, got %d", fs.CurrentStep)
	}

	fs2 := NewFlowState("permit2")
	if fs2.TransferMethod != "permit2" {
		t.Errorf("expected transfer method 'permit2', got %q", fs2.TransferMethod)
	}
}
