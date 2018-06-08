package eth_block_transactions

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/eth-block-extractor/pkg/db"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/util"
)

const (
	EthBlockTransactionCode = 0x93
)

type BlockTransactionsDagPutter struct {
	adder   ipfs.Adder
	decoder db.Decoder
}

func NewBlockTransactionsDagPutter(adder ipfs.Adder, decoder db.Decoder) *BlockTransactionsDagPutter {
	return &BlockTransactionsDagPutter{adder: adder, decoder: decoder}
}

func (bbdp *BlockTransactionsDagPutter) DagPut(raw []byte) ([]string, error) {
	var blockBody types.Body
	out, err := bbdp.decoder.Decode(raw, &blockBody)
	if err != nil {
		return nil, err
	}
	body := out.(*types.Body)
	transactions := body.Transactions
	var cids []string
	for _, transaction := range transactions {
		buffer := new(bytes.Buffer)
		err = transaction.EncodeRLP(buffer)
		if err != nil {
			return nil, err
		}
		transactionCid, err := util.RawToCid(EthBlockTransactionCode, buffer.Bytes())
		if err != nil {
			return nil, err
		}
		transactionNode := &EthTransactionNode{
			Transaction: transaction,
			cid:         transactionCid,
			rawdata:     buffer.Bytes(),
		}
		err = bbdp.adder.Add(transactionNode)
		if err != nil {
			return nil, err
		}
		cids = append(cids, transactionCid.String())
	}
	return cids, nil
}
