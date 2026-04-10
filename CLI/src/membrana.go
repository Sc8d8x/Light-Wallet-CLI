package src

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tyler-smith/go-bip39"
)

// генерируем мнемонику
func GenerateEntropy() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

//	ручная деривацию BIP-44 для Ethereum
//
// Путь: m/44'/60'/0'/0/0
func DerivePrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	// Мнемоника (BIP-39)
	seed := bip39.NewSeed(mnemonic, "") // passphrase пустой

	// Создаем мастер-ключ из Seed (BIP-32).
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	// деривация по пути m/44'/60'/0'/0/0
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}

	//  coin_type = 60' (Ethereum)
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 60)
	if err != nil {
		return nil, err
	}

	// account = 0'
	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	// change = 0
	change, err := account.Derive(0)
	if err != nil {
		return nil, err
	}

	//  address_index = 0
	addressIndex, err := change.Derive(0)
	if err != nil {
		return nil, err
	}

	// получаем приватный ключ в формате ECDSA
	privateKeyECDSA, err := addressIndex.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get EC private key: %w", err)
	}

	// конвертируем в тип *ecdsa.PrivateKey, который понимает go-ethereum
	return privateKeyECDSA.ToECDSA(), nil
}

// создает хранилище, импортируя ключ из мнемоники
func GenerateKeyStore(mnemonic string, password string) (string, error) {
	//деривируем ключ нашей надежной функцией
	privateKey, err := DerivePrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return "", fmt.Errorf("derivation failed: %w", err)
	}

	//настраиваем путь
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	keyDir := filepath.Join(homeDir, ".crypto_wallet", "keys")

	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return "", err
	}

	// создаем хранилище
	ks := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// импортируем ключ
	account, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		return "", err
	}

	return account.Address.Hex(), nil
}

// утилита
func MnemonicToSeed(mnemonic string, passphrase string) (string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", errors.New("invalid BIP-39 mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, passphrase)
	return hex.EncodeToString(seed), nil
}

// расшифровывка приватного ключа из хранилища
func LoadPrivateKeyFromKeystore(keyDir, addressHex, password string) (*ecdsa.PrivateKey, error) {
	// инициализируем хранилище
	ks := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// создаем объект аккаунта для поиска
	accountToFind := accounts.Account{Address: common.HexToAddress(addressHex)}

	// ищем аккаунт в хранилище (получаем полный путь к файлу UTC)
	account, err := ks.Find(accountToFind)
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	return getPrivateKeyFromJSON(account, password)
}

// вспомогательная функция для ручного декодирования JSON файла ключа
func getPrivateKeyFromJSON(account accounts.Account, password string) (*ecdsa.PrivateKey, error) {
	// находим файл ключа вручную (путь есть в account.URL.Path)
	keyJSON, err := os.ReadFile(account.URL.Path)
	if err != nil {
		return nil, err
	}

	// используем внутренний метод DecryptKey из пакета keystore
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, err
	}

	// key.PrivateKey - это и есть *ecdsa.PrivateKey
	return key.PrivateKey, nil
}
