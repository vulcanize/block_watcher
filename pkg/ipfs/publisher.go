package ipfs

import "fmt"

type Error struct {
	msg string
	err error
}

func (ie Error) Error() string {
	return fmt.Sprintf("%s: %s", ie.msg, ie.err.Error())
}

type Publisher interface {
	Write(blockData []byte) ([]string, error)
}

type BlockDataPublisher struct {
	DagPutter
}

func NewIpfsPublisher(dagPutter DagPutter) *BlockDataPublisher {
	return &BlockDataPublisher{DagPutter: dagPutter}
}

func (ip *BlockDataPublisher) Write(blockData []byte) ([]string, error) {
	cids, err := ip.DagPutter.DagPut(blockData)
	if err != nil {
		return nil, err
	}
	return cids, nil
}
