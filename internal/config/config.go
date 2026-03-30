package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		// Only warn on actual parse errors; missing .env file is fine
		if !os.IsNotExist(err) {
			slog.Warn("failed to parse .env file", "error", err)
		}
	}
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

// LoadFacilitator loads facilitator server configuration from environment variables.
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
	if err := validateNetwork(cfg.Network); err != nil {
		errs = append(errs, err)
	}
	if err := validatePort(cfg.Port); err != nil {
		errs = append(errs, err)
	}
	if err := validateTransferMethod(cfg.AssetTransferMethod); err != nil {
		errs = append(errs, err)
	}
	if err := validateURL(cfg.RPCURL, "RPC_URL"); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

// LoadResource loads resource server configuration from environment variables.
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
	} else if err := validateURL(cfg.FacilitatorURL, "FACILITATOR_URL"); err != nil {
		errs = append(errs, err)
	}
	if cfg.PayToAddress == "" {
		errs = append(errs, fmt.Errorf("PAY_TO_ADDRESS is required"))
	} else if err := validateEthAddress(cfg.PayToAddress, "PAY_TO_ADDRESS"); err != nil {
		errs = append(errs, err)
	}
	if err := validateNetwork(cfg.Network); err != nil {
		errs = append(errs, err)
	}
	if err := validatePort(cfg.Port); err != nil {
		errs = append(errs, err)
	}
	if err := validateTransferMethod(cfg.AssetTransferMethod); err != nil {
		errs = append(errs, err)
	}
	if err := validateEthAddress(cfg.USDCAddress, "USDC_ADDRESS"); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

// LoadClient loads client CLI configuration from environment variables.
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
	} else if err := validateURL(cfg.ResourceURL, "RESOURCE_URL"); err != nil {
		errs = append(errs, err)
	}
	if err := validateNetwork(cfg.Network); err != nil {
		errs = append(errs, err)
	}
	if err := validateTransferMethod(cfg.AssetTransferMethod); err != nil {
		errs = append(errs, err)
	}
	if err := validateEthAddress(cfg.USDCAddress, "USDC_ADDRESS"); err != nil {
		errs = append(errs, err)
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

// LoadExplorer loads TUI explorer configuration from environment variables.
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
	} else if err := validateURL(cfg.ResourceURL, "RESOURCE_URL"); err != nil {
		errs = append(errs, err)
	}
	if cfg.FacilitatorURL == "" {
		errs = append(errs, fmt.Errorf("FACILITATOR_URL is required"))
	} else if err := validateURL(cfg.FacilitatorURL, "FACILITATOR_URL"); err != nil {
		errs = append(errs, err)
	}
	if err := validateNetwork(cfg.Network); err != nil {
		errs = append(errs, err)
	}
	if err := validateTransferMethod(cfg.AssetTransferMethod); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return cfg, nil
}

func validateEthAddress(addr, name string) error {
	if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
		return fmt.Errorf("%s must be 0x-prefixed 40-char hex address, got %q", name, addr)
	}
	if _, err := hex.DecodeString(addr[2:]); err != nil {
		return fmt.Errorf("%s contains invalid hex: %w", name, err)
	}
	return nil
}

func validateNetwork(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] != "eip155" {
		return fmt.Errorf("NETWORK must be in format eip155:<chainId>, got %q", s)
	}
	if _, err := strconv.ParseUint(parts[1], 10, 64); err != nil {
		return fmt.Errorf("NETWORK chain ID must be numeric, got %q", parts[1])
	}
	return nil
}

func validateURL(s, name string) error {
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("%s is not a valid URL: %w", name, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%s must use http or https scheme, got %q", name, s)
	}
	return nil
}

func validatePort(s string) error {
	p, err := strconv.Atoi(s)
	if err != nil || p < 1 || p > 65535 {
		return fmt.Errorf("port must be 1-65535, got %q", s)
	}
	return nil
}

func validateTransferMethod(s string) error {
	if s != "eip3009" && s != "permit2" {
		return fmt.Errorf("ASSET_TRANSFER_METHOD must be eip3009 or permit2, got %q", s)
	}
	return nil
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
