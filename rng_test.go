package main

import (
	"context"
	"os"
	"testing"

	"DRNG/config"
	"DRNG/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

func Test_NewRNG(t *testing.T) {
	log.SetDefault(log.NewLogger(log.NewTerminalHandler(os.Stdout, true)))

	ast := assert.New(t)
	rng, err := core.NewRNG(context.Background(), &config.Config{})
	ast.NoError(err)
	rng.Start()
}
