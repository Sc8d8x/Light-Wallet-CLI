package main

import (
	"fmt"

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

	privateKeyHex, err := src.DerivePrivateKeyFromMnemonic(seedHex)
	if err != nil {
		panic(err)
	}

	fmt.Println("Mnemonic:", mnemonic)
	fmt.Println("Seed (hex):", seedHex)
	fmt.Println("Private key (hex):", privateKeyHex)
}
