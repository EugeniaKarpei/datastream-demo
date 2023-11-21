package processor

import "bytes"

func NewTrieNode() *TrieNode {
	return &TrieNode{
		chars: make(map[byte]*TrieNode),
	}
}

type TrieNode struct {
	chars  map[byte]*TrieNode
	isLeaf bool
}

func (node *TrieNode) AddWord(word string) {
	if len(word) == 0 {
		node.isLeaf = true
		return
	}

	char := word[0]
	nextNode, found := node.chars[char]
	if !found {
		nextNode = NewTrieNode()
	}
	node.chars[char] = nextNode
	nextNode.AddWord(word[1:])
}

func (node *TrieNode) GetWordsInOrder(searchTerm string) []string {
	return traverse(node, "", searchTerm)
}

func traverse(node *TrieNode, partialWord string, path string) []string {
	if len(path) == 0 {
		// return all subtrie nodes
		words := []string{}
		listWords(node, partialWord, &words)
		return words
	}
	char := path[0]
	nextNode, found := node.chars[char]
	if !found {
		return []string{}
	}
	newWord := bytes.NewBufferString(partialWord)
	newWord.WriteByte(char)
	return traverse(nextNode, newWord.String(), path[1:])
}

func listWords(node *TrieNode, partialWord string, words *[]string) {
	if node.isLeaf {
		*words = append(*words, partialWord)
	}

	for char, nextNode := range node.chars {
		newWord := bytes.NewBufferString(partialWord)
		newWord.WriteByte(char)
		listWords(nextNode, newWord.String(), words)
	}
}
