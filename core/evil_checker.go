package core

import (
	"sync"

	"DRNG/common/blsSignatures"
	_types "DRNG/common/types"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/slog"
)

type EvilCheckerInterface interface {
	EvilCheckerOuterInterface
	AddEvilRecord(node string)
	AddBlackList(node string)
	CheckOneByOne(sigs []blsSignatures.Signature, plains [][]byte, pubs []blsSignatures.PublicKey, nodes []string) []string
}

type EvilCheckerOuterInterface interface {
	IsEvil(node string) bool
	GetEvils() []string
	GetBlackListLen() int
	FreeFromBlackList(node string)
	GetFailureTimesThreshold() int
}

type EvilChecker struct {
	BlackList        map[string]struct{} // records suspected of malicious nodes, nodeID:struct{}
	FailureTimes     sync.Map            // node:failureTimes
	FailureThreshold int
	logger           log.Logger
}

func NewEvilChecker() EvilCheckerInterface {
	return &EvilChecker{
		BlackList:        make(map[string]struct{}), // check when leader receives requests
		FailureTimes:     sync.Map{},
		FailureThreshold: _types.FailureThreshold,
		logger:           log.Root().With(slog.String("module", "EvilChecker")),
	}
}

func (e *EvilChecker) AddEvilRecord(node string) {
	var times int
	value, loaded := e.FailureTimes.Load(node)
	if loaded {
		// 如果 key 已存在，则对应的 value 加 1
		if v, ok := value.(int); ok {
			times = v + 1
		} else {
			e.logger.Error("failureTimes map load value not type int")
			return
		}
	} else {
		times = 1
	}
	e.logger.Debug("addEvilRecord", "node", node, "fail times", times)
	e.FailureTimes.Store(node, times)
	if times >= e.FailureThreshold {
		e.AddBlackList(node)
		e.logger.Debug("Determine evil node", "node", node, "failureTimes", times)
	}
}

func (e *EvilChecker) AddBlackList(node string) {
	e.BlackList[node] = struct{}{}
}

func (e *EvilChecker) GetBlackListLen() int {
	return len(e.BlackList)
}

func (e *EvilChecker) FreeFromBlackList(node string) {
	delete(e.BlackList, node)
	e.FailureTimes.Store(node, 0)
}

func (e *EvilChecker) IsEvil(node string) bool {
	if _, ok := e.BlackList[node]; ok {
		return true
	}
	return false
}

func (e *EvilChecker) GetEvils() []string {
	l := len(e.BlackList)
	res := make([]string, l)
	for n, _ := range e.BlackList {
		res = append(res, n)
	}
	return res
}

func (e *EvilChecker) GetFailureTimesThreshold() int {
	return e.FailureThreshold
}

func (e *EvilChecker) CheckOneByOne(sigs []blsSignatures.Signature, plains [][]byte, pubs []blsSignatures.PublicKey, nodes []string) []string {
	if sigs == nil || plains == nil || pubs == nil || nodes == nil {
		return nil
	}
	l := len(sigs)

	if l == 0 {
		return nil
	}

	if !(l == len(plains) && l == len(pubs) && l == len(nodes)) {
		e.logger.Error("args length not equal")
		return nil
	}

	evilNodes := make([]string, 0)
	for i := 0; i < l; i++ {
		sig, message, pub, node := sigs[i], plains[i], pubs[i], nodes[i]
		ok, err := blsSignatures.VerifySignature(sig, message, pub)
		if !ok {
			evilNodes = append(evilNodes, node)
			e.AddEvilRecord(node)
			if err == nil {
				e.logger.Error("check signature one by one failed with no err")
			} else {
				e.logger.Error("check signature one by one failed", "err", err)
			}
		}
	}
	return evilNodes
}
