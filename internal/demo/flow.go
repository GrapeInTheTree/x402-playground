package demo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FlowExecutor runs the 10-step payment flow.
type FlowExecutor struct {
	FacilitatorURL string
	ResourceURL    string
	EndpointPath   string
	State          *FlowState
}

// StepResult captures the output of a single step execution.
type StepResult struct {
	Step        int
	Description string
	Data        json.RawMessage // key data produced by this step
	Error       error
}

// Step2_FacilitatorSupported calls GET /supported on the facilitator.
func (f *FlowExecutor) Step2_FacilitatorSupported(ctx context.Context) (*StepResult, error) {
	resp, err := http.Get(f.FacilitatorURL + "/supported")
	if err != nil {
		return nil, fmt.Errorf("GET /supported: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return &StepResult{Step: 2, Description: StepDescription(2), Data: body}, nil
}

// Step3_NaiveAPICall calls the resource server without payment headers.
func (f *FlowExecutor) Step3_NaiveAPICall(ctx context.Context) (*http.Response, []byte, error) {
	targetURL := f.ResourceURL + f.EndpointPath
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, nil, fmt.Errorf("GET %s: %w", targetURL, err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, body, nil
}

// Step7_Verify calls the facilitator /verify endpoint directly.
func (f *FlowExecutor) Step7_Verify(ctx context.Context, payloadBytes, requirementsBytes []byte) (json.RawMessage, error) {
	verifyBody := map[string]interface{}{
		"x402Version":         2,
		"paymentPayload":      json.RawMessage(payloadBytes),
		"paymentRequirements": json.RawMessage(requirementsBytes),
	}
	verifyJSON, err := json.Marshal(verifyBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(f.FacilitatorURL+"/verify", "application/json", bytes.NewReader(verifyJSON))
	if err != nil {
		return nil, fmt.Errorf("POST /verify: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

// Step8_PaidRequest makes the actual request with payment headers.
func (f *FlowExecutor) Step8_PaidRequest(ctx context.Context, headerName, headerValue string) (*http.Response, []byte, error) {
	targetURL := f.ResourceURL + f.EndpointPath
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set(headerName, headerValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("GET %s (paid): %w", targetURL, err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, body, nil
}
