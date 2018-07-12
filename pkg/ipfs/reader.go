package ipfs

import (
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type Reader interface {
	Read(cid *cid.Cid) (format.Node, error)
}

type IpldReader struct {
	getter   Getter
	resolver Resolver
}

func NewIpldReader(getter Getter, resolver Resolver) *IpldReader {
	return &IpldReader{
		getter:   getter,
		resolver: resolver,
	}
}

func (reader *IpldReader) Read(cid *cid.Cid) (format.Node, error) {
	block, err := reader.getter.Get(cid)
	if err != nil {
		return nil, err
	}
	return reader.resolver.Resolve(block)
}
