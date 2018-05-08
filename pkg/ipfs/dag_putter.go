package ipfs

type DagPutter interface {
	DagPut(raw []byte) ([]string, error)
}
