package core

import (
	_common "DRNG/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type Reviewer interface {
	ReviewerOuter
	ReviewLeaderLog(reqID string, num uint64) (bool, error)
	ReviewGeneratorsLog(reqID string, num uint64) ([]common.Address, error)
}

type ReviewerOuter interface {
	ExamineRandNumWithLeaderAndGeneratorsLogs(reqID string, num uint64) (bool, []common.Address, error)
}

type ReviewerImpl struct {
	inner *RNG
}

func NewReviewer(rng *RNG) Reviewer {
	return &ReviewerImpl{
		rng,
	}
}

// ExamineRandNumWithLeaderAndGeneratorsLogs check whether leader saved randNum on chain,
// check which generators have recorded this random number, and return their address
func (r *ReviewerImpl) ExamineRandNumWithLeaderAndGeneratorsLogs(reqID string, num uint64) (bool, []common.Address, error) {
	log.Debug("Reviewer start reviewing leader logs", "reqID", reqID, "randNum", num)
	leaderHas, err := r.ReviewLeaderLog(reqID, num)
	if err != nil {
		log.Error("Reviewer check leader log failed", "err", err, "reqID", reqID, "randNum", num)
		return leaderHas, nil, err
	}

	log.Debug("Reviewer start reviewing generators logs", "reqID", reqID, "randNum", num)
	generators, err := r.ReviewGeneratorsLog(reqID, num)
	if err != nil {
		log.Error("Reviewer check generators logs failed", "err", err, "reqID", reqID, "randNum", num)
	} else {
		log.Debug("Reviewer check generators logs", "reqID", reqID, "randNum", num, "witnessNum", len(generators), "witness", generators)
	}
	return leaderHas, generators, err
}

func (r *ReviewerImpl) ReviewLeaderLog(reqID string, num uint64) (bool, error) {
	iter, err := r.inner.randHashStorageClient.FilterHashListStored(&bind.FilterOpts{}, []string{reqID})
	if err != nil {
		return false, err
	}
	defer iter.Close()

	numHash := _common.Uint64sToStrings(r.inner.randHashStorageClient.GetHasher(), []uint64{num})[0]
	for iter.Next() {
		for _, h := range iter.Event.HashList {
			if h == numHash {
				return true, nil
			}
		}
	}
	return false, nil
}

func (r *ReviewerImpl) ReviewGeneratorsLog(reqID string, num uint64) ([]common.Address, error) {
	iter, err := r.inner.randHashStorageClient.FilterCollectionSaved(&bind.FilterOpts{}, []string{reqID}, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	res := make([]common.Address, 0)
	numHash := _common.Uint64sToStrings(r.inner.randHashStorageClient.GetHasher(), []uint64{num})[0]
	for iter.Next() {
		for _, h := range iter.Event.Collection {
			if h == numHash {
				res = append(res, iter.Event.Sender)
			}
		}
	}
	return res, nil
}
