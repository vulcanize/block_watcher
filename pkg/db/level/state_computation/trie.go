package state_computation

import "github.com/ethereum/go-ethereum/core/state"

type ITrie interface {
	NodeIterator(startKey []byte) IIterator
}

type Trie struct {
	trie state.Trie
}

func NewTrie(trie state.Trie) ITrie {
	return &Trie{trie: trie}
}

func (t *Trie) NodeIterator(startKey []byte) IIterator {
	iterator := t.trie.NodeIterator(startKey)
	return NewStateIterator(iterator)
}
