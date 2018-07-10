package transformers

import (
	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"log"
)

type EthBlockTransactionsTransformer struct {
	database  db.Database
	publisher ipfs.Publisher
}

func NewEthBlockTransactionsTransformer(db db.Database, publisher ipfs.Publisher) *EthBlockTransactionsTransformer {
	return &EthBlockTransactionsTransformer{database: db, publisher: publisher}
}

func (t EthBlockTransactionsTransformer) Execute(startingBlockNumber int64, endingBlockNumber int64) error {
	if endingBlockNumber < startingBlockNumber {
		return ErrInvalidRange
	}
	for i := startingBlockNumber; i <= endingBlockNumber; i++ {
		blockData, err := t.database.GetBlockBodyByBlockNumber(i)
		if err != nil {
			return NewExecuteError(GetBlockRlpErr, err)
		}
		res, err := t.publisher.Write(blockData)
		if err != nil {
			return NewExecuteError(PutIpldErr, err)
		}
		log.Println("Created CIDs: ", res)
	}
	return nil
}
