package ipfs

import (
	"context"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type IPFS struct {
	n   *core.IpfsNode
	ctx context.Context
}

func InitIPFSNode(repoPath string) (*IPFS, error) {
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	cfg := &core.BuildCfg{
		Online: false,
		Repo:   r,
	}
	ipfsNode, err := core.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &IPFS{n: ipfsNode, ctx: ctx}, nil
}

func (ipfs IPFS) Add(node format.Node) error {
	return ipfs.n.DAG.Add(ipfs.n.Context(), node)
}

func (ipfs IPFS) Get(cid *cid.Cid) (blocks.Block, error) {
	return ipfs.n.Blocks.GetBlock(context.Background(), cid)
}
