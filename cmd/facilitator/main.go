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
	evmfacilitator "github.com/coinbase/x402/go/mechanisms/evm/exact/facilitator"
	"github.com/gin-gonic/gin"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/facilserver"
	"github.com/GrapeInTheTree/x402-playground/internal/signer"
	"github.com/GrapeInTheTree/x402-playground/pkg/health"
)

func main() {
	cfg, err := config.LoadFacilitator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	// Create facilitator EVM signer
	evmSigner, err := signer.NewFacilitatorSigner(cfg.PrivateKey, cfg.RPCURL, logger)
	if err != nil {
		logger.Error("failed to create signer", "error", err)
		os.Exit(1)
	}

	defer evmSigner.Close()

	logger.Info("facilitator initialized",
		"address", evmSigner.Address(),
		"network", cfg.Network,
		"rpcURL", cfg.RPCURL,
	)

	// Create x402 facilitator and register EVM exact scheme
	facilitator := x402.Newx402Facilitator()

	evmScheme := evmfacilitator.NewExactEvmScheme(evmSigner, &evmfacilitator.ExactEvmSchemeConfig{
		DeployERC4337WithEIP6492: false,
	})

	networks := []x402.Network{x402.Network(cfg.Network)}
	facilitator.Register(networks, evmScheme)

	// Set up HTTP server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	srv := facilserver.New(facilitator, logger)
	r.POST("/verify", srv.HandleVerify)
	r.POST("/settle", srv.HandleSettle)
	r.GET("/supported", srv.HandleSupported)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, health.Response{
			Status:  "ok",
			Service: "facilitator",
			Network: cfg.Network,
			Address: evmSigner.Address(),
		})
	})

	// Graceful shutdown
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("facilitator server starting", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
