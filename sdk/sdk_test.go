package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SDK(t *testing.T) {
	ast := assert.New(t)
	addr := ":8081"
	sdk, err := NewDRNGSdk(addr)
	ast.NoError(err)
	ast.NotNil(sdk)

	// GetPoolUpperLimit
	upper, err := sdk.GetPoolUpperLimit()
	ast.NoError(err)
	t.Logf("upper limit: %d", upper)

	poolLatestSeq, err := sdk.GetPoolLatestSeq()
	ast.NoError(err)
	t.Logf("poolLatestSeq: %d", poolLatestSeq)

	latestProcessingSeq, err := sdk.GetLatestProcessingSeq()
	ast.NoError(err)
	t.Logf("latestProcessingSeq: %d", latestProcessingSeq)

	generatorsNum, err := sdk.GetGeneratorsNum()
	ast.NoError(err)
	t.Logf("generatorsNum: %d", generatorsNum)

	lastDiscardSeq, err := sdk.GetLastDiscardSeq()
	ast.NoError(err)
	t.Logf("lastDiscardSeq: %d", lastDiscardSeq)

	sizeofPool, err := sdk.GetSizeofPool()
	ast.NoError(err)
	t.Logf("sizeofPool: %d", sizeofPool)
}

func Test_GetSingleRandomNumber(t *testing.T) {
	ast := assert.New(t)
	addr := ":8081"
	sdk, err := NewDRNGSdk(addr)
	ast.NoError(err)
	ast.NotNil(sdk)

	singleRandNum, err := sdk.GetSingleRandomNumber()
	ast.NoError(err)
	t.Logf("single random number %v", singleRandNum)
}

func Test_GetMultiRandomNumbers(t *testing.T) {
	ast := assert.New(t)
	addr := ":8081"
	sdk, err := NewDRNGSdk(addr)
	ast.NoError(err)
	ast.NotNil(sdk)

	multiRandNums, err := sdk.GetMultiRandomNumbers(2)
	ast.NoError(err)
	t.Logf("multi random numbers %v", multiRandNums)
}

func Test_Evil(t *testing.T) {
	ast := assert.New(t)
	addr := ":8081"
	sdk, err := NewDRNGSdk(addr)
	ast.NoError(err)
	ast.NotNil(sdk)

	t.Log(sdk.GetBlackList())
}

func Test_Reviewer(t *testing.T) {
	ast := assert.New(t)
	addr := ":8081"
	sdk, err := NewDRNGSdk(addr)
	ast.NoError(err)
	ast.NotNil(sdk)

	reqID, randNum := "c696a649-ee9e-4e58-b3c5-94d2832ad6ee", uint64(13763515586510503104)
	t.Log(sdk.ExamineRandNumWithLeaderAndGeneratorsLogs(reqID, randNum))
}
