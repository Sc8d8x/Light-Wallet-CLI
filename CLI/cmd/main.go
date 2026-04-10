package main

import (
	"Light-Wallet-CLI/CLI/src"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	amount := "0.01"
	mnemonic, err := src.GenerateEntropy()
	if err != nil {
		panic(err)
	}

	seedHex, err := src.MnemonicToSeed(mnemonic, "")
	if err != nil {
		panic(err)
	}

	privateKey, err := src.DerivePrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	address, err := src.GenerateKeyStore(mnemonic, "test-password")
	if err != nil {
		panic(err)
	}

	fmt.Println("Mnemonic:", mnemonic)
	fmt.Println("Seed (hex):", seedHex)
	fmt.Println("Private key (hex):", hex.EncodeToString(crypto.FromECDSA(privateKey)))
	fmt.Println("Address:", address)

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Println("RPC_URL is empty, skip balance check")
		return
	}

	txHash, err := src.Checktraction(rpcURL, privateKey, address, amount)
	if err != nil {
		log.Fatalf("SendTransaction failed: %v", err)
	}

	if err := src.CheckBalance(address, rpcURL); err != nil {
		log.Printf("balance check failed: %v\n", err)
	}

	fmt.Println(txHash)
}
