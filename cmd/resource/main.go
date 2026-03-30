package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	ginmw "github.com/coinbase/x402/go/http/gin"
	evmserver "github.com/coinbase/x402/go/mechanisms/evm/exact/server"
	"github.com/gin-gonic/gin"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/server"
)

func main() {
	cfg, err := config.LoadResource()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	network := x402.Network(cfg.Network)

	// Create facilitator client
	facilitatorClient := x402http.NewHTTPFacilitatorClient(&x402http.FacilitatorConfig{
		URL:     cfg.FacilitatorURL,
		Timeout: 60 * time.Second,
	})

	// Build payment-protected routes
	routes := server.BuildRoutes(cfg.PayToAddress, network)

	logger.Info("resource server initialized",
		"network", cfg.Network,
		"payTo", cfg.PayToAddress,
		"facilitatorURL", cfg.FacilitatorURL,
		"assetTransferMethod", cfg.AssetTransferMethod,
	)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Health endpoint (no payment required)
	r.GET("/health", server.HealthHandler("resource", cfg.Network))

	// Register EVM exact scheme with money parser.
	// In eip3009 mode: SDK-supported networks use SDK defaults; unknown networks use custom parser.
	// In permit2 mode: ALL networks use custom parser to inject assetTransferMethod=permit2.
	evmScheme := evmserver.NewExactEvmScheme()
	evmScheme.RegisterMoneyParser(func(amount float64, net x402.Network) (*x402.AssetAmount, error) {
		knownNetworks := map[string]struct {
			address string
			name    string
			version string
		}{
			"eip155:8453":  {address: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", name: "USD Coin", version: "2"},
			"eip155:84532": {address: "0x036CbD53842c5426634e7929541eC2318f3dCF7e", name: "USDC", version: "2"},
			"eip155:137":   {address: "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359", name: "USD Coin", version: "2"},
			"eip155:42161": {address: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831", name: "USD Coin", version: "2"},
		}

		if cfg.AssetTransferMethod == "permit2" {
			// Permit2 mode: override all networks to inject assetTransferMethod
			asset := cfg.USDCAddress
			name := "USDC"
			version := "2"
			if info, ok := knownNetworks[string(net)]; ok {
				asset = info.address
				name = info.name
				version = info.version
			}
			atomicAmount := int64(amount * 1_000_000)
			return &x402.AssetAmount{
				Asset:  asset,
				Amount: fmt.Sprintf("%d", atomicAmount),
				Extra: map[string]interface{}{
					"name":                name,
					"version":             version,
					"assetTransferMethod": "permit2",
				},
			}, nil
		}

		// EIP-3009 mode: delegate SDK-supported networks to SDK default
		if _, ok := knownNetworks[string(net)]; ok {
			return nil, nil
		}

		// Custom parser for unknown networks (Chiliz, etc.)
		if cfg.USDCAddress != "" {
			atomicAmount := int64(amount * 1_000_000)
			return &x402.AssetAmount{
				Asset:  cfg.USDCAddress,
				Amount: fmt.Sprintf("%d", atomicAmount),
				Extra: map[string]interface{}{
					"name":    "USDC",
					"version": "2",
				},
			}, nil
		}
		return nil, nil
	})

	// Apply x402 payment middleware
	r.Use(ginmw.X402Payment(ginmw.Config{
		Routes:      routes,
		Facilitator: facilitatorClient,
		Schemes: []ginmw.SchemeConfig{
			{Network: network, Server: evmScheme},
		},
		SyncFacilitatorOnStart: true,
		Timeout:                60 * time.Second,
		SettlementHandler: func(c *gin.Context, resp *x402.SettleResponse) {
			logger.Info("payment settled",
				"txHash", resp.Transaction,
				"network", resp.Network,
				"payer", resp.Payer,
			)
		},
	}))

	// Protected endpoints
	r.GET("/weather", server.WeatherHandler)
	r.GET("/joke", server.JokeHandler)
	r.GET("/premium-data", server.PremiumDataHandler)

	// Graceful shutdown
	httpServer := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("resource server starting", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case err := <-errCh:
		logger.Error("server failed to start", "error", err)
	}

	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
