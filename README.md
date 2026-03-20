<p align="center">
  <h1 align="center">x402-demo</h1>
  <p align="center">
    <strong>HTTP-native micropayments for any EVM blockchain.</strong>
  </p>
  <p align="center">
    A production-grade Go implementation of the <a href="https://x402.org/">x402 payment protocol</a> — three independently deployable components that demonstrate the full payment lifecycle with real on-chain USDC transfers.
  </p>
</p>

---

## Why x402-demo?

Most x402 examples are minimal snippets. This project is a **complete, working reference implementation** with real money tested on Base Sepolia:

- **Full Lifecycle** — Facilitator, Resource Server, and Client CLI working end-to-end
- **Dual Transfer Methods** — EIP-3009 and Permit2, switchable via one environment variable
- **Interactive Demo** — 10-step walkthrough that explains every protocol detail as it executes
- **Chain-Agnostic** — Configure any EVM chain via environment variables
- **28 Unit Tests** — Config, handlers, signer, all covered

## How It Works

```
Client CLI ──HTTP──> Resource Server ──HTTP──> Facilitator Server ──RPC──> EVM Chain
cmd/client           cmd/resource              cmd/facilitator
```

| Step | What happens |
|:---:|-------------|
| 1 | Client sends a normal HTTP request to a protected endpoint |
| 2 | Resource Server responds with **HTTP 402** + `PAYMENT-REQUIRED` header |
| 3 | Client creates an EIP-712 signature (EIP-3009 or Permit2) and retries with `PAYMENT-SIGNATURE` |
| 4-5 | Resource Server delegates to Facilitator via `/verify` and `/settle` |
| 6 | Facilitator submits the tx on-chain, pays gas, returns tx hash |
| 7 | Client receives API response + `PAYMENT-RESPONSE` with settlement hash |

<details>
<summary>Detailed flow diagram</summary>

```
                          x402 Payment Flow
                          ─────────────────

   ┌────────────┐                ┌──────────────────┐                ┌────────────────────┐
   │ Client CLI │  1. GET /api   │  Resource Server  │                │ Facilitator Server │
   │            │ ──────────────>│                   │                │                    │
   │            │                │  "No payment      │                │                    │
   │            │  2. HTTP 402   │   header found"   │                │                    │
   │            │ <──────────────│                   │                │                    │
   │            │  + PAYMENT-    │                   │                │                    │
   │            │    REQUIRED    │                   │                │                    │
   │            │                │                   │                │                    │
   │  Signs     │                │                   │                │                    │
   │  EIP-3009  │  3. GET /api   │                   │  4. POST       │                    │
   │  or        │ ──────────────>│  Parses header,   │ ──/verify────> │  Recovers signer,  │
   │  Permit2   │  + PAYMENT-    │  forwards to      │                │  checks balance,   │
   │  payload   │    SIGNATURE   │  facilitator      │  5. POST       │  simulates call    │  ┌───────────┐
   │            │                │                   │ ──/settle────> │                    │  │ EVM Chain │
   │            │                │                   │                │  Builds EIP-1559   │  │           │
   │            │                │                   │                │  tx, settles via   │──│ USDC      │
   │            │                │                   │  6. tx hash    │  EIP-3009 or       │  │ transfer  │
   │            │  7. HTTP 200   │  Returns API data │ <─────────────│  Permit2 Proxy     │  │           │
   │            │ <──────────────│  + PAYMENT-       │                │                    │  │           │
   │            │  + response    │    RESPONSE       │                │                    │  └───────────┘
   └────────────┘                └──────────────────┘                └────────────────────┘
```

</details>

## Wallet Roles

Three wallets, three roles. `PAY_TO_ADDRESS` does **not** need a private key:

```
Client Wallet                          PAY_TO Address
(signs payment)                        (receives USDC)
      │                                      ▲
      │  0.1 USDC                            │
      │  (EIP-3009: transferWithAuth)        │
      │  (Permit2: via x402Permit2Proxy)     │
      └──────────────────────────────────────┘
                        │
               Facilitator Wallet
               (pays gas only, never touches USDC)
```

| Wallet | Private Key? | Role |
|--------|:---:|------|
| **Facilitator** | Yes | Pays gas, submits payment tx on-chain |
| **Client** | Yes | Holds USDC, signs payment authorizations |
| **PAY_TO** | **No** | Receives USDC — any EVM address works |

## Quick Start

### 1. Clone and build

```bash
git clone https://github.com/GrapeInTheTree/x402-demo.git
cd x402-demo
make build
```

### 2. Configure

```bash
cp .env.example .env
```

Edit `.env` with your values:

