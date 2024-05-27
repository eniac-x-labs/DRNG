package api

import (
	"net/http"
	"strconv"

	"DRNG/common"
	_types "DRNG/common/types"
	"DRNG/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"

	"golang.org/x/exp/slog"
)

type RNGOuterInterface interface {
	GetSingleRandomNumber() (common.PoolElem, error)
	GetMultiRandomNumbers(n int) ([]common.PoolElem, error)
	UpdateConfig(configPath string) error // update config with the configPath file
}

type Api struct {
	inner     core.RNGInterface
	evilOuter core.EvilCheckerOuterInterface
	reviewer  core.Reviewer
	logger    log.Logger
}

// ErrorResponse api response for abnormal situation
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse api response for success situation
type SuccessResponse struct {
	Result []common.PoolElem
}

// NewApi
func NewApi(rng *core.RNG) *Api {
	return &Api{
		inner:     rng,
		evilOuter: rng.EvilChecker,
		reviewer:  core.NewReviewer(rng),
		logger:    log.Root().With(slog.String("module", "api")),
	}
}

func PrepareWebServer(addrConf string, rng *core.RNG) *http.Server {
	apiHandler := NewApi(rng)
	r := gin.Default()
	//r.Use(RateLimitMiddleware(offchain.Config.RateLimiting.FillInterval, offchain.Config.RateLimiting.Cap, offchain.Config.RateLimiting.Quantum))

	// 定义路由和处理函数
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "200"})
	})
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, decentralised random number generator"})
	})

	r.GET("/getSingleRandomNumber", apiHandler.GetSingleRandomNumber)
	r.GET("/getMultiRandomNumbers/:num", apiHandler.GetMultiRandomNumbers)
	r.GET("/getPoolUpperLimit", apiHandler.GetPoolUpperLimit)
	r.GET("/getPoolLatestSeq", apiHandler.GetPoolLatestSeq)
	r.GET("/getLatestProcessingSeq", apiHandler.GetLatestProcessingSeq)
	r.GET("/getGeneratorsNum", apiHandler.GetGeneratorsNum)
	r.GET("/getLastDiscardSeq", apiHandler.GetLastDiscardSeq)
	r.GET("/getSizeofPool", apiHandler.GetSizeofPool)
	r.GET("/isEvil/:node", apiHandler.IsEvil)
	r.GET("/getEvils", apiHandler.GetEvils)
	r.GET("/getBlackListLen", apiHandler.GetBlackListLen)
	r.GET("/freeFromBlackList/:node", apiHandler.FreeFromBlackList)
	r.GET("/getFailureTimesThreshold", apiHandler.GetFailureTimesThreshold)

	r.POST("/examineRandNumWithLeaderAndGeneratorsLogs", apiHandler.ExamineRandNumWithLeaderAndGeneratorsLogs)

	addr := ":8080" // default localhost:8080
	if len(addrConf) != 0 {
		addr = addrConf
		log.Info("api service address set by config", "address", addr)
	} else {
		log.Info("api service address set by default", "address", addr)
	}
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

