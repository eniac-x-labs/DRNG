package types

const (
	StartType          = "StartType"
	SingleRandType     = "SingleRandType"
	CollectionRandType = "CollectionRandType"
	MergedRandType     = "MergedRandType"
	UnknownType        = "UnknownType"
)

const (
	Suffix_1 = "suffix_1" // send single
	Suffix_2 = "suffix_2" // send collection
)

type Status uint8

// status
const (
	Idle Status = 1 << iota
	Generated
	Collecting
	Mergeing
	Done
	Safe      // onchain
	Finalized // pass checkpoint
)

// response msg string
const (
	Success = "200"
	Fail    = "400"
)

// errType for errMessage through errCh
const (
	VerifyAggrSignErr = "VerifyAggrSignErr"
	InvalidAggrSign   = "InvalidAggrSign"
)

const RpcServiceName = "RpcServiceName"

const (
	DefaultRandomNumbersPoolUpperLimit = 100
	FailureThreshold                   = 3
)