```bash
FACILITATOR_PRIVATE_KEY=0x...   # Pays gas for on-chain settlement
CLIENT_PRIVATE_KEY=0x...        # Holds USDC, signs authorizations
PAY_TO_ADDRESS=0x...            # Receives USDC (no private key needed)
ASSET_TRANSFER_METHOD=eip3009   # or permit2
```

### 3. Run

```bash
# Terminal 1 — Facilitator (port 4022)
make run-facilitator

# Terminal 2 — Resource Server (port 4021)
make run-resource

# Terminal 3 — Interactive demo
make run-demo            # EIP-3009 mode
make run-demo-permit2    # Permit2 mode
```

### 4. Or use the simple client

```bash
make run-client
```

```
→ GET http://localhost:4021/weather
← 200 OK

💰 Payment Settlement:
   Success:     true
   Transaction: 0x99e49093...faedc24
   Network:     eip155:84532

Response:
{
  "city": "New York",
  "temperature": 25,
  "condition": "Windy"
}
```

## Transfer Methods

Two methods, switchable via `ASSET_TRANSFER_METHOD`:

| | EIP-3009 (default) | Permit2 |
|---|---|---|
| **On-chain call** | `USDC.transferWithAuthorization(...)` | `x402Permit2Proxy.settle(...)` |
| **Token requirement** | Must implement EIP-3009 | Any ERC-20 |
| **Setup** | None | Client must `approve(Permit2)` once |
| **Gas cost** | Lower (direct call) | Slightly higher (proxy hop) |
| **EIP-712 domain** | Token contract | Permit2 contract |

Permit2 contracts use CREATE2 — same address on all EVM chains:

| Contract | Address |
|----------|---------|
| Permit2 (Uniswap) | `0x000000000022D473030F116dDEE9F6B43aC78BA3` |
| x402Permit2Proxy (Coinbase) | `0x402085c248EeA27D92E8b30b2C58ed07f9E20001` |

## Chain Compatibility

| Chain | EIP-3009 | Permit2 | Status |
|-------|:---:|:---:|--------|
| Base Sepolia | Yes | Yes | **Verified working** |
| Base Mainnet | Yes | Yes | SDK built-in |
| Polygon | Yes | Yes | SDK built-in |
| Arbitrum | Yes | Yes | SDK built-in |
| Chiliz Mainnet | **No** | **Unverified** | Bridged USDC, no EIP-3009 |
| Chiliz Spicy | **No** | **No** | No USDC deployed |

To switch chains, change three env vars:

```bash
NETWORK=eip155:8453
RPC_URL=https://mainnet.base.org
USDC_ADDRESS=0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913
```

## Protected Endpoints

| Endpoint | Price | Description |
|----------|:-----:|-------------|
| `GET /weather` | $0.10 | Random city weather data |
| `GET /joke` | $0.10 | Programming jokes |
| `GET /premium-data` | $0.10 | Mock analytics report |
| `GET /health` | Free | Server health check |

## Development

```bash
make build           # Compile all binaries
make test            # Run all 28 unit tests
make lint            # Run golangci-lint
make clean           # Remove compiled binaries
```

| Package | Tests | What's tested |
|---------|:-----:|---------------|
| `internal/config` | 10 | Env loading, defaults, validation, Permit2 config |
| `internal/facilserver` | 9 | Verify/settle/supported handlers with mock |
| `internal/server` | 4 | API response structure, health |
| `internal/signer` | 5 | Key parsing, address derivation, Close() zeroing |
| **Total** | **28** | |

## Docker

```bash
docker compose up --build                        # Facilitator + Resource
RESOURCE_URL=http://localhost:4021 make run-client # Client against Docker
docker compose down
```

<details>
<summary>Facilitator API Reference</summary>

### `POST /verify`

Validates a payment payload off-chain. Checks signature, balance, timestamps, nonce.

```json
// Request
{ "x402Version": 2, "paymentPayload": { ... }, "paymentRequirements": { ... } }

// Response
{ "isValid": true, "payer": "0x47322Ca2..." }
```

### `POST /settle`

Executes payment on-chain. Auto-detects EIP-3009 vs Permit2 payload type.

```json
{ "success": true, "transaction": "0x99e49093...", "network": "eip155:84532", "payer": "0x47322Ca2..." }
```

### `GET /supported`

```json
{ "kinds": [{ "x402Version": 2, "scheme": "exact", "network": "eip155:84532" }], "signers": { "eip155:*": ["0x23fbdE5A..."] } }
```

### `GET /health`

```json
{ "status": "ok", "service": "facilitator", "network": "eip155:84532", "address": "0x23fbdE5A..." }
```

</details>

<details>
<summary>Protocol Details</summary>

### HTTP Headers (V2)

