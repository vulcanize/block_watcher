package state

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/trie"
)

type GethStateTrie interface {
	NodeIterator(startKey []byte) trie.GethTrieNodeIterator
}

type Trie struct {
	trie state.Trie
}

func NewTrie(trie state.Trie) GethStateTrie {
	return &Trie{trie: trie}
}

func (t *Trie) NodeIterator(startKey []byte) trie.GethTrieNodeIterator {
	iterator := t.trie.NodeIterator(startKey)
	return trie.NewNodeIterator(iterator)
}
