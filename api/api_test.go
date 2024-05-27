package api

import (
	"context"
	"time"

	"DRNG/common"
	_conf "DRNG/config"
	"DRNG/core"
)

func newMockLeader() (*core.RNG, error) {
	conf := &_conf.Config{
		Local:                       "node-local",
		Leader:                      "node-local",
		NodeURLTable:                nil,
		NodeKeyTable:                nil,
		ConnectionMap:               nil,
		NodeRpcClientTable:          nil,
		Role:                        common.LEADER,
		PrivKey:                     nil,
		BlsThreshold:                0,
		RpcChanLen:                  10,
		GrpcServerPort:              "",
		TimeoutDuration:             10 * time.Second,
		DiscardDuration:             10 * time.Second,
		LastDiscardSeq:              0,
		RandomNumbersPoolUpperLimit: 20,
		ApiServiceAddress:           "",
		RandStorageCli:              nil,
	}
	return core.NewRNG(context.Background(), conf)
}
