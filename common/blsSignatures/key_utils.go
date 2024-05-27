package blsSignatures

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

// Note for Decode functions
// Ethereum's BLS library doesn't like the byte slice containing the BLS keys to be
// any larger than necessary, so we need to create a Decoder to avoid returning any padding.

func DecodeBase64BLSPublicKey(pubKeyEncodedBytes []byte) (*PublicKey, error) {
	pubKeyDecoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(pubKeyEncodedBytes))
	pubKeyBytes, err := io.ReadAll(pubKeyDecoder)
	if err != nil {
		return nil, err
	}
	pubKey, err := PublicKeyFromBytes(pubKeyBytes, false)
	if err != nil {
		return nil, err
	}
	return &pubKey, nil
}

func DecodeBase64BLSPrivateKey(privKeyEncodedBytes []byte) (PrivateKey, error) {
	privKeyDecoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(privKeyEncodedBytes))
	privKeyBytes, err := io.ReadAll(privKeyDecoder)
	if err != nil {
		return nil, err
	}
	privKey, err := PrivateKeyFromBytes(privKeyBytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

const DefaultPubKeyFilename = "das_bls.pub"
const DefaultPrivKeyFilename = "das_bls"

func GenerateAndStoreKeys(keyDir string) (*PublicKey, *PrivateKey, error) {
	pubKey, privKey, err := GenerateKeys()
	if err != nil {
		return nil, nil, err
	}
	pubKeyPath := keyDir + "/" + DefaultPubKeyFilename
	pubKeyBytes := PublicKeyToBytes(pubKey)
	encodedPubKey := make([]byte, base64.StdEncoding.EncodedLen(len(pubKeyBytes)))
	base64.StdEncoding.Encode(encodedPubKey, pubKeyBytes)
	err = os.WriteFile(pubKeyPath, encodedPubKey, 0o600)
	if err != nil {
		return nil, nil, err
	}

	privKeyPath := keyDir + "/" + DefaultPrivKeyFilename
	privKeyBytes := PrivateKeyToBytes(privKey)
	encodedPrivKey := make([]byte, base64.StdEncoding.EncodedLen(len(privKeyBytes)))
	base64.StdEncoding.Encode(encodedPrivKey, privKeyBytes)
	err = os.WriteFile(privKeyPath, encodedPrivKey, 0o600)
	if err != nil {
		return nil, nil, err
	}
	return &pubKey, &privKey, nil
}

func ReadKeysFromFile(keyDir string) (*PublicKey, PrivateKey, error) {
	pubKey, err := ReadPubKeyFromFile(keyDir + "/" + DefaultPubKeyFilename)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := ReadPrivKeyFromFile(keyDir + "/" + DefaultPrivKeyFilename)
	if err != nil {
		return nil, nil, err
	}
	return pubKey, privKey, nil
}

func ReadPubKeyFromFile(pubKeyPath string) (*PublicKey, error) {
	pubKeyEncodedBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	pubKey, err := DecodeBase64BLSPublicKey(pubKeyEncodedBytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func ReadPrivKeyFromFile(privKeyPath string) (PrivateKey, error) {
	privKeyEncodedBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}
	privKey, err := DecodeBase64BLSPrivateKey(privKeyEncodedBytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func GenerateAndStoreECDSAKeys(dir string) error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	err = crypto.SaveECDSA(dir+"/ecdsa", privateKey)
	if err != nil {
		return err
	}
	encodedPubKey := hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey))
	return os.WriteFile(dir+"/ecdsa.pub", []byte(encodedPubKey), 0o600)
}
