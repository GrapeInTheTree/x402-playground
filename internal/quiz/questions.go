package quiz

// AllQuestions returns all quiz questions organized by progressive difficulty.
func AllQuestions() []Question {
	var all []Question
	all = append(all, level1Basics()...)
	all = append(all, level2Standards()...)
	all = append(all, level3Protocol()...)
	all = append(all, level4Advanced()...)
	all = append(all, level5Agents()...)
	all = append(all, SolidityQuestions()...)
	return all
}

// ============================================================
// LEVEL 1: Basics — Go fundamentals for blockchain
// ============================================================

func level1Basics() []Question {
	return []Question{
		{
			ID: "hex-address", Title: "Hex Address Validation",
			Difficulty: "easy", Category: "Basics",
			Description: `Ethereum addresses are 40 hex characters prefixed with "0x".
Write a function to validate and normalize an Ethereum address.`,
			Template: `package x402quiz

import (
	"fmt"
	"strings"
)

// IsValidAddress checks if the given string is a valid Ethereum address.
// Valid: starts with "0x" and has exactly 40 hex characters after prefix.
func IsValidAddress(addr string) bool {
	// TODO: Check prefix and length
	// TODO: Check all characters after "0x" are hex (0-9, a-f, A-F)
	_ = strings.HasPrefix
	return false
}

// NormalizeAddress lowercases an Ethereum address for comparison.
// Example: "0xABCD1234..." → "0xabcd1234..."
func NormalizeAddress(addr string) string {
	// TODO: Convert to lowercase
	_ = strings.ToLower
	_ = fmt.Sprintf
	return ""
}
`,
			TestCode: `package x402quiz

import "testing"

func TestIsValidAddress(t *testing.T) {
	valid := []string{
		"0x036CbD53842c5426634e7929541eC2318f3dCF7e",
		"0x0000000000000000000000000000000000000000",
		"0xABCDEF1234567890abcdef1234567890ABCDEF12",
	}
	for _, a := range valid {
		if !IsValidAddress(a) {
			t.Errorf("expected valid: %s", a)
		}
	}
	invalid := []string{"", "0x", "0xGGGG", "036CbD53842c5426634e7929541eC2318f3dCF7e", "0x036CbD"}
	for _, a := range invalid {
		if IsValidAddress(a) {
			t.Errorf("expected invalid: %s", a)
		}
	}
}

func TestNormalizeAddress(t *testing.T) {
	got := NormalizeAddress("0xABCD1234EF567890abcd1234ef567890ABCD1234")
	want := "0xabcd1234ef567890abcd1234ef567890abcd1234"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
`,
			Hints: []string{
				`Check strings.HasPrefix(addr, "0x") and len(addr) == 42`,
				"Hex characters: 0-9, a-f, A-F",
				"strings.ToLower() converts the whole string",
			},
		},
		{
			ID: "usdc-amount", Title: "USDC Amount Conversion",
			Difficulty: "easy", Category: "ERC-20",
			Description: `USDC uses 6 decimal places. In smart contracts, amounts are in
the smallest unit. $0.10 = 100,000 units. $1.00 = 1,000,000 units.

Write conversion functions between dollar amounts and USDC units.`,
			Template: `package x402quiz

import "fmt"

// DollarsToUSDC converts a dollar amount to USDC smallest units (6 decimals).
// Example: 0.10 → 100000, 1.00 → 1000000
func DollarsToUSDC(dollars float64) uint64 {
	// TODO: Multiply by 10^6
	return 0
}

// USDCToDollars converts USDC units back to dollars.
// Example: 100000 → 0.10
func USDCToDollars(units uint64) float64 {
	// TODO: Divide by 10^6
	return 0.0
}

// FormatUSDC returns "$X.XX" for the given USDC units.
// Example: 100000 → "$0.10", 1500000 → "$1.50"
func FormatUSDC(units uint64) string {
	// TODO: Convert and format
	_ = fmt.Sprintf
	return ""
}
`,
			TestCode: `package x402quiz

import "testing"

func TestDollarsToUSDC(t *testing.T) {
	tests := []struct{ in float64; want uint64 }{
		{0.10, 100000}, {1.00, 1000000}, {0.01, 10000}, {100.0, 100000000},
	}
	for _, tt := range tests {
		if got := DollarsToUSDC(tt.in); got != tt.want {
			t.Errorf("DollarsToUSDC(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestUSDCToDollars(t *testing.T) {
	tests := []struct{ in uint64; want float64 }{
		{100000, 0.10}, {1000000, 1.00}, {10000, 0.01},
	}
	for _, tt := range tests {
		got := USDCToDollars(tt.in)
		if diff := got - tt.want; diff < -0.001 || diff > 0.001 {
			t.Errorf("USDCToDollars(%d) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestFormatUSDC(t *testing.T) {
	tests := []struct{ in uint64; want string }{
		{100000, "$0.10"}, {1000000, "$1.00"}, {1500000, "$1.50"},
	}
	for _, tt := range tests {
		if got := FormatUSDC(tt.in); got != tt.want {
			t.Errorf("FormatUSDC(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
`,
			Hints: []string{
				"1 USDC = 1,000,000 units (6 decimals)",
				"uint64(dollars * 1_000_000)",
				`fmt.Sprintf("$%.2f", float64(units)/1_000_000)`,
			},
		},
		{
			ID: "base64-json", Title: "Base64 + JSON Decode",
			Difficulty: "easy", Category: "Basics",
			Description: `x402 protocol headers use base64-encoded JSON. This is the most
fundamental operation: decode base64 → parse JSON → extract fields.

Write a generic decode function.`,
			Template: `package x402quiz

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// DecodeBase64JSON decodes a base64 string and unmarshals the JSON
// into the provided target (pointer to struct or map).
// Returns an error if base64 decoding or JSON parsing fails.
func DecodeBase64JSON(encoded string, target any) error {
	// TODO: 1. base64.StdEncoding.DecodeString(encoded)
	// TODO: 2. json.Unmarshal into target
	_ = base64.StdEncoding.DecodeString
	_ = json.Unmarshal
	_ = fmt.Errorf
	return nil
}

// EncodeToBase64JSON marshals the value to JSON and base64 encodes it.
func EncodeToBase64JSON(value any) (string, error) {
	// TODO: 1. json.Marshal(value)
	// TODO: 2. base64.StdEncoding.EncodeToString
	_ = json.Marshal
	_ = base64.StdEncoding.EncodeToString
	return "", nil
}
`,
			TestCode: `package x402quiz

import "testing"

func TestDecodeBase64JSON(t *testing.T) {
	encoded := "eyJuYW1lIjoiVVNEQyIsInZhbHVlIjoxMDB9"
	var result struct {
		Name  string ` + "`json:\"name\"`" + `
		Value int    ` + "`json:\"value\"`" + `
	}
	if err := DecodeBase64JSON(encoded, &result); err != nil {
		t.Fatal(err)
	}
	if result.Name != "USDC" || result.Value != 100 {
		t.Errorf("got %+v", result)
	}
}

func TestDecodeBase64JSON_Invalid(t *testing.T) {
	var m map[string]any
	if err := DecodeBase64JSON("not-valid!!!", &m); err == nil {
		t.Error("expected error")
	}
}

func TestEncodeToBase64JSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	encoded, err := EncodeToBase64JSON(data)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]string
	if err := DecodeBase64JSON(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["key"] != "value" {
		t.Errorf("roundtrip failed: %v", decoded)
	}
}
`,
			Hints: []string{
				"base64.StdEncoding.DecodeString returns ([]byte, error)",
				"json.Unmarshal(bytes, target) — target is already a pointer",
				"json.Marshal returns ([]byte, error), then base64.StdEncoding.EncodeToString(bytes)",
			},
		},

		// ============================================================
		// ERC-20 ABI encoding
		// ============================================================
		{
			ID: "erc20-abi", Title: "ERC-20 Function Selectors",
			Difficulty: "easy", Category: "ERC-20",
			Description: `In Ethereum, function calls are encoded using the first 4 bytes
of the keccak256 hash of the function signature.

For ERC-20, the key functions are:
  balanceOf(address) → 0x70a08231
  transfer(address,uint256) → 0xa9059cbb
  approve(address,uint256) → 0x095ea7b3
  allowance(address,address) → 0xdd62ed3e

Write a function that returns the correct selector for each function name.`,
			Template: `package x402quiz

// ERC20Selector returns the 4-byte function selector (as hex string)
// for common ERC-20 functions.
//
// Known selectors:
//   "balanceOf"  → "0x70a08231"
//   "transfer"   → "0xa9059cbb"
//   "approve"    → "0x095ea7b3"
//   "allowance"  → "0xdd62ed3e"
//   "transferFrom" → "0x23b872dd"
func ERC20Selector(functionName string) string {
	// TODO: Return the correct selector for each function
	// Hint: use a map or switch
	return ""
}

// IsERC20Function checks if the given 4-byte hex selector
// belongs to a standard ERC-20 function.
func IsERC20Function(selector string) bool {
	// TODO: Check if selector matches any known ERC-20 function
	return false
}
`,
			TestCode: `package x402quiz

import "testing"

func TestERC20Selector(t *testing.T) {
	tests := map[string]string{
		"balanceOf": "0x70a08231", "transfer": "0xa9059cbb",
		"approve": "0x095ea7b3", "allowance": "0xdd62ed3e",
		"transferFrom": "0x23b872dd",
	}
	for name, want := range tests {
		if got := ERC20Selector(name); got != want {
			t.Errorf("ERC20Selector(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestERC20Selector_Unknown(t *testing.T) {
	if got := ERC20Selector("unknown"); got != "" {
		t.Errorf("expected empty for unknown, got %q", got)
	}
}

func TestIsERC20Function(t *testing.T) {
	if !IsERC20Function("0x70a08231") { t.Error("balanceOf should be ERC20") }
	if !IsERC20Function("0xa9059cbb") { t.Error("transfer should be ERC20") }
	if IsERC20Function("0x12345678")  { t.Error("random should not be ERC20") }
}
`,
			Hints: []string{
				"Use a switch statement or map[string]string",
				`selectors := map[string]string{"balanceOf": "0x70a08231", ...}`,
				"For IsERC20Function, check if the selector exists in the known set",
			},
		},
	}
}

