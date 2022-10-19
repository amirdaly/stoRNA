package CommitDAG

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

type CommitDAG struct {
	Root         *Node
	dagRoot      []byte
	Nodes        []*Node
	Levels       map[int][]*Node
	Leafs        []*Node
	hashStrategy func() hash.Hash
}

type Node struct {
	DAG     *CommitDAG
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

// This function Create a new DAG pointer and push first node as initialized node in it.
func NewDAGGenesis(cs Content) (*CommitDAG, error) {
	var defaultHashStrategy = sha256.New               // Default hash strategy for calculation
	t := &CommitDAG{hashStrategy: defaultHashStrategy} // New DAG call by refrence
	var nodes []*Node                                  // Array for nodes of DAG
	hash, err := cs.CalculateHash()                    // Calculate hash
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, &Node{ // Assign data to new node
		Hash:   hash,         // Field of Hash in node structure
		C:      cs,           // Field of Content in node structure
		Data:   cs.GetData(), // Field of Data in node structure
		leaf:   true,         // Is this node leaf?
		DAG:    t,            // Pointer of DAG of this node
		Index:  "0",          // Index of node in binary string based of binary DAG
		done:   true,         // Is calculation completely Done?
		Number: 1,            // Traversing Number of Node - first node is 1
	})
	t.Root = nodes[0]
	t.Nodes = nodes
	t.dagRoot = hash
	emptyNode := &Node{DAG: t}
	level := make(map[int][]*Node)
	level[0] = append(level[0], emptyNode)
	level[1] = append(level[1], t.Nodes[0])
	t.Nodes[0].Number = 1 //add navigation Number 1 to first node of DAG
	t.Levels = level
	if !setParentsToNode(t.Nodes[0], t) {
		return nil, err
	}
	return t, nil
}

// This function add new Leaf to DAG.
func AddNewLeafToDAG(cs Content, t *CommitDAG, depth int) (*Node, error) {
	updateNodesIndex(t, depth)                        // generate or update binary indexes of nodes
	leafsCount := countLeafs(t)                       // count leafs of the DAG
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
		DAG:    t,                // Pointer of DAG of this node
		Index:  index,            // Index of node in binary string based of binary DAG
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

// Tis function add and Intermediate Node to DAG that is not a leaf.
func AddIntermediateNode(cs Content, t *CommitDAG, depth int, index string) (*Node, error) {
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
		DAG:    t,
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

// This Function totaly add new Node to CommitDAG. It means that Proofs of Sequential Work
// and Binary MerkleTree are both Supported.
// There are 4 if check for adding a new Node to DAG
// 1: If lastNode of DAG is leaf and leafs count is odd. So we must add another leaf to DAG.
// 2: If lastNode of DAG is leaf and leafs count is even. So we must add an upper parent to
// last 2 leafs. Its an intermediate node to DAG.
// 3: If lastNode of DAG is an intermediate node and count of nodes in that level is even.
// So we must add an other intermediate Node to DAG.
// 4: If lastNode of DAG is an intermediate node and count of nodes in that level is odd.
// So we must add new leaf to DAG.
func AddNodeToDAG(cs Content, t *CommitDAG) ([]byte, *CommitDAG, error) {
	depth := int(math.Log2(float64(len(t.Nodes) + 2))) // calculate depth of DAG bu log(n) + 2
	lastNode := t.Nodes[len(t.Nodes)-1]
	if lastNode.leaf == true {
		if countLeafs(t)%2 != 0 {
			AddNewLeafToDAG(cs, t, depth)
		} else if countLeafs(t)%2 == 0 {
			var upperLevelCount int
			if depth == 2 && len(t.Nodes) <= 3 {
				upperLevelCount = 0
			} else {
				upperLevelCount = len(t.Levels[depth-1])
			}
			index := integerToBinaryString(upperLevelCount, depth-1)
			updateNodesIndex(t, depth) // generate or update binary indexes of nodes
			_, err := AddIntermediateNode(cs, t, depth, index)
			if err != nil {
				return nil, nil, err
			}
			return t.dagRoot, t, nil
		}
	} else if lastNode.leaf == false {
		lastNodeLevelCount := len(t.Levels[len(t.Nodes[len(t.Nodes)-1].Index)])
		if lastNodeLevelCount%2 == 0 {
			updateNodesIndex(t, depth) // generate or update binary indexes of nodes
			index := lastNode.Index[:len(lastNode.Index)-1]
			AddIntermediateNode(cs, t, depth, index)
		} else if lastNodeLevelCount%2 != 0 {
			AddNewLeafToDAG(cs, t, depth)
		}
	}
	return t.dagRoot, t, nil
}

// This function update the state of node hash be calling
func Update(cs Content, node *Node, t *CommitDAG) {

}

// This function retrun an integer number for Leafs count of the DAG.
func countLeafs(t *CommitDAG) int {
	count := 0
	for _, n := range t.Nodes {
		if n.leaf == true {
			count++
		}
	}
	return count
}

// This function get the DAG pointer and new Depth of DAG, then update each
// nodes Lable index in binary string format.
func updateNodesIndex(t *CommitDAG, depth int) bool {
	T := false
	lastDepth := len(t.Nodes[0].Index)
	if (depth - lastDepth) >= 1 {
		t.Levels = nil
		emptyNode := &Node{DAG: t}
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

// This function get DAG pointer and reorder all arrays of levels of DAG.
func updateLevelsEntry(t *CommitDAG) bool {
	T := false
	t.Levels = nil
	emptyNode := &Node{DAG: t}
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

// This function get an integer number for convert to its binary string with length as second argument.
func integerToBinaryString(num int, length int) string {
	binaryString := strconv.FormatInt(int64(num), 2)
	if len(binaryString) >= length {
		return binaryString
	}
	return strings.Repeat("0", length-len(binaryString)) + binaryString
}

// This function check each node, if node is a leaf then add a
func setParentsToNode(node *Node, t *CommitDAG) bool {
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

// This function checks if a node is in DAG or not.
func IsNodeInDAG(Index string, t *CommitDAG) *Node {
	for _, i := range t.Nodes {
		if i.Index == Index {
			return i
		}
	}
	return nil
}

// This is a helper function for converting *Node pointer to string with some data of it.
func (n *Node) String() string {
	return fmt.Sprintf(
		"Number: %d | Index: %s | leaf: %t | hash: %x data: %s",
		n.Number, n.Index, n.leaf, n.Hash, n.Data)
}

// This function is a helper function for export string for CommitDAG struct that print each nodes data.
func (m *CommitDAG) String() string {
	s := ""
	for _, l := range m.Nodes {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}
