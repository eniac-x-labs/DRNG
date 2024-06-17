package core

import (
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

// 所有人都要处理timeout，因为任何人在任何阶段的err都没有广播
func (n *RNG) leaderDiscard() {
	log.Info("start leader discard")
	defer n.wg.Done()
	//ticker := time.NewTicker(n.discardDuration)
	timer := time.NewTimer(n.discardDuration)
	for {
		select {
		case <-timer.C:

			log.Debug("start discarding")

			// for from n.lastDiscardSeq+1,
			seq := n.lastDiscardSeq + 1

			n.leaderRecords.Lock()
			//log.Warn("进for前")

			//log.Debug("看全部", "all", n.leaderRecords.seqTimeoutQueue.Queue)
			var ti *time.Time
			for {
				if ti = n.leaderRecords.seqTimeoutQueue.FrontTime(); ti == nil {
					log.Debug("discarding, seqTimeoutQueue is empty")
					break
				}
				if time.Now().After(*ti) {
					seq = n.leaderRecords.seqTimeoutQueue.Dequeue().Seq

					log.Debug("discarding timeout record", "seq", seq)
					delete(n.leaderRecords.msgBuffer, seq)
					delete(n.leaderRecords.seqStatusMap, seq)
					delete(n.leaderRecords.mergedRandNums, seq)
					atomic.StoreInt64(&n.lastDiscardSeq, seq)
				} else {
					log.Debug("record not timeout yet, break loop", "seq", seq)
					break
				}
			}
			//log.Warn("break 了")

			n.leaderRecords.Unlock()
			timer.Reset(n.discardDuration)
			log.Debug("finish discarding")

		case <-n.ctx.Done():
			return
		}
	}
}

func (n *RNG) generatorDiscard() {
	log.Info("start generator discard")
	n.wg.Done()
	//ticker := time.NewTicker(n.discardDuration)
	timer := time.NewTimer(n.discardDuration)
	for {
		select {
		case <-timer.C:
			log.Debug("start discarding")
			// for from n.lastDiscardSeq+1,
			seq := n.lastDiscardSeq + 1
			n.generatorRecords.Lock()
			log.Warn("进for前")
			var ti *time.Time
			for {
				if ti = n.generatorRecords.seqTimeoutQueue.FrontTime(); ti == nil {
					log.Debug("discarding, seqTimeoutQueue is empty")
					break
				}
				if time.Now().After(*ti) {
					seq = n.generatorRecords.seqTimeoutQueue.Dequeue().Seq

					// store

					//del
					delete(n.generatorRecords.msgBuffer, seq)
					delete(n.generatorRecords.seqStatusMap, seq)

					// update lastDiscardSeq
					atomic.StoreInt64(&n.lastDiscardSeq, seq)
				} else {
					log.Debug("record not timeout yet, break loop", "seq", seq)
					break
				}
			}
			log.Warn("break 了")
			n.generatorRecords.Unlock()
			timer.Reset(n.discardDuration)
			log.Debug("finish discarding")
		case <-n.ctx.Done():
			return
		}
	}
}
