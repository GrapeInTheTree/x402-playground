# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

Go-based x402 payment protocol demo — tested and verified on **Base Sepolia** with real USDC transfers. Supports both **EIP-3009** and **Permit2** transfer methods. Three independently deployable components:
- **Facilitator Server** — Verifies and settles payments on-chain (EIP-3009 or Permit2, auto-detected)
- **Resource Server** — Protected APIs that return HTTP 402 with payment requirements
- **Client CLI** — Signs payment payloads (EIP-3009 or Permit2) and handles automatic payment flow

Chain-agnostic: configure via environment variables. Verified working on Base Sepolia (eip155:84532).

## Build & Run

```bash
make build                    # Build all three binaries
make test                     # Run all 28 unit tests
make run-facilitator          # go run ./cmd/facilitator
make run-resource             # go run ./cmd/resource
make run-client               # go run ./cmd/client
make run-demo                 # Interactive 10-step payment flow demo (EIP-3009)
make run-demo-permit2         # Interactive demo with Permit2 transfer method
go test ./internal/config -run TestLoadFacilitator -v        # Single test
go test ./internal/facilserver -run TestHandleVerify -v     # Facilserver tests
go test ./internal/signer -run TestFacilitatorSigner -v     # Signer tests

# Utilities
go run ./cmd/balance          # Check wallet balances on current network

# Docker
docker compose up             # Facilitator + Resource server
```

## Architecture

```
Client CLI ──HTTP──> Resource Server ──HTTP──> Facilitator Server ──RPC──> EVM Chain
cmd/client           cmd/resource              cmd/facilitator
```

### Wallet Roles

Three distinct roles — `PAY_TO_ADDRESS` does NOT need a private key:

| Wallet | Private Key? | Role |
|--------|:---:|------|
| `FACILITATOR_PRIVATE_KEY` | Required | Pays gas, submits payment tx on-chain (EIP-3009 or Permit2) |
| `CLIENT_PRIVATE_KEY` | Required | Holds USDC, signs payment authorizations |
| `PAY_TO_ADDRESS` | **Not needed** | Receives USDC payments — any EVM address works |

USDC flows directly from Client → PAY_TO. The Facilitator never touches USDC — it only relays the signed transaction and pays gas.

### Key Code Locations

- `internal/signer/facilitator.go` — Custom `FacilitatorEvmSigner` implementation (~330 lines). Implements the SDK's `evm.FacilitatorEvmSigner` interface with `Close()` for key zeroing. The SDK does NOT provide a facilitator signer constructor.
- `internal/facilserver/iface.go` — `Facilitator` interface for testability (decouples handlers from SDK)
- `internal/facilserver/server.go` — Facilitator HTTP handlers (`/verify`, `/settle`, `/supported`)
- `internal/facilserver/errors.go` — Sentinel errors for request validation
- `internal/server/routes.go` — Payment-protected route definitions with pricing (currently $0.1 per endpoint)
- `internal/server/handlers.go` — Demo API handlers (weather, joke, premium-data)
- `internal/config/config.go` — Environment variable loading for all three components
- `cmd/facilitator/main.go` — Wires SDK facilitator + EVM exact scheme + Gin router
- `cmd/resource/main.go` — Wires SDK Gin middleware + facilitator HTTP client + custom MoneyParser
- `cmd/client/main.go` — Wires SDK client signer + HTTP RoundTripper for auto-payment
- `cmd/demo/main.go` — Interactive 10-step demo: 402 → sign → verify → settle → balance diff
- `cmd/balance/main.go` — Utility to check ETH/USDC balances on current network

### SDK Usage Pattern

The project uses the official **Coinbase x402 Go SDK** (`github.com/coinbase/x402/go` v2.6.0).

Key SDK types:
- `x402.Newx402Facilitator()` → `*x402.X402Facilitator`
- `evmfacilitator.NewExactEvmScheme(signer, config)` — EVM exact scheme for facilitator
- `evmserver.NewExactEvmScheme()` — EVM exact scheme for resource server (no signer needed)
- `evmclient.NewExactEvmScheme(signer, config)` — EVM exact scheme for client
- `evmsigner.NewClientSignerFromPrivateKey(key)` — Client-side EIP-712 signer
- `x402http.NewHTTPFacilitatorClient(config)` — HTTP client for calling facilitator
- `x402http.WrapHTTPClientWithPayment(httpClient, x402Client)` — Auto-payment RoundTripper
- `ginmw.X402Payment(config)` — Gin middleware for payment-gated routes

### MoneyParser: Network-Aware Price Resolution

`cmd/resource/main.go` registers a custom `MoneyParser` on `evmserver.NewExactEvmScheme()`. **Dual-mode behavior:**

