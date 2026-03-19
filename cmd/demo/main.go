package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	evmclient "github.com/coinbase/x402/go/mechanisms/evm/exact/client"
	evmsigner "github.com/coinbase/x402/go/signers/evm"

	"github.com/GrapeInTheTree/x402-demo/internal/config"
)

func main() {
	cfg, err := config.LoadClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	resCfg, err := config.LoadResource()
	if err != nil {
		fmt.Fprintf(os.Stderr, "resource config error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	transferMethod := envOr("ASSET_TRANSFER_METHOD", "eip3009")
	banner("x402 Payment Protocol — Full Flow Demo")
	fmt.Printf("  Transfer Method: %s\n", transferMethod)
	if transferMethod == "permit2" {
		fmt.Println("  → Permit2 모드: x402Permit2Proxy 경유, 모든 ERC-20 토큰 지원")
	} else {
		fmt.Println("  → EIP-3009 모드: transferWithAuthorization 직접 호출 (USDC 전용)")
	}
	fmt.Println()

	// ─────────────────────────────────────────────────────────
	// STEP 1: Wallet Overview
	// ─────────────────────────────────────────────────────────
	step(1, "지갑 주소 & 잔액 확인")

	clientSigner, _ := evmsigner.NewClientSignerFromPrivateKey(cfg.PrivateKey)

	facilitatorAddr := "0x23fbdE5A14dFB508502f5A2622f66c0D3B0ab37A"
	// Get facilitator address from /health
	if healthResp, err := http.Get(resCfg.FacilitatorURL + "/health"); err == nil {
		var health struct {
			Address string `json:"address"`
		}
		json.NewDecoder(healthResp.Body).Decode(&health)
		healthResp.Body.Close()
		if health.Address != "" {
			facilitatorAddr = health.Address
		}
	}

	fmt.Printf("  Facilitator (가스비 지불):  %s\n", facilitatorAddr)
	fmt.Printf("  Client      (USDC 결제):   %s\n", clientSigner.Address())
	fmt.Printf("  PAY_TO      (USDC 수신):   %s\n", resCfg.PayToAddress)
	fmt.Printf("  Network:                   %s\n", cfg.Network)
	fmt.Printf("  USDC Contract:             %s\n", cfg.USDCAddress)
	fmt.Println()

	ethClient, _ := ethclient.Dial(cfg.RPCURL)
	showBalances(ctx, ethClient, cfg.USDCAddress, "BEFORE",
		wallet{"Facilitator", facilitatorAddr},
		wallet{"Client", clientSigner.Address()},
		wallet{"PAY_TO", resCfg.PayToAddress},
	)
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 2: Facilitator /supported
	// ─────────────────────────────────────────────────────────
	step(2, "Resource Server → Facilitator /supported (서비스 디스커버리)")

	fmt.Printf("  → GET %s/supported\n\n", resCfg.FacilitatorURL)
	supportedResp, _ := http.Get(resCfg.FacilitatorURL + "/supported")
	supportedBody, _ := io.ReadAll(supportedResp.Body)
	supportedResp.Body.Close()
	prettyPrint("  Facilitator 지원 현황", supportedBody)
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 3: Naive API Call → 402
	// ─────────────────────────────────────────────────────────
	step(3, "Client → Resource Server: 결제 없이 API 호출")

	targetURL := cfg.ResourceURL + "/weather"
	fmt.Printf("  → GET %s\n", targetURL)
	fmt.Println("    (PAYMENT-SIGNATURE 헤더 없음 — 일반 HTTP 요청)")
	fmt.Println()

	resp402, _ := http.Get(targetURL)
	body402, _ := io.ReadAll(resp402.Body)
	resp402.Body.Close()

	fmt.Printf("  ← HTTP %d %s\n", resp402.StatusCode, http.StatusText(resp402.StatusCode))
	fmt.Printf("    Body: %s\n", string(body402))
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 4: Decode 402 PAYMENT-REQUIRED
	// ─────────────────────────────────────────────────────────
	step(4, "402 응답의 PAYMENT-REQUIRED 헤더 디코딩")

	paymentRequiredHeader := resp402.Header.Get("Payment-Required")
	fmt.Printf("  Raw Header (base64, %d chars):\n", len(paymentRequiredHeader))
	fmt.Printf("    %s...%s\n\n", paymentRequiredHeader[:60], paymentRequiredHeader[len(paymentRequiredHeader)-20:])

	prBytes, _ := base64.StdEncoding.DecodeString(paymentRequiredHeader)
	prettyPrint("  디코딩된 결제 요구사항", prBytes)

	var paymentRequired struct {
		Accepts []json.RawMessage `json:"accepts"`
	}
	json.Unmarshal(prBytes, &paymentRequired)
	fmt.Println()
	fmt.Println("  서버가 요구하는 것:")
	fmt.Println("    • scheme: exact (정확한 금액 결제)")
	fmt.Println("    • network: eip155:84532 (Base Sepolia)")
	fmt.Println("    • amount: 100000 (0.1 USDC, 6 decimals)")
	fmt.Println("    • payTo: " + resCfg.PayToAddress)
	fmt.Println("    • extra.name: USDC (EIP-712 도메인 name)")
	if transferMethod == "permit2" {
		fmt.Println("    • extra.assetTransferMethod: permit2 ← Client가 Permit2 서명 생성")
	}
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 5: Create Payment Signature
	// ─────────────────────────────────────────────────────────
	if transferMethod == "permit2" {
		step(5, "Client: Permit2 서명 생성 (오프체인)")
		fmt.Println("  1. 402 응답의 accepts에서 (scheme=exact, assetTransferMethod=permit2) 선택")
		fmt.Println("  2. EIP-712 Typed Data 구성:")
		fmt.Println("     Domain: Permit2 (0x000000000022D473030F116dDEE9F6B43aC78BA3)")
		fmt.Println("     Message: PermitWitnessTransferFrom { permitted, spender, nonce, deadline, witness }")
		fmt.Println("     - spender: x402Permit2Proxy (0x402085c248EeA27D92E8b30b2C58ed07f9E20001)")
		fmt.Println("     - witness: { to, validAfter }")
		fmt.Println("  3. Client private key로 ECDSA 서명 (v, r, s)")
		fmt.Println("  4. Permit2Payload JSON 생성 → base64 인코딩 → PAYMENT-SIGNATURE 헤더")
		fmt.Println()
	} else {
		step(5, "Client: EIP-3009 서명 생성 (오프체인)")
		fmt.Println("  1. 402 응답의 accepts에서 (scheme=exact, network=eip155:84532) 선택")
		fmt.Println("  2. EIP-712 Typed Data 구성:")
		fmt.Println("     Domain: { name: \"USDC\", version: \"2\", chainId: 84532, verifyingContract: 0x036C... }")
		fmt.Println("     Message: TransferWithAuthorization { from, to, value, validAfter, validBefore, nonce }")
		fmt.Println("  3. Client private key로 ECDSA 서명 (v, r, s)")
		fmt.Println("  4. PaymentPayload JSON 생성 → base64 인코딩 → PAYMENT-SIGNATURE 헤더")
		fmt.Println()
	}

	// Actually create the payment using SDK
	x402Client := x402.Newx402Client().
		Register("eip155:*", evmclient.NewExactEvmScheme(clientSigner, &evmclient.ExactEvmSchemeConfig{
			RPCURL: cfg.RPCURL,
		}))
	httpx402Client := x402http.Newx402HTTPClient(x402Client)

	// Parse payment required to get requirements
	headers402 := make(map[string]string)
	for k, v := range resp402.Header {
		if len(v) > 0 {
			headers402[k] = v[0]
		}
	}
	paymentReq, _ := httpx402Client.GetPaymentRequiredResponse(headers402, body402)

	// Select requirements and create payload
	selectedReq, _ := x402Client.SelectPaymentRequirements(paymentReq.Accepts)
	payloadTyped, _ := x402Client.CreatePaymentPayload(ctx, selectedReq, paymentReq.Resource, paymentReq.Extensions)
	payloadBytes, _ := json.Marshal(payloadTyped)

	prettyPrint("  생성된 PaymentPayload", payloadBytes)

	// Encode as header
	paymentHeaders, _ := httpx402Client.EncodePaymentSignatureHeader(payloadBytes)
	headerName := ""
	headerValue := ""
	for k, v := range paymentHeaders {
		headerName = k
		headerValue = v
	}
	fmt.Printf("\n  PAYMENT-SIGNATURE 헤더 (base64, %d chars):\n", len(headerValue))
	fmt.Printf("    %s: %s...%s\n", headerName, headerValue[:50], headerValue[len(headerValue)-20:])
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 6: Client → Resource Server (with payment)
	// ─────────────────────────────────────────────────────────
	step(6, "Client → Resource Server: PAYMENT-SIGNATURE 포함 재요청")

	fmt.Printf("  → GET %s\n", targetURL)
	fmt.Printf("    %s: <base64 payload>\n\n", headerName)

	req, _ := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	req.Header.Set(headerName, headerValue)

	fmt.Println("  Resource Server 내부 처리:")
	fmt.Println("    1. PAYMENT-SIGNATURE 헤더 파싱")
	fmt.Println("    2. PaymentPayload 디코딩")
	fmt.Println("    3. Facilitator에 /verify 요청 전달")
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 7: Facilitator /verify (off-chain)
	// ─────────────────────────────────────────────────────────
	step(7, "Resource Server → Facilitator /verify (오프체인 검증)")

	// Call verify directly to show the response
	verifyBody := map[string]interface{}{
		"x402Version":         2,
		"paymentPayload":      json.RawMessage(payloadBytes),
		"paymentRequirements": json.RawMessage(must(json.Marshal(selectedReq))),
	}
	verifyJSON, _ := json.Marshal(verifyBody)

	fmt.Printf("  → POST %s/verify\n\n", resCfg.FacilitatorURL)
	prettyPrint("  요청 본문 (paymentPayload + paymentRequirements)", verifyJSON)
	fmt.Println()

	if transferMethod == "permit2" {
		fmt.Println("  Facilitator 검증 항목 (Permit2):")
		fmt.Println("    ✓ EIP-712 서명 복원 → from 주소 일치 확인")
		fmt.Println("    ✓ spender == x402Permit2Proxy 확인")
		fmt.Println("    ✓ witness.to == payTo 확인")
		fmt.Println("    ✓ permitted.amount >= amount 확인")
		fmt.Println("    ✓ permitted.token == asset 확인")
		fmt.Println("    ✓ deadline >= now, witness.validAfter <= now 시간 유효성")
		fmt.Println("    ✓ Client USDC 잔액 >= value 온체인 확인")
		fmt.Println("    ✓ Client → Permit2 approve 확인 (allowance)")
		fmt.Println("    ✓ eth_call로 x402Permit2Proxy.settle 시뮬레이션")
	} else {
		fmt.Println("  Facilitator 검증 항목:")
		fmt.Println("    ✓ EIP-712 서명 복원 → from 주소 일치 확인")
		fmt.Println("    ✓ authorization.to == payTo 확인")
		fmt.Println("    ✓ authorization.value >= amount 확인")
		fmt.Println("    ✓ validAfter <= now <= validBefore 시간 유효성")
		fmt.Println("    ✓ nonce 미사용 확인 (이중 결제 방지)")
		fmt.Println("    ✓ Client USDC 잔액 >= value 온체인 확인")
		fmt.Println("    ✓ eth_call로 transferWithAuthorization 시뮬레이션")
	}
	fmt.Println()

	verifyResp, _ := http.Post(resCfg.FacilitatorURL+"/verify", "application/json", bytes.NewReader(verifyJSON))
	verifyRespBody, _ := io.ReadAll(verifyResp.Body)
	verifyResp.Body.Close()

	prettyPrint("  ← Facilitator /verify 응답", verifyRespBody)
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 8: Resource Server → Client (200 + data)
	// ─────────────────────────────────────────────────────────
	step(8, "검증 성공 → Resource Server가 Client에게 데이터 반환 + /settle 요청")

	fmt.Println("  검증 통과 → Resource Server가 핸들러 실행 → Client에 200 OK + 데이터 반환")
	fmt.Println("  동시에 Facilitator에 /settle 요청 (온체인 정산)")
	fmt.Println()

	// Now do the actual request through the resource server
	paidResp, _ := http.DefaultClient.Do(req)
	paidBody, _ := io.ReadAll(paidResp.Body)
	paidResp.Body.Close()

	fmt.Printf("  ← HTTP %d %s\n\n", paidResp.StatusCode, http.StatusText(paidResp.StatusCode))
	prettyPrint("  API 응답 데이터", paidBody)
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 9: Settlement + PAYMENT-RESPONSE
	// ─────────────────────────────────────────────────────────
	step(9, "Facilitator /settle → 온체인 정산 + PAYMENT-RESPONSE")

	if transferMethod == "permit2" {
		fmt.Println("  Facilitator 정산 과정 (Permit2):")
		fmt.Println("    1. Permit2Authorization + 서명 추출")
		fmt.Println("    2. EIP-1559 트랜잭션 빌드 (gas estimation + 20% 버퍼)")
		fmt.Println("    3. x402Permit2Proxy.settle(owner, permitted, nonce, deadline, witness, signature)")
		fmt.Println("       → Proxy가 Permit2.permitWitnessTransferFrom() 호출")
		fmt.Println("       → USDC가 Client → PAY_TO로 이동")
		fmt.Println("    4. Facilitator 지갑이 가스비 지불")
		fmt.Println("    5. eth_sendRawTransaction → Base Sepolia 제출")
		fmt.Println("    6. Receipt 대기 (2초 간격 polling)")
	} else {
		fmt.Println("  Facilitator 정산 과정:")
		fmt.Println("    1. 서명에서 v, r, s 추출")
		fmt.Println("    2. EIP-1559 트랜잭션 빌드 (gas estimation + 20% 버퍼)")
		fmt.Println("    3. USDC.transferWithAuthorization(from, to, value, validAfter, validBefore, nonce, v, r, s)")
		fmt.Println("    4. Facilitator 지갑이 가스비 지불")
		fmt.Println("    5. eth_sendRawTransaction → Base Sepolia 제출")
		fmt.Println("    6. Receipt 대기 (2초 간격 polling)")
	}
	fmt.Println()

	prHeader := paidResp.Header.Get("Payment-Response")
	if prHeader == "" {
		prHeader = paidResp.Header.Get("PAYMENT-RESPONSE")
	}

	if prHeader != "" {
		prDecoded, _ := base64.StdEncoding.DecodeString(prHeader)
		prettyPrint("  PAYMENT-RESPONSE 헤더 (정산 결과)", prDecoded)

		var settleResult struct {
			Success     bool   `json:"success"`
			Transaction string `json:"transaction"`
			Network     string `json:"network"`
			Payer       string `json:"payer"`
		}
		json.Unmarshal(prDecoded, &settleResult)

		fmt.Println()
		if settleResult.Success {
			fmt.Println("  ✅ 온체인 정산 완료!")
			fmt.Printf("     Tx: https://sepolia.basescan.org/tx/%s\n", settleResult.Transaction)
		} else {
			fmt.Println("  ❌ 정산 실패")
		}
	} else {
		fmt.Println("  ⏳ PAYMENT-RESPONSE 헤더 없음 (정산이 아직 진행 중이거나 별도 확인 필요)")
	}
	pause()

	// ─────────────────────────────────────────────────────────
	// STEP 10: Final Balances
	// ─────────────────────────────────────────────────────────
	step(10, "정산 후 잔액 확인")

	// Wait a moment for chain to settle
	time.Sleep(3 * time.Second)

	showBalances(ctx, ethClient, cfg.USDCAddress, "AFTER",
		wallet{"Facilitator", facilitatorAddr},
		wallet{"Client", clientSigner.Address()},
		wallet{"PAY_TO", resCfg.PayToAddress},
	)

	fmt.Println()
	fmt.Println("  자금 이동 요약:")
	if transferMethod == "permit2" {
		fmt.Println("    Client  → PAY_TO:       0.1 USDC (Permit2 → x402Permit2Proxy.settle)")
	} else {
		fmt.Println("    Client  → PAY_TO:       0.1 USDC (transferWithAuthorization)")
	}
	fmt.Println("    Facilitator:             가스비만 소모 (USDC 변동 없음)")
	fmt.Println()

	banner("Demo Complete")
}

// ─── Helpers ─────────────────────────────────────────────────

type wallet struct {
	name string
	addr string
}

func showBalances(ctx context.Context, client *ethclient.Client, usdcAddr string, label string, wallets ...wallet) {
	usdc := common.HexToAddress(usdcAddr)
	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`))
	ethDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	usdcDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))

	fmt.Printf("  ┌─────────── Balances (%s) ───────────┐\n", label)
	for _, w := range wallets {
		addr := common.HexToAddress(w.addr)
		ethBal, _ := client.BalanceAt(ctx, addr, nil)
		data, _ := erc20ABI.Pack("balanceOf", addr)
		result, _ := client.CallContract(ctx, ethereum.CallMsg{To: &usdc, Data: data}, nil)
		out, _ := erc20ABI.Unpack("balanceOf", result)

		ethF := new(big.Float).Quo(new(big.Float).SetInt(ethBal), ethDiv)
		usdcF := new(big.Float)
		if len(out) > 0 {
			usdcF = new(big.Float).Quo(new(big.Float).SetInt(out[0].(*big.Int)), usdcDiv)
		}
		fmt.Printf("  │ %-12s %s  │\n", w.name+":", w.addr[:10]+"..."+w.addr[len(w.addr)-4:])
		fmt.Printf("  │   ETH:  %-14s USDC: %-12s │\n", ethF.Text('f', 6), usdcF.Text('f', 6))
	}
	fmt.Println("  └──────────────────────────────────────────┘")
}

func prettyPrint(label string, data []byte) {
	fmt.Println(label + ":")
	var v interface{}
	if json.Unmarshal(data, &v) == nil {
		formatted, _ := json.MarshalIndent(v, "    ", "  ")
		fmt.Printf("    %s\n", formatted)
	} else {
		fmt.Printf("    %s\n", data)
	}
}

func banner(text string) {
	line := strings.Repeat("═", len(text)+4)
	fmt.Printf("\n  ╔%s╗\n", line)
	fmt.Printf("  ║  %s  ║\n", text)
	fmt.Printf("  ╚%s╝\n\n", line)
}

func step(n int, desc string) {
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  STEP %d: %s\n", n, desc)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
}

func pause() {
	fmt.Println()
	fmt.Print("  [Enter를 누르면 다음 단계로 →] ")
	fmt.Scanln()
	fmt.Println()
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func must(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}
