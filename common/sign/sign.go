package sign

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// hasherPool holds LegacyKeccak256 hashers for RlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

func rlpHash(x interface{}) (h common.Hash) {
	fmt.Printf("RlpHash: elem: %+v\n", x)
	var err error
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)

	sha.Reset()
	if err = rlp.Encode(sha, x); err != nil {
		fmt.Printf("RlpHash ERR1: %s\n", err.Error())
		return common.Hash{}
	}
	if _, err = sha.Read(h[:]); err != nil {
		fmt.Printf("RlpHash ERR2: %s\n", err.Error())

		return common.Hash{}
	}
	return h
}

//func SignHash(msgHash common.Hash, priv *ecdsa.PrivateKey) ([]byte, error) {
//	return crypto.Sign(msgHash.Bytes(), priv)
//}
//
//func VerifySignature(hash []byte, sig []byte, pub []byte) bool {
//	return crypto.VerifySignature(pub, hash, sig)
//}
//func GetRSVFromSignature(sig []byte) (*big.Int, *big.Int, *big.Int, error) {
//	if len(sig) != crypto.SignatureLength {
//		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
//	}
//	r := new(big.Int).SetBytes(sig[:32])
//	s := new(big.Int).SetBytes(sig[32:64])
//	v := new(big.Int).SetBytes([]byte{sig[64] + 27})
//	return r, s, v, nil
//}

func SignHash(data interface{}, priv *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
	hash := rlpHash(data)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash.Bytes())
	return r, s, err
}
func VerifySignature(pub *ecdsa.PublicKey, data interface{}, r, s *big.Int) (bool, error) {
	hash := rlpHash(data)
	//pub, err := crypto.DecompressPubkey(pubByte)
	//if err != nil {
	//	return false, err
	//}
	ok := ecdsa.Verify(pub, hash.Bytes(), r, s)
	return ok, nil
}

func EcdsaPubKeyToByte(key *ecdsa.PublicKey) []byte {
	return crypto.CompressPubkey(key)
}

func EcdsaPubKeyFromByte(pubByte []byte) (*ecdsa.PublicKey, error) {
	return crypto.DecompressPubkey(pubByte)
}
