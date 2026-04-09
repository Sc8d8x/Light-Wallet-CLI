package main

import (
	"fmt"
	"log"
	"os"

	"Light-Wallet-CLI/CLI/src"
)

func main() {
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
	fmt.Println("Private key:", privateKey)
	fmt.Println("Address:", address)

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Println("RPC_URL is empty, skip balance check")
		return
	}

	if err := src.CheckBalance(address, rpcURL); err != nil {
		log.Printf("balance check failed: %v\n", err)
	}
}
