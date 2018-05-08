package db

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
)

type Decoder interface {
	Decode(raw []byte, out interface{}) (interface{}, error)
}

type RlpDecoder struct{}

func (RlpDecoder) Decode(raw []byte, out interface{}) (interface{}, error) {
	err := rlp.Decode(bytes.NewBuffer(raw), out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
