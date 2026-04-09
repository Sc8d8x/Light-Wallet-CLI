package src

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckBalance(addstr string, rpcUrl string) error {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}

	defer client.Close()

	adress := common.HexToAddress(addstr)

	balance, err := client.BalanceAt(context.Background(), adress, nil) // nil означает последний блок

	if err != nil {
		return fmt.Errorf("failed to fetch balance: %w", err)
	}
	// conversion from  Wei в Ether (1 ETH = 10^18 Wei)
	ethbalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1.e18))

	fmt.Printf("Баланс адреса %s: %s ETH\n", adress, ethbalance.Text('f', 18))
	return nil

}
