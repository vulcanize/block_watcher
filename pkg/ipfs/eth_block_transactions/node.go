package eth_block_transactions

import (
	"github.com/ethereum/go-ethereum/core/types"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

type EthTransactionNode struct {
	*types.Transaction

	cid     *cid.Cid
	rawdata []byte
}

func (etn *EthTransactionNode) RawData() []byte {
	return etn.rawdata
}

func (etn *EthTransactionNode) Cid() *cid.Cid {
	return etn.cid
}

func (EthTransactionNode) String() string {
	return ""
}

func (EthTransactionNode) Loggable() map[string]interface{} {
	panic("implement me")
}

func (EthTransactionNode) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (EthTransactionNode) Tree(path string, depth int) []string {
	panic("implement me")
}

func (EthTransactionNode) ResolveLink(path []string) (*ipld.Link, []string, error) {
	panic("implement me")
}

func (EthTransactionNode) Copy() ipld.Node {
	panic("implement me")
}

func (EthTransactionNode) Links() []*ipld.Link {
	panic("implement me")
}

func (EthTransactionNode) Stat() (*ipld.NodeStat, error) {
	panic("implement me")
}

func (EthTransactionNode) Size() (uint64, error) {
	panic("implement me")
}
