package transformers

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"log"
)

const (
	GenesisBlockNumber  = int64(0)
	FirstBlockToCompute = int64(1)
)

type ComputeEthStateTrieTransformer struct {
	database  db.Database
	decoder   db.Decoder
	publisher ipfs.Publisher
}

func NewComputeEthStateTrieTransformer(database db.Database, decoder db.Decoder, publisher ipfs.Publisher) *ComputeEthStateTrieTransformer {
	return &ComputeEthStateTrieTransformer{database: database, decoder: decoder, publisher: publisher}
}

func (t ComputeEthStateTrieTransformer) Execute(endingBlockNumber int64) error {
	root, err := t.getStateRootForBlock(GenesisBlockNumber)
	if err != nil {
		return err
	}
	stateTrieNodes, err := t.database.GetStateTrieNodes(root.Bytes())
	if err != nil {
		return fmt.Errorf("Error fetching state trie for genesis block: %s\n", err)
	}
	err = t.writeStateTrieNodesToIpfs(stateTrieNodes)
	if err != nil {
		return err
	}

	for n := FirstBlockToCompute; n <= endingBlockNumber; n++ {
		currentBlock := t.database.GetBlockByBlockNumber(n)
		parentBlock := t.database.GetBlockByBlockNumber(n - 1)
		nextStateTrieNodes, err := t.database.ComputeBlockStateTrie(currentBlock, parentBlock)
		if err != nil {
			return err
		}
		err = t.writeStateTrieNodesToIpfs(nextStateTrieNodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t ComputeEthStateTrieTransformer) getStateRootForBlock(blockNumber int64) (common.Hash, error) {
	var header types.Header
	rawHeader, err := t.database.GetBlockHeaderByBlockNumber(blockNumber)
	if err != nil {
		return common.Hash{}, fmt.Errorf("Error fetching header for block %d: %s\n", blockNumber, err)
	}
	out, err := t.decoder.Decode(rawHeader, &header)
	parsedHeader := out.(*types.Header)
	return parsedHeader.Root, nil
}

func (t ComputeEthStateTrieTransformer) writeStateTrieNodesToIpfs(stateTrieNodes [][]byte) error {
	for _, node := range stateTrieNodes {
		output, err := t.publisher.Write(node)
		if err != nil {
			return fmt.Errorf("Error writing state trie node to ipfs: %s\n", err)
		}
		log.Println("Created ipld: ", output)
	}
	return nil
}
