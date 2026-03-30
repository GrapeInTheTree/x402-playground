package demo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	evmclient "github.com/coinbase/x402/go/mechanisms/evm/exact/client"
	evmsigner "github.com/coinbase/x402/go/signers/evm"
)

// LiveExecutor executes the 10-step payment flow with real HTTP/on-chain calls.
type LiveExecutor struct {
	FacilitatorURL string
	ResourceURL    string
	RPCURL         string
	USDCAddress    string
	PayToAddress   string
	ClientKey      string
	EndpointPath   string
	TransferMethod string

	// HTTP client with timeout for all outbound calls
	httpClient *http.Client

	// Accumulated state across steps
	ethClient        *ethclient.Client
	facilitatorAddr  string
	clientAddr       string
	resp402Headers   map[string]string
	body402          []byte
	payloadBytes     []byte
	headerName       string
	headerValue      string
	selectedReqBytes []byte
	paidRespHeaders  map[string]string
}

// NewLiveExecutor creates a new executor. Returns error if required config is missing.
func NewLiveExecutor(facilitatorURL, resourceURL, rpcURL, usdcAddr, payToAddr, clientKey, transferMethod string) (*LiveExecutor, error) {
	if clientKey == "" {
		return nil, fmt.Errorf("CLIENT_PRIVATE_KEY is required for live execution")
	}
	if resourceURL == "" {
		return nil, fmt.Errorf("RESOURCE_URL is required for live execution")
	}
	if facilitatorURL == "" {
		return nil, fmt.Errorf("FACILITATOR_URL is required for live execution")
	}
	if rpcURL == "" {
		rpcURL = "https://sepolia.base.org"
	}
	if transferMethod == "" {
		transferMethod = "eip3009"
	}
	return &LiveExecutor{
		FacilitatorURL: facilitatorURL,
		ResourceURL:    resourceURL,
		RPCURL:         rpcURL,
		USDCAddress:    usdcAddr,
		PayToAddress:   payToAddr,
		ClientKey:      clientKey,
		EndpointPath:   "/weather",
		TransferMethod: transferMethod,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// RunStep executes the given step (0-indexed) and returns display text.
func (e *LiveExecutor) RunStep(ctx context.Context, step int) (string, error) {
	switch step {
	case 0:
		return e.step1WalletOverview(ctx)
	case 1:
		return e.step2Supported(ctx)
	case 2:
		return e.step3NaiveCall(ctx)
	case 3:
		return e.step4DecodeHeader()
	case 4:
		return e.step5CreateSignature(ctx)
	case 5:
		return e.step6SendPayment()
	case 6:
		return e.step7Verify(ctx)
	case 7:
		return e.step8PaidRequest(ctx)
	case 8:
		return e.step9Settlement()
	case 9:
		return e.step10FinalBalances(ctx)
	default:
		return "", fmt.Errorf("unknown step %d", step)
	}
}

func (e *LiveExecutor) step1WalletOverview(ctx context.Context) (string, error) {
	clientSigner, err := evmsigner.NewClientSignerFromPrivateKey(e.ClientKey)
	if err != nil {
		return "", fmt.Errorf("invalid client key: %w", err)
	}
	e.clientAddr = clientSigner.Address()

	// Try to get facilitator address from /health
	e.facilitatorAddr = "(unknown)"
	resp, err := e.httpClient.Get(e.FacilitatorURL + "/health")
	if err == nil {
		var health struct {
			Address string `json:"address"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&health)
		resp.Body.Close()
		if health.Address != "" {
			e.facilitatorAddr = health.Address
		}
	}

	// Query balances
	var balanceInfo string
	client, err := ethclient.Dial(e.RPCURL)
	if err == nil {
		e.ethClient = client
		wallets := []WalletInfo{
			{Name: "Facilitator", Address: e.facilitatorAddr},
			{Name: "Client", Address: e.clientAddr},
		}
		if e.PayToAddress != "" {
			wallets = append(wallets, WalletInfo{Name: "PAY_TO", Address: e.PayToAddress})
		}
		bals, err := QueryBalances(ctx, client, e.USDCAddress, wallets)
		if err == nil {
			var sb strings.Builder
			for _, b := range bals {
				sb.WriteString(fmt.Sprintf("  %s: ETH %s  USDC %s\n", b.Wallet.Name, b.ETH, b.USDC))
			}
			balanceInfo = sb.String()
		}
	}

	return fmt.Sprintf("Facilitator: %s\nClient:      %s\nPAY_TO:      %s\nMethod:      %s\n\n%s",
		e.facilitatorAddr, e.clientAddr, e.PayToAddress, e.TransferMethod, balanceInfo), nil
}

func (e *LiveExecutor) step2Supported(_ context.Context) (string, error) {
	resp, err := e.httpClient.Get(e.FacilitatorURL + "/supported")
	if err != nil {
		return "", fmt.Errorf("GET /supported failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	return FormatJSON(body), nil
}

func (e *LiveExecutor) step3NaiveCall(_ context.Context) (string, error) {
	targetURL := e.ResourceURL + e.EndpointPath
	resp, err := e.httpClient.Get(targetURL)
	if err != nil {
		return "", fmt.Errorf("GET %s failed: %w", targetURL, err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// Store for later steps
	e.body402 = body
	e.resp402Headers = make(map[string]string)
	for k := range resp.Header {
		e.resp402Headers[k] = resp.Header.Get(k)
	}

	return fmt.Sprintf("HTTP %d %s\n\nBody: %s", resp.StatusCode, http.StatusText(resp.StatusCode), string(body)), nil
}

func (e *LiveExecutor) step4DecodeHeader() (string, error) {
	headerVal := e.resp402Headers["Payment-Required"]
	if headerVal == "" {
		return "", fmt.Errorf("no PAYMENT-REQUIRED header in 402 response")
	}

	decoded, err := base64.StdEncoding.DecodeString(headerVal)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	return fmt.Sprintf("Raw: %d chars (base64)\n\n%s", len(headerVal), FormatJSON(decoded)), nil
}

func (e *LiveExecutor) step5CreateSignature(ctx context.Context) (string, error) {
	clientSigner, err := evmsigner.NewClientSignerFromPrivateKey(e.ClientKey)
	if err != nil {
		return "", fmt.Errorf("create signer: %w", err)
	}

	x402Client := x402.Newx402Client().
		Register("eip155:*", evmclient.NewExactEvmScheme(clientSigner, &evmclient.ExactEvmSchemeConfig{
			RPCURL: e.RPCURL,
		}))
	httpx402Client := x402http.Newx402HTTPClient(x402Client)

	paymentReq, err := httpx402Client.GetPaymentRequiredResponse(e.resp402Headers, e.body402)
	if err != nil {
		return "", fmt.Errorf("parse payment required: %w", err)
	}

	selectedReq, err := x402Client.SelectPaymentRequirements(paymentReq.Accepts)
	if err != nil {
		return "", fmt.Errorf("select requirements: %w", err)
	}

	payloadTyped, err := x402Client.CreatePaymentPayload(ctx, selectedReq, paymentReq.Resource, paymentReq.Extensions)
	if err != nil {
		return "", fmt.Errorf("create payload: %w", err)
	}

	e.payloadBytes, _ = json.Marshal(payloadTyped)
	e.selectedReqBytes, _ = json.Marshal(selectedReq)

	headers, err := httpx402Client.EncodePaymentSignatureHeader(e.payloadBytes)
	if err != nil {
		return "", fmt.Errorf("encode header: %w", err)
	}
	for k, v := range headers {
		e.headerName = k
		e.headerValue = v
	}

	return fmt.Sprintf("Payload:\n%s\n\nHeader: %s (%d chars)",
		FormatJSON(e.payloadBytes), e.headerName, len(e.headerValue)), nil
}

func (e *LiveExecutor) step6SendPayment() (string, error) {
	if e.headerName == "" {
		return "", fmt.Errorf("no payment signature created (run step 5 first)")
	}
	return fmt.Sprintf("-> GET %s%s\n  %s: <base64 payload>\n\nResource Server:\n  1. PAYMENT-SIGNATURE header parsing\n  2. PaymentPayload decoding\n  3. Forward /verify request to Facilitator",
		e.ResourceURL, e.EndpointPath, e.headerName), nil
}

func (e *LiveExecutor) step7Verify(_ context.Context) (string, error) {
	if e.payloadBytes == nil {
		return "", fmt.Errorf("no payload (run step 5 first)")
	}

	verifyBody := map[string]any{
		"x402Version":         2,
		"paymentPayload":      json.RawMessage(e.payloadBytes),
		"paymentRequirements": json.RawMessage(e.selectedReqBytes),
	}
	verifyJSON, _ := json.Marshal(verifyBody)

	resp, err := e.httpClient.Post(e.FacilitatorURL+"/verify", "application/json", strings.NewReader(string(verifyJSON)))
	if err != nil {
		return "", fmt.Errorf("POST /verify: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return fmt.Sprintf("POST %s/verify\n\n%s", e.FacilitatorURL, FormatJSON(body)), nil
}

func (e *LiveExecutor) step8PaidRequest(ctx context.Context) (string, error) {
	if e.headerName == "" {
		return "", fmt.Errorf("no payment header (run step 5 first)")
	}

	targetURL := e.ResourceURL + e.EndpointPath
	req, _ := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	req.Header.Set(e.headerName, e.headerValue)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("paid request: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	e.paidRespHeaders = make(map[string]string)
	for k := range resp.Header {
		e.paidRespHeaders[k] = resp.Header.Get(k)
	}

	return fmt.Sprintf("HTTP %d %s\n\n%s", resp.StatusCode, http.StatusText(resp.StatusCode), FormatJSON(body)), nil
}

func (e *LiveExecutor) step9Settlement() (string, error) {
	prHeader := e.paidRespHeaders["Payment-Response"]

	if prHeader == "" {
		return "PAYMENT-RESPONSE header not found\n(Settlement may still be in progress or requires separate confirmation)", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(prHeader)
	if err != nil {
		return "", fmt.Errorf("decode PAYMENT-RESPONSE: %w", err)
	}

	var settle struct {
		Success     bool   `json:"success"`
		Transaction string `json:"transaction"`
		Network     string `json:"network"`
	}
	_ = json.Unmarshal(decoded, &settle)

	result := FormatJSON(decoded)
	if settle.Success && settle.Transaction != "" {
		result += fmt.Sprintf("\n\nTx: https://sepolia.basescan.org/tx/%s", settle.Transaction)
	}
	return result, nil
}

func (e *LiveExecutor) step10FinalBalances(ctx context.Context) (string, error) {
	if e.ethClient == nil {
		client, err := ethclient.Dial(e.RPCURL)
		if err != nil {
			return "", fmt.Errorf("connect RPC: %w", err)
		}
		e.ethClient = client
	}

	wallets := []WalletInfo{
		{Name: "Facilitator", Address: e.facilitatorAddr},
		{Name: "Client", Address: e.clientAddr},
	}
	if e.PayToAddress != "" {
		wallets = append(wallets, WalletInfo{Name: "PAY_TO", Address: e.PayToAddress})
	}

	bals, err := QueryBalances(ctx, e.ethClient, e.USDCAddress, wallets)
	if err != nil {
		return "", fmt.Errorf("query balances: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("AFTER balances:\n\n")
	for _, b := range bals {
		sb.WriteString(fmt.Sprintf("  %s: ETH %s  USDC %s\n", b.Wallet.Name, b.ETH, b.USDC))
	}
	sb.WriteString(fmt.Sprintf("\nClient -> PAY_TO: 0.1 USDC (%s)", e.TransferMethod))
	return sb.String(), nil
}

// Close releases resources held by the executor.
func (e *LiveExecutor) Close() {
	if e.ethClient != nil {
		e.ethClient.Close()
	}
}
