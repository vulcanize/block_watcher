package level

type Database struct {
	reader
}

func NewLevelDatabase(ldbReader reader) *Database {
	return &Database{
		reader: ldbReader,
	}
}

func (l Database) GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return l.reader.GetBodyRLP(h, n), nil
}

func (l Database) GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return l.reader.GetHeaderRLP(h, n), nil
}
