package errors

import "errors"

// error message
const (
	ErrWaitingForGeneratingMsg   = "generating random numbers, please call later"
	LeftRandoNumbersNotEnoughMsg = "left random numbers not enough"
)

var (
	ErrWaitingForGenerating      = errors.New(ErrWaitingForGeneratingMsg)
	ErrLeftRandoNumbersNotEnough = errors.New(LeftRandoNumbersNotEnoughMsg)
)
