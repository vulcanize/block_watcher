package db

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/vulcanize/block_watcher/pkg/db/level"
	"github.com/vulcanize/block_watcher/pkg/db/level/state_computation"
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
	ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error)
	GetBlockByBlockNumber(blockNumber int64) *types.Block
	GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error)
	GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error)
	GetStateTrieNodes(root []byte) ([][]byte, error)
}

func CreateDatabase(config DatabaseConfig) (Database, error) {
	switch config.Type {
	case Level:
		levelDBConnection, err := ethdb.NewLDBDatabase(config.Path, 128, 1024)
		if err != nil {
			return nil, ReadError{msg: "Failed to connect to LevelDB", err: err}
		}
		levelDBReader := level.NewLevelDatabaseReader(levelDBConnection)
		stateComputer, err := createStateComputer(levelDBConnection)
		if err != nil {
			return nil, err
		}
		levelDB := level.NewLevelDatabase(levelDBReader, stateComputer)
		return levelDB, nil
	default:
		return nil, ReadError{msg: "Unknown database not implemented", err: ErrNoSuchDb}
	}
}

func createStateComputer(databaseConnection ethdb.Database) (state_computation.Computer, error) {
	blockChain, err := state_computation.NewStateBlockChain(databaseConnection)
	if err != nil {
		return nil, err
	}
	db := state_computation.NewDatabase(databaseConnection)
	processor := state_computation.NewStateProcessor(*blockChain)
	trieFactory := state_computation.NewStateTrieFactory()
	validator := state_computation.NewStateValidator(*blockChain)
	computer := state_computation.NewStateComputer(blockChain, db, processor, trieFactory, validator)
	return computer, nil
}
