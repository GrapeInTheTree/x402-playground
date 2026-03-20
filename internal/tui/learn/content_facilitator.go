package learn

// ContentFacilitator returns the markdown content explaining the facilitator role and gas sponsorship.
func ContentFacilitator() string {
	return `# Facilitator Role & Gas Sponsorship

## Overview

**Facilitator**는 x402 프로토콜의 핵심 중개자이다.
클라이언트의 오프체인 서명을 검증하고, 온체인에 제출하여 정산한다.

## 역할

1. **서명 검증** (/verify) — 오프체인에서 결제 유효성 확인
2. **온체인 정산** (/settle) — 실제 트랜잭션 제출
3. **가스비 대납** — Facilitator 지갑이 가스비 부담
4. **서비스 디스커버리** (/supported) — 지원 네트워크/방식 공개

## 자금 흐름

` + "```" + `
USDC:  Client ────────────────────> PayTo
       (서명자)                     (수신자)
         ↑                            ↑
         └── Facilitator는 중간에     └── USDC 직접 수신
             USDC를 건드리지 않음

ETH:   Facilitator ───> Network (가스비)
       (가스 지불자)
` + "```" + `

**중요**: Facilitator는 USDC를 절대 터치하지 않는다.
서명된 트랜잭션을 릴레이할 뿐이며, 가스비만 소모한다.

## API Endpoints

### POST /verify

오프체인 검증. 서명, 잔액, 시간, 시뮬레이션을 확인.

` + "```json" + `
// Request
{
  "x402Version": 2,
  "paymentPayload": "<json.RawMessage>",
  "paymentRequirements": "<json.RawMessage>"
}

// Response
{
  "isValid": true,
  "invalidReason": ""
}
` + "```" + `

### POST /settle

온체인 정산. 트랜잭션을 빌드하고 제출.

` + "```json" + `
// Response
{
  "success": true,
  "transaction": "0x99e4...dc24",
  "network": "eip155:84532",
  "payer": "0xFacilitator..."
}
` + "```" + `

### GET /supported

지원 정보 공개.

` + "```json" + `
{
  "supportedSchemes": {
    "exact": {
      "eip155:84532": { ... }
    }
  }
}
` + "```" + `

## 트랜잭션 빌드 과정

1. EIP-1559 트랜잭션 구성 (Dynamic fee)
2. Base fee 조회 → Gas fee cap = 2 * baseFee + tipCap
3. Gas 추정 + 20% 버퍼
4. Facilitator private key로 서명
5. eth_sendRawTransaction 제출
6. Receipt 대기 (2초 간격 polling, 60초 타임아웃)

## Wallet Requirements

| Wallet | Private Key | 역할 |
|--------|:-----------:|------|
| Facilitator | 필요 | 가스비 지불, 트랜잭션 제출 |
| Client | 필요 | USDC 보유, 결제 서명 |
| PayTo | **불필요** | USDC 수신만 — 아무 EVM 주소 |

## Security Considerations

- Facilitator는 **Client 자금에 접근 불가**
- 서명은 특정 금액/수신자/시간에만 유효
- Replay protection: nonce (EIP-3009) 또는 Permit2 nonce
- 시뮬레이션: eth_call로 사전 실행 확인
`
}