// ============================================================
// LEVEL 2: Standards — EIP-712, EIP-2612, EIP-3009
// ============================================================

func level2Standards() []Question {
	return []Question{
		{
			ID: "eip712-domain", Title: "EIP-712 Domain Separator",
			Difficulty: "medium", Category: "EIP-712",
			Description: `EIP-712 requires a domain separator to prevent signature replay.
For USDC on Base Sepolia:
  - name must match token's name() return: "USDC" (NOT "USD Coin")
  - version: "2" (FiatTokenV2)
  - chainId: 84532 (Base Sepolia)

For Permit2 (same address on all chains via CREATE2):
  - address: 0x000000000022D473030F116dDEE9F6B43aC78BA3`,
			Template: `package x402quiz

// EIP712Domain represents the domain separator fields.
type EIP712Domain struct {
	Name              string
	Version           string
	ChainID           uint64
	VerifyingContract string
}

// USDCDomain returns the EIP-712 domain for USDC on Base Sepolia.
func USDCDomain() EIP712Domain {
	return EIP712Domain{
		// TODO: Fill in correct values
		// IMPORTANT: Name must match token contract's name() exactly!
		Name:              "",
		Version:           "",
		ChainID:           0,
		VerifyingContract: "",
	}
}

// Permit2Domain returns the EIP-712 domain for Permit2 on Base Sepolia.
func Permit2Domain() EIP712Domain {
	return EIP712Domain{
		// TODO: Fill in correct values
		Name:              "",
		Version:           "",
		ChainID:           0,
		VerifyingContract: "",
	}
}
`,
			TestCode: `package x402quiz

import "testing"

func TestUSDCDomain(t *testing.T) {
	d := USDCDomain()
	if d.Name != "USDC" { t.Errorf("Name = %q, want \"USDC\"", d.Name) }
	if d.Version != "2" { t.Errorf("Version = %q, want \"2\"", d.Version) }
	if d.ChainID != 84532 { t.Errorf("ChainID = %d, want 84532", d.ChainID) }
	if d.VerifyingContract != "0x036CbD53842c5426634e7929541eC2318f3dCF7e" {
		t.Errorf("Contract = %q", d.VerifyingContract)
	}
}

func TestPermit2Domain(t *testing.T) {
	d := Permit2Domain()
	if d.Name != "Permit2" { t.Errorf("Name = %q, want \"Permit2\"", d.Name) }
	if d.ChainID != 84532 { t.Errorf("ChainID = %d, want 84532", d.ChainID) }
	if d.VerifyingContract != "0x000000000022D473030F116dDEE9F6B43aC78BA3" {
		t.Errorf("Contract = %q", d.VerifyingContract)
	}
}
`,
			Hints: []string{
				`Base Sepolia USDC returns "USDC" from name(), not "USD Coin"`,
				"FiatTokenV2 uses version \"2\"",
				"Permit2 uses CREATE2: same address on all EVM chains",
			},
		},
		{
			ID: "eip712-typehash", Title: "EIP-712 Type Hash Construction",
			Difficulty: "medium", Category: "EIP-712",
			Description: `In EIP-712, each struct type has a "type hash" computed as:
  keccak256("TypeName(type1 name1,type2 name2,...)")

For EIP-3009's TransferWithAuthorization:
  "TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"

Write a function that builds the type string from field definitions.`,
			Template: `package x402quiz

import "strings"

// TypeField represents a single field in an EIP-712 type.
type TypeField struct {
	Type string // e.g., "address", "uint256", "bytes32"
	Name string // e.g., "from", "to", "value"
}

// BuildTypeString constructs the EIP-712 type encoding string.
// Example: BuildTypeString("Transfer", fields) → "Transfer(address from,address to,uint256 value)"
func BuildTypeString(typeName string, fields []TypeField) string {
	// TODO: Build "TypeName(type1 name1,type2 name2,...)"
	_ = strings.Join
	return ""
}
`,
			TestCode: `package x402quiz

import "testing"

func TestBuildTypeString_TransferWithAuth(t *testing.T) {
	fields := []TypeField{
		{"address", "from"}, {"address", "to"}, {"uint256", "value"},
		{"uint256", "validAfter"}, {"uint256", "validBefore"}, {"bytes32", "nonce"},
	}
	got := BuildTypeString("TransferWithAuthorization", fields)
	want := "TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"
	if got != want {
		t.Errorf("got:\n  %s\nwant:\n  %s", got, want)
	}
}

func TestBuildTypeString_Simple(t *testing.T) {
	fields := []TypeField{{"address", "owner"}, {"uint256", "amount"}}
	got := BuildTypeString("Approve", fields)
	if got != "Approve(address owner,uint256 amount)" {
		t.Errorf("got %q", got)
	}
}

func TestBuildTypeString_Empty(t *testing.T) {
	got := BuildTypeString("Empty", nil)
	if got != "Empty()" {
		t.Errorf("got %q, want \"Empty()\"", got)
	}
}
`,
			Hints: []string{
				"Build each field as \"type name\", then join with \",\"",
				`parts := make([]string, len(fields)); for i, f := range fields { parts[i] = f.Type + " " + f.Name }`,
				`return typeName + "(" + strings.Join(parts, ",") + ")"`,
			},
		},
		{
			ID: "eip3009-fields", Title: "EIP-3009 Authorization Fields",
			Difficulty: "medium", Category: "EIP-3009",
			Description: `EIP-3009 transferWithAuthorization requires specific fields:
  from, to, value, validAfter, validBefore, nonce

Write a struct and constructor that validates these fields.
- from and to must be valid addresses
- value must be > 0
- validBefore must be > validAfter
- nonce must be 32 bytes (64 hex chars + "0x" prefix)`,
			Template: `package x402quiz

import (
	"fmt"
	"strings"
)

// TransferAuth holds the fields for EIP-3009 transferWithAuthorization.
type TransferAuth struct {
	From        string
	To          string
	Value       uint64
	ValidAfter  uint64
	ValidBefore uint64
	Nonce       string // "0x" + 64 hex chars
}

// NewTransferAuth creates and validates a TransferAuth.
// Returns an error if any field is invalid.
func NewTransferAuth(from, to string, value, validAfter, validBefore uint64, nonce string) (*TransferAuth, error) {
	// TODO: Validate all fields
	// 1. from and to must start with "0x" and be 42 chars
	// 2. value must be > 0
	// 3. validBefore must be > validAfter
	// 4. nonce must be "0x" + 64 hex chars (66 total)
	_ = fmt.Errorf
	_ = strings.HasPrefix
	return nil, nil
}
`,
			TestCode: `package x402quiz

import "testing"

func TestNewTransferAuth_Valid(t *testing.T) {
	auth, err := NewTransferAuth(
		"0x036CbD53842c5426634e7929541eC2318f3dCF7e",
		"0x000000000022D473030F116dDEE9F6B43aC78BA3",
		100000, 0, 1718000000,
		"0x"+strings.Repeat("ab", 32),
	)
	if err != nil { t.Fatal(err) }
	if auth.Value != 100000 { t.Errorf("Value = %d", auth.Value) }
}

func TestNewTransferAuth_ZeroValue(t *testing.T) {
	_, err := NewTransferAuth("0x"+strings.Repeat("a", 40), "0x"+strings.Repeat("b", 40), 0, 0, 100, "0x"+strings.Repeat("c", 64))
	if err == nil { t.Error("expected error for zero value") }
}

func TestNewTransferAuth_BadTime(t *testing.T) {
	_, err := NewTransferAuth("0x"+strings.Repeat("a", 40), "0x"+strings.Repeat("b", 40), 100, 200, 100, "0x"+strings.Repeat("c", 64))
	if err == nil { t.Error("expected error: validBefore <= validAfter") }
}

func TestNewTransferAuth_BadNonce(t *testing.T) {
	_, err := NewTransferAuth("0x"+strings.Repeat("a", 40), "0x"+strings.Repeat("b", 40), 100, 0, 100, "0xshort")
	if err == nil { t.Error("expected error for short nonce") }
}

import "strings"
`,
			Hints: []string{
				"Check len(from) == 42 && strings.HasPrefix(from, \"0x\")",
				"if value == 0 { return nil, fmt.Errorf(\"value must be > 0\") }",
				"Nonce: len(nonce) == 66 && strings.HasPrefix(nonce, \"0x\")",
			},
		},
		{
			ID: "eip2612-permit", Title: "EIP-2612 Permit Concept",
			Difficulty: "medium", Category: "EIP-2612",
			Description: `EIP-2612 "permit" allows gasless token approvals via signatures.
Instead of calling approve() on-chain, the owner signs a message
off-chain, and anyone can submit it.

The permit message contains: owner, spender, value, nonce, deadline.

Write a function to build a permit message and check if it's expired.`,
			Template: `package x402quiz

import (
	"fmt"
	"time"
)

// PermitMessage represents an EIP-2612 permit.
type PermitMessage struct {
	Owner    string
	Spender  string
	Value    uint64
	Nonce    uint64
	Deadline uint64 // Unix timestamp
}

// NewPermit creates a new permit message.
// Deadline is set to the given duration from now.
func NewPermit(owner, spender string, value, nonce uint64, validFor time.Duration) *PermitMessage {
	// TODO: Create permit with deadline = now + validFor
	_ = fmt.Sprintf
	_ = time.Now
	return nil
}

// IsExpired returns true if the permit's deadline has passed.
func (p *PermitMessage) IsExpired() bool {
	// TODO: Compare deadline with current time
	return false
}

// String returns a human-readable representation.
func (p *PermitMessage) String() string {
	// TODO: Format as "permit(owner→spender, value, deadline: <time>)"
	return ""
}
`,
			TestCode: `package x402quiz

import (
	"strings"
	"testing"
	"time"
)

func TestNewPermit(t *testing.T) {
	p := NewPermit("0xOwner", "0xSpender", 100000, 0, 1*time.Hour)
	if p == nil { t.Fatal("got nil") }
	if p.Owner != "0xOwner" { t.Errorf("Owner = %q", p.Owner) }
	if p.Value != 100000 { t.Errorf("Value = %d", p.Value) }
	if p.Deadline <= uint64(time.Now().Unix()) { t.Error("deadline should be in the future") }
}

func TestPermit_IsExpired(t *testing.T) {
	p := &PermitMessage{Deadline: uint64(time.Now().Unix() - 100)}
	if !p.IsExpired() { t.Error("should be expired") }

	p2 := &PermitMessage{Deadline: uint64(time.Now().Unix() + 3600)}
	if p2.IsExpired() { t.Error("should not be expired") }
}

func TestPermit_String(t *testing.T) {
	p := &PermitMessage{Owner: "0xA", Spender: "0xB", Value: 100}
	s := p.String()
	if !strings.Contains(s, "0xA") || !strings.Contains(s, "0xB") {
		t.Errorf("String() = %q, missing addresses", s)
	}
}
`,
			Hints: []string{
				"uint64(time.Now().Add(validFor).Unix()) for deadline",
				"time.Now().Unix() > int64(p.Deadline) means expired",
				`fmt.Sprintf("permit(%s→%s, %d, deadline: %d)", p.Owner, p.Spender, p.Value, p.Deadline)`,
			},
		},
	}
}

