package eth_block_transactions

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/go-ethereum/rlp"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
)

type TransactionResolver struct {
	decoder rlp.Decoder
}

func NewTransactionResolver(decoder rlp.Decoder) TransactionResolver {
	return TransactionResolver{decoder: decoder}
}

func (resolver TransactionResolver) Resolve(block blocks.Block) (format.Node, error) {
	var transaction types.Transaction
	raw := block.RawData()
	err := resolver.decoder.Decode(raw, &transaction)
	if err != nil {
		return nil, err
	}
	node := &EthTransactionNode{
		Transaction: &transaction,
		cid:         block.Cid(),
		rawdata:     raw,
	}
	return node, nil
}
