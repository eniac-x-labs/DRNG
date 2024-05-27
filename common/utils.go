package common

import (
	"hash"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

func TimeToBytes(t time.Time) []byte {
	// 这里使用 time.Time 的字节表示，通常是 time.FixedZone 的字节表示
	return []byte(t.Format(time.RFC3339))
}

// BytesToTime
func BytesToTime(b []byte) (time.Time, error) {
	return time.Parse(time.RFC3339, string(b))
}

func Uint64sToStrings(hasher hash.Hash, randNums []uint64) []string {
	if len(randNums) <= 0 {
		return nil
	}
	randHashs := make([]string, len(randNums), len(randNums))

	//buf := make([]byte, 8)
	//var hashHex string
	for i, n := range randNums {
		hasher.Reset()
		hasher.Write([]byte(hexutil.EncodeUint64(n)))
		//binary.BigEndian.PutUint64(buf, n)
		//hasher.Write(buf)
		hashByte := hasher.Sum(nil)
		hashHex := hexutil.Encode(hashByte)
		randHashs[i] = hashHex
	}
	return randHashs
}

func SetLogLevel(level string) {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelTrace, true)))
	case "debug":
		log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelDebug, true)))
	case "info":
		log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelInfo, true)))
	case "warn":
		log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelWarn, true)))
	case "error":
		log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelError, true)))
	case "crit":
		log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelCrit, true)))
	}
}
