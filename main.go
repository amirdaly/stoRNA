package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"merkletree/merkletree"
	"strings"
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

	// var list []merkletree.Content
	// list = append(list, TestContent{x: "amir"})
	// list = append(list, TestContent{x: "ali"})
	// list = append(list, TestContent{x: "zahra"})
	// list = append(list, TestContent{x: "sara"})

	// t, err := merkletree.NewTree(list)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var x, y, z []merkletree.Content

	// z = append(z, TestContent{x: "a"})
	// t, err := merkletree.NewTreeGenesis(z, 2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(t)
	// fmt.Println("---------------------------------")

	// x = append(x, TestContent{x: "b"})
	// x = append(x, TestContent{x: "c"})
	// x = append(x, TestContent{x: "d"})
	// merkletree.AddNode(x, t)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(t)
	// fmt.Println("---------------------------------")

	// y = append(y, TestContent{x: "e"})
	// y = append(y, TestContent{x: "f"})
	// y = append(y, TestContent{x: "g"})
	// y = append(y, TestContent{x: "h"})
	// y = append(y, TestContent{x: "i"})
	// y = append(y, TestContent{x: "j"})
	// merkletree.AddNode(y, t)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	t1 := TestContent{x: "a"}
	t, err := merkletree.NewTreeGenesis(t1, 1)
	if err != nil {
		log.Fatal(err)
	}
	t2 := TestContent{x: "b"}
	merkletree.AddNodeToTree(t2, t, 2)

	t3 := TestContent{x: "c"}
	merkletree.AddNodeToTree(t3, t, 2)
	t4 := TestContent{x: "d"}
	merkletree.AddNodeToTree(t4, t, 2)
	fmt.Println(t)

	fmt.Println(t.Levels)
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
