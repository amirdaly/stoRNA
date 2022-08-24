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
	Nodes        []*Node
	Levels       map[int][]*Node
	hashStrategy func() hash.Hash
}

type Levels map[int][]*Node

type Node struct {
	Tree   *MerkleTree
	Parent []*Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   []byte
	C      Content
	Index  string
	Number int
}

func NewTreeGenesis(cs Content, length int) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{hashStrategy: defaultHashStrategy}
	var nodes []*Node
	//
	hash, err := cs.CalculateHash()
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, &Node{
		Hash: hash,
		C:    cs,
		leaf: true,
		Tree: t,
		// Index: integerToBinaryString(0, length),
		Index: "0",
	})
	t.Root = nodes[0]
	t.Nodes = nodes
	t.merkleRoot = hash
	emptyNode := &Node{Tree: t}
	level := make(map[int][]*Node)
	level[0] = append(level[0], emptyNode)
	level[1] = append(level[1], t.Nodes[0])
	t.Levels = level
	t.Nodes[0].Number = 1 //add navigation Number 1 to first node of tree
	return t, nil
}

func AddNodeToTree(cs Content, t *MerkleTree, depth int) (*Node, []*Node, error) {
	// // if there is nothing to add to tree
	// if len(cs) == 0 {
	// 	return t.Root, t.Nodes, nil
	// }

	// N := int(math.Pow(2, float64(depth)+1) - 1) // N = 2^n+1 − 1 for a tree of depth n

	// treeLastNodeCount := len(t.Nodes)           // Last node Count that tree had
	// treeNewNodesToAddCount := 1                 // count of new nodes that will be added to Tree
	treeDepthLength := depth // depth of Tree and Number Of edge Index Bytes

	// for N (Max Number of Merkle Tree Nodes) start to navigate whole tree
	// At first we will check the tree depth. if depth is growing up we must update Node Index strings
	// if depth of tree is not changing, we must just Add new nodes to their places
	//update nodes Index to new strings

	// is leaf or not
	// traversing Number
	// Index string

	if len(t.Nodes[0].Index) < depth {
		updateNodesIndex(t, treeDepthLength) //"00"
	}

	// add second node
	if len(t.Nodes) == 1 {
		hash, err := cs.CalculateHash()
		if err != nil {
			return nil, nil, err
		}
		x := integerToBinaryString(1, 2) // "01"
		newNode := &Node{
			Hash:   hash,
			C:      cs,
			leaf:   true,
			Tree:   t,
			Number: 2,
			Index:  x,
		}
		t.Nodes = append(t.Nodes, newNode)
		t.Levels[2] = append(t.Levels[2], newNode)
		t.Nodes[1].Parent = append(t.Nodes[1].Parent, newNode)
		return t.Root, t.Nodes, nil
	}

	// leafsCount := 1
	// navigate Tree Nodes to add new nodes (N is all of tree nodes count)
	i := len(t.Nodes)
	traversingNumber := i + 1

	beforeNode := t.Nodes[i-1]

	if beforeNode.leaf == true {
		if countLeafs(t)%2 == 1 { //added as new leaf
			hash, err := cs.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			x := integerToBinaryString(i, depth)
			newNode := &Node{
				Hash:   hash,
				C:      cs,
				leaf:   true,
				Tree:   t,
				Number: traversingNumber,
				Index:  x,
			}
			t.Nodes = append(t.Nodes, newNode)
			t.Nodes[i-1].Parent = append(t.Nodes[i-1].Parent, newNode)
			return t.Root, t.Nodes, nil
		} else { // added as parent node
			hash, err := cs.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			x := beforeNode.Index[:depth-1]
			newNode := &Node{
				Hash:   hash,
				C:      cs,
				leaf:   false,
				Tree:   t,
				Number: traversingNumber,
				Index:  x,
				Left:   t.Nodes[i-2],
				Right:  t.Nodes[i-1],
			}
			t.Nodes = append(t.Nodes, newNode)
			t.Levels[len(x)] = append(t.Levels[len(x)], newNode)
			t.Nodes[1].Parent = append(t.Nodes[1].Parent, newNode)
			return t.Root, t.Nodes, nil
		}
	}
	return nil, nil, nil
}

func countLeafs(t *MerkleTree) int {
	count := 0
	for _, n := range t.Nodes {
		if n.leaf == true {
			count++
		}
	}
	return count
}

