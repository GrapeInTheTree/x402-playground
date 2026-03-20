package demo

import (
	"encoding/base64"
	"encoding/json"
)

// DecodePaymentRequiredHeader decodes a base64-encoded PAYMENT-REQUIRED header.
func DecodePaymentRequiredHeader(header string) (*DecodedPaymentRequired, []byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return nil, nil, err
	}

	var pr DecodedPaymentRequired
	if err := json.Unmarshal(decoded, &pr); err != nil {
		return nil, decoded, err
	}

	return &pr, decoded, nil
}

// DecodeBase64JSON decodes a base64-encoded JSON string and returns the raw JSON.
func DecodeBase64JSON(encoded string) (json.RawMessage, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(decoded), nil
}

// FormatJSON pretty-prints JSON bytes.
func FormatJSON(data []byte) string {
	var v interface{}
	if json.Unmarshal(data, &v) != nil {
		return string(data)
	}
	formatted, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(data)
	}
	return string(formatted)
}

// ParseAcceptItem parses a single accept item from the PAYMENT-REQUIRED header.
func ParseAcceptItem(raw json.RawMessage) (*AcceptItem, error) {
	var item AcceptItem
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, err
	}
	return &item, nil
}
