package db

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	DefaultLevelDBPathFormat = "./db/meta-%s.db"
)

var (
	PendingSeqKey = []byte("pendingSeq")
)

type MetaDB interface {
	SetPendingSeq(seq int64) error
	GetPendingSeq() (int64, error)
	CloseDB()
}

type MetaDBImpl struct {
	*leveldb.DB
}

func NewDB(path string) (MetaDB, error) {
	// open or create level db
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Error("open or create level db failed", "err", err, "path", path)
		return nil, err
	}
	return &MetaDBImpl{db}, nil
}

func (db *MetaDBImpl) SetPendingSeq(seq int64) error {
	return db.Put(PendingSeqKey, []byte(fmt.Sprintf("%d", seq)), nil)
}

func (db *MetaDBImpl) GetPendingSeq() (int64, error) {
	seqByte, err := db.Get(PendingSeqKey, nil)
	if err != nil || seqByte == nil {
		return 0, err
	}

	newSeq, err := strconv.ParseInt(string(seqByte), 10, 64)
	if err != nil {
		return 0, err
	}
	return newSeq, nil
}

func (db *MetaDBImpl) CloseDB() {
	db.Close()
}
