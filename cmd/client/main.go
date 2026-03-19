package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	evmclient "github.com/coinbase/x402/go/mechanisms/evm/exact/client"
	evmsigner "github.com/coinbase/x402/go/signers/evm"

	"github.com/GrapeInTheTree/x402-demo/internal/config"
)

func main() {
	verbose := flag.Bool("v", false, "verbose output")
	endpoint := flag.String("endpoint", "", "API endpoint (overrides ENDPOINT_PATH)")
	url := flag.String("url", "", "resource server URL (overrides RESOURCE_URL)")
	flag.Parse()

	cfg, err := config.LoadClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	if *endpoint != "" {
		cfg.EndpointPath = *endpoint
	}
	if *url != "" {
		cfg.ResourceURL = *url
	}

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	// Create client signer
	clientSigner, err := evmsigner.NewClientSignerFromPrivateKey(cfg.PrivateKey)
	if err != nil {
		logger.Error("failed to create signer", "error", err)
		os.Exit(1)
	}

	logger.Info("client initialized",
		"address", clientSigner.Address(),
		"resourceURL", cfg.ResourceURL,
		"endpoint", cfg.EndpointPath,
		"assetTransferMethod", cfg.AssetTransferMethod,
	)

	// Create x402 client with EVM exact scheme
	x402Client := x402.Newx402Client().
		Register("eip155:*", evmclient.NewExactEvmScheme(clientSigner, &evmclient.ExactEvmSchemeConfig{
			RPCURL: cfg.RPCURL,
		}))

	// Wrap HTTP client with payment handling
	httpx402Client := x402http.Newx402HTTPClient(x402Client)
	httpClient := x402http.WrapHTTPClientWithPayment(&http.Client{Timeout: 120 * time.Second}, httpx402Client)

	// Make request
	targetURL := cfg.ResourceURL + cfg.EndpointPath
	fmt.Printf("→ GET %s\n", targetURL)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		logger.Error("failed to create request", "error", err)
		os.Exit(1)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("request failed", "error", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("← %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))

	// Check for payment response header
	paymentResponse := resp.Header.Get("PAYMENT-RESPONSE")
	if paymentResponse == "" {
		paymentResponse = resp.Header.Get("X-PAYMENT-RESPONSE")
	}

	if paymentResponse != "" {
		settleResp, err := httpx402Client.GetPaymentSettleResponse(headerMap(resp.Header))
		if err == nil && settleResp != nil {
			fmt.Printf("\n💰 Payment Settlement:\n")
			fmt.Printf("   Success:     %v\n", settleResp.Success)
			fmt.Printf("   Transaction: %s\n", settleResp.Transaction)
			fmt.Printf("   Network:     %s\n", settleResp.Network)
			fmt.Printf("   Payer:       %s\n", settleResp.Payer)
		}
	}

	// Pretty-print response body
	if len(body) > 0 {
		var prettyJSON json.RawMessage
		if json.Unmarshal(body, &prettyJSON) == nil {
			formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
			fmt.Printf("\nResponse:\n%s\n", formatted)
		} else {
			fmt.Printf("\nResponse:\n%s\n", body)
		}
	}
}

func headerMap(h http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}
