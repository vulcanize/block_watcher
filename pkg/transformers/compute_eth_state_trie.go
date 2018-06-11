package transformers

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
	"log"
)

const (
	GenesisBlockNumber  = int64(0)
	FirstBlockToCompute = int64(1)
)

type ComputeEthStateTrieTransformer struct {
	database             db.Database
	decoder              rlp.Decoder
	stateTriePublisher   ipfs.Publisher
	storageTriePublisher ipfs.Publisher
}

func NewComputeEthStateTrieTransformer(database db.Database, decoder rlp.Decoder, stateTriePublisher, storageTriePublisher ipfs.Publisher) *ComputeEthStateTrieTransformer {
	return &ComputeEthStateTrieTransformer{
		database:             database,
		decoder:              decoder,
		stateTriePublisher:   stateTriePublisher,
		storageTriePublisher: storageTriePublisher,
	}
}

func (t ComputeEthStateTrieTransformer) Execute(endingBlockNumber int64) error {
	root, err := t.getStateRootForBlock(GenesisBlockNumber)
	if err != nil {
		return err
	}
	// ignore storage trie node return val for genesis block
	stateTrieNodes, _, err := t.database.GetStateAndStorageTrieNodes(root)
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
		stateRoot, err := t.database.ComputeBlockStateTrie(currentBlock, parentBlock)
		if err != nil {
			return err
		}
		nextStateTrieNodes, nextStorageTrieNodes, err := t.database.GetStateAndStorageTrieNodes(stateRoot)
		if err != nil {
			return err
		}
		err = t.writeStateTrieNodesToIpfs(nextStateTrieNodes)
		if err != nil {
			return err
		}
		err = t.writeStorageTrieNodesToIpfs(nextStorageTrieNodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t ComputeEthStateTrieTransformer) getStateRootForBlock(blockNumber int64) (root common.Hash, err error) {
	var header types.Header
	rawHeader, err := t.database.GetBlockHeaderByBlockNumber(blockNumber)
	if err != nil {
		return root, fmt.Errorf("Error fetching header for block %d: %s\n", blockNumber, err)
	}
	err = t.decoder.Decode(rawHeader, &header)
	if err != nil {
		return root, err
	}
	return header.Root, err
}

func (t ComputeEthStateTrieTransformer) writeStateTrieNodesToIpfs(stateTrieNodes [][]byte) error {
	for _, node := range stateTrieNodes {
		output, err := t.stateTriePublisher.Write(node)
		if err != nil {
			return fmt.Errorf("Error writing state trie node to ipfs: %s\n", err)
		}
		log.Println("Created ipld: ", output)
	}
	return nil
}

func (t ComputeEthStateTrieTransformer) writeStorageTrieNodesToIpfs(storageTrieNodes [][]byte) error {
	for _, node := range storageTrieNodes {
		output, err := t.storageTriePublisher.Write(node)
		if err != nil {
			return fmt.Errorf("Error writing storage trie node to ipfs: %s\n", err.Error())
		}
		log.Println("Created ipld: ", output)
	}
	return nil
}