// ============================================================
// LEVEL 3: x402 Protocol — Payment flow
// ============================================================

func level3Protocol() []Question {
	return []Question{
		{
			ID: "decode-header", Title: "Decode PAYMENT-REQUIRED Header",
			Difficulty: "medium", Category: "x402",
			Description: `When a resource server returns HTTP 402, it includes a
PAYMENT-REQUIRED header containing base64-encoded JSON.

Decode and extract the payTo address from the first accepts entry.`,
			Template: `package x402quiz

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// DecodePayTo decodes a base64-encoded PAYMENT-REQUIRED header
// and returns the payTo address from the first accepts entry.
func DecodePayTo(headerValue string) (string, error) {
	// TODO: 1. Base64 decode
	// TODO: 2. JSON unmarshal — struct should have Accepts []struct{ PayTo string }
	// TODO: 3. Check len(accepts) > 0
	// TODO: 4. Return first payTo
	_ = base64.StdEncoding.DecodeString
	_ = json.Unmarshal
	_ = fmt.Errorf
	return "", nil
}
`,
			TestCode: `package x402quiz

import (
	"encoding/base64"
	"testing"
)

func TestDecodePayTo_Basic(t *testing.T) {
	raw := ` + "`" + `{"accepts":[{"scheme":"exact","payTo":"0xABCD1234","network":"eip155:84532"}]}` + "`" + `
	got, err := DecodePayTo(base64.StdEncoding.EncodeToString([]byte(raw)))
	if err != nil { t.Fatal(err) }
	if got != "0xABCD1234" { t.Errorf("got %q", got) }
}

func TestDecodePayTo_InvalidBase64(t *testing.T) {
	if _, err := DecodePayTo("not-valid!!!"); err == nil { t.Error("expected error") }
}

func TestDecodePayTo_EmptyAccepts(t *testing.T) {
	raw := ` + "`" + `{"accepts":[]}` + "`" + `
	if _, err := DecodePayTo(base64.StdEncoding.EncodeToString([]byte(raw))); err == nil {
		t.Error("expected error for empty accepts")
	}
}
`,
			Hints: []string{
				"Define: var pr struct { Accepts []struct { PayTo string `json:\"payTo\"` } `json:\"accepts\"` }",
				"if len(pr.Accepts) == 0 { return \"\", fmt.Errorf(\"empty\") }",
				"return pr.Accepts[0].PayTo, nil",
			},
		},
		{
			ID: "build-verify", Title: "Build /verify Request Body",
			Difficulty: "medium", Category: "x402",
			Description: `The facilitator's /verify endpoint expects JSON with:
  { "x402Version": 2, "paymentPayload": <raw JSON>, "paymentRequirements": <raw JSON> }

IMPORTANT: payload and requirements must be embedded as raw JSON,
not re-encoded as strings. Use json.RawMessage.`,
			Template: `package x402quiz

import "encoding/json"

// BuildVerifyBody constructs the /verify request body.
// payload and requirements are raw JSON bytes that must be
// embedded directly (not re-encoded as strings).
func BuildVerifyBody(payload, requirements []byte) ([]byte, error) {
	// TODO: Build the request body with json.RawMessage
	// Hint: json.RawMessage preserves raw JSON bytes
	_ = json.RawMessage{}
	_ = json.Marshal
	return nil, nil
}
`,
			TestCode: `package x402quiz

import (
	"encoding/json"
	"testing"
)

func TestBuildVerifyBody(t *testing.T) {
	payload := []byte(` + "`" + `{"from":"0xClient"}` + "`" + `)
	reqs := []byte(` + "`" + `{"scheme":"exact"}` + "`" + `)
	body, err := BuildVerifyBody(payload, reqs)
	if err != nil { t.Fatal(err) }

	var result struct {
		Version int             ` + "`" + `json:"x402Version"` + "`" + `
		Payload json.RawMessage ` + "`" + `json:"paymentPayload"` + "`" + `
		Reqs    json.RawMessage ` + "`" + `json:"paymentRequirements"` + "`" + `
	}
	if err := json.Unmarshal(body, &result); err != nil { t.Fatal(err) }
	if result.Version != 2 { t.Errorf("version = %d", result.Version) }
	if result.Payload[0] == '"' { t.Error("payload should be raw JSON, not string") }
	if string(result.Payload) != string(payload) { t.Errorf("payload mismatch") }
}
`,
			Hints: []string{
				"json.RawMessage is just []byte that won't be re-encoded",
				`body := map[string]any{"x402Version": 2, "paymentPayload": json.RawMessage(payload), ...}`,
				"return json.Marshal(body)",
			},
		},
		{
			ID: "parse-settlement", Title: "Parse PAYMENT-RESPONSE",
			Difficulty: "easy", Category: "x402",
			Description: `After settlement, the PAYMENT-RESPONSE header contains base64 JSON:
  { "success": true, "transaction": "0xabc...", "network": "eip155:84532", "payer": "0x..." }

Parse it and validate: if success is true, transaction must not be empty.`,
			Template: `package x402quiz

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Settlement holds the parsed settlement result.
type Settlement struct {
	Success     bool   ` + "`json:\"success\"`" + `
	Transaction string ` + "`json:\"transaction\"`" + `
	Network     string ` + "`json:\"network\"`" + `
	Payer       string ` + "`json:\"payer\"`" + `
}

// ParsePaymentResponse decodes a PAYMENT-RESPONSE header.
// Validates that successful settlements have a transaction hash.
func ParsePaymentResponse(headerValue string) (*Settlement, error) {
	// TODO: 1. Base64 decode
	// TODO: 2. JSON unmarshal into Settlement
	// TODO: 3. If success==true && transaction=="" → error
	_ = base64.StdEncoding.DecodeString
	_ = json.Unmarshal
	_ = fmt.Errorf
	return nil, nil
}
`,
			TestCode: `package x402quiz

import (
	"encoding/base64"
	"testing"
)

func TestParsePaymentResponse_Success(t *testing.T) {
	raw := ` + "`" + `{"success":true,"transaction":"0xABC","network":"eip155:84532","payer":"0xDEF"}` + "`" + `
	s, err := ParsePaymentResponse(base64.StdEncoding.EncodeToString([]byte(raw)))
	if err != nil { t.Fatal(err) }
	if !s.Success { t.Error("expected success") }
	if s.Transaction != "0xABC" { t.Errorf("tx = %q", s.Transaction) }
}

func TestParsePaymentResponse_SuccessNoTx(t *testing.T) {
	raw := ` + "`" + `{"success":true,"transaction":""}` + "`" + `
	_, err := ParsePaymentResponse(base64.StdEncoding.EncodeToString([]byte(raw)))
	if err == nil { t.Error("expected error: success without tx") }
}

func TestParsePaymentResponse_Invalid(t *testing.T) {
	if _, err := ParsePaymentResponse("bad!!!"); err == nil { t.Error("expected error") }
}
`,
			Hints: []string{
				"Same base64→JSON pattern as DecodePayTo",
				"Use the Settlement struct with json tags (already defined)",
				`if s.Success && s.Transaction == "" { return nil, fmt.Errorf("...") }`,
			},
		},
	}
}

