package config

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"DRNG/common"
	"DRNG/common/blsSignatures"
	"DRNG/common/eth"
	"DRNG/common/types"
	"DRNG/contract"
	_proto "DRNG/core/grpc/proto"
	"DRNG/db"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	EnvVarPrefix = "DRNG"
	DRNGCategory = "DRNGCategory"
	ConfigDir    = "./config"
	ConfigName   = "config"
	ConfigType   = "toml"
)

func prefixEnvVars(name string) []string {
	return []string{EnvVarPrefix + "_" + name}
}

var (
	/* Required Flags */
	BlsPrivKey = &cli.StringFlag{
		Name:     "privKey",
		Usage:    "bls12381 privateKey, to sign message bls signature",
		Value:    "",
		EnvVars:  prefixEnvVars("BLS_PRIV_KEY"),
		Category: DRNGCategory,
	}
)

type Config struct {
	Local  string
	Leader string

	//NodeInfoFromConfig   map[string]NodeInfo
	//UnConnectedNodeTable map[string]struct{}

	NodeURLTable       map[string]string
	NodeKeyTable       map[string]blsSignatures.PublicKey
	ConnectionMap      map[string]*grpc.ClientConn
	NodeRpcClientTable map[string]_proto.RandomNumberGeneratorClient

	Role            common.RoleType
	PrivKey         blsSignatures.PrivateKey
	BlsThreshold    int
	RpcChanLen      int
	GrpcServerPort  string
	TimeoutDuration time.Duration
	DiscardDuration time.Duration
	LastDiscardSeq  int64

	//RetryUnconnectedNodeDuration time.Duration

	RandomNumbersPoolUpperLimit int

	LogLevel string

	ApiServiceAddress string

	MustOnChain    bool
	RandStorageCli *eth.RandHashStorage

	EnableSDK         bool
	ExposedRpcAddress string

	MetaDBPath string
}

type fileConf struct {
	Local           string        `toml:"local"`
	Leader          string        `toml:"leader"`
	Role            string        `toml:"role"`
	PrivKey         string        `toml:"privKey"`
	BlsThreshold    int           `toml:"blsThreshold"`
	RpcChanLen      int           `toml:"rpcChanLen"`
	GrpcServerPort  string        `toml:"grpcServerPort"`
	ConnMaxAttempts int           `toml:"connMaxAttempts"`
	RetryDelay      time.Duration `toml:"retryDelay"`
	//RetryUnconnectedNodeDuration time.Duration `toml:"retryUnconnectedNodeDuration"`
	TimeoutDuration             time.Duration `toml:"timeoutDuration"`
	DiscardDuration             time.Duration `toml:"discardDuration"`
	LastDiscardSeq              int64         `toml:"lastDiscardSeq"`
	RandomNumbersPoolUpperLimit int           `toml:"randomNumbersPoolUpperLimit"`

	LogLevel string `toml:"logLevel"`

	ApiServiceAddress string `toml:"apiServiceAddress"`

	EnableSDK         bool   `toml:"enableSDK"`
	ExposedRpcAddress string `toml:"exposedRpcAddress"`

	MetaDBPath string    `toml:"metaDBPath"`
	EthConfig  EthConfig `toml:"ethConfig"`

	NodeInfos []NodeInfo `toml:"nodeInfos"`
}

type NodeInfo struct {
	NodeID string `toml:"nodeID"`
	Url    string `toml:"url"`
	PubKey string `toml:"pubKey"`
}

type EthConfig struct {
	Must            bool   `toml:"must"`
	ChainUrl        string `toml:"chainUrl"`
	ChainID         int64  `toml:"chainID"`
	ContractAddress string `toml:"contractAddress"` // RandomHashListStorageContractAddress
	PrivKey         string `toml:"privKey"`
}

