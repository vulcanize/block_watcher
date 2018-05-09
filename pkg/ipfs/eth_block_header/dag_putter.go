package eth_block_header

import (
	"github.com/ethereum/go-ethereum/core/types"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"

	"github.com/vulcanize/block_watcher/pkg/db"
	"github.com/vulcanize/block_watcher/pkg/ipfs"
	"github.com/vulcanize/block_watcher/pkg/ipfs/util"
)

const (
	EthBlockHeaderCode = 0x90
)

type BlockHeaderDagPutter struct {
	adder   ipfs.Adder
	decoder db.Decoder
}

func NewBlockHeaderDagPutter(adder ipfs.Adder, decoder db.Decoder) *BlockHeaderDagPutter {
	return &BlockHeaderDagPutter{adder: adder, decoder: decoder}
}

func (bhdp *BlockHeaderDagPutter) DagPut(raw []byte) ([]string, error) {
	nd, err := bhdp.getNodeForBlockHeader(raw)
	if err != nil {
		return nil, err
	}
	err = bhdp.adder.Add(nd)
	if err != nil {
		return nil, err
	}
	return []string{nd.Cid().String()}, nil
}

func (bhdp *BlockHeaderDagPutter) getNodeForBlockHeader(raw []byte) (ipld.Node, error) {
	var blockHeader types.Header
	out, err := bhdp.decoder.Decode(raw, &blockHeader)
	if err != nil {
		return nil, err
	}
	blockHeaderCid, err := util.RawToCid(EthBlockHeaderCode, raw)
	if err != nil {
		return nil, err
	}
	header := out.(*types.Header)
	return &EthBlockHeaderNode{
		Header:  header,
		cid:     blockHeaderCid,
		rawdata: raw,
	}, nil
}
