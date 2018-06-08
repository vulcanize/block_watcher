package test_helpers

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	FakeError = errors.New("failed")
	FakeHash  = common.HexToHash("0x123")
)
