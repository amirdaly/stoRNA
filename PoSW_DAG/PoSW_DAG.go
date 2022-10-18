package PoSW_DAG

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"math"
	"strconv"
	"strings"
)

type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
	GetData() string
}

type PoSW_DAG struct {
	Root         *Node
	merkleRoot   []byte
	Nodes        []*Node
	Levels       map[int][]*Node
	Leafs        []*Node
	hashStrategy func() hash.Hash
}

type Node struct {
	Tree    *PoSW_DAG
	Parents []*Node
	Left    *Node
	Right   *Node
	leaf    bool
	Hash    []byte
	C       Content
	Data    string
	Index   string
	Number  int
	done    bool
	verify  []byte
}

func NewTreeGenesis(cs Content, length int) (*PoSW_DAG, error) {
	var defaultHashStrategy = sha256.New              // Default hash strategy for calculation
	t := &PoSW_DAG{hashStrategy: defaultHashStrategy} // New tree call by refrence
	var nodes []*Node                                 // Array for nodes of tree
	hash, err := cs.CalculateHash()                   // Calculate hash
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, &Node{ // Assign data to new node
		Hash:   hash,         // Field of Hash in node structure
		C:      cs,           // Field of Content in node structure
		Data:   cs.GetData(), // Field of Data in node structure
		leaf:   true,         // Is this node leaf?
		Tree:   t,            // Pointer of tree of this node
		Index:  "0",          // Index of node in binary string based of binary tree
		done:   true,         // Is calculation completely Done?
		Number: 1,            // Traversing Number of Node - first node is 1
	})
	t.Root = nodes[0]
	t.Nodes = nodes
	t.merkleRoot = hash
	emptyNode := &Node{Tree: t}
	level := make(map[int][]*Node)
	level[0] = append(level[0], emptyNode)
	level[1] = append(level[1], t.Nodes[0])
	t.Nodes[0].Number = 1 //add navigation Number 1 to first node of tree
	t.Levels = level
	if !setParentsToNode(t.Nodes[0], t) {
		return nil, err
	}
	return t, nil
}

func AddNewLeafToTree(cs Content, t *PoSW_DAG, depth int) (*Node, error) {
	updateNodesIndex(t, depth)                        // generate or update binary indexes of nodes
	leafsCount := countLeafs(t)                       // count leafs of the tree
	traversingNumber := len(t.Nodes) + 1              // travering number of node is count of nodes + 1
	index := integerToBinaryString(leafsCount, depth) // export string binary index of leaf

	hash, err := cs.CalculateHash() // calculate hash of entry
	if err != nil {
		return nil, err
	}
	newNode := &Node{
		Hash:   hash,             // Field of Hash in node structure
		C:      cs,               // Field of Content in node structure
		Data:   cs.GetData(),     // Field of Data in node structure
		leaf:   true,             // Is this node leaf?
		Tree:   t,                // Pointer of tree of this node
		Index:  index,            // Index of node in binary string based of binary tree
		Number: traversingNumber, // Traversing Number of Node
		done:   true,             // Is calculation completely Done?
	}
	t.Nodes = append(t.Nodes, newNode)
	t.Levels[depth] = append(t.Levels[depth], newNode) // Add new leaf to leafs level nodes
	if !setParentsToNode(newNode, t) {
		fmt.Printf("error in adding parent to %s\n", newNode.Data)
		return nil, err
	}
	return t.Root, nil
}

func AddIntermediateNode(cs Content, t *PoSW_DAG, depth int, index string) (*Node, error) {
	traversingNumber := len(t.Nodes) + 1 // travering number of node is count of nodes + 1
	hash, err := cs.CalculateHash()
	if err != nil {
		return nil, err
	}
	newNode := &Node{
		Hash:   hash,
		C:      cs,
		Data:   cs.GetData(),
		leaf:   false,
		Tree:   t,
		Number: traversingNumber,
		Index:  index,
	}
	t.Nodes = append(t.Nodes, newNode)
	newNodeLevel := len(index)
	t.Levels[len(index)] = append(t.Levels[newNodeLevel], newNode) // Add new Node to its level nodes
	if !setParentsToNode(newNode, t) {
		fmt.Printf("error in adding parent to %s\n", newNode.Data)
		return nil, err
	}
	return t.Root, nil
}

