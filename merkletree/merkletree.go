package merkletree

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"math"
	"strconv"
	"strings"
)

type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
}

type MerkleTree struct {
	Root         *Node
	merkleRoot   []byte
	Leafs        []*Node
	hashStrategy func() hash.Hash
}

type Node struct {
	Tree   *MerkleTree
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   []byte
	C      Content
	index  string
}

func AddNode(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return t.Root, t.Leafs, nil
	}

	var nodeCount int
	nodeCount = len(t.Leafs)
	inCount := len(cs)
	leafCount := nodeCount + inCount
	indexLength := int(math.Round(math.Log2(float64(leafCount) + 1)))

	if t.Leafs[0].index == "0" { // just have one node
		index := integerToBinaryString(0, indexLength)
		t.Leafs[0].index = index
	}

	if nodeCount >= 1 {
		for i := 0; i < nodeCount; i++ {

			iindex := integerToBinaryString(i, indexLength)
			t.Leafs[i].index = iindex
		}
		in := nodeCount
		for _, c := range cs {
			hash, err := c.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			iindex := integerToBinaryString(in, indexLength)
			t.Leafs = append(t.Leafs, &Node{
				Hash:  hash,
				C:     c,
				leaf:  true,
				Tree:  t,
				index: iindex,
			})
			in += 1
		}
	}

	return t.Leafs[0], t.Leafs, nil
}

func buildWithContent(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("Error: Cannot onstruct new tree with no content")
	}

	var leafs []*Node
	for _, c := range cs {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Hash: hash,
			C:    c,
			leaf: true,
			Tree: t,
		})
	}

	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			C:    leafs[len(leafs)-1].C,
			leaf: true,
			dup:  true,
			Tree: t,
		}
		leafs = append(leafs, duplicate)
	}

	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

func buildIntermediate(nl []*Node, t *MerkleTree) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := t.hashStrategy()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		chash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
			Tree:  t,
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes, t)
}

func NewTree(cs []Content) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{hashStrategy: defaultHashStrategy}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

func NewTreeGenesis(cs []Content, length int) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{hashStrategy: defaultHashStrategy}
	var leafs []*Node

	hash, err := cs[0].CalculateHash()
	if err != nil {
		return nil, err
	}
	leafs = append(leafs, &Node{
		Hash: hash,
		C:    cs[0],
		leaf: true,
		Tree: t,
		// index: integerToBinaryString(0, length),
		index: "0",
	})
	t.Root = leafs[0]
	t.Leafs = leafs
	t.merkleRoot = hash
	return t, nil
}

func (n *Node) String() string {
	return fmt.Sprintf("index: %s | leaf: %t | hash: %x data: %s", n.index, n.leaf, n.Hash, n.C)
}

func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}

func integerToBinaryString(num int, length int) string {
	binaryString := strconv.FormatInt(int64(num), 2)
	if len(binaryString) >= length {
		return binaryString
	}
	return strings.Repeat("0", length-len(binaryString)) + binaryString
}
