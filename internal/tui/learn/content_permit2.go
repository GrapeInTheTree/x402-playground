package learn

func ContentPermit2() string {
	return `# Permit2: Universal ERC-20 Payments

## Overview

**Permit2**는 Uniswap이 개발한 토큰 승인 프로토콜이다.
EIP-3009와 달리, **모든 ERC-20 토큰**에서 오프체인 서명 기반 전송이 가능하다.

## EIP-3009 vs Permit2

| | EIP-3009 | Permit2 |
|---|---|---|
| 토큰 요구사항 | EIP-3009 구현 필수 | 모든 ERC-20 |
| 사전 설정 | 없음 | approve(Permit2, amount) 필요 |
| 중개 컨트랙트 | 없음 (토큰 직접) | Permit2 + x402Permit2Proxy |
| 서명 대상 | USDC 도메인 | Permit2 도메인 |
| 수수료 | 가스비만 | 가스비만 |

## Architecture

` + "```" + `
Client (서명자)
  │
  │ 1. approve(Permit2, amount)     ← 1회만 필요
  │
  ▼
Permit2 (0x000...22D4)
  │
  │ 2. permitWitnessTransferFrom()  ← Facilitator가 호출
  │
  ▼
x402Permit2Proxy (0x4020...0001)
  │
  │ 3. settle()                     ← 토큰 이동 실행
  │
  ▼
Token: Client → PayTo              ← USDC 전송 완료
` + "```" + `

## Contract Addresses (CREATE2 — 모든 EVM 체인 동일)

| Contract | Address |
|----------|---------|
| Permit2 | ` + "`0x000000000022D473030F116dDEE9F6B43aC78BA3`" + ` |
| x402Permit2Proxy | ` + "`0x402085c248EeA27D92E8b30b2C58ed07f9E20001`" + ` |

CREATE2 deterministic deployment — 주소는 체인 불문이나, 실제 배포 여부는 다를 수 있음.

## EIP-712 TypedData 구조

` + "```" + `
Domain {
  name:              "Permit2"
  chainId:           84532
  verifyingContract: 0x000000000022D473...  ← Permit2 컨트랙트
}

Message: PermitWitnessTransferFrom {
  permitted: {
    token:   0x036CbD...    ← USDC 주소
    amount:  100000          ← 0.1 USDC
  }
  spender:   0x402085...     ← x402Permit2Proxy
  nonce:     12345           ← Permit2 nonce
  deadline:  1718000000      ← 만료 시간
  witness: {
    to:         0xPayTo...   ← 수신자
    validAfter: 0            ← 유효 시작
  }
}
` + "```" + `

## Prerequisites

1. Client가 USDC에 대해 Permit2 contract에 approve() 필요
2. Permit2와 x402Permit2Proxy가 대상 체인에 배포되어 있어야 함
3. ` + "`ASSET_TRANSFER_METHOD=permit2`" + ` 환경변수 설정

## MoneyParser 동작

Resource Server의 MoneyParser가 ` + "`extra.assetTransferMethod = \"permit2\"`" + `를 주입.
이를 통해 Client SDK가 EIP-3009 대신 Permit2 서명을 생성한다.
`
}
