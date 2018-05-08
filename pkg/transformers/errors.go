package transformers

import "fmt"

const (
	GetBlockRlpErr = "Error fetching block RLP data"
	PutIpldErr     = "Error writing to IPFS"
)

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