- **EIP-3009 mode** (`ASSET_TRANSFER_METHOD=eip3009`, default):
  - SDK-supported networks (Base Sepolia, Base Mainnet, etc.): returns `nil` to delegate to SDK defaults.
  - Unknown networks: returns custom `AssetAmount` using `USDC_ADDRESS` from config.
- **Permit2 mode** (`ASSET_TRANSFER_METHOD=permit2`):
  - ALL networks use custom parser to inject `extra.assetTransferMethod = "permit2"`.
  - This overrides SDK defaults because the SDK's built-in configs default to EIP-3009.
  - The `assetTransferMethod` field in `Extra` tells the Client SDK to create a Permit2 payload instead of EIP-3009.

**Lesson learned:** The EIP-712 domain `name` must match the token contract's actual `name()` return value exactly. Base Sepolia USDC returns `"USDC"` (not `"USD Coin"`). A mismatch causes `FiatTokenV2: invalid signature` on-chain.

### Protocol Version

This project uses **x402 V2 protocol** exclusively:
- Payment header: `PAYMENT-SIGNATURE` (V2), not `X-PAYMENT` (V1)
- Requirements header: `PAYMENT-REQUIRED` (V2)
- Response header: `PAYMENT-RESPONSE` (V2)
- SDK methods: `Register()` (V2), not `RegisterV1()`

### Payload Forwarding

Facilitator endpoints receive payloads as `json.RawMessage` and pass `[]byte` directly to the SDK. Do NOT re-marshal payloads — this breaks signature verification.

## Environment Variables

| Variable | Component | Default |
|---|---|---|
| `FACILITATOR_PRIVATE_KEY` | facilitator | required |
| `CLIENT_PRIVATE_KEY` | client | required |
| `RPC_URL` | facilitator, client | `https://sepolia.base.org` |
| `NETWORK` | all | `eip155:84532` |
| `USDC_ADDRESS` | resource, client | Base Sepolia USDC |
| `FACILITATOR_URL` | resource | required |
| `FACILITATOR_PORT` | facilitator | `4022` |
| `RESOURCE_PORT` | resource | `4021` |
| `PAY_TO_ADDRESS` | resource | required (no private key needed) |
| `RESOURCE_URL` | client | required |
| `ASSET_TRANSFER_METHOD` | all | `eip3009` |
| `LOG_LEVEL` | all | `info` |

## Verified Test Results (Base Sepolia)

Successfully tested on Base Sepolia with real USDC transfers:
- Transaction: `0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24`
- Transaction: `0x6d3a230de24f0650703fc87fd9b3f0cb19cc914e6530aca4512d5956f4fb2445`

## Transfer Methods

Two transfer methods supported, switchable via `ASSET_TRANSFER_METHOD`:

| Method | Env Value | How it works | Token requirement |
|--------|-----------|-------------|-------------------|
| **EIP-3009** | `eip3009` (default) | `transferWithAuthorization` on USDC contract directly | Token must implement EIP-3009 |
| **Permit2** | `permit2` | `permitWitnessTransferFrom` via x402Permit2Proxy | Any ERC-20 (needs Permit2 approve) |

### Permit2 Contract Addresses (CREATE2 — same on all EVM chains)

| Contract | Address | Deployed by |
|----------|---------|-------------|
| Permit2 | `0x000000000022D473030F116dDEE9F6B43aC78BA3` | Uniswap |
| x402Permit2Proxy | `0x402085c248EeA27D92E8b30b2C58ed07f9E20001` | Coinbase |

Both use CREATE2 deterministic deployment. Address is chain-agnostic but actual deployment status varies — verify with `cast code <address> --rpc-url <rpc>` before use.

### Permit2 Prerequisites

- Client must `approve(Permit2, amount)` the token to the Permit2 contract
- Both Permit2 and x402Permit2Proxy must be deployed on the target chain
- SDK handles all EIP-712 signing and on-chain settlement automatically

## Chain Compatibility Notes

| Chain | EIP-3009 | Permit2 | Status |
|-------|:---:|:---:|--------|
| Base Sepolia (`eip155:84532`) | Supported | Supported | Verified working (EIP-3009) |
| Base Mainnet (`eip155:8453`) | Supported | Supported | SDK built-in |
| Polygon (`eip155:137`) | Supported | Supported | SDK built-in |
| Chiliz Mainnet (`eip155:88888`) | **No** | **Unverified** | Bridged USDC, no EIP-3009. Permit2/Proxy deployment status unknown |
| Chiliz Spicy (`eip155:88882`) | **No** | **No** | No USDC deployed |

## External References

- [x402 Protocol](https://x402.org/) | [GitHub](https://github.com/coinbase/x402)
- [EIP-3009](https://eips.ethereum.org/EIPS/eip-3009)
- [Coinbase x402 Go SDK](https://pkg.go.dev/github.com/coinbase/x402/go)
