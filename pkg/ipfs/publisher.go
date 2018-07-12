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
	Write(input interface{}) ([]string, error)
}

type IpldPublisher struct {
	DagPutter
}

func NewIpldPublisher(dagPutter DagPutter) *IpldPublisher {
	return &IpldPublisher{DagPutter: dagPutter}
}

func (ip *IpldPublisher) Write(input interface{}) ([]string, error) {
	cids, err := ip.DagPutter.DagPut(input)
	if err != nil {
		return nil, err
	}
	return cids, nil
}
