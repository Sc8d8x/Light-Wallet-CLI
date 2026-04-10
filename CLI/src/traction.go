package src

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Checktraction(rpcURL string, privateKey *ecdsa.PrivateKey, toAddress string, amountInEther string) (string, error) {
	client, err := ethclient.Dial(rpcURL)

	if err != nil {
		return "", fmt.Errorf("Error connect user: %v", err)
	}
	defer client.Close()
	// Get address sender
	// получение адреса отправителя

	fromAdress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Getting the current transaction number (Nonce) for this address
	// получение текущего номера транзакции (Nonce) для этого адреса
	// Nonce must be unique and sequential for each address
	// Nonce должен быть уникальным и последовательным для каждого адреса

	nonce, err := client.PendingNonceAt(context.Background(), fromAdress)

	if err != nil {
		return "", fmt.Errorf("Ошибка получения nonce: %v", err)
	}

	// конвертация без потери точности при отправке данных

	value := new(big.Float)
	value.SetString(amountInEther)

	weiMultiplier := new(big.Float).SetInt(big.NewInt(1e18))
	value.Mul(value, weiMultiplier)

	valueInWei := new(big.Int)
	value.Int(valueInWei) // конвертация в целое число

	// Setting the gas limit and gas price (Gas Tip Cap / Gas Fee Cap for EIP-1559)
	// установка лимита газа и цены газа (Gas Tip Cap / Gas Fee Cap для EIP-1559)
	gasLimit := uint64(21000)

	// получаем текущие цены газа
	// Getting current gas prices

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("Ошибка получения цены газа: %v", err)
	}

	// Creating a transaction structure
	// создание структуры транзакции

	toAddr := common.HexToAddress(toAddress)
	tx := types.NewTransaction(nonce, toAddr, valueInWei, gasLimit, gasPrice, nil)

	// You need to find out the network ChainID (1 for Mainnet, 11155111 for Sepolia, etc.)
	// нужно узнать ChainID сети (1 для Mainnet, 11155111 для Sepolia и т.д.)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("Ошибка получения ChainID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("Ошибка подписания транзакции: %v", err)
	}
	// sending a raw transaction to the network
	//  отправка сырой транзакции в сеть
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("Ошибка отправки транзакции: %v", err)
	}

	fmt.Printf("Транзакция отправлена! Хеш: %s\n", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}
