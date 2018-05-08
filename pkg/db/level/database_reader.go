package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rlp"
)

// Wraps go-ethereum db operations
type reader interface {
	GetCanonicalHash(number uint64) common.Hash
	GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
	GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue
}

type LDBReader struct {
	core.DatabaseReader
}

func NewLevelDatabaseReader(reader core.DatabaseReader) *LDBReader {
	return &LDBReader{DatabaseReader: reader}
}

func (ldbr *LDBReader) GetCanonicalHash(number uint64) common.Hash {
	return core.GetCanonicalHash(ldbr.DatabaseReader, number)
}

func (ldbr *LDBReader) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	return core.GetHeaderRLP(ldbr.DatabaseReader, hash, number)
}

func (ldbr *LDBReader) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	return core.GetBodyRLP(ldbr.DatabaseReader, hash, number)
}
