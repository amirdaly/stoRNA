package merkletree

import (
	"crypto/sha256"
	"fmt"
	"hash"
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
	done   bool
	verify []byte
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
		done:  true,
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

func AddNodeToTree(cs Content, t *MerkleTree) (*Node, []*Node, error) {
	// if there is nothing to add to tree
	// if len(cs) == 0 {
	// 	return t.Root, t.Nodes, nil
	// }

	// N := int(math.Pow(2, float64(depth)+1) - 1) // N = 2^n+1 − 1 for a tree of depth n

	// treeLastNodeCount := len(t.Nodes)           // Last node Count that tree had
	// treeNewNodesToAddCount := 1                 // count of new nodes that will be added to Tree
	//treeDepthLength := depth // depth of Tree and Number Of edge Index Bytes

	// for N (Max Number of Merkle Tree Nodes) start to navigate whole tree
	// At first we will check the tree depth. if depth is growing up we must update Node Index strings
	// if depth of tree is not changing, we must just Add new nodes to their places
	//update nodes Index to new strings

	// is leaf or not
	// traversing Number
	// Index string

	var depth int

	// add second node
	if len(t.Nodes) == 1 {
		depth = len(t.Nodes) + 1
		updateNodesIndex(t, depth)

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
			done:   true,
		}
		t.Nodes = append(t.Nodes, newNode)
		t.Levels[2] = append(t.Levels[2], newNode)
		t.Nodes[1].Parent = append(t.Nodes[1].Parent, t.Nodes[0])
		return t.Root, t.Nodes, nil
	}

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
			x := integerToBinaryString(countLeafs(t), depth)
			newNode := &Node{
				Hash:   hash,
				C:      cs,
				leaf:   true,
				Tree:   t,
				Number: traversingNumber,
				Index:  x,
				done:   true,
			}
			t.Nodes = append(t.Nodes, newNode)
			t.Levels[depth] = append(t.Levels[depth], newNode)
			setParentsToNode(newNode, t)
			return t.Root, t.Nodes, nil
		} else { // added as parent node
			hash, err := cs.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			x := beforeNode.Index[:len(beforeNode.Index)-1]
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
	} else if beforeNode.leaf == false {
		beforeNodeLevel := len(beforeNode.Index)
		if len(t.Levels[beforeNodeLevel])%2 != 0 { // before node is a parent up to the leafs. so new node will be leaf
			hash, err := cs.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			x := integerToBinaryString(countLeafs(t), depth)
			newNode := &Node{
				Hash:   hash,
				C:      cs,
				leaf:   true,
				Tree:   t,
				Number: traversingNumber,
				Index:  x,
				Left:   t.Nodes[i-1],
			}
			t.Nodes = append(t.Nodes, newNode)
			t.Levels[len(x)] = append(t.Levels[len(x)], newNode)
			setParentsToNode(newNode, t)
			return t.Root, t.Nodes, nil
		} else if len(t.Levels[beforeNodeLevel])%2 == 0 { // before node is a parent and new node will upper node of them
			hash, err := cs.CalculateHash()
			if err != nil {
				return nil, nil, err
			}
			x := beforeNode.Index[:len(beforeNode.Index)-1]
			leftLevelEntry := t.Levels[len(beforeNode.Index)][len(beforeNode.Index)-2]
			newNode := &Node{
				Hash:   hash,
				C:      cs,
				leaf:   false,
				Tree:   t,
				Number: traversingNumber,
				Index:  x,
				Right:  t.Nodes[i-1],
				Left:   leftLevelEntry,
			}
			t.Nodes = append(t.Nodes, newNode)
			t.Levels[len(x)] = append(t.Levels[len(x)], newNode)
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
	lastDepth := len(t.Nodes[0].Index)
	if (depth - lastDepth) >= 1 {
		t.Levels = nil
		emptyNode := &Node{Tree: t}
		level := make(map[int][]*Node)
		level[0] = append(level[0], emptyNode)
		for _, i := range t.Nodes {
			tmpIndex := i.Index
			newIndex := strings.Repeat("0", depth-lastDepth) + tmpIndex
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

func exportParentsIndex(index string, length int) []string {
	var parentsIndexString []string
	if len(index) == length && !strings.Contains(index, "1") {
		parentsIndexString = append(parentsIndexString, index)
		return parentsIndexString
	} else if len(index) == length && strings.Contains(index, "1") {
		for i := length - 1; i >= 0; i-- {
			newstr := index
			if index[i] == '1' {
				newstr = index[:i] + string('0')
				parentsIndexString = append(parentsIndexString, newstr)
			}
		}
	}
	return parentsIndexString
}

func setParentsToNode(node *Node, t *MerkleTree) bool {
	var parentsIndexString []string
	index := node.Index
	if node.leaf == true { // this is leaf
		if !strings.Contains(index, "1") { // this is first node
			parentsIndexString = append(parentsIndexString, index)
			node.Parent = append(node.Parent, node)
			return true
		} else if strings.Contains(index, "1") {
			for i := len(index) - 1; i >= 0; i-- {
				newstr := index
				if index[i] == '1' {
					newstr = index[:i] + string('0')
					parentsIndexString = append(parentsIndexString, newstr)
					for _, n := range t.Nodes {
						if n.Index == newstr {
							node.Parent = append(node.Parent, n)
						}
					}
				}
			}

		}

	} else { // this is not leaf
		return false
	}
	return true
}

func IsNodeInTree(Index string, t *MerkleTree) *Node {
	for _, i := range t.Nodes {
		if i.Index == Index {
			return i
		}
	}
	return nil
}
