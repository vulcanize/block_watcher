package db

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/vulcanize/block_watcher/pkg/db/level"
)

var ErrNoSuchDb = errors.New("no such database")

type ReadError struct {
	msg string
	err error
}

func (re ReadError) Error() string {
	return fmt.Sprintf("%s: %s", re.msg, re.err.Error())
}

type Database interface {
	GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error)
	GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error)
}

func CreateDatabase(config DatabaseConfig) (Database, error) {
	switch config.Type {
	case Level:
		levelDBConnection, err := ethdb.NewLDBDatabase(config.Path, 128, 1024)
		if err != nil {
			return nil, ReadError{msg: "Failed to connect to LevelDB", err: err}
		}
		levelDBReader := level.NewLevelDatabaseReader(levelDBConnection)
		levelDB := level.NewLevelDatabase(levelDBReader)
		return levelDB, nil
	default:
		return nil, ReadError{msg: "Unknown database not implemented", err: ErrNoSuchDb}
	}
}