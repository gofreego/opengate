package utils

import (
	"strings"
)

// TrieNode represents a node in the prefix trie
type TrieNode[k any] struct {
	children map[string]*TrieNode[k]
	route    k
	isEnd    bool
}

// Trie represents a prefix tree for efficient route matching
type Trie[k any] struct {
	root *TrieNode[k]
}

// NewTrie creates a new trie
func NewTrie[k any]() *Trie[k] {
	return &Trie[k]{
		root: &TrieNode[k]{
			children: make(map[string]*TrieNode[k]),
		},
	}
}

// Insert adds a route to the trie
func (t *Trie[k]) Insert(pathPrefix string, route k) {
	node := t.root
	parts := strings.Split(strings.Trim(pathPrefix, "/"), "/")

	for _, part := range parts {
		if part == "" {
			continue
		}
		if node.children[part] == nil {
			node.children[part] = &TrieNode[k]{
				children: make(map[string]*TrieNode[k]),
			}
		}
		node = node.children[part]
	}

	node.route = route
	node.isEnd = true
}

// FindLongestMatch finds the route with the longest matching prefix
func (t *Trie[k]) FindLongestMatch(path string) k {
	node := t.root
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var lastMatch k

	// Check root for exact "/" match
	if node.isEnd {
		lastMatch = node.route
	}

	for _, part := range parts {
		if part == "" {
			continue
		}
		if node.children[part] == nil {
			break
		}
		node = node.children[part]
		if node.isEnd {
			lastMatch = node.route
		}
	}

	return lastMatch
}
