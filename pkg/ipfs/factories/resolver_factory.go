package factories

import (
	"fmt"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/go-ethereum/rlp"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type ResolverFactory struct{}

func (factory ResolverFactory) GetResolver(decodedCid *cid.Cid) (ipfs.Resolver, error) {
	rlpDecoder := rlp.RlpDecoder{}
	switch decodedCid.Prefix().Codec {
	case cid.EthBlock:
		return eth_block_header.NewBlockHeaderResolver(rlpDecoder), nil
	case cid.EthTxReceipt:
		return eth_block_receipts.NewReceiptResolver(rlpDecoder), nil
	case cid.EthTx:
		return eth_block_transactions.NewTransactionResolver(rlpDecoder), nil
	default:
		return nil, fmt.Errorf("resolver not found for codec: %d\n", decodedCid.Prefix().Codec)
	}
}
