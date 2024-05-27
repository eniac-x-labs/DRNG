package config

import (
	"os"
	"testing"

	"DRNG/common/blsSignatures"
	"github.com/stretchr/testify/assert"
)

func Test_config(t *testing.T) {
	ast := assert.New(t)
	privString, err := blsSignatures.GeneratePrivKeyString()
	ast.NoError(err)
	t.Logf("privKeyStr: %s", privString)
	ast.NoError(os.Setenv("DRNG_PRIV_KEY", privString))

	conf, err := PrepareAndGetConfig("", "")
	if err != nil {
		t.Logf("getConfig err: %s", err.Error())
		return
	}
	t.Logf("%+v", conf)
}
