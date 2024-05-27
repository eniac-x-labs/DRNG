package core

import (
	"context"
	"math/big"
	"testing"

	"DRNG/common/eth"
	"DRNG/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

var (
	url             = "http://127.0.0.1:8545"
	contractAddrHex = "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"
	accoutPivHex    = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	chainID         = 31337
)

func Test_Filter(t *testing.T) {
	ast := assert.New(t)
	cli, err := ethclient.Dial(url)
	ast.NoError(err)
	contractAddr := common.HexToAddress(contractAddrHex)
	storage, err := contract.NewRandomHashListStorage(contractAddr, cli)
	chainPrivKey, err := crypto.HexToECDSA(accoutPivHex)
	ast.NoError(err)
	randStorage := eth.NewRandHashStorage(storage, chainPrivKey, big.NewInt(int64(chainID)), cli)

	iter, err := randStorage.FilterHashListStored(&bind.FilterOpts{}, []string{"reqID1"})
	ast.NoError(err)
	defer iter.Close()
	for iter.Next() {
		t.Log(iter.Event.ReqID, iter.Event.HashList)
	}

	reqID1Hash := common.HexToHash(hexutil.Encode([]byte("reqID1")))
	ast.NoError(err)
	t.Log(reqID1Hash.Hex())

	iter2, err := randStorage.FilterCollectionSaved(&bind.FilterOpts{}, []string{"reqID1"}, nil)
	ast.NoError(err)
	defer iter2.Close()
	for iter2.Next() {
		t.Log(iter2.Event.ReqID, iter2.Event.Sender, iter2.Event.Collection)
		p, _ := randStorage.ParseCollectionSaved(iter2.Event.Raw)
		t.Logf("%+v", p)

	}

	suggGasPrice, err := cli.SuggestGasPrice(context.Background())
	ast.NoError(err)
	t.Logf("suggest gas price is %d", suggGasPrice)

	crypto.PubkeyToAddress(chainPrivKey.PublicKey)
	accBal, err := cli.BalanceAt(context.Background(), crypto.PubkeyToAddress(chainPrivKey.PublicKey), nil)
	ast.NoError(err)
	t.Logf("balance is %d", accBal)
}
