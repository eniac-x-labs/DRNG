package types

import (
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type SeqID uint64
type LifeCycleMeta struct {
	Timeout       time.Time
	ReqID         string
	Status        Status
	SelfRandNum   uint64
	remoteRandNum map[string]uint64 // target:randomNum
	mergedRandNum []uint64
}
type RandomNumberRecords map[SeqID]*LifeCycleMeta

type RpcMsg struct {
	MsgType  string
	ProtoMsg interface{}
}

type ReqInfo struct {
	SeqNum  int64
	ReqID   string
	Timeout time.Time
}
type LeaderReqs struct {
	//ProcessingReqs []ReqInfo
	PendingSeqNum atomic.Int64
}

type ErrMessage struct {
	Err     error
	ErrType string
	From    string
	SeqNum  int64
	ReqID   string
}

type ChainResult struct {
	TxHash common.Hash
	Err    error
}

type ExamineRandNumRequest struct {
	ReqID   string `json:"reqID"`
	RandNum uint64 `json:"randNum"`
}

type ExamineResponse struct {
	LeaderSaved       bool             `json:"leaderSaved"`
	GeneratorsAddress []common.Address `json:"GeneratorsAddress"`
}