// ============================================================
// LEVEL 4: Advanced — Integration challenges
// ============================================================

func level4Advanced() []Question {
	return []Question{
		{
			ID: "permit2-flow", Title: "Permit2 Approval Check",
			Difficulty: "hard", Category: "Permit2",
			Description: `Before using Permit2, the token owner must approve the Permit2
contract. Write a function that determines the required approval steps.

Permit2 address: 0x000000000022D473030F116dDEE9F6B43aC78BA3
x402Permit2Proxy: 0x402085c248EeA27D92E8b30b2C58ed07f9E20001`,
			Template: `package x402quiz

import "fmt"

const (
	Permit2Address     = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
	X402ProxyAddress   = "0x402085c248EeA27D92E8b30b2C58ed07f9E20001"
)

// ApprovalStep describes one step the user needs to take.
type ApprovalStep struct {
	Action      string // "approve" or "ready"
	Contract    string // which contract to call
	Description string
}

// CheckPermit2Readiness returns the steps needed to use Permit2.
// - If currentAllowance >= requiredAmount, return a single "ready" step
// - Otherwise, return an "approve" step telling the user to approve Permit2
func CheckPermit2Readiness(currentAllowance, requiredAmount uint64) []ApprovalStep {
	// TODO: Check if allowance is sufficient
	// TODO: Return appropriate steps
	_ = fmt.Sprintf
	return nil
}

// Permit2TransferPath returns the contract call chain for a Permit2 payment.
// The chain is: Client → Permit2 → x402Permit2Proxy → Token → PayTo
func Permit2TransferPath() []string {
	// TODO: Return the ordered list of contracts/actors in the transfer
	return nil
}
`,
			TestCode: `package x402quiz

import "testing"

func TestCheckPermit2Readiness_Sufficient(t *testing.T) {
	steps := CheckPermit2Readiness(1000000, 100000)
	if len(steps) != 1 { t.Fatalf("expected 1 step, got %d", len(steps)) }
	if steps[0].Action != "ready" { t.Errorf("action = %q, want \"ready\"", steps[0].Action) }
}

func TestCheckPermit2Readiness_NeedApprove(t *testing.T) {
	steps := CheckPermit2Readiness(0, 100000)
	if len(steps) != 1 { t.Fatalf("expected 1 step, got %d", len(steps)) }
	if steps[0].Action != "approve" { t.Errorf("action = %q, want \"approve\"", steps[0].Action) }
	if steps[0].Contract != Permit2Address { t.Errorf("contract = %q", steps[0].Contract) }
}

func TestPermit2TransferPath(t *testing.T) {
	path := Permit2TransferPath()
	if len(path) < 4 { t.Fatalf("expected at least 4 entries, got %d", len(path)) }
	// Permit2 and Proxy must be in the path
	hasPermit2, hasProxy := false, false
	for _, p := range path {
		if p == Permit2Address { hasPermit2 = true }
		if p == X402ProxyAddress { hasProxy = true }
	}
	if !hasPermit2 { t.Error("missing Permit2 in path") }
	if !hasProxy { t.Error("missing x402Permit2Proxy in path") }
}
`,
			Hints: []string{
				"if currentAllowance >= requiredAmount → return []ApprovalStep{{Action: \"ready\", ...}}",
				"Otherwise → return []ApprovalStep{{Action: \"approve\", Contract: Permit2Address, ...}}",
				`Path: []string{"Client", Permit2Address, X402ProxyAddress, "Token", "PayTo"}`,
			},
		},
		{
			ID: "payment-flow", Title: "x402 Payment Flow State Machine",
			Difficulty: "hard", Category: "x402",
			Description: `Model the x402 payment flow as a state machine.
States: idle → requesting → got402 → signing → sending → verifying → settled → done
Each transition should validate the previous state.`,
			Template: `package x402quiz

import "fmt"

// FlowState represents the current state of an x402 payment.
type FlowState string

const (
	StateIdle       FlowState = "idle"
	StateRequesting FlowState = "requesting"
	StateGot402     FlowState = "got402"
	StateSigning    FlowState = "signing"
	StateSending    FlowState = "sending"
	StateVerifying  FlowState = "verifying"
	StateSettled    FlowState = "settled"
	StateDone       FlowState = "done"
	StateError      FlowState = "error"
)

// PaymentFlow tracks the state of a payment.
type PaymentFlow struct {
	State FlowState
	Error error
}

// NewPaymentFlow creates a flow in idle state.
func NewPaymentFlow() *PaymentFlow {
	return &PaymentFlow{State: StateIdle}
}

// Transition moves to the next state. Returns error if the
// transition is invalid (e.g., can't go from idle to settled).
//
// Valid transitions:
//   idle → requesting → got402 → signing → sending → verifying → settled → done
//   any state → error
func (f *PaymentFlow) Transition(next FlowState) error {
	// TODO: Validate the transition is allowed
	// TODO: Update f.State
	_ = fmt.Errorf
	return nil
}

// IsComplete returns true if the flow reached "done" state.
func (f *PaymentFlow) IsComplete() bool {
	// TODO
	return false
}
`,
			TestCode: `package x402quiz

import "testing"

func TestPaymentFlow_HappyPath(t *testing.T) {
	f := NewPaymentFlow()
	states := []FlowState{StateRequesting, StateGot402, StateSigning, StateSending, StateVerifying, StateSettled, StateDone}
	for _, s := range states {
		if err := f.Transition(s); err != nil {
			t.Fatalf("transition to %s failed: %v", s, err)
		}
	}
	if !f.IsComplete() { t.Error("should be complete") }
}

func TestPaymentFlow_InvalidTransition(t *testing.T) {
	f := NewPaymentFlow()
	if err := f.Transition(StateSettled); err == nil {
		t.Error("idle → settled should be invalid")
	}
}

func TestPaymentFlow_ErrorFromAnyState(t *testing.T) {
	f := NewPaymentFlow()
	f.Transition(StateRequesting)
	if err := f.Transition(StateError); err != nil {
		t.Errorf("any → error should be valid: %v", err)
	}
}

func TestPaymentFlow_NotComplete(t *testing.T) {
	f := NewPaymentFlow()
	f.Transition(StateRequesting)
	if f.IsComplete() { t.Error("should not be complete") }
}
`,
			Hints: []string{
				"Define valid transitions as a map[FlowState]FlowState",
				`valid := map[FlowState]FlowState{StateIdle: StateRequesting, StateRequesting: StateGot402, ...}`,
				"Allow StateError from any state: if next == StateError { f.State = next; return nil }",
			},
		},
	}
}

