package learn

func ContentEIP712() string {
	return `# EIP-712: Typed Structured Data Signing

## Overview

**EIP-712**는 이더리움의 구조화된 데이터 서명 표준이다.
바이트 배열이 아닌, 사람이 읽을 수 있는 형태의 데이터에 서명한다.

## Why EIP-712?

| 방식 | 표시 | 안전성 |
|------|------|--------|
| eth_sign | ` + "`0x1a2b3c...`" + ` (해시) | 위험: 임의 데이터 서명 가능 |
| personal_sign | ` + "`\"Hello World\"`" + ` (문자열) | 중간: 텍스트만 |
| **EIP-712** | 구조화된 필드 표시 | 안전: 각 필드를 확인 가능 |

## 구조

EIP-712 서명은 세 가지 요소로 구성된다:

### 1. Domain Separator

서명이 어떤 컨트랙트에서 사용되는지 식별:

` + "```" + `
{
  name:              "USDC"          ← 컨트랙트 이름
  version:           "2"             ← 컨트랙트 버전
  chainId:           84532           ← 체인 ID
  verifyingContract: 0x036CbD...     ← 컨트랙트 주소
}
` + "```" + `

### 2. Type Definition

메시지의 구조를 정의:

` + "```" + `
TransferWithAuthorization(
  address from,
  address to,
  uint256 value,
  uint256 validAfter,
  uint256 validBefore,
  bytes32 nonce
)
` + "```" + `

### 3. Message

실제 서명할 데이터:

` + "```" + `
{
  from:        "0xClient...",
  to:          "0xPayTo...",
  value:       100000,
  validAfter:  0,
  validBefore: 1718000000,
  nonce:       "0xabcd..."
}
` + "```" + `

## 해싱 과정

` + "```" + `
1. domainSeparator = hashStruct(domain)
2. messageHash     = hashStruct(primaryType, message)
3. digest          = keccak256(0x19, 0x01, domainSeparator, messageHash)
4. signature       = ECDSA.sign(privateKey, digest)
` + "```" + `

## x402에서의 사용

### EIP-3009 (USDC 직접)
- Domain: USDC 컨트랙트
- PrimaryType: TransferWithAuthorization
- 서명 → v, r, s 추출 → USDC.transferWithAuthorization() 호출

### Permit2 (범용)
- Domain: Permit2 컨트랙트
- PrimaryType: PermitWitnessTransferFrom
- Witness: x402 확장 필드 (to, validAfter)
- 서명 → Permit2.permitWitnessTransferFrom() 호출

## 보안

- **Replay protection**: chainId + verifyingContract + nonce
- **Domain separation**: 다른 컨트랙트에서 서명 재사용 불가
- **Time bounds**: validAfter / validBefore로 유효 기간 제한
`
}
