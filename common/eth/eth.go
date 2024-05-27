package eth

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"hash"
	"math/big"

	_common "DRNG/common"
	"DRNG/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

var EmptyHash = common.Hash{}

type RandHashStorage struct {
	*contract.RandomHashListStorage
	hasher  hash.Hash
	privKey *ecdsa.PrivateKey
	chainID *big.Int
	*ethclient.Client
}

func NewRandHashStorage(storageCli *contract.RandomHashListStorage, privKey *ecdsa.PrivateKey, chainID *big.Int, client *ethclient.Client) *RandHashStorage {
	return &RandHashStorage{
		RandomHashListStorage: storageCli,
		hasher:                sha256.New(),
		privKey:               privKey,
		chainID:               chainID,
		Client:                client,
	}
}

func (r *RandHashStorage) LeaderSaveResultsOnChain(reqID string, randNums []uint64) (common.Hash, error) {
	randHashs := _common.Uint64sToStrings(r.hasher, randNums)
	if randHashs == nil {
		return EmptyHash, nil
	}
	opt, err := bind.NewKeyedTransactorWithChainID(r.privKey, r.chainID)
	if err != nil {
		log.Error("NewKeyedTransactorWithChainID failed", "err", err)
		return EmptyHash, err
	}

	opt.GasPrice, err = r.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("leader SuggestGasPrice failed", "err", err)
		return EmptyHash, err
	}

	log.Debug("leader storing on chain", "reqID", reqID, "randNum", randNums, "randomHashs", randHashs)
	tx, err := r.StoreHashList(opt, reqID, randHashs)
	if err != nil {
		log.Error("leader StoreHashList failed", "err", err)
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (r *RandHashStorage) GeneratorSaveCollectionOnChain(reqID string, randNums []uint64) {
	if randNums == nil || len(randNums) == 0 {
		log.Error("GeneratorSaveCollectionOnChain got no randNums")
		return
	}
	randHashs := _common.Uint64sToStrings(r.hasher, randNums)
	if randHashs == nil {
		log.Error("turn randNums to hashs get nil result")
		return
	}

	log.Debug("generators storing on chain", "randNums", randNums, "randomHashs", randHashs)
	opt, err := bind.NewKeyedTransactorWithChainID(r.privKey, r.chainID)
	if err != nil {
		log.Error("NewKeyedTransactorWithChainID failed", "err", err)
		return
	}

	opt.GasPrice, err = r.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("generator SuggestGasPrice failed", "err", err)
		return
	}

	tx, err := r.GeneratorSaveCollection(opt, reqID, randHashs)
	if err != nil {
		log.Error("GeneratorSaveCollection failed", "err", err, "reqID", reqID)
		return
	}
	log.Debug("GeneratorSaveCollectionOnChain success", "txHash", tx.Hash().Hex())
}

//func (r *RandHashStorage) HasRandNumHashOnChain(reqID string, hash string) (bool, error) {
//	iter, err := r.FilterHashListStored(&bind.FilterOpts{}, []string{reqID})
//	if err != nil {
//		log.Error("FilterHashListStored", "err", err)
//		return false, err
//	}
//
//	hashList := make([]string, 0)
//	for iter.Next() {
//		hashList = append(hashList, iter.Event.HashList...)
//	}
//	for _, h := range hashList {
//		if h == hash {
//			return true, nil
//		}
//	}
//	return false, nil
//}

func (r *RandHashStorage) GetHasher() hash.Hash {
	return r.hasher
}
