package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	_proto "DRNG/core/grpc/proto"
	"github.com/ethereum/go-ethereum/log"
)

type RNGRPCServer struct {
	RandomServer _proto.RandomNumberGeneratorServer
}

func StartRNGRPCServer(ctx context.Context, addr string, portNum uint64, rng _proto.RandomNumberGeneratorServer) (*http.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, portNum))
	if err != nil {
		return nil, err
	}
	return StartRNGRPCServerOnListener(ctx, listener, rng)
}

func StartRNGRPCServerOnListener(ctx context.Context, listener net.Listener, rng _proto.RandomNumberGeneratorServer) (*http.Server, error) {
	rpcServer := rpc.NewServer()
	err := rpcServer.RegisterName("rng", &RNGRPCServer{
		rng,
	})
	if err != nil {
		log.Error("rngServer register failed", "err", err)
		return nil, err
	}

	srv := &http.Server{
		Handler: rpcServer,
	}

	log.Info("rngServer start")
	go func() {
		err := srv.Serve(listener)
		if err != nil {
			log.Error("rngServer serve listener failed", "err", err)
			return
		}
	}()
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()
	return srv, nil
}

func (s *RNGRPCServer) CheckLive(ctx context.Context) (string, error) {
	return "succ", nil
}

func (s *RNGRPCServer) RequestNewRandomNumber(ctx context.Context) (int64, string, error) {
	resp, err := s.RandomServer.RequestRandomNumber(ctx, nil)
	return resp.GetSeqNum(), resp.GetReqID(), err
}
