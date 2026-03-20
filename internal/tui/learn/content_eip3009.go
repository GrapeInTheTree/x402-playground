package learn

// ContentEIP3009 returns the markdown content explaining EIP-3009 transferWithAuthorization.
func ContentEIP3009() string {
	return `# EIP-3009: transferWithAuthorization

## Overview

**EIP-3009**는 USDC 등 특정 토큰에 내장된 메타 트랜잭션 표준이다.
토큰 보유자가 오프체인 서명을 생성하면, 제3자(Facilitator)가 이를 온체인에 제출하여
토큰을 이동시킬 수 있다.

## 핵심 개념

- **토큰 보유자**는 가스비를 지불하지 않음
- **Facilitator**가 서명을 온체인에 제출하고 가스비 대납
- **USDC, EURC** 등 Circle 토큰에 기본 내장

## EIP-712 TypedData 구조

` + "```" + `
Domain {
  name:              "USDC"          ← 토큰 컨트랙트의 name() 반환값
  version:           "2"             ← 토큰 컨트랙트의 EIP-712 version
  chainId:           84532           ← Base Sepolia
  verifyingContract: 0x036CbD...     ← USDC 컨트랙트 주소
}

Message: TransferWithAuthorization {
  from:        0xClient...          ← 서명자 주소
  to:          0xPayTo...           ← 수신자 주소
  value:       100000               ← 0.1 USDC (6 decimals)
  validAfter:  0                    ← 즉시 유효
  validBefore: 1718000000           ← 만료 시간 (Unix timestamp)
  nonce:       0x1234...            ← 랜덤 32바이트 (이중 사용 방지)
}
` + "```" + `

## 온체인 실행

` + "```solidity" + `
// FiatTokenV2.sol (USDC contract)
function transferWithAuthorization(
    address from,
    address to,
    uint256 value,
    uint256 validAfter,
    uint256 validBefore,
    bytes32 nonce,
    uint8 v,
    bytes32 r,
    bytes32 s
) external
` + "```" + `

Facilitator가 이 함수를 호출 → USDC가 from → to로 이동.

## 검증 항목 (Facilitator /verify)

1. EIP-712 서명 복원 → from 주소 일치 확인
2. to == payTo (수신자 주소 일치)
3. value >= maxAmountRequired (금액 충분)
4. validAfter <= now <= validBefore (시간 유효)
5. nonce 미사용 (이중 결제 방지)
6. Client USDC 잔액 >= value
7. eth_call 시뮬레이션 (실제 실행 가능 확인)

## 제한사항

- EIP-3009를 구현한 토큰만 사용 가능 (USDC, EURC 등)
- 일반 ERC-20 토큰은 사용 불가 → Permit2 필요

## Domain Name 주의사항

EIP-712 domain의 ` + "`name`" + ` 필드는 토큰 컨트랙트의 ` + "`name()`" + ` 반환값과 정확히 일치해야 한다.
Base Sepolia USDC는 ` + "`\"USDC\"`" + `를 반환한다 (` + "`\"USD Coin\"`" + `이 아님).
불일치 시 ` + "`FiatTokenV2: invalid signature`" + ` 오류 발생.
`
}
