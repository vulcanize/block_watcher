package state_computation

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

type Validator interface {
	ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error
}

type StateValidator struct {
	validator core.Validator
}

func NewStateValidator(blockChain StateBlockChain) *StateValidator {
	validator := core.NewBlockValidator(blockChain.Config(), blockChain.BlockChain(), blockChain.Engine())
	return &StateValidator{validator: validator}
}

func (sv *StateValidator) ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	return sv.validator.ValidateState(block, parent, state, receipts, usedGas)
}