func AddNodeToTree(cs Content, t *PoSW_DAG) (*Node, error) {
	depth := int(math.Log2(float64(len(t.Nodes) + 2))) // calculate depth of tree bu log(n) + 2
	lastNode := t.Nodes[len(t.Nodes)-1]
	// There are 4 if check for adding a new Node to tree
	// 1: If lastNode of tree is leaf and leafs count is odd. So we must add another leaf to tree.
	// 2: If lastNode of tree is leaf and leafs count is even. So we must add an upper parent to last 2 leafs. Its an intermediate node to tree.
	// 3: If lastNode of tree is an intermediate node and count of nodes in that level is even. So we must add an other intermediate Node to tree.
	// 4: If lastNode of tree is an intermediate node and count of nodes in that level is odd. So we must add new leaf to tree.
	if lastNode.leaf == true {
		if countLeafs(t)%2 != 0 {
			AddNewLeafToTree(cs, t, depth)
		} else if countLeafs(t)%2 == 0 {
			var upperLevelCount int
			if depth == 2 && len(t.Nodes) <= 3 {
				upperLevelCount = 0
			} else {
				upperLevelCount = len(t.Levels[depth-1])
			}
			index := integerToBinaryString(upperLevelCount, depth-1)
			updateNodesIndex(t, depth) // generate or update binary indexes of nodes
			newNodeAdded, err := AddIntermediateNode(cs, t, depth, index)
			if err != nil {
				return nil, err
			}
			return newNodeAdded, nil
		}
	} else if lastNode.leaf == false {
		lastNodeLevelCount := len(t.Levels[len(t.Nodes[len(t.Nodes)-1].Index)])
		if lastNodeLevelCount%2 == 0 {
			updateNodesIndex(t, depth) // generate or update binary indexes of nodes
			index := lastNode.Index[:len(lastNode.Index)-1]
			AddIntermediateNode(cs, t, depth, index)
		} else if lastNodeLevelCount%2 != 0 {
			AddNewLeafToTree(cs, t, depth)
		}
	}
	return t.Root, nil
}

func countLeafs(t *PoSW_DAG) int {
	count := 0
	for _, n := range t.Nodes {
		if n.leaf == true {
			count++
		}
	}
	return count
}

func updateNodesIndex(t *PoSW_DAG, depth int) bool {
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

func updateLevelsEntry(t *PoSW_DAG) bool {
	T := false
	t.Levels = nil
	emptyNode := &Node{Tree: t}
	level := make(map[int][]*Node)
	level[0] = append(level[0], emptyNode)
	for _, i := range t.Nodes {
		tmpIndex := i.Index
		l := len(tmpIndex)
		level[l] = append(level[l], i)
		T = true
	}
	t.Levels = level
	return T
}

func integerToBinaryString(num int, length int) string {
	binaryString := strconv.FormatInt(int64(num), 2)
	if len(binaryString) >= length {
		return binaryString
	}
	return strings.Repeat("0", length-len(binaryString)) + binaryString
}

func setParentsToNode(node *Node, t *PoSW_DAG) bool {
	var parentsIndexString []string
	index := node.Index
	if node.leaf == true { // this is leaf
		if !strings.Contains(index, "1") { // this is first node
			node.Parents = nil
			node.Parents = append(node.Parents, node)
			return true
		} else if strings.Contains(index, "1") {
			for i := len(index) - 1; i >= 0; i-- {
				newstr := index
				if index[i] == '1' {
					newstr = index[:i] + string('0')
					parentsIndexString = append(parentsIndexString, newstr)
					for _, n := range t.Nodes {
						if n.Index == newstr {
							node.Parents = append(node.Parents, n)
						}
					}
				}
			}

		}

	} else if node.leaf == false { // this is intermediate node
		leftString := index + "0"
		rightString := index + "1"
		for _, n := range t.Nodes {
			if n.Index == leftString {
				node.Left = n
			}
			if n.Index == rightString {
				node.Right = n
			}
		}
		return true
	}
	return true
}

func IsNodeInTree(Index string, t *PoSW_DAG) *Node {
	for _, i := range t.Nodes {
		if i.Index == Index {
			return i
		}
	}
	return nil
}

func (n *Node) String() string {
	return fmt.Sprintf("Number: %d | Index: %s | leaf: %t | hash: %x data: %s", n.Number, n.Index, n.leaf, n.Hash, n.Data)
}

func (m *PoSW_DAG) String() string {
	s := ""
	for _, l := range m.Nodes {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}

func printLevels(t *PoSW_DAG) {
	for i := 0; i < len(t.Levels); i++ {
		fmt.Printf("Level %d counted nodes are: %d\n", i, len(t.Levels[i]))
		for t, j := range t.Levels[i] {

			fmt.Println(t, j)
		}
	}
}

func printLevel(t *PoSW_DAG, level int) {
	fmt.Printf("Level %d counted nodes are: %d\n", level, len(t.Levels[level]))
	for t, j := range t.Levels[level] {

		fmt.Println(t, j)
	}
}

type NewContent struct {
	x string
}

func (t NewContent) CalculateHash() ([]byte, error) {

	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
func (t NewContent) Equals(other Content) (bool, error) {
	return false, nil
}

func (t NewContent) GetData() string {
	return t.x
}
