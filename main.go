package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"DRNG/api"
	"DRNG/common"
	_conf "DRNG/config"
	"DRNG/core"
	_rpc "DRNG/rpc"
	"github.com/ethereum/go-ethereum/log"
	flag "github.com/spf13/pflag"
)

// ./DRNG --confdir ./tests/conf --confname config-1
// ./DRNG --confdir ./tests/conf --confname config-2
// ./DRNG --confdir ./tests/conf --confname config-3
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true)))

	dirPtr := flag.String("confdir", "./config", "file dir of config file")
	namePtr := flag.String("confname", "config", "file name of config file")
	flag.Parse()

	conf, err := _conf.PrepareAndGetConfig(*dirPtr, *namePtr)
	if err != nil {
		log.Error("prepare and get config failed", "err", err)
		return
	}
	if len(conf.LogLevel) != 0 {
		common.SetLogLevel(conf.LogLevel)
	}

	rng, err := core.NewRNG(ctx, conf)
	if err != nil {
		log.Error("new RNG failed", "err", err)
		return
	}
	log.Debug("get config", "detail", conf)

	log.Info("RNG is trying starting")
	var webServer *http.Server
	if conf.Role.String() == common.GENERATORSTR {
		rng.Start()
	} else { // Leader
		go rng.Start()
		if conf.EnableSDK && len(conf.ExposedRpcAddress) != 0 {
			log.Info("starting exposed rpc server for sdk")
			go _rpc.NewAndStartDRNGRpcServer(ctx, conf.ExposedRpcAddress, rng)
		}

		// =================================== go-gin ===========================================
		// ================================= web server =========================================
		log.Info("preparing DRNG web server")
		webServer = api.PrepareWebServer(conf.ApiServiceAddress, rng)
		go func() {
			if err = webServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Println("Error starting server:", err)
			}
		}()
		log.Info("DRNG web server started successfully")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")
	cancel()

	if webServer != nil {
		if err := webServer.Shutdown(ctx); err != nil {
			fmt.Println("Error shutting down server:", err)
		}
	}

	fmt.Println("Server gracefully stopped")
}
