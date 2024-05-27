package core

import _common "DRNG/common"

type RNGInterface interface {
	GetPoolUpperLimit() int
	GetPoolLatestSeq() int64
	GetLatestProcessingSeq() int64
	GetGeneratorsNum() int
	GetLastDiscardSeq() int64
	GetSizeofPool() int
	GetSingleRandomNumber() (_common.PoolElem, error)
	GetMultiRandomNumbers(num int) ([]_common.PoolElem, error)
}
