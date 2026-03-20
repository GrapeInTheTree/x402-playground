package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

// FacilitatorConfig holds configuration for the facilitator server.
type FacilitatorConfig struct {
	PrivateKey          string
	RPCURL              string
	Network             string
	Port                string
	AssetTransferMethod string // "eip3009" (default) or "permit2"
	LogLevel            slog.Level
}

// ResourceConfig holds configuration for the resource server.
type ResourceConfig struct {
	FacilitatorURL      string
	PayToAddress        string
	Network             string
	USDCAddress         string
	Port                string
	AssetTransferMethod string // "eip3009" (default) or "permit2"
	LogLevel            slog.Level
}

// ClientConfig holds configuration for the client CLI.
type ClientConfig struct {
	PrivateKey          string
	ResourceURL         string
	EndpointPath        string
	RPCURL              string
	Network             string
	USDCAddress         string
	AssetTransferMethod string // "eip3009" (default) or "permit2"
	LogLevel            slog.Level
}

func LoadFacilitator() (*FacilitatorConfig, error) {
	cfg := &FacilitatorConfig{
		PrivateKey:          os.Getenv("FACILITATOR_PRIVATE_KEY"),
		RPCURL:              envOr("RPC_URL", "https://sepolia.base.org"),
		Network:             envOr("NETWORK", "eip155:84532"),
		Port:                envOr("FACILITATOR_PORT", "4022"),
		AssetTransferMethod: envOr("ASSET_TRANSFER_METHOD", "eip3009"),
		LogLevel:            parseLogLevel(os.Getenv("LOG_LEVEL")),
	}

	var errs []error
	if cfg.PrivateKey == "" {
		errs = append(errs, fmt.Errorf("FACILITATOR_PRIVATE_KEY is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

func LoadResource() (*ResourceConfig, error) {
	cfg := &ResourceConfig{
		FacilitatorURL:      os.Getenv("FACILITATOR_URL"),
		PayToAddress:        os.Getenv("PAY_TO_ADDRESS"),
		Network:             envOr("NETWORK", "eip155:84532"),
		USDCAddress:         envOr("USDC_ADDRESS", "0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
		Port:                envOr("RESOURCE_PORT", "4021"),
		AssetTransferMethod: envOr("ASSET_TRANSFER_METHOD", "eip3009"),
		LogLevel:            parseLogLevel(os.Getenv("LOG_LEVEL")),
	}

	var errs []error
	if cfg.FacilitatorURL == "" {
		errs = append(errs, fmt.Errorf("FACILITATOR_URL is required"))
	}
	if cfg.PayToAddress == "" {
		errs = append(errs, fmt.Errorf("PAY_TO_ADDRESS is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

func LoadClient() (*ClientConfig, error) {
	cfg := &ClientConfig{
		PrivateKey:          os.Getenv("CLIENT_PRIVATE_KEY"),
		ResourceURL:         os.Getenv("RESOURCE_URL"),
		EndpointPath:        envOr("ENDPOINT_PATH", "/weather"),
		RPCURL:              envOr("RPC_URL", "https://sepolia.base.org"),
		Network:             envOr("NETWORK", "eip155:84532"),
		USDCAddress:         envOr("USDC_ADDRESS", "0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
		AssetTransferMethod: envOr("ASSET_TRANSFER_METHOD", "eip3009"),
		LogLevel:            parseLogLevel(os.Getenv("LOG_LEVEL")),
	}

	var errs []error
	if cfg.PrivateKey == "" {
		errs = append(errs, fmt.Errorf("CLIENT_PRIVATE_KEY is required"))
	}
	if cfg.ResourceURL == "" {
		errs = append(errs, fmt.Errorf("RESOURCE_URL is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

// ExplorerConfig holds configuration for the TUI explorer.
type ExplorerConfig struct {
	ClientPrivateKey    string
	FacilitatorURL      string
	ResourceURL         string
	PayToAddress        string
	RPCURL              string
	Network             string
	USDCAddress         string
	AssetTransferMethod string // "eip3009" (default) or "permit2"
	LogLevel            slog.Level
}

func LoadExplorer() (*ExplorerConfig, error) {
	cfg := &ExplorerConfig{
		ClientPrivateKey:    os.Getenv("CLIENT_PRIVATE_KEY"),
		FacilitatorURL:      os.Getenv("FACILITATOR_URL"),
		ResourceURL:         os.Getenv("RESOURCE_URL"),
		PayToAddress:        os.Getenv("PAY_TO_ADDRESS"),
		RPCURL:              envOr("RPC_URL", "https://sepolia.base.org"),
		Network:             envOr("NETWORK", "eip155:84532"),
		USDCAddress:         envOr("USDC_ADDRESS", "0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
		AssetTransferMethod: envOr("ASSET_TRANSFER_METHOD", "eip3009"),
		LogLevel:            parseLogLevel(os.Getenv("LOG_LEVEL")),
	}

	var errs []error
	if cfg.ClientPrivateKey == "" {
		errs = append(errs, fmt.Errorf("CLIENT_PRIVATE_KEY is required"))
	}
	if cfg.ResourceURL == "" {
		errs = append(errs, fmt.Errorf("RESOURCE_URL is required"))
	}
	if cfg.FacilitatorURL == "" {
		errs = append(errs, fmt.Errorf("FACILITATOR_URL is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