// @Summary Handle single random number request, return result or error meessage
// @Description get and return one random number from leader's randomNumbersPool, if pool is empty, maybe start a generate round
// @Accept json
// @Produce json
// @Success 200 {object} common.PoolElem "random number and corresponding reqID"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getSingleRandomNumber [get]
func (a *Api) GetSingleRandomNumber(c *gin.Context) {
	res, err := a.inner.GetSingleRandomNumber()
	if err != nil {
		a.logger.Debug("GetSingleRandomNumber failed", "err", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	} else {
		a.logger.Debug("GetSingleRandomNumber success", "result", res)
		c.JSON(http.StatusOK, res)
	}
}

// @Summary Handle multi random numbers request, return results or error meessage
// @Description get and return specified number of random numbers from leader's randomNumbersPool, if pool does not have enough numbers, maybe start a generate round
// @Accept json
// @Produce json
// @Param num path int true "number of random numbers requested" Minimum(1)
// @Success 200 {object} []common.PoolElem "random number and corresponding reqID"
// @Failure 400 {object} ErrorResponse "Not getting expected results"
// @Router /getMultiRandomNumbers/{num} [get]
func (a *Api) GetMultiRandomNumbers(c *gin.Context) {
	nStr := c.Param("num")
	a.logger.Debug("request GetMultiRandomNumbers", "askNum", nStr)
	n, err := strconv.Atoi(nStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := a.inner.GetMultiRandomNumbers(n)
	if err != nil {
		a.logger.Debug("GetMultiRandomNumbers failed", "err", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	} else {
		a.logger.Debug("GetMultiRandomNumbers success", "result", res)
		c.JSON(http.StatusOK, res)
	}
}

// @Summary Handle GetPoolUpperLimit request, return result or error meessage
// @Description get random numbers pool's upper limit which is set by config
// @Accept json
// @Produce json
// @Success 200 {object} int "random numbers' pool upper limit"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getPoolUpperLimit [get]
func (a *Api) GetPoolUpperLimit(c *gin.Context) {
	res := a.inner.GetPoolUpperLimit()
	c.JSON(http.StatusOK, res)
}

// @Summary Handle GetPoolLatestSeq request, return result or error meessage
// @Description get the latest seq in random numbers pool
// @Accept json
// @Produce json
// @Success 200 {object} int "latest seq in pool"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getPoolLatestSeq [get]
func (a *Api) GetPoolLatestSeq(c *gin.Context) {
	res := a.inner.GetPoolLatestSeq()
	c.JSON(http.StatusOK, res)
}

// @Summary Handle GetLatestProcessingSeq request, return result or error meessage
// @Description get the latest seq rng node is processing
// @Accept json
// @Produce json
// @Success 200 {object} int "latest seq rng node process"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getLatestProcessingSeq [get]
func (a *Api) GetLatestProcessingSeq(c *gin.Context) {
	res := a.inner.GetLatestProcessingSeq()
	c.JSON(http.StatusOK, res)
}

// @Summary Handle GetGeneratorsNum request, return result or error meessage
// @Description get the number of generator nodes
// @Accept json
// @Produce json
// @Success 200 {object} int "number of generators"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getGeneratorsNum [get]
func (a *Api) GetGeneratorsNum(c *gin.Context) {
	res := a.inner.GetGeneratorsNum()
	c.JSON(http.StatusOK, res)
}

// @Summary Handle GetLastDiscardSeq request, return result or error meessage
// @Description get the last seq rng node has discarded
// @Accept json
// @Produce json
// @Success 200 {object} int "last seq has discarded"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getLastDiscardSeq [get]
func (a *Api) GetLastDiscardSeq(c *gin.Context) {
	res := a.inner.GetLastDiscardSeq()
	c.JSON(http.StatusOK, res)
}

// @Summary GetSizeofPool request, return result or error meessage
// @Description get the random numbers pool size
// @Accept json
// @Produce json
// @Success 200 {object} int "random numbers pool size"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getSizeofPool [get]
func (a *Api) GetSizeofPool(c *gin.Context) {
	res := a.inner.GetSizeofPool()
	c.JSON(http.StatusOK, res)
}

// @Summary IsEvil, return whether the given node is evil
// @Description return given node is evil or not
// @Accept json
// @Produce json
// @Param node path string true "nodeID for checking evil"
// @Success 200 {object} bool " whether the given node is evil"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /isEvil/{node} [get]
func (a *Api) IsEvil(c *gin.Context) {
	node := c.Param("node")
	res := a.evilOuter.IsEvil(node)
	c.JSON(http.StatusOK, res)
}

// @Summary GetEvils get all evil nodes
// @Description return a list of evil nodes
// @Accept json
// @Produce json
// @Success 200 {object} []string "list of evil nodes"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getEvils [get]
func (a *Api) GetEvils(c *gin.Context) {
	res := a.evilOuter.GetEvils()
	c.JSON(http.StatusOK, res)
}

// @Summary GetBlackListLen request, return length of block list
// @Description get the length of block list
// @Accept json
// @Produce json
// @Success 200 {object} int "length of block list"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getBlackListLen [get]
func (a *Api) GetBlackListLen(c *gin.Context) {
	res := a.evilOuter.GetBlackListLen()
	c.JSON(http.StatusOK, res)
}

// @Summary FreeFromBlackList request, remove node from blacklist
// @Description remove node from blacklist
// @Accept json
// @Produce json
// @Param node path string true "nodeID waiting to remove from blacklist"
// @Success 200
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /freeFromBlackList/{node} [get]
func (a *Api) FreeFromBlackList(c *gin.Context) {
	node := c.Param("node")
	a.evilOuter.FreeFromBlackList(node)
	c.JSON(http.StatusOK, gin.H{"message": "200"})
}

// @Summary GetFailureTimesThreshold request, return threshold of number of failures for a node to become evil
// @Description get the threshold of number of failures for a node to become evil
// @Accept json
// @Produce json
// @Success 200 {object} int "threshold for node to become evil"
// @Failure 400 {object} ErrorResponse "Not get expected results"
// @Router /getFailureTimesThreshold [get]
func (a *Api) GetFailureTimesThreshold(c *gin.Context) {
	res := a.evilOuter.GetFailureTimesThreshold()
	c.JSON(http.StatusOK, res)
}

func (a *Api) ExamineRandNumWithLeaderAndGeneratorsLogs(c *gin.Context) {
	var req _types.ExamineRandNumRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	leaderSaved, generators, err := a.reviewer.ExamineRandNumWithLeaderAndGeneratorsLogs(req.ReqID, req.RandNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &_types.ExamineResponse{
		LeaderSaved:       leaderSaved,
		GeneratorsAddress: generators,
	})
}

func (a *Api) UpdateConfig(configPath string) error {
	return nil
}
