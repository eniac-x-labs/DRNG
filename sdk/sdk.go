package sdk

import (
	"net/rpc"

	_common "DRNG/common"
	_types "DRNG/common/types"
	"DRNG/core"
	"github.com/ethereum/go-ethereum/common"
)

// drng rpc client impl
type DRNGSDKInterface interface {
	GetPoolUpperLimit() (int, error)
	GetPoolLatestSeq() (int64, error)
	GetLatestProcessingSeq() (int64, error)
	GetGeneratorsNum() (int, error)
	GetLastDiscardSeq() (int64, error)
	GetSizeofPool() (int, error)
	GetSingleRandomNumber() (_common.PoolElem, error)
	GetMultiRandomNumbers(num int) ([]_common.PoolElem, error)
	GetBlackList() ([]string, error)
	GetBlackListLen() (int, error)
	FreeNodeFromBlackList(node string) error
	core.ReviewerOuter
}

type DRNGSDK struct {
	*rpc.Client
}

func NewDRNGSdk(addr string) (DRNGSDKInterface, error) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &DRNGSDK{client}, nil
}

func (d *DRNGSDK) GetPoolUpperLimit() (int, error) {
	var res int
	err := d.Call("DRNGRpcServer.GetPoolUpperLimit", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetPoolLatestSeq() (int64, error) {
	var res int64
	err := d.Call("DRNGRpcServer.GetPoolLatestSeq", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetLatestProcessingSeq() (int64, error) {
	var res int64
	err := d.Call("DRNGRpcServer.GetLatestProcessingSeq", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetGeneratorsNum() (int, error) {
	var res int
	err := d.Call("DRNGRpcServer.GetGeneratorsNum", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetLastDiscardSeq() (int64, error) {
	var res int64
	err := d.Call("DRNGRpcServer.GetLastDiscardSeq", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetSizeofPool() (int, error) {
	var res int
	err := d.Call("DRNGRpcServer.GetSizeofPool", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetSingleRandomNumber() (_common.PoolElem, error) {
	var res _common.PoolElem
	err := d.Call("DRNGRpcServer.GetSingleRandomNumber", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetMultiRandomNumbers(num int) ([]_common.PoolElem, error) {
	//res := make([]common.PoolElem, num)
	var res []_common.PoolElem
	err := d.Call("DRNGRpcServer.GetMultiRandomNumbers", &num, &res)
	return res, err
}

func (d *DRNGSDK) GetBlackList() ([]string, error) {
	var res []string
	err := d.Call("DRNGRpcServer.GetBlackList", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) GetBlackListLen() (int, error) {
	var res int
	err := d.Call("DRNGRpcServer.GetBlackListLen", struct{}{}, &res)
	return res, err
}

func (d *DRNGSDK) FreeNodeFromBlackList(node string) error {
	return d.Call("DRNGRpcServer.FreeNodeFromBlackList", &node, struct{}{})
}

func (d *DRNGSDK) ExamineRandNumWithLeaderAndGeneratorsLogs(reqID string, num uint64) (bool, []common.Address, error) {
	var resp _types.ExamineResponse
	err := d.Call("DRNGRpcServer.ExamineRandNumWithLeaderAndGeneratorsLogs", &_types.ExamineRandNumRequest{
		ReqID:   reqID,
		RandNum: num,
	}, &resp)
	if err != nil {
		return false, nil, err
	}
	return resp.LeaderSaved, resp.GeneratorsAddress, nil
}
