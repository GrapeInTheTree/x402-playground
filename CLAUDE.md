# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

Go-based x402 payment protocol demo — tested and verified on **Base Sepolia** with real USDC transfers. Supports both **EIP-3009** and **Permit2** transfer methods. Four components:
- **Facilitator Server** — Verifies and settles payments on-chain (EIP-3009 or Permit2, auto-detected)
- **Resource Server** — Protected APIs that return HTTP 402 with payment requirements
- **Client CLI** — Signs payment payloads (EIP-3009 or Permit2) and handles automatic payment flow
- **Explorer TUI** — Bubbletea-based interactive learning tool for the x402 protocol (Learn, Explore, Practice, Dashboard)

Chain-agnostic: configure via environment variables. Verified working on Base Sepolia (eip155:84532).

## Build & Run

```bash
make build                    # Build all four binaries
make test                     # Run all unit tests
make run-facilitator          # go run ./cmd/facilitator
make run-resource             # go run ./cmd/resource
make run-client               # go run ./cmd/client
make run-explorer             # Interactive TUI explorer (Home menu)
make run-demo                 # Practice mode: EIP-3009 flow
make run-demo-permit2         # Practice mode: Permit2 flow
make run-learn                # Learn mode: interactive coding quiz (Go + Solidity)
make run-dashboard            # Dashboard: wallet balances

# Explorer with flags
go run ./cmd/explorer --mode=learn           # Jump directly to coding quiz
go run ./cmd/explorer --mode=practice --flow=eip3009  # EIP-3009 practice
go run ./cmd/explorer --mode=dashboard       # Dashboard only

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

Explorer TUI (cmd/explorer) — Interactive learning & practice tool
  ├── Learn     — x402 protocol concepts (6 topics, markdown)
  ├── Explore   — Data structure inspector (headers, EIP-712, on-chain)
  ├── Practice  — 10-step payment flow (EIP-3009, Permit2, side-by-side)
  └── Dashboard — Wallet balances (live from chain)
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
- `internal/config/config.go` — Environment variable loading for all four components (Facilitator, Resource, Client, Explorer)
- `cmd/facilitator/main.go` — Wires SDK facilitator + EVM exact scheme + Gin router
- `cmd/resource/main.go` — Wires SDK Gin middleware + facilitator HTTP client + custom MoneyParser
- `cmd/client/main.go` — Wires SDK client signer + HTTP RoundTripper for auto-payment
- `cmd/explorer/main.go` — Bubbletea TUI entry point with `--mode` and `--flow` flags
- `cmd/balance/main.go` — Utility to check ETH/USDC balances on current network
- `internal/demo/` — Extracted protocol logic: types, balance queries, header decoding, flow execution
- `internal/quiz/` — Quiz engine: questions (Go + Solidity), runner (`go test` + `forge test`), types
- `internal/tui/` — TUI framework: app routing, components, pages (home, learn, explore, practice, dashboard)

### Explorer TUI Architecture

The TUI uses [bubbletea](https://github.com/charmbracelet/bubbletea) (Elm architecture). `RootModel` in `internal/tui/app.go` routes between pages:

- **Layout**: OpenCode-style — header bar (app name + page tab) + full-screen content + status bar. No outer border. Pages fill the available space. RootModel adds PaddingLeft(2)+PaddingRight(2), so pages must account for 4 chars less width. Pages using bordered panels must subtract border width (2 per box) from available space: `lipgloss.Width(n)` renders outer width `n+2`.
- `SubModel` interface — each page implements `Init()`, `Update()`, `View()`, `SetSize()`
- `SubModelFactory` — lazy initialization of pages on first `WindowSizeMsg`
- Navigation: `NavigateMsg` to go to a page, `BackMsg` to return home
- Help overlay: `?` toggles keyboard shortcuts (rendered as centered overlay)
- Minimum terminal size: 60x20 — shows resize prompt if too small
- Home page: title + menu centered vertically and horizontally (landing screen)
- Other pages: content fills from top, responsive to terminal width
- CLI flags `--mode` (learn/explore/practice/dashboard) and `--flow` (eip3009/permit2/sidebyside)

Key TUI packages:
- `internal/tui/components/` — Reusable: Menu (with highlight bar), Panel, TriPanel, FieldExplorer, JSONView, Progress, StatusBar
- `internal/tui/learn/` — Interactive coding quiz with split-panel layout (left: scrollable question list, right: description preview). 48 problems (17 Go + 31 Solidity) across 12 groups (5 Go levels + 7 Solidity modules). Pass/fail animation (2s pulsing) after test submission. Viewport scrolling with group header awareness. Uses `tea.ExecProcess` to launch `$EDITOR` (nvim→vim→nano fallback), then auto-grades via `go test` or `forge test`. Shared `*quiz.QuizProgress` synced to Dashboard.
- `internal/quiz/` — Quiz engine:
  - `types.go` — `Question` (with Lang, Category, Difficulty), `Result`, `Score`, `ModuleProgress`, `QuizProgress`
  - `runner.go` — `Runner` supports Go (`go test` in temp module) and Solidity (`forge test` in temp Foundry project with forge-std)
  - `questions.go` — 17 Go questions across 5 levels: hex validation, USDC conversion, base64/JSON, ERC-20 selectors, EIP-712 domains/type hashes, EIP-3009 fields, EIP-2612 permits, x402 headers/verify/settlement, Permit2 flow, payment state machine, ERC-8004 agent registration/global ID/WAD encoding/x402 integration
  - `questions_solidity.go` — 23 Solidity questions across 6 modules:
    - M1 Foundations (4): types/vars, functions/visibility, control flow/events, msg/block globals
    - M2 ERC-20 (4): basic token, approval system, events, metadata
    - M3 Signatures (4): keccak256, ecrecover, EIP-712 domain separator, struct hash
    - M4 Gasless (3): EIP-2612 permit, EIP-3009 transferWithAuth, nonce management
    - M5 Advanced (5): CREATE2, Permit2 SignatureTransfer, minimal proxy, access control, reentrancy guard
    - M6 x402 (3): settlement contract, Permit2 proxy, full verify+settle
    - M7 ERC-8004 (8): agent identity registry, metadata storage, reputation feedback, self-feedback prevention, validation request/response, EIP-712 wallet verification, WAD reputation summary, x402+reputation combined
- `internal/tui/explore/` — PAYMENT-REQUIRED field explorer, EIP-712 TypedData inspector (Tab to switch EIP-3009/Permit2), EIP-3009 vs Permit2 comparison, on-chain state viewer
- `internal/tui/practice/` — **Live execution** of 10-step flow via `LiveExecutor`. Each step fires async `tea.Cmd` with real HTTP/SDK calls. 3-column panel (Client/Resource/Facilitator) with spinner animation during execution. `stepManager` tracks state (pending/running/done/error). Press `n` to execute, `p` to review previous.
- `internal/tui/dashboard/` — Side-by-side layout: left panel shows live wallet balances from chain via RPC (spinner, `r` to refresh), right panel shows quiz progress with uniform `━`/`─` bars per module and overall percentage. Shared `*quiz.QuizProgress` pointer updated by Learn page.

Shared protocol logic extracted to `internal/demo/`:
- `types.go` — `FlowState`, `WalletInfo`, `WalletBalance`, `DecodedPaymentRequired`, `AcceptItem`
- `balance.go` — `QueryBalances()`, `QueryAllowance()` via ethclient + ERC-20 ABI (exported `ERC20BalanceOfABI`/`ERC20AllowanceABI` to avoid duplication)
- `decoder.go` — `DecodePaymentRequiredHeader()`, `DecodeBase64JSON()`, `FormatJSON()`, `ParseAcceptItem()`
- `flow.go` — `FlowExecutor` with step methods for HTTP calls to facilitator/resource
- `live.go` — `LiveExecutor` wraps the full 10-step flow with x402 SDK integration. Creates `Newx402Client` + `NewClientSignerFromPrivateKey` for payment signature generation. Accumulates state across steps (402 headers, payload bytes, selected requirements). Each `RunStep(ctx, step)` returns display text or error.

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
| `CLIENT_PRIVATE_KEY` | client, explorer | required |
| `RPC_URL` | facilitator, client, explorer | `https://sepolia.base.org` |
| `NETWORK` | all | `eip155:84532` |
| `USDC_ADDRESS` | resource, client, explorer | Base Sepolia USDC |
| `FACILITATOR_URL` | resource, explorer | required |
| `FACILITATOR_PORT` | facilitator | `4022` |
| `RESOURCE_PORT` | resource | `4021` |
| `PAY_TO_ADDRESS` | resource, explorer | required (no private key needed) |
| `RESOURCE_URL` | client, explorer | required |
| `ASSET_TRANSFER_METHOD` | all | `eip3009` |
| `LOG_LEVEL` | all | `info` |

Explorer's Learn and Explore modes work without private keys or server URLs (best-effort config).
Learn mode requires `$EDITOR` (nvim/vim/nano) and `go` for Go quizzes; additionally `forge` (Foundry) for Solidity quizzes.
Practice mode requires `CLIENT_PRIVATE_KEY`, `FACILITATOR_URL`, `RESOURCE_URL` + running servers for live execution.
Dashboard requires `RPC_URL` and `PAY_TO_ADDRESS` for balance queries (no servers needed).

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
- [ERC-8004: Trustless Agents](https://eips.ethereum.org/EIPS/eip-8004)
- [Coinbase x402 Go SDK](https://pkg.go.dev/github.com/coinbase/x402/go)
