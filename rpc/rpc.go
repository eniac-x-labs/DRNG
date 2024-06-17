package rpc

import (
	"context"
	"net"
	"net/rpc"

	_common "DRNG/common"
	_errors "DRNG/common/errors"
	_types "DRNG/common/types"
	"DRNG/core"
	"github.com/ethereum/go-ethereum/log"
)

type DRNGRpcInterface interface {
	GetPoolUpperLimit(_ struct{}, reply *int) error
	GetPoolLatestSeq(_ struct{}, reply *int64) error
	GetLatestProcessingSeq(_ struct{}, reply *int64) error
	GetGeneratorsNum(_ struct{}, reply *int) error
	GetLastDiscardSeq(_ struct{}, reply *int64) error
	GetSizeofPool(_ struct{}, reply *int) error
	GetSingleRandomNumber(_ struct{}, reply *_common.PoolElem) error
	GetMultiRandomNumbers(num *int, reply *[]_common.PoolElem) error
	GetBlackList(_ struct{}, reply *[]string) error
	GetBlackListLen(_ struct{}, reply *int) error
	FreeNodeFromBlackList(node *string, _ struct{}) error
	IsEvil(node *string, reply *bool) error
	GetFailureTimesThreshold(_ struct{}, reply *int) error
	ExamineRandNumWithLeaderAndGeneratorsLogs(req *_types.ExamineRandNumRequest, reply *_types.ExamineResponse) error
}

type DRNGRpcServer struct {
	rng         *core.RNG
	evilChecker core.EvilCheckerOuterInterface
	reviewer    core.Reviewer
}

func NewAndStartDRNGRpcServer(ctx context.Context, address string, rng *core.RNG) {
	if err := rpc.Register(&DRNGRpcServer{
		rng:         rng,
		evilChecker: rng.EvilChecker,
		reviewer:    core.NewReviewer(rng),
	}); err != nil {
		log.Error("RpcServer Register failed", "err", err)
	}
	log.Debug("RpcServer Register finished")

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("RpcServer Listen failed", "err", err, "address", address)
		return
	}
	log.Debug("RpcServer listen address finished", "address", address)

	for {
		select {
		case <-ctx.Done():
			listener.Close()
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Error("RpcServer listener.Accept failed", "err", err)
			}

			go rpc.ServeConn(conn)
		}

	}
}

func (s *DRNGRpcServer) GetPoolUpperLimit(_ struct{}, reply *int) error {
	*reply = s.rng.GetPoolUpperLimit()
	return nil
}

func (s *DRNGRpcServer) GetPoolLatestSeq(_ struct{}, reply *int64) error {
	*reply = s.rng.GetPoolLatestSeq()
	return nil
}

func (s *DRNGRpcServer) GetLatestProcessingSeq(_ struct{}, reply *int64) error {
	*reply = s.rng.GetLatestProcessingSeq()
	return nil
}

func (s *DRNGRpcServer) GetGeneratorsNum(_ struct{}, reply *int) error {
	*reply = s.rng.GetGeneratorsNum()
	return nil
}

func (s *DRNGRpcServer) GetLastDiscardSeq(_ struct{}, reply *int64) error {
	*reply = s.rng.GetLastDiscardSeq()
	return nil
}

func (s *DRNGRpcServer) GetSizeofPool(_ struct{}, reply *int) error {
	*reply = s.rng.GetSizeofPool()
	return nil
}

func (s *DRNGRpcServer) GetSingleRandomNumber(_ struct{}, reply *_common.PoolElem) error {
	res, poolSize := s.rng.GetRandNumFromPool(1)
	s.rng.JudgeAndStartNewGenerateRound(poolSize)
	if res == nil || len(res) == 0 {
		*reply = _common.PoolElem{}
		return _errors.ErrLeftRandoNumbersNotEnough
	}
	*reply = res[0]
	return nil
}

func (s *DRNGRpcServer) GetMultiRandomNumbers(num *int, reply *[]_common.PoolElem) error {
	res, pooSize := s.rng.GetRandNumFromPool(*num)
	if err := s.rng.JudgeAndStartNewGenerateRound(pooSize); err != nil {
		log.Warn("JudgeAndStartNewGenerateRound has error", "err", err)
	}

	if res == nil || len(res) == 0 {
		return _errors.ErrLeftRandoNumbersNotEnough
	}

	*reply = res
	return nil
}

func (s *DRNGRpcServer) GetBlackList(_ struct{}, reply *[]string) error {
	*reply = s.evilChecker.GetEvils()
	return nil
}

func (s *DRNGRpcServer) GetBlackListLen(_ struct{}, reply *int) error {
	*reply = s.evilChecker.GetBlackListLen()
	return nil
}

func (s *DRNGRpcServer) FreeNodeFromBlackList(node *string, _ struct{}) error {
	s.evilChecker.FreeFromBlackList(*node)
	return nil
}

func (s *DRNGRpcServer) IsEvil(node *string, reply *bool) error {
	*reply = s.evilChecker.IsEvil(*node)
	return nil
}

func (s *DRNGRpcServer) GetFailureTimesThreshold(_ struct{}, reply *int) error {
	*reply = s.evilChecker.GetFailureTimesThreshold()
	return nil
}

func (s *DRNGRpcServer) ExamineRandNumWithLeaderAndGeneratorsLogs(req *_types.ExamineRandNumRequest, reply *_types.ExamineResponse) error {
	leaderSaved, generators, err := s.reviewer.ExamineRandNumWithLeaderAndGeneratorsLogs(req.ReqID, req.RandNum)
	if err != nil {
		return err
	}
	*reply = _types.ExamineResponse{
		LeaderSaved:       leaderSaved,
		GeneratorsAddress: generators,
	}
	return nil
}
