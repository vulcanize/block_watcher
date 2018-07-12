package eth_block_header

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/go-ethereum/rlp"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
)

type BlockHeaderResolver struct {
	decoder rlp.Decoder
}

func NewBlockHeaderResolver(decoder rlp.Decoder) BlockHeaderResolver {
	return BlockHeaderResolver{decoder: decoder}
}

func (resolver BlockHeaderResolver) Resolve(block blocks.Block) (format.Node, error) {
	raw := block.RawData()
	var header types.Header
	err := resolver.decoder.Decode(raw, &header)
	if err != nil {
		return nil, err
	}
	headerNode := &EthBlockHeaderNode{
		Header:  &header,
		cid:     block.Cid(),
		rawdata: raw,
	}
	return headerNode, nil
}
