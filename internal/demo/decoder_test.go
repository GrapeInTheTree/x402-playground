package demo

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestDecodePaymentRequiredHeader(t *testing.T) {
	payload := map[string]interface{}{
		"accepts": []map[string]interface{}{
			{
				"scheme":  "exact",
				"network": "eip155:84532",
				"payTo":   "0x1234",
				"asset":   "0x036CbD53842c5426634e7929541eC2318f3dCF7e",
				"maxAmountRequired": "100000",
			},
		},
	}
	raw, _ := json.Marshal(payload)
	encoded := base64.StdEncoding.EncodeToString(raw)

	pr, decoded, err := DecodePaymentRequiredHeader(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded == nil {
		t.Fatal("expected decoded bytes")
	}
	if pr == nil {
		t.Fatal("expected parsed result")
	}
	if len(pr.Accepts) != 1 {
		t.Fatalf("expected 1 accept item, got %d", len(pr.Accepts))
	}
}

func TestDecodePaymentRequiredHeader_InvalidBase64(t *testing.T) {
	_, _, err := DecodePaymentRequiredHeader("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecodeBase64JSON(t *testing.T) {
	original := `{"key":"value"}`
	encoded := base64.StdEncoding.EncodeToString([]byte(original))

	result, err := DecodeBase64JSON(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m map[string]string
	if err := json.Unmarshal(result, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if m["key"] != "value" {
		t.Errorf("expected 'value', got %q", m["key"])
	}
}

func TestFormatJSON(t *testing.T) {
	input := `{"a":1,"b":"hello"}`
	result := FormatJSON([]byte(input))

	var v map[string]interface{}
	if err := json.Unmarshal([]byte(result), &v); err != nil {
		t.Fatalf("FormatJSON produced invalid JSON: %v", err)
	}
}

func TestFormatJSON_InvalidJSON(t *testing.T) {
	input := "not json"
	result := FormatJSON([]byte(input))
	if result != input {
		t.Errorf("expected passthrough for invalid JSON, got %q", result)
	}
}

func TestParseAcceptItem(t *testing.T) {
	raw := json.RawMessage(`{
		"scheme": "exact",
		"network": "eip155:84532",
		"maxAmountRequired": "100000",
		"payTo": "0xABCD",
		"asset": "0x036CbD53842c5426634e7929541eC2318f3dCF7e",
		"extra": {"name": "USDC", "version": "2"}
	}`)

	item, err := ParseAcceptItem(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Scheme != "exact" {
		t.Errorf("expected scheme 'exact', got %q", item.Scheme)
	}
	if item.Network != "eip155:84532" {
		t.Errorf("expected network 'eip155:84532', got %q", item.Network)
	}
	if item.PayTo != "0xABCD" {
		t.Errorf("expected payTo '0xABCD', got %q", item.PayTo)
	}
	if item.Extra["name"] != "USDC" {
		t.Errorf("expected extra.name 'USDC', got %v", item.Extra["name"])
	}
}

func TestStepDescription(t *testing.T) {
	for i := 1; i <= 10; i++ {
		desc := StepDescription(i)
		if desc == "" {
			t.Errorf("step %d has empty description", i)
		}
	}
	// Unknown step
	if desc := StepDescription(99); desc != "" {
		t.Errorf("expected empty description for unknown step, got %q", desc)
	}
}
