package demo

import "encoding/json"

// WalletInfo identifies a named wallet address.
type WalletInfo struct {
	Name    string
	Address string
}

// WalletBalance holds ETH and USDC balances for a wallet.
type WalletBalance struct {
	Wallet  WalletInfo
	ETH     string // formatted with 6 decimal places
	USDC    string // formatted with 6 decimal places
	ETHRaw  string // raw wei value
	USDCRaw string // raw smallest-unit value
}

// DecodedPaymentRequired is the decoded PAYMENT-REQUIRED header.
type DecodedPaymentRequired struct {
	Accepts     []json.RawMessage `json:"accepts"`
	Resource    string            `json:"resource,omitempty"`
	Description string            `json:"description,omitempty"`
	MimeType    string            `json:"mimeType,omitempty"`
}

// AcceptItem is a single entry in the accepts array.
type AcceptItem struct {
	Scheme             string                 `json:"scheme"`
	Network            string                 `json:"network"`
	MaxAmountRequired  string                 `json:"maxAmountRequired"`
	Resource           string                 `json:"resource"`
	Description        string                 `json:"description"`
	MimeType           string                 `json:"mimeType"`
	PayTo              string                 `json:"payTo"`
	Asset              string                 `json:"asset"`
	Extra              map[string]interface{} `json:"extra"`
}

// FlowState tracks the state of a payment flow execution.
type FlowState struct {
	TransferMethod  string
	Wallets         []WalletInfo
	BalancesBefore  []WalletBalance
	BalancesAfter   []WalletBalance
	PaymentRequired *DecodedPaymentRequired
	PaymentPayload  json.RawMessage // raw payload JSON
	VerifyResponse  json.RawMessage // raw verify response
	SettleResponse  json.RawMessage // raw settle response
	TxHash          string
	CurrentStep     int
	TotalSteps      int
	Error           error
}

// NewFlowState creates a new flow state with default values.
func NewFlowState(transferMethod string) *FlowState {
	return &FlowState{
		TransferMethod: transferMethod,
		TotalSteps:     10,
	}
}

// StepDescription returns the description for each step.
func StepDescription(step int) string {
	descriptions := map[int]string{
		1:  "지갑 주소 & 잔액 확인",
		2:  "Facilitator /supported (서비스 디스커버리)",
		3:  "Client → Resource Server: 결제 없이 API 호출",
		4:  "402 응답의 PAYMENT-REQUIRED 헤더 디코딩",
		5:  "Client: 결제 서명 생성 (오프체인)",
		6:  "Client → Resource Server: PAYMENT-SIGNATURE 포함 재요청",
		7:  "Resource Server → Facilitator /verify (오프체인 검증)",
		8:  "검증 성공 → 데이터 반환 + /settle 요청",
		9:  "Facilitator /settle → 온체인 정산 + PAYMENT-RESPONSE",
		10: "정산 후 잔액 확인",
	}
	return descriptions[step]
}
