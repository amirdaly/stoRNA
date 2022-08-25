package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"merkletree/merkletree"
)

//------------------------------ Testbed Area ----------------------

type TestContent struct {
	x string
}

func (t TestContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
func (t TestContent) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}

func main() {
	t1 := TestContent{x: "a"}
	t, err := merkletree.NewTreeGenesis(t1, 1)
	if err != nil {
		log.Fatal(err)
	}
	//-----------------------------------
	t2 := TestContent{x: "b"}
	merkletree.AddNodeToTree(t2, t, 2)
	t3 := TestContent{x: "c"}
	merkletree.AddNodeToTree(t3, t, 2)
	t4 := TestContent{x: "d"}
	merkletree.AddNodeToTree(t4, t, 2)
	t5 := TestContent{x: "e"}
	merkletree.AddNodeToTree(t5, t, 2)
	t6 := TestContent{x: "f"}
	merkletree.AddNodeToTree(t6, t, 2)
	//-----------------------------------
	t7 := TestContent{x: "g"}
	merkletree.AddNodeToTree(t7, t, 3)
	t8 := TestContent{x: "h"}
	merkletree.AddNodeToTree(t8, t, 3)
	t9 := TestContent{x: "i"}
	merkletree.AddNodeToTree(t9, t, 3)

	fmt.Println(t)

}
