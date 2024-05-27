package core

import (
	"sync"

	"DRNG/common"
	"DRNG/common/blsSignatures"
	"DRNG/common/types"
)

type collectionElem struct {
	randNum   uint64
	blsSig    blsSignatures.Signature
	signedMsg []byte
}
type GeneratorMsgBuffer struct {
	reqID string
	//status            types.Status
	singleRandomNum   uint64
	collectRandomNums map[string]*collectionElem // remote:{randomNum, blsSig}

}

func NewGeneratorMsgBuffer() *GeneratorMsgBuffer {
	return &GeneratorMsgBuffer{
		reqID:             "",
		singleRandomNum:   0,
		collectRandomNums: make(map[string]*collectionElem),
	}
}

type GeneratorRandomNumberRecord struct {
	sync.RWMutex
	//seqTimeoutMap map[int64]time.Time
	seqTimeoutQueue common.TimeoutQueue // FIFO, {seq, timeout}
	seqStatusMap    map[int64]types.Status
	msgBuffer       map[int64]*GeneratorMsgBuffer // seqNum:struct
}

func NewGeneratorRandomNumberRecord() GeneratorRandomNumberRecord {
	return GeneratorRandomNumberRecord{
		RWMutex: sync.RWMutex{},
		//seqTimeoutMap: make(map[int64]time.Time),
		seqTimeoutQueue: common.TimeoutQueue{},
		seqStatusMap:    make(map[int64]types.Status),
		msgBuffer:       make(map[int64]*GeneratorMsgBuffer),
	}
}

type RandTimes map[uint64]uint64 // rand:times

type CollectionBuffer struct {
	remote2Collection map[string][]uint64  // remote: remoteCollectionRandNums
	nodeRandTimes     map[string]RandTimes // from:{rand: times}, avoid evil node send different random numbers to different nodes
	plainAndSig       map[string]struct {
		plain []byte
		sig   blsSignatures.Signature
	}
}

func NewCollectionBuffer() *CollectionBuffer {
	return &CollectionBuffer{
		remote2Collection: make(map[string][]uint64),
		nodeRandTimes:     make(map[string]RandTimes),
		plainAndSig: make(map[string]struct {
			plain []byte
			sig   blsSignatures.Signature
		}),
	}
}

type LeaderRandomNumberRecord struct {
	sync.RWMutex
	//Req2Seq        map[string]int64           // reqID: seqNum
	//seq2Req        map[int64]string           // seqNum: reqID
	//seqTimeoutMap   map[int64]time.Time
	seqTimeoutQueue common.TimeoutQueue         // FIFO, {seq, timeout}
	seqStatusMap    map[int64]types.Status      // seqNum: status
	msgBuffer       map[int64]*CollectionBuffer // seqNum: remotesCollections

	mergedRandNums map[int64][]uint64 // seqNUm:randoms
}

func NewLeaderRandomNumberRecord() LeaderRandomNumberRecord {
	return LeaderRandomNumberRecord{
		RWMutex: sync.RWMutex{},

		//seqTimeoutMap:  make(map[int64]time.Time),
		seqTimeoutQueue: common.TimeoutQueue{},
		seqStatusMap:    make(map[int64]types.Status),
		msgBuffer:       make(map[int64]*CollectionBuffer),
		mergedRandNums:  make(map[int64][]uint64),
	}
}
