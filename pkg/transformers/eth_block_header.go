package transformers

import (
	"log"

	"github.com/8thlight/block_watcher/pkg/db"
	"github.com/8thlight/block_watcher/pkg/ipfs"
)

type EthBlockHeaderTransformer struct {
	database  db.Database
	publisher ipfs.Publisher
}

func NewEthBlockHeaderTransformer(ethDB db.Database, publisher ipfs.Publisher) *EthBlockHeaderTransformer {
	return &EthBlockHeaderTransformer{database: ethDB, publisher: publisher}
}

func (t EthBlockHeaderTransformer) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		blockData, err := t.database.GetBlockHeaderByBlockNumber(i)
		if err != nil {
			return NewExecuteError(GetBlockRlpErr, err)
		}
		output, err := t.publisher.Write(blockData)
		if err != nil {
			return NewExecuteError(PutIpldErr, err)
		}
		log.Printf("Created IPLD: %s", output)
	}
	return nil
}
