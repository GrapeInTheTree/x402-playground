package demo

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20BalanceOfABI is the minimal ABI for the ERC-20 balanceOf function.
var ERC20BalanceOfABI = `[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// ERC20AllowanceABI is the minimal ABI for the ERC-20 allowance function.
var ERC20AllowanceABI = `[{"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// QueryBalances fetches ETH and USDC balances for the given wallets.
func QueryBalances(ctx context.Context, client *ethclient.Client, usdcAddr string, wallets []WalletInfo) ([]WalletBalance, error) {
	usdc := common.HexToAddress(usdcAddr)
	erc20ABI, err := abi.JSON(strings.NewReader(ERC20BalanceOfABI))
	if err != nil {
		return nil, err
	}

	ethDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	usdcDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))

	balances := make([]WalletBalance, 0, len(wallets))

	for _, w := range wallets {
		addr := common.HexToAddress(w.Address)

		ethBal, err := client.BalanceAt(ctx, addr, nil)
		if err != nil {
			ethBal = big.NewInt(0)
		}

		data, _ := erc20ABI.Pack("balanceOf", addr)
		result, err := client.CallContract(ctx, ethereum.CallMsg{To: &usdc, Data: data}, nil)

		usdcBal := big.NewInt(0)
		if err == nil && len(result) > 0 {
			out, err := erc20ABI.Unpack("balanceOf", result)
			if err == nil && len(out) > 0 {
				if val, ok := out[0].(*big.Int); ok {
					usdcBal = val
				}
			}
		}

		ethF := new(big.Float).Quo(new(big.Float).SetInt(ethBal), ethDiv)
		usdcF := new(big.Float).Quo(new(big.Float).SetInt(usdcBal), usdcDiv)

		balances = append(balances, WalletBalance{
			Wallet:  w,
			ETH:     ethF.Text('f', 6),
			USDC:    usdcF.Text('f', 6),
			ETHRaw:  ethBal.String(),
			USDCRaw: usdcBal.String(),
		})
	}

	return balances, nil
}

// QueryAllowance fetches the ERC-20 allowance from owner to spender.
func QueryAllowance(ctx context.Context, client *ethclient.Client, tokenAddr, ownerAddr, spenderAddr string) (string, error) {
	token := common.HexToAddress(tokenAddr)
	owner := common.HexToAddress(ownerAddr)
	spender := common.HexToAddress(spenderAddr)

	allowanceABI, err := abi.JSON(strings.NewReader(ERC20AllowanceABI))
	if err != nil {
		return "0", err
	}

	data, err := allowanceABI.Pack("allowance", owner, spender)
	if err != nil {
		return "0", err
	}

	result, err := client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return "0", err
	}

	out, err := allowanceABI.Unpack("allowance", result)
	if err != nil || len(out) == 0 {
		return "0", err
	}

	val, ok := out[0].(*big.Int)
	if !ok {
		return "0", fmt.Errorf("unexpected type from allowance call: %T", out[0])
	}
	usdcDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))
	formatted := new(big.Float).Quo(new(big.Float).SetInt(val), usdcDiv)

	return formatted.Text('f', 6), nil
}
