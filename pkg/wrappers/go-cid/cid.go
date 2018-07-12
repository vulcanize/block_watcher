package go_cid

import "gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"

type ICidDecoder interface {
	Decode(v string) (*cid.Cid, error)
}

type CidDecoder struct{}

func NewCidDecoder() CidDecoder {
	return CidDecoder{}
}

func (CidDecoder) Decode(v string) (*cid.Cid, error) {
	return cid.Decode(v)
}
