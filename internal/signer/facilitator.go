package signer

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	evm "github.com/coinbase/x402/go/mechanisms/evm"
)

// FacilitatorSigner implements evm.FacilitatorEvmSigner using a local ECDSA private key.
type FacilitatorSigner struct {
	client     *ethclient.Client
	privateKey *[32]byte
	address    common.Address
	logger     *slog.Logger
}

// NewFacilitatorSigner creates a new FacilitatorEvmSigner from a hex-encoded private key and RPC URL.
func NewFacilitatorSigner(privateKeyHex, rpcURL string, logger *slog.Logger) (*FacilitatorSigner, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	pk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	address := crypto.PubkeyToAddress(pk.PublicKey)

	// Store key as fixed-size array and zero the intermediate slice
	rawKey := crypto.FromECDSA(pk)
	var keyBytes [32]byte
	copy(keyBytes[:], rawKey)
	for i := range rawKey {
		rawKey[i] = 0
	}

	return &FacilitatorSigner{
		client:     client,
		privateKey: &keyBytes,
		address:    address,
		logger:     logger,
	}, nil
}

func (s *FacilitatorSigner) ecdsaKey() *[32]byte {
	return s.privateKey
}

// GetAddresses returns the list of addresses controlled by this signer.
func (s *FacilitatorSigner) GetAddresses() []string {
	return []string{s.address.Hex()}
}

// ReadContract calls a read-only function on a smart contract and returns the result.
func (s *FacilitatorSigner) ReadContract(ctx context.Context, contractAddress string, abiBytes []byte, functionName string, args ...interface{}) (interface{}, error) {
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := contractABI.Pack(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call: %w", err)
	}

	addr := common.HexToAddress(contractAddress)
	result, err := s.client.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	outputs, err := contractABI.Unpack(functionName, result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	if len(outputs) == 0 {
		return nil, nil
	}
	if len(outputs) == 1 {
		return outputs[0], nil
	}
	return outputs, nil
}

// VerifyTypedData recovers the signer from an EIP-712 typed data signature and checks it matches the expected address.
func (s *FacilitatorSigner) VerifyTypedData(
	ctx context.Context,
	address string,
	domain evm.TypedDataDomain,
	typesMap map[string][]evm.TypedDataField,
	primaryType string,
	message map[string]interface{},
	signature []byte,
) (bool, error) {
	typedData := apitypes.TypedData{
		Types:       make(apitypes.Types),
		PrimaryType: primaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              domain.Name,
			Version:           domain.Version,
			ChainId:           (*math.HexOrDecimal256)(domain.ChainID),
			VerifyingContract: domain.VerifyingContract,
		},
		Message: message,
	}

	for typeName, fields := range typesMap {
		typedFields := make([]apitypes.Type, len(fields))
		for i, f := range fields {
			typedFields[i] = apitypes.Type{Name: f.Name, Type: f.Type}
		}
		typedData.Types[typeName] = typedFields
	}

	if _, exists := typedData.Types["EIP712Domain"]; !exists {
		typedData.Types["EIP712Domain"] = []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		}
	}

	dataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return false, fmt.Errorf("failed to hash struct: %w", err)
	}

	domainSep, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return false, fmt.Errorf("failed to hash domain: %w", err)
	}

	rawData := append([]byte{0x19, 0x01}, domainSep...)
	rawData = append(rawData, dataHash...)
	digest := crypto.Keccak256(rawData)

	// Normalize v value (27/28 → 0/1 for recovery)
	sig := make([]byte, len(signature))
	copy(sig, signature)
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	pubKey, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	expectedAddr := common.HexToAddress(address)

	return recoveredAddr == expectedAddr, nil
}

// WriteContract calls a state-changing function on a smart contract and returns the transaction hash.
func (s *FacilitatorSigner) WriteContract(ctx context.Context, contractAddress string, abiBytes []byte, functionName string, args ...interface{}) (string, error) {
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := contractABI.Pack(functionName, args...)
	if err != nil {
		return "", fmt.Errorf("failed to pack tx data: %w", err)
	}

	return s.sendRawTx(ctx, contractAddress, data)
}