// ============================================================
// LEVEL 5: Agents — ERC-8004 on-chain agent registry
// ============================================================

func level5Agents() []Question {
	return []Question{
		{
			ID: "agent-registration", Title: "Agent Registration File Parser",
			Difficulty: "easy", Category: "ERC-8004",
			Description: `ERC-8004 defines an on-chain registry for autonomous agents. Each agent
publishes a JSON registration file describing its capabilities, endpoints,
and payment configuration.

Parse the registration JSON, check whether the agent has x402 payment
support enabled, and extract all service endpoints.`,
			Template: `package x402quiz

import "encoding/json"

// X402Config holds the x402 payment configuration for an agent.
type X402Config struct {
	Enabled bool   ` + "`json:\"enabled\"`" + `
	Network string ` + "`json:\"network\"`" + `
	Asset   string ` + "`json:\"asset\"`" + `
}

// AgentService describes a single service endpoint.
type AgentService struct {
	Name     string ` + "`json:\"name\"`" + `
	Endpoint string ` + "`json:\"endpoint\"`" + `
	Price    string ` + "`json:\"price\"`" + `
}

// AgentRegistration is the top-level agent registration structure.
type AgentRegistration struct {
	AgentID     uint64         ` + "`json:\"agentId\"`" + `
	Name        string         ` + "`json:\"name\"`" + `
	Description string         ` + "`json:\"description\"`" + `
	Services    []AgentService ` + "`json:\"services\"`" + `
	X402        X402Config     ` + "`json:\"x402\"`" + `
}

// ParseRegistration parses a JSON registration file into an AgentRegistration.
func ParseRegistration(data []byte) (*AgentRegistration, error) {
	// TODO: json.Unmarshal data into AgentRegistration
	_ = json.Unmarshal
	return nil, nil
}

// HasX402Support returns true if the agent has x402 payments enabled.
func HasX402Support(reg *AgentRegistration) bool {
	// TODO: Check reg.X402.Enabled
	return false
}

// ServiceEndpoints collects all endpoint URLs from the agent's services.
func ServiceEndpoints(reg *AgentRegistration) []string {
	// TODO: Iterate reg.Services and collect Endpoint fields
	return nil
}
`,
			TestCode: `package x402quiz

import "testing"

func TestParseRegistration_Valid(t *testing.T) {
	data := []byte(` + "`" + `{
		"agentId": 1,
		"name": "WeatherBot",
		"description": "Provides weather data",
		"services": [
			{"name": "weather", "endpoint": "https://agent.example.com/weather", "price": "0.10"},
			{"name": "forecast", "endpoint": "https://agent.example.com/forecast", "price": "0.25"}
		],
		"x402": {"enabled": true, "network": "eip155:84532", "asset": "0x036CbD53842c5426634e7929541eC2318f3dCF7e"}
	}` + "`" + `)
	reg, err := ParseRegistration(data)
	if err != nil { t.Fatal(err) }
	if reg.AgentID != 1 { t.Errorf("AgentID = %d, want 1", reg.AgentID) }
	if reg.Name != "WeatherBot" { t.Errorf("Name = %q", reg.Name) }
	if len(reg.Services) != 2 { t.Errorf("Services count = %d, want 2", len(reg.Services)) }
	if reg.X402.Network != "eip155:84532" { t.Errorf("Network = %q", reg.X402.Network) }
}

func TestHasX402Support(t *testing.T) {
	enabled := &AgentRegistration{X402: X402Config{Enabled: true}}
	if !HasX402Support(enabled) { t.Error("expected x402 support enabled") }

	disabled := &AgentRegistration{X402: X402Config{Enabled: false}}
	if HasX402Support(disabled) { t.Error("expected x402 support disabled") }
}

func TestServiceEndpoints(t *testing.T) {
	reg := &AgentRegistration{
		Services: []AgentService{
			{Endpoint: "https://a.example.com/one"},
			{Endpoint: "https://a.example.com/two"},
			{Endpoint: "https://a.example.com/three"},
		},
	}
	eps := ServiceEndpoints(reg)
	if len(eps) != 3 { t.Fatalf("expected 3 endpoints, got %d", len(eps)) }
	if eps[0] != "https://a.example.com/one" { t.Errorf("eps[0] = %q", eps[0]) }
	if eps[2] != "https://a.example.com/three" { t.Errorf("eps[2] = %q", eps[2]) }
}

func TestParseRegistration_Invalid(t *testing.T) {
	_, err := ParseRegistration([]byte("not json"))
	if err == nil { t.Error("expected error for invalid JSON") }
}
`,
			Hints: []string{
				"var reg AgentRegistration; err := json.Unmarshal(data, &reg)",
				"HasX402Support simply returns reg.X402.Enabled",
				"for _, s := range reg.Services { eps = append(eps, s.Endpoint) }",
			},
		},
		{
			ID: "agent-global-id", Title: "Global Agent ID Format",
			Difficulty: "easy", Category: "ERC-8004",
			Description: `ERC-8004 agents are identified by a global ID that encodes chain,
contract, and token information in a single string:
  eip155:{chainId}:{contractAddress}:{tokenId}

Parse this format into its components, format it back, and validate
that the prefix is "eip155", the contract starts with "0x" and is
42 characters long.`,
			Template: `package x402quiz

import (
	"fmt"
	"strconv"
	"strings"
)

// GlobalAgentID holds the parsed components of a global agent identifier.
type GlobalAgentID struct {
	ChainID  uint64
	Contract string
	TokenID  uint64
}

// ParseGlobalAgentID parses "eip155:{chainId}:{contract}:{tokenId}".
// Returns an error if the format is invalid.
func ParseGlobalAgentID(id string) (*GlobalAgentID, error) {
	// TODO: Split by ":", validate 4 parts, first part must be "eip155"
	// TODO: Parse chainId and tokenId as uint64
	_ = strings.Split
	_ = strconv.ParseUint
	_ = fmt.Errorf
	return nil, nil
}

// FormatGlobalAgentID formats the ID back to "eip155:{chainId}:{contract}:{tokenId}".
func FormatGlobalAgentID(gid *GlobalAgentID) string {
	// TODO: Format the global ID string
	_ = fmt.Sprintf
	return ""
}

// ValidateGlobalAgentID checks that the ID components are well-formed:
// - Contract starts with "0x" and is 42 characters
func ValidateGlobalAgentID(gid *GlobalAgentID) error {
	// TODO: Validate contract address format
	_ = strings.HasPrefix
	return nil
}
`,
			TestCode: `package x402quiz

import "testing"

func TestParseGlobalAgentID_Valid(t *testing.T) {
	gid, err := ParseGlobalAgentID("eip155:84532:0x036CbD53842c5426634e7929541eC2318f3dCF7e:42")
	if err != nil { t.Fatal(err) }
	if gid.ChainID != 84532 { t.Errorf("ChainID = %d, want 84532", gid.ChainID) }
	if gid.Contract != "0x036CbD53842c5426634e7929541eC2318f3dCF7e" { t.Errorf("Contract = %q", gid.Contract) }
	if gid.TokenID != 42 { t.Errorf("TokenID = %d, want 42", gid.TokenID) }
}

func TestGlobalAgentID_Roundtrip(t *testing.T) {
	original := "eip155:84532:0x036CbD53842c5426634e7929541eC2318f3dCF7e:42"
	gid, err := ParseGlobalAgentID(original)
	if err != nil { t.Fatal(err) }
	formatted := FormatGlobalAgentID(gid)
	if formatted != original { t.Errorf("roundtrip failed: got %q, want %q", formatted, original) }
}

func TestParseGlobalAgentID_InvalidFormat(t *testing.T) {
	_, err := ParseGlobalAgentID("eip155:84532:0xABC")
	if err == nil { t.Error("expected error for wrong number of parts") }
}

func TestParseGlobalAgentID_InvalidChainID(t *testing.T) {
	_, err := ParseGlobalAgentID("eip155:notanumber:0x036CbD53842c5426634e7929541eC2318f3dCF7e:42")
	if err == nil { t.Error("expected error for invalid chain ID") }
}

func TestValidateGlobalAgentID(t *testing.T) {
	valid := &GlobalAgentID{ChainID: 84532, Contract: "0x036CbD53842c5426634e7929541eC2318f3dCF7e", TokenID: 1}
	if err := ValidateGlobalAgentID(valid); err != nil { t.Errorf("unexpected error: %v", err) }

	bad := &GlobalAgentID{ChainID: 84532, Contract: "notanaddress", TokenID: 1}
	if err := ValidateGlobalAgentID(bad); err == nil { t.Error("expected error for invalid contract") }
}
`,
			Hints: []string{
				`parts := strings.Split(id, ":"); if len(parts) != 4 { return error }`,
				`fmt.Sprintf("eip155:%d:%s:%d", gid.ChainID, gid.Contract, gid.TokenID)`,
				`if !strings.HasPrefix(gid.Contract, "0x") || len(gid.Contract) != 42 { return error }`,
			},
		},
		{
			ID: "agent-wad-encoding", Title: "Feedback Value Encoding with WAD Math",
			Difficulty: "medium", Category: "ERC-8004",
			Description: `ERC-8004 feedback values are stored on-chain using WAD encoding:
fixed-point math with 18 decimal places (1 WAD = 1e18). This is the
standard precision format in DeFi for avoiding floating-point errors.

Implement WAD conversion functions using math/big for arbitrary
precision: convert to/from WAD scale, compute averages, and
normalize between different decimal representations (e.g., USDC's
6 decimals to WAD's 18 decimals).`,
			Template: `package x402quiz

import "math/big"

// WAD is 1e18 — the standard fixed-point scale.
var WAD = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// ToWAD converts an integer value to WAD scale (value * 1e18).
func ToWAD(value int64) *big.Int {
	// TODO: Multiply value by WAD
	_ = new(big.Int).Mul
	return nil
}

// FromWAD converts a WAD-scaled value back to a plain integer (wad / 1e18).
func FromWAD(wad *big.Int) int64 {
	// TODO: Divide wad by WAD and return as int64
	_ = new(big.Int).Div
	return 0
}

// WADAverage computes the average of WAD-scaled values: sum / count.
// Returns zero if the slice is empty.
func WADAverage(values []*big.Int) *big.Int {
	// TODO: Sum all values, divide by count
	return nil
}

// NormalizeDecimals converts a value from one decimal scale to another.
// Example: 1000000 (6 decimals) → 1000000000000000000 (18 decimals)
func NormalizeDecimals(value int64, fromDecimals, toDecimals int) *big.Int {
	// TODO: Scale value by 10^(toDecimals - fromDecimals)
	// If toDecimals > fromDecimals, multiply; otherwise divide
	return nil
}
`,
			TestCode: `package x402quiz

import (
	"math/big"
	"testing"
)

func TestToWAD(t *testing.T) {
	got := ToWAD(100)
	want := new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	if got.Cmp(want) != 0 { t.Errorf("ToWAD(100) = %s, want %s", got, want) }
}

func TestFromWAD_Roundtrip(t *testing.T) {
	wad := ToWAD(42)
	got := FromWAD(wad)
	if got != 42 { t.Errorf("FromWAD(ToWAD(42)) = %d, want 42", got) }
}

func TestWADAverage(t *testing.T) {
	wad18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	v1 := new(big.Int).Mul(big.NewInt(80), wad18)
	v2 := new(big.Int).Mul(big.NewInt(-20), wad18)
	v3 := new(big.Int).Mul(big.NewInt(60), wad18)
	avg := WADAverage([]*big.Int{v1, v2, v3})
	want := new(big.Int).Mul(big.NewInt(40), wad18)
	if avg.Cmp(want) != 0 { t.Errorf("WADAverage = %s, want %s", avg, want) }
}

func TestWADAverage_Empty(t *testing.T) {
	avg := WADAverage(nil)
	if avg.Sign() != 0 { t.Errorf("expected zero for empty, got %s", avg) }
}

func TestNormalizeDecimals_6to18(t *testing.T) {
	got := NormalizeDecimals(1_000_000, 6, 18)
	want := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	if got.Cmp(want) != 0 { t.Errorf("NormalizeDecimals(1e6, 6, 18) = %s, want %s", got, want) }
}

func TestNormalizeDecimals_18to6(t *testing.T) {
	got := NormalizeDecimals(1_000_000_000_000_000_000, 18, 6)
	want := big.NewInt(1_000_000)
	if got.Cmp(want) != 0 { t.Errorf("NormalizeDecimals(1e18, 18, 6) = %s, want %s", got, want) }
}
`,
			Hints: []string{
				"ToWAD: return new(big.Int).Mul(big.NewInt(value), WAD)",
				"WADAverage: sum with Add in a loop, then Div by big.NewInt(int64(len(values)))",
				"NormalizeDecimals: compute 10^|toDecimals-fromDecimals| then multiply or divide accordingly",
			},
		},
		{
			ID: "agent-x402-integration", Title: "x402 + ERC-8004 Integration",
			Difficulty: "medium", Category: "ERC-8004",
			Description: `When an agent performs a service paid via x402, the settlement receipt
can be linked to on-chain feedback. This creates a verifiable connection:
proof-of-payment tied to an agent rating.

Build an AgentFeedback struct that embeds a ProofOfPayment from an x402
settlement. Validate that the proof contains a well-formed transaction
hash and positive payment amount. The feedback tag must be "x402-payment"
to identify payment-linked reviews.`,
			Template: `package x402quiz

import (
	"fmt"
	"strings"
)

// ProofOfPayment links feedback to an on-chain x402 settlement.
type ProofOfPayment struct {
	TxHash  string
	Network string
	Amount  uint64
	Payer   string
}

// AgentFeedback represents a payment-linked agent review.
type AgentFeedback struct {
	AgentID  uint64
	Provider string
	Value    int64 // rating score (can be negative)
	Tag      string
	Proof    *ProofOfPayment
}

// BuildFeedbackFromSettlement creates an AgentFeedback from x402 settlement data.
// Validates: txHash starts with "0x" and len >= 66, amount > 0.
// Sets Tag to "x402-payment" and Provider to the payer address.
func BuildFeedbackFromSettlement(txHash, network string, amount uint64, payer string, agentID uint64, rating int64) (*AgentFeedback, error) {
	// TODO: Validate txHash format (starts with "0x", length >= 66)
	// TODO: Validate amount > 0
	// TODO: Build AgentFeedback with embedded ProofOfPayment
	_ = fmt.Errorf
	_ = strings.HasPrefix
	return nil, nil
}

// ValidateProofOfPayment checks that all fields are non-empty and txHash is well-formed.
func ValidateProofOfPayment(p *ProofOfPayment) error {
	// TODO: Check TxHash, Network, Payer are non-empty
	// TODO: Check TxHash starts with "0x" and len >= 66
	_ = fmt.Errorf
	return nil
}
`,
			TestCode: `package x402quiz

import "testing"

func TestBuildFeedbackFromSettlement_Valid(t *testing.T) {
	fb, err := BuildFeedbackFromSettlement(
		"0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24",
		"eip155:84532", 100000, "0xClientAddr", 1, 85,
	)
	if err != nil { t.Fatal(err) }
	if fb.AgentID != 1 { t.Errorf("AgentID = %d", fb.AgentID) }
	if fb.Tag != "x402-payment" { t.Errorf("Tag = %q, want \"x402-payment\"", fb.Tag) }
	if fb.Proof == nil { t.Fatal("Proof is nil") }
	if fb.Proof.Amount != 100000 { t.Errorf("Proof.Amount = %d", fb.Proof.Amount) }
	if fb.Value != 85 { t.Errorf("Value = %d, want 85", fb.Value) }
	if fb.Provider != "0xClientAddr" { t.Errorf("Provider = %q", fb.Provider) }
}

func TestBuildFeedbackFromSettlement_BadTxHash(t *testing.T) {
	_, err := BuildFeedbackFromSettlement("short", "eip155:84532", 100000, "0xPayer", 1, 50)
	if err == nil { t.Error("expected error for invalid txHash") }
}

func TestBuildFeedbackFromSettlement_ZeroAmount(t *testing.T) {
	_, err := BuildFeedbackFromSettlement(
		"0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24",
		"eip155:84532", 0, "0xPayer", 1, 50,
	)
	if err == nil { t.Error("expected error for zero amount") }
}

func TestValidateProofOfPayment_Valid(t *testing.T) {
	p := &ProofOfPayment{
		TxHash:  "0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24",
		Network: "eip155:84532", Amount: 100000, Payer: "0xAddr",
	}
	if err := ValidateProofOfPayment(p); err != nil { t.Errorf("unexpected error: %v", err) }
}

func TestValidateProofOfPayment_MissingTxHash(t *testing.T) {
	p := &ProofOfPayment{TxHash: "", Network: "eip155:84532", Amount: 100000, Payer: "0xAddr"}
	if err := ValidateProofOfPayment(p); err == nil { t.Error("expected error for empty txHash") }
}
`,
			Hints: []string{
				`if !strings.HasPrefix(txHash, "0x") || len(txHash) < 66 { return nil, fmt.Errorf("invalid txHash") }`,
				`Tag must be set to "x402-payment" and Provider to payer`,
				`ValidateProofOfPayment: check p.TxHash != "" && p.Network != "" && p.Payer != "" then check txHash format`,
			},
		},
	}
}