func PrepareAndGetConfig(dir, fileName string) (*Config, error) {
	log.Info("preparing config", "config file dir", dir, "config file name", fileName)

	// 1. set viper args and read config file into 'fileConf'
	if len(dir) != 0 && len(fileName) != 0 {
		viper.SetConfigName(fileName)
		viper.AddConfigPath(dir)
		log.Debug("node config", "dir", dir, "file name", fileName)
	} else {
		viper.SetConfigName(ConfigName)
		viper.AddConfigPath(ConfigDir)
		log.Debug("node config", "dir", ConfigDir, "file name", ConfigName)

	}

	viper.SetConfigType(ConfigType)

	// privKey env flag: DRNG_PRIV_KEY
	viper.SetEnvPrefix(EnvVarPrefix)
	viper.BindEnv("PRIV_KEY")

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}

	confFromFile := &fileConf{}
	if err := viper.Unmarshal(confFromFile); err != nil {
		fmt.Println("Error parsing config file:", err)
		return nil, err
	}

	// 2. check self bls priv key
	if len(confFromFile.PrivKey) == 0 {
		// read env
		confFromFile.PrivKey = viper.GetString("PRIV_KEY")
		if len(confFromFile.PrivKey) == 0 {
			return nil, errors.New("priv key is required")
		}
	}
	privKey, err := blsSignatures.DecodeBase64BLSPrivateKey([]byte(confFromFile.PrivKey))
	if err != nil {
		return nil, err
	}
	pubKey, err := blsSignatures.PublicKeyFromPrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	// 3. fill maps
	urlMap := make(map[string]string)
	keyMap := make(map[string]blsSignatures.PublicKey) // include self
	clientMap := make(map[string]_proto.RandomNumberGeneratorClient)
	connectionMap := make(map[string]*grpc.ClientConn)
	//nodeInfosMap := make(map[string]NodeInfo)
	//unConnectedNodeMap := make(map[string]struct{})

	keyMap[confFromFile.Local] = pubKey

	fmt.Printf("confFromFile detail: %+v", confFromFile)
	succConn := 0
	for _, info := range confFromFile.NodeInfos {
		node, url, pubKeyStr := info.NodeID, info.Url, info.PubKey
		pubKey, err := blsSignatures.DecodeBase64BLSPublicKey([]byte(pubKeyStr))
		if err != nil {
			return nil, err
		}

		keyMap[node] = *pubKey
		//nodeInfosMap[node] = info

		attempt := 1
		var conn *grpc.ClientConn
		for attempt <= confFromFile.ConnMaxAttempts {
			log.Debug("Attempting to connect to remote node", "nodeID", node, "attempt", attempt, "maxAttempts", confFromFile.ConnMaxAttempts)
			conn, err = grpc.Dial(
				url,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithConnectParams(grpc.ConnectParams{
					Backoff: backoff.DefaultConfig,
				}),
			)
			if err != nil {
				log.Warn("Failed to connect to remote node", "nodeID", node, "err", err, "remote url", url)
				attempt++
				time.Sleep(confFromFile.RetryDelay)
				continue
			}
			break
		}

		if conn == nil {
			log.Warn("tried connecting remote grpc server failed", "times", confFromFile.ConnMaxAttempts)
			//unConnectedNodeMap[node] = struct{}{}
			continue
		}

		succConn++

		connectionMap[node] = conn
		clientMap[node] = _proto.NewRandomNumberGeneratorClient(conn)
		urlMap[node] = url

	}

	//if len(clientMap) < confFromFile.BlsThreshold {
	//	log.Error("connected grpc server not enough", "at least", confFromFile.BlsThreshold)
	//	return nil, errors.New(fmt.Sprintf("connected grpc server not enough, at least %d", confFromFile.BlsThreshold))
	//}
	log.Info("have connected remote grpc server", "num", len(clientMap))

	// prepare chain client for saving results on chain with contract
	var randStorage *eth.RandHashStorage
	url, addrHex, chainPrivHex, chainID := confFromFile.EthConfig.ChainUrl, confFromFile.EthConfig.ContractAddress, confFromFile.EthConfig.PrivKey, confFromFile.EthConfig.ChainID
	contractAddr := common2.HexToAddress(addrHex)
	chainPrivKey, err := crypto.HexToECDSA(chainPrivHex)
	if err != nil {
		log.Error("get chainPrivKey from hex failed, cannot save on chain", "err", err, "chainPrivKeyHex", chainPrivHex)
	} else {
		if len(url) == 0 || len(addrHex) == 0 || chainPrivKey == nil {
			log.Warn("chainUrl, contractAddress, and chainPrivKey are necessary for storing random numbers on chain", "url", url, "address", contractAddr, "privKeyHex", chainPrivHex)
		} else {
			log.Debug("parse eth config", "url", url, "address", contractAddr, "privKeyHex", chainPrivHex, "chainID", chainID)
			cli, err := ethclient.Dial(url)
			if err != nil {
				log.Error("ethClient dial failed", "err", err, "url", url, "address", addrHex)
			} else {
				storage, err := contract.NewRandomHashListStorage(contractAddr, cli)
				if err != nil {
					log.Error("contract.NewRandomHashListStorage failed", "err", err, "addr", addrHex)
				} else {
					randStorage = eth.NewRandHashStorage(storage, chainPrivKey, big.NewInt(chainID), cli)
				}
			}
		}
	}
	if randStorage == nil {
		log.Warn("the RandStorageCli is nil, will not save random numbers' hash an chain")
	}

	if confFromFile.RandomNumbersPoolUpperLimit <= 0 {
		confFromFile.RandomNumbersPoolUpperLimit = types.DefaultRandomNumbersPoolUpperLimit
	}

	if len(confFromFile.MetaDBPath) == 0 {
		confFromFile.MetaDBPath = fmt.Sprintf(db.DefaultLevelDBPathFormat, confFromFile.Local)
	}

	return &Config{
		Local:  confFromFile.Local,
		Leader: confFromFile.Leader,

		//NodeInfoFromConfig:   nodeInfosMap,
		//UnConnectedNodeTable: unConnectedNodeMap,

		NodeURLTable:       urlMap,
		NodeKeyTable:       keyMap,
		ConnectionMap:      connectionMap,
		NodeRpcClientTable: clientMap,
		Role:               common.String2RoleType(confFromFile.Role),
		PrivKey:            privKey,
		BlsThreshold:       confFromFile.BlsThreshold,
		RpcChanLen:         confFromFile.RpcChanLen,
		GrpcServerPort:     confFromFile.GrpcServerPort,
		TimeoutDuration:    confFromFile.TimeoutDuration,
		DiscardDuration:    confFromFile.DiscardDuration,
		LastDiscardSeq:     confFromFile.LastDiscardSeq,
		//RetryUnconnectedNodeDuration: confFromFile.RetryUnconnectedNodeDuration,
		RandomNumbersPoolUpperLimit: confFromFile.RandomNumbersPoolUpperLimit,
		LogLevel:                    confFromFile.LogLevel,
		ApiServiceAddress:           confFromFile.ApiServiceAddress,
		MustOnChain:                 confFromFile.EthConfig.Must,
		RandStorageCli:              randStorage,
		EnableSDK:                   confFromFile.EnableSDK,
		ExposedRpcAddress:           confFromFile.ExposedRpcAddress,
		MetaDBPath:                  confFromFile.MetaDBPath,
	}, nil
}
