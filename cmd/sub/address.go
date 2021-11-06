package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shengdoushi/base58"
	"golang.org/x/crypto/sha3"
	"os"
)

func DecodeAddress (privateKey string) string {
	// Instructions from https://secretscan.org/PrivateKeyTron

	publicKey := privateKeyToPublicKey(privateKey)

	hexAddress, err := publicKeyToHexAddress(publicKey)
	if err != nil {
		handlerError(err)
	}

	address, err := hexToAddress(hexAddress)
	if err != nil {
		handlerError(err)
	}

	fmt.Println(address)

	return address
}

func privateKeyToPublicKey(privateKey string) []byte {
	var pk *ecdsa.PrivateKey

	pk, _ = crypto.HexToECDSA(privateKey)
	publicKey := elliptic.Marshal(crypto.S256(), pk.X, pk.Y)

	return publicKey
}

func publicKeyToHexAddress(publicKey []byte) (string, error) {
	var dataKeccak []byte

	keccak := sha3.NewLegacyKeccak256()
	_, err := keccak.Write([]byte(publicKey[1:]))
	if err != nil {
		return "", fmt.Errorf("encountered an error while encoding to keccak: %v", err)
	}

	sumKeccak := keccak.Sum(dataKeccak)
	formattedKeccak := "41" + hex.EncodeToString(sumKeccak[len(sumKeccak)-20:])

	hash, err := hashSHA256(formattedKeccak)
	if err != nil {
		return "", err
	}

	hash, err = hashSHA256(hash)
	if err != nil {
		return "", err
	}

	hexAddress := formattedKeccak + hash[:8]

	return hexAddress, nil
}

func hashSHA256(data string) (string, error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("encountered an error while decoding hex: %v", err)
	}

	hasher := sha256.New()
	hasher.Write(decoded)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func hexToAddress(hexAddress string) (string, error) {
	fmt.Println(hexAddress)
	addressDecoded, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", fmt.Errorf("encountered an error while decoding from hex: %v", err)
	}

	address := string(base58.Encode(addressDecoded, base58.BitcoinAlphabet))

	asd, _ := DecodeCheck(address)
	CheckBalance(string(asd))
	return address, nil
}

func handlerError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func Decode(input string) ([]byte, error) {
	fmt.Println("DECODE")
	fmt.Println(base58.Decode(input, base58.BitcoinAlphabet))
	return base58.Decode(input, base58.BitcoinAlphabet)
}

func DecodeCheck(input string) ([]byte, error) {
	decodeCheck, err := Decode(input)

	if err != nil {
		return nil, err
	}

	if len(decodeCheck) < 4 {
		return nil, fmt.Errorf("b58 check error")
	}

	decodeData := decodeCheck[:len(decodeCheck)-4]

	h256h0 := sha256.New()
	h256h0.Write(decodeData)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	if h1[0] == decodeCheck[len(decodeData)] &&
		h1[1] == decodeCheck[len(decodeData)+1] &&
		h1[2] == decodeCheck[len(decodeData)+2] &&
		h1[3] == decodeCheck[len(decodeData)+3] {
		fmt.Println("DecodeData")
		fmt.Println(decodeData)
		fmt.Println("DecodeData")
		return decodeData, nil
	}
	return nil, fmt.Errorf("b58 check error")
}
