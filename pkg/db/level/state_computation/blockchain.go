package state_computation

import (
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

type BlockChain interface {
	BlockChain() *core.BlockChain
	Config() *params.ChainConfig
	Engine() consensus.Engine
}

type StateBlockChain struct {
	blockChain *core.BlockChain
}

func NewStateBlockChain(databaseConnection ethdb.Database) (*StateBlockChain, error) {
	blockchain, err := core.NewBlockChain(databaseConnection, nil, params.MainnetChainConfig, ethash.NewFaker(), vm.Config{})
	if err != nil {
		return nil, err
	}
	return &StateBlockChain{blockChain: blockchain}, nil
}

func (sb *StateBlockChain) BlockChain() *core.BlockChain {
	return sb.blockChain
}

func (sb *StateBlockChain) Config() *params.ChainConfig {
	return sb.blockChain.Config()
}

func (sb *StateBlockChain) Engine() consensus.Engine {
	return sb.blockChain.Engine()
}
