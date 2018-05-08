package eth_block_header

import (
	"github.com/ethereum/go-ethereum/core/types"

	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

type EthBlockHeaderNode struct {
	*types.Header

	cid     *cid.Cid
	rawdata []byte
}

func (ebh *EthBlockHeaderNode) RawData() []byte {
	return ebh.rawdata
}

func (ebh *EthBlockHeaderNode) Cid() *cid.Cid {
	return ebh.cid
}

func (EthBlockHeaderNode) String() string {
	return ""
}

func (EthBlockHeaderNode) Loggable() map[string]interface{} {
	panic("implement me")
}

func (EthBlockHeaderNode) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (EthBlockHeaderNode) Tree(path string, depth int) []string {
	panic("implement me")
}

func (EthBlockHeaderNode) ResolveLink(path []string) (*ipld.Link, []string, error) {
	panic("implement me")
}

func (EthBlockHeaderNode) Copy() ipld.Node {
	panic("implement me")
}

func (EthBlockHeaderNode) Links() []*ipld.Link {
	panic("implement me")
}

func (EthBlockHeaderNode) Stat() (*ipld.NodeStat, error) {
	panic("implement me")
}

func (EthBlockHeaderNode) Size() (uint64, error) {
	panic("implement me")
}
