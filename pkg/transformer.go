package pkg

import (
	"fmt"
	"log"

	"github.com/8thlight/block_watcher/pkg/db"
	"github.com/8thlight/block_watcher/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

const (
	GetVulcanizeBlockErr = "Error fetching vulcanize block"
	GetBlockRlpErr       = "Error fetching block RLP data"
	PutIpldErr           = "Error writing to IPFS"
)

type Transformer struct {
	db.Database
	datastore.BlockRepository
	ipfs.Publisher
}

type ExecuteError struct {
	msg string
	err error
}

func NewExecuteError(msg string, err error) *ExecuteError {
	return &ExecuteError{msg: msg, err: err}
}

func (ee ExecuteError) Error() string {
	return fmt.Sprintf("%s: %s", ee.msg, ee.err.Error())
}

func (t Transformer) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		block, err := t.BlockRepository.GetBlock(i)
		if err != nil {
			return NewExecuteError(GetVulcanizeBlockErr, err)
		}
		blockData, err := t.Database.Get(block)
		if err != nil {
			return NewExecuteError(GetBlockRlpErr, err)
		}
		output, err := t.Publisher.Write(blockData)
		if err != nil {
			return NewExecuteError(PutIpldErr, err)
		}
		log.Printf("Created IPLD: %s", output)
	}
	return nil
}

func NewTransformer(blockRepository datastore.BlockRepository, ethDB db.Database, publisher ipfs.Publisher) *Transformer {
	return &Transformer{BlockRepository: blockRepository, Database: ethDB, Publisher: publisher}
}
