package eth_block_receipts

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/go-ethereum/rlp"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
)

type ReceiptResolver struct {
	decoder rlp.Decoder
}

func NewReceiptResolver(decoder rlp.Decoder) ReceiptResolver {
	return ReceiptResolver{decoder: decoder}
}

func (resolver ReceiptResolver) Resolve(block blocks.Block) (format.Node, error) {
	var receipt types.Receipt
	raw := block.RawData()
	err := resolver.decoder.Decode(raw, &receipt)
	if err != nil {
		return nil, err
	}
	node := &EthReceiptNode{
		Receipt: &receipt,
		raw:     raw,
		cid:     block.Cid(),
	}
	return node, nil
}
