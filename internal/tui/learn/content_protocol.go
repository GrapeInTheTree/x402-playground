package learn

func ContentProtocol() string {
	return `# What is x402?

## Overview

**x402**는 HTTP 402 상태 코드를 활용한 **머신-투-머신 결제 프로토콜**이다.

기존에는 사실상 사용되지 않던 HTTP 402 (Payment Required)를 실제 결제 프로토콜로 구현해,
API 서버가 직접 결제를 요구하고, 클라이언트가 자동으로 결제하는 흐름을 만든다.

## Why x402?

| 기존 방식 | x402 방식 |
|-----------|-----------|
| API Key 발급 → 월정액 | 요청마다 마이크로 결제 |
| 사전 계약 필요 | 익명 결제 가능 |
| 중앙화된 빌링 시스템 | 온체인 정산, 즉시 확인 |
| 사람이 결제 | AI 에이전트가 자동 결제 |

## Architecture (3-party model)

` + "```" + `
Client ──HTTP──> Resource Server ──HTTP──> Facilitator ──RPC──> EVM Chain
` + "```" + `

1. **Client** — API를 호출하고, 결제 서명을 생성
2. **Resource Server** — 보호된 API를 제공, 402 응답으로 결제 요구
3. **Facilitator** — 서명을 검증하고 온체인 정산 (가스비 대납)

## Payment Flow (10 steps)

1. Client가 API 호출 (결제 헤더 없음)
2. Resource Server가 **HTTP 402** + ` + "`PAYMENT-REQUIRED`" + ` 헤더 반환
3. Client가 결제 요구사항 파싱 (금액, 자산, 네트워크)
4. Client가 **EIP-712 서명** 생성 (오프체인, 가스비 없음)
5. Client가 ` + "`PAYMENT-SIGNATURE`" + ` 헤더와 함께 재요청
6. Resource Server가 Facilitator에 **/verify** 요청
7. Facilitator가 서명 + 잔액 + 시뮬레이션 검증
8. 검증 통과 → Resource Server가 데이터 반환
9. Resource Server가 Facilitator에 **/settle** 요청
10. Facilitator가 온체인 트랜잭션 제출 (가스비 자체 부담)

## Protocol Version

이 프로젝트는 **x402 V2** 프로토콜을 사용한다:

- ` + "`PAYMENT-SIGNATURE`" + ` (V2) — V1은 ` + "`X-PAYMENT`" + `
- ` + "`PAYMENT-REQUIRED`" + ` (V2)
- ` + "`PAYMENT-RESPONSE`" + ` (V2)

## Supported Networks

- Base Sepolia (eip155:84532) — 테스트넷, 검증 완료
- Base Mainnet (eip155:8453)
- Polygon (eip155:137)
`
}