// SendTransaction sends a raw transaction to the given address and returns the transaction hash.
func (s *FacilitatorSigner) SendTransaction(ctx context.Context, to string, data []byte) (string, error) {
	return s.sendRawTx(ctx, to, data)
}

func (s *FacilitatorSigner) sendRawTx(ctx context.Context, to string, data []byte) (string, error) {
	chainID, err := s.client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	nonce, err := s.client.PendingNonceAt(ctx, s.address)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasTipCap, err := s.client.SuggestGasTipCap(ctx)
	if err != nil {
		gasTipCap = big.NewInt(1_000_000_000) // 1 gwei fallback
	}

	header, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get latest block: %w", err)
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		baseFee = big.NewInt(1_000_000_000)
	}
	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

	toAddr := common.HexToAddress(to)
	gasLimit, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
		From:      s.address,
		To:        &toAddr,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Data:      data,
	})
	if err != nil {
		gasLimit = 300_000 // fallback
		s.logger.Warn("gas estimation failed, using fallback", "error", err, "gasLimit", gasLimit)
	}

	// Add 20% buffer
	gasLimit = gasLimit * 120 / 100

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     big.NewInt(0),
		Data:      data,
	})

	pk, err := crypto.ToECDSA(s.ecdsaKey()[:])
	if err != nil {
		return "", fmt.Errorf("failed to reconstruct key: %w", err)
	}

	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, pk)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	s.logger.Info("transaction sent", "txHash", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

// WaitForTransactionReceipt polls for a transaction receipt until it is mined or a timeout is reached.
func (s *FacilitatorSigner) WaitForTransactionReceipt(ctx context.Context, txHash string) (*evm.TransactionReceipt, error) {
	hash := common.HexToHash(txHash)

	timeout := 60 * time.Second
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		receipt, err := s.client.TransactionReceipt(ctx, hash)
		if err == nil {
			return &evm.TransactionReceipt{
				Status:      receipt.Status,
				BlockNumber: receipt.BlockNumber.Uint64(),
				TxHash:      receipt.TxHash.Hex(),
			}, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return nil, fmt.Errorf("timeout waiting for transaction receipt: %s", txHash)
}

// GetBalance returns the native or ERC-20 token balance for the given address.
func (s *FacilitatorSigner) GetBalance(ctx context.Context, address string, tokenAddress string) (*big.Int, error) {
	if tokenAddress == "" || tokenAddress == "0x0000000000000000000000000000000000000000" {
		balance, err := s.client.BalanceAt(ctx, common.HexToAddress(address), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get native balance: %w", err)
		}
		return balance, nil
	}

	erc20ABI := `[{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"}]`
	result, err := s.ReadContract(ctx, tokenAddress, []byte(erc20ABI), "balanceOf", common.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	if balance, ok := result.(*big.Int); ok {
		return balance, nil
	}
	return nil, fmt.Errorf("unexpected balance type: %T", result)
}

// GetChainID returns the chain ID from the connected RPC node.
func (s *FacilitatorSigner) GetChainID(ctx context.Context) (*big.Int, error) {
	return s.client.ChainID(ctx)
}

// GetCode returns the contract bytecode deployed at the given address.
func (s *FacilitatorSigner) GetCode(ctx context.Context, address string) ([]byte, error) {
	return s.client.CodeAt(ctx, common.HexToAddress(address), nil)
}

// Address returns the facilitator's Ethereum address (convenience method).
func (s *FacilitatorSigner) Address() string {
	return s.address.Hex()
}

// Close zeroes the private key from memory and closes the RPC client.
// Call this when the signer is no longer needed.
func (s *FacilitatorSigner) Close() {
	if s.privateKey != nil {
		for i := range s.privateKey {
			s.privateKey[i] = 0
		}
	}
	if s.client != nil {
		s.client.Close()
	}
}