func updateNodesIndex(t *MerkleTree, depth int) bool {
	T := false
	if (depth - len(t.Nodes[0].Index)) >= 1 {
		t.Levels = nil
		emptyNode := &Node{Tree: t}
		level := make(map[int][]*Node)
		level[0] = append(level[0], emptyNode)
		for _, i := range t.Nodes {
			tmpIndex := i.Index
			newIndex := strings.Repeat("0", depth-len(t.Nodes[0].Index)) + tmpIndex
			i.Index = newIndex
			l := len(newIndex)
			level[l] = append(level[l], i)
			T = true
		}
		t.Levels = level
		return T
	}
	return T
}

func AddNode(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return t.Root, t.Nodes, nil
	}

	var leafsCountHad int
	leafsCountHad = len(t.Nodes)
	leafsInCount := len(cs)
	newNodesCount := leafsCountHad + leafsInCount
	IndexLength := int(math.Round(math.Log2(float64(newNodesCount)) + 1))

	// node zero Index reimpliment
	Index := integerToBinaryString(0, IndexLength)
	t.Nodes[0].Index = Index
	t.Nodes[0].Parent = append(t.Nodes[0].Parent, t.Nodes[0])
	fmt.Println("added zero node", Index)
	//-----------------------------------

	if leafsCountHad >= 1 {
		for i := 1; i < leafsCountHad; i++ {
			iIndex := integerToBinaryString(i, IndexLength)
			t.Nodes[i].Index = iIndex
			parentsIndexString := exportParentsIndex(iIndex, IndexLength)

			for _, x := range parentsIndexString {
				if IsNodeInTree(x, t) != nil {
					fmt.Println("is not in tree ", i, " : ", x)
					var tmpParent *Node
					tmpParent = IsNodeInTree(x, t)
					t.Nodes[i].Parent = append(t.Nodes[i].Parent, tmpParent)
				} else {
					newNode := &Node{
						Tree:  t,
						Index: x,
					}
					t.Nodes = append(t.Nodes, newNode)
					t.Nodes[i].Parent = append(t.Nodes[i].Parent, newNode)
					fmt.Println("is in tree ", i, " : ", x)
				}
			}

		}

		in := leafsCountHad
		for _, c := range cs {
			var newNode *Node
			hash, err := c.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			iIndex := integerToBinaryString(in, IndexLength)
			newNode = &Node{
				Hash:  hash,
				C:     c,
				leaf:  true,
				Tree:  t,
				Index: iIndex,
			}
			parentsIndexString := exportParentsIndex(iIndex, IndexLength)
			for _, x := range parentsIndexString {
				if IsNodeInTree(x, t) != nil {
					var tmpParent *Node
					tmpParent = IsNodeInTree(x, t)
					newNode.Parent = append(newNode.Parent, tmpParent)
				} else {
					newNode2 := &Node{
						Tree:  t,
						Index: x,
					}
					t.Nodes = append(t.Nodes, newNode2)
					newNode.Parent = append(newNode.Parent, newNode2)
				}
			}

			t.Nodes = append(t.Nodes, newNode)
			in += 1
		}
	}

	return t.Nodes[0], t.Nodes, nil
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
		// nl[left].Parent = n
		// nl[right].Parent = n
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
	t.Nodes = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

func (n *Node) String() string {
	return fmt.Sprintf("Number: %d | Index: %s | leaf: %t | hash: %x data: %s", n.Number, n.Index, n.leaf, n.Hash, n.C)
}

func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Nodes {
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

func exportParentsIndex(Index string, length int) []string {
	var parentsIndexString []string
	if len(Index) == length && !strings.Contains(Index, "1") {
		parentsIndexString = append(parentsIndexString, Index)
		return parentsIndexString
	} else if len(Index) == length && strings.Contains(Index, "1") {
		for i := length - 1; i >= 0; i-- {
			newstr := Index
			if Index[i] == '1' {
				newstr = Index[:i] + string('0')
				parentsIndexString = append(parentsIndexString, newstr)
			}
		}
	}
	return parentsIndexString
}

func checkParent(node *Node, length int) []*Node {
	var parents []*Node
	if len(node.Index) == length && !strings.Contains(node.Index, "1") {
		parents = append(parents, node)
		return parents
	} else if len(node.Index) == length && strings.Contains(node.Index, "1") {
		for i := len(node.Index) - 1; i >= 0; i-- {
			str := node.Index
			var parentsStrings []string
			if str[i] == '1' {
				newstr := str[:i] + string('0')
				parentsStrings = append(parentsStrings, newstr)
			}
		}
	}
	return parents
}

func IsNodeInTree(Index string, t *MerkleTree) *Node {
	for _, i := range t.Nodes {
		if i.Index == Index {
			return i
		}
	}
	return nil
}
