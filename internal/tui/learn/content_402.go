package learn

// Content402 returns the markdown content explaining HTTP 402 and payment headers.
func Content402() string {
	return `# HTTP 402 & Payment Headers

## HTTP 402 Payment Required

RFC 7231에 정의된 HTTP 상태 코드. 오랫동안 "reserved for future use"였으나,
x402 프로토콜이 이를 실제 결제 메커니즘으로 구현했다.

## Request / Response Flow

` + "```" + `
Client                    Resource Server              Facilitator
  |                            |                           |
  |──GET /weather──────────────>|                           |
  |                            |                           |
  |<──402 + PAYMENT-REQUIRED───|                           |
  |                            |                           |
  |──GET /weather──────────────>|                           |
  |  + PAYMENT-SIGNATURE       |──POST /verify────────────>|
  |                            |<─────────────────────────|
  |<──200 + data───────────────|                           |
  |  + PAYMENT-RESPONSE        |──POST /settle────────────>|
  |                            |<─────────────────────────|
` + "```" + `

## PAYMENT-REQUIRED Header (402 응답)

Base64 인코딩된 JSON. 서버가 수락하는 결제 조건을 명시한다.

` + "```json" + `
{
  "accepts": [
    {
      "scheme": "exact",
      "network": "eip155:84532",
      "maxAmountRequired": "100000",
      "resource": "https://api.example.com/weather",
      "description": "Current weather data",
      "mimeType": "application/json",
      "payTo": "0x1234...abcd",
      "asset": "0x036CbD53842c5426634e7929541eC2318f3dCF7e",
      "extra": {
        "name": "USDC",
        "version": "2"
      }
    }
  ]
}
` + "```" + `

### Key Fields

| Field | Description |
|-------|-------------|
| ` + "`scheme`" + ` | 결제 방식 — ` + "`exact`" + ` (정확한 금액) |
| ` + "`network`" + ` | CAIP-2 네트워크 ID (예: eip155:84532) |
| ` + "`maxAmountRequired`" + ` | 최대 결제 금액 (USDC: 6 decimals, 100000 = $0.1) |
| ` + "`payTo`" + ` | 결제 수신 주소 |
| ` + "`asset`" + ` | 토큰 컨트랙트 주소 (USDC) |
| ` + "`extra.name`" + ` | EIP-712 도메인 name (서명에 필요) |
| ` + "`extra.assetTransferMethod`" + ` | ` + "`permit2`" + `이면 Permit2 방식 사용 |

## PAYMENT-SIGNATURE Header (재요청)

Client가 생성한 결제 서명. Base64 인코딩된 PaymentPayload JSON.

## PAYMENT-RESPONSE Header (200 응답)

Facilitator의 정산 결과. Base64 인코딩.

` + "```json" + `
{
  "success": true,
  "transaction": "0x99e4...dc24",
  "network": "eip155:84532",
  "payer": "0xABCD...1234"
}
` + "```" + `
`
}
