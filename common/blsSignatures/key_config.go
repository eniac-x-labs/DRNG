package blsSignatures

import (
	"errors"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type KeyConfig struct {
	KeyDir  string `koanf:"key-dir"`
	PrivKey string `koanf:"priv-key"`
}

func (c *KeyConfig) BLSPrivKey() (PrivateKey, error) {
	var privKeyBytes []byte
	if len(c.PrivKey) != 0 {
		privKeyBytes = []byte(c.PrivKey)
	} else if len(c.KeyDir) != 0 {
		var err error
		privKeyBytes, err = os.ReadFile(c.KeyDir + "/" + DefaultPrivKeyFilename)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("required BLS keypair did not exist at %s", c.KeyDir)
			}
			return nil, err
		}
	} else {
		return nil, errors.New("must specify privKey or KeyDir")
	}
	privKey, err := DecodeBase64BLSPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("'priv-key' was invalid: %w", err)
	}
	return privKey, nil
}

var DefaultKeyConfig = KeyConfig{}

func KeyConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.String(prefix+".key-dir", DefaultKeyConfig.KeyDir, fmt.Sprintf("the directory to read the bls keypair ('%s' and '%s') from; if using any of the DAS storage types exactly one of key-dir or priv-key must be specified", DefaultPubKeyFilename, DefaultPrivKeyFilename))
	f.String(prefix+".priv-key", DefaultKeyConfig.PrivKey, "the base64 BLS private key to use for signing DAS certificates; if using any of the DAS storage types exactly one of key-dir or priv-key must be specified")
}