| Header | Direction | Purpose |
|--------|-----------|---------|
| `PAYMENT-REQUIRED` | Server → Client | Base64-encoded payment requirements (in 402 response) |
| `PAYMENT-SIGNATURE` | Client → Server | Base64-encoded signed payment payload |
| `PAYMENT-RESPONSE` | Server → Client | Base64-encoded settlement result |

### EIP-3009 Mode (default)

1. Client signs EIP-712 typed data authorizing a token transfer
2. Facilitator calls `transferWithAuthorization(from, to, value, validAfter, validBefore, nonce, v, r, s)`
3. Nonces are random 32-byte values (not sequential)
4. Facilitator pays gas; USDC goes directly Client → PAY_TO

### Permit2 Mode

1. Client signs EIP-712 for Permit2 `PermitWitnessTransferFrom` (domain = Permit2 contract)
2. Facilitator calls `x402Permit2Proxy.settle()` → `Permit2.permitWitnessTransferFrom()`
3. Requires one-time `approve(Permit2, amount)` from Client wallet
4. Works with any ERC-20 token

</details>

<details>
<summary>Configuration Reference</summary>

| Variable | Used by | Default | Description |
|----------|---------|---------|-------------|
| `FACILITATOR_PRIVATE_KEY` | facilitator | *required* | Wallet that pays gas for settlement |
| `CLIENT_PRIVATE_KEY` | client | *required* | Wallet that holds USDC and signs authorizations |
| `RPC_URL` | facilitator, client | `https://sepolia.base.org` | JSON-RPC endpoint |
| `NETWORK` | all | `eip155:84532` | CAIP-2 network identifier |
| `USDC_ADDRESS` | resource, client | `0x036CbD53...` | Token contract address |
| `FACILITATOR_URL` | resource | *required* | Facilitator base URL |
| `PAY_TO_ADDRESS` | resource | *required* | Payment recipient (no private key needed) |
| `FACILITATOR_PORT` | facilitator | `4022` | HTTP listen port |
| `RESOURCE_PORT` | resource | `4021` | HTTP listen port |
| `RESOURCE_URL` | client | *required* | Resource Server base URL |
| `ENDPOINT_PATH` | client | `/weather` | Default API endpoint |
| `ASSET_TRANSFER_METHOD` | all | `eip3009` | `eip3009` or `permit2` |
| `LOG_LEVEL` | all | `info` | `debug`, `info`, `warn`, `error` |

</details>

<details>
<summary>Project Structure</summary>

```
x402-demo/
├── cmd/
│   ├── facilitator/main.go    Facilitator HTTP server
│   ├── resource/main.go       Resource HTTP server
│   ├── client/main.go         Client CLI
│   ├── demo/main.go           Interactive 10-step walkthrough
│   └── balance/main.go        Wallet balance checker
├── internal/
│   ├── config/                Environment variable loading + validation
│   ├── facilserver/           Facilitator HTTP handlers (/verify, /settle, /supported)
│   ├── server/                Resource Server routes + API handlers
│   └── signer/                FacilitatorEvmSigner (EIP-1559, EIP-712)
├── pkg/health/                Shared health check type
├── .env.example               Environment variable template
├── Dockerfile                 Multi-stage build
├── docker-compose.yml         Facilitator + Resource orchestration
└── Makefile                   Build, test, run targets
```

</details>

## Technology Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.24+ |
| x402 SDK | [coinbase/x402/go](https://github.com/coinbase/x402) v2.6.0 |
| EVM Client | [go-ethereum](https://github.com/ethereum/go-ethereum) v1.17 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| Payment Scheme | EIP-3009 / Permit2 (exact scheme) |
| Signatures | EIP-712 Typed Structured Data |
| Transactions | EIP-1559 (dynamic fee) |

## Verified Transactions

| Tx Hash | From | To | Amount |
|---------|------|-----|--------|
| [`0x99e4...dc24`](https://sepolia.basescan.org/tx/0x99e49093d0bb2805b2e1097a6c71336c73f5871a4e51ec2dacc733f51faedc24) | `0x4732...037b` | `0x23fb...b37A` | 0.1 USDC |
| [`0x6d3a...2445`](https://sepolia.basescan.org/tx/0x6d3a230de24f0650703fc87fd9b3f0cb19cc914e6530aca4512d5956f4fb2445) | `0x4732...037b` | `0xDBCb...07F5` | 0.1 USDC |

## Further Reading

- [x402 Protocol](https://x402.org/) | [Documentation](https://docs.x402.org/) | [GitHub](https://github.com/coinbase/x402)
- [Coinbase x402 Go SDK](https://pkg.go.dev/github.com/coinbase/x402/go)
- [EIP-3009: Transfer With Authorization](https://eips.ethereum.org/EIPS/eip-3009) | [EIP-712: Typed Structured Data](https://eips.ethereum.org/EIPS/eip-712)

## License

[MIT](LICENSE)
