package common

import (
	"sync/atomic"

	"DRNG/common/types"
)

// RandomNumbersPool is a random numbers pool for leader.
// If RandomNumbersPool is not empty, return {reqID, randNum} from pool when a NewRandomNumberRequest is coming.
// As the same time, if the length of the pool is not larger than upperLimit,
// leader and generators also run a new round of random number generating, results will append into queue.
// Overwise, if the length of the pool is larger than upperLimit, new generating round will not start.
type RandomNumbersPool struct {
	Queue           Queue
	upperLimit      int // if length of queue is larger than upperLimit, do not start a new generate round
	latestSeqInPool int64
}

// PoolElem is the type of RandomNumbersPool's Queue
type PoolElem struct {
	ReqID     string `json:"reqId"`
	RandomNum uint64 `json:"randomNum"`
}

func NewRandomNumbersPool(upperLimit int, latestSeq int64) *RandomNumbersPool {
	if upperLimit < 0 {
		upperLimit = types.DefaultRandomNumbersPoolUpperLimit
	}
	return &RandomNumbersPool{
		Queue:           Queue{},
		upperLimit:      upperLimit,
		latestSeqInPool: latestSeq,
	}
}

func (p *RandomNumbersPool) GetUpperLimit() int {
	return p.upperLimit
}

func (p *RandomNumbersPool) GetLatestSeqInPool() int64 {
	return atomic.LoadInt64(&p.latestSeqInPool)
}

func (p *RandomNumbersPool) SetLatestSeqInPool(s int64) {
	atomic.StoreInt64(&p.latestSeqInPool, s)
}

func (p *RandomNumbersPool) PutRNGResults(reqID string, randNums []uint64) {
	p.Queue.Lock()
	defer p.Queue.Unlock()

	if randNums == nil || len(randNums) == 0 {
		return
	}
	for _, n := range randNums {
		p.Queue.items = append(p.Queue.items, PoolElem{
			ReqID:     reqID,
			RandomNum: n,
		})
	}

}

func (p *RandomNumbersPool) Size() int {
	p.Queue.RLock()
	defer p.Queue.RUnlock()
	return p.Queue.Size()
}

func (p *RandomNumbersPool) DequeueMultiRandNums(num int) ([]PoolElem, int) {
	p.Queue.Lock()
	defer p.Queue.Unlock()

	inters := p.Queue.DequeueMulti(num)
	if inters == nil {
		return nil, p.Queue.Size()
	}
	results := make([]PoolElem, num)
	for i, e := range inters {
		results[i] = e.(PoolElem)
	}
	return results, p.Queue.Size()
}
