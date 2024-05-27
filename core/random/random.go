package random

import (
	"crypto/rand"
	"encoding/binary"
	"sync"

	"github.com/ethereum/go-ethereum/log"
)

var byteArrayPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 8)
	},
}

func GenerateRandomNum() (uint64, error) {
	var randNum uint64
	for randNum == 0 {
		buf := byteArrayPool.Get().([]byte)
		defer byteArrayPool.Put(buf)

		_, err := rand.Read(buf)
		if err != nil {
			log.Error("crypto/rand generate random number:", "err", err)
			return 0, err
		}
		randNum = binary.BigEndian.Uint64(buf)
	}
	return randNum, nil
}
