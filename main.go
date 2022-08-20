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

	// var list []merkletree.Content
	// list = append(list, TestContent{x: "amir"})
	// list = append(list, TestContent{x: "ali"})
	// list = append(list, TestContent{x: "zahra"})
	// list = append(list, TestContent{x: "sara"})

	// t, err := merkletree.NewTree(list)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var y, z []merkletree.Content

	z = append(z, TestContent{x: "amir"})
	t, err := merkletree.NewTreeGenesis(z, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)
	fmt.Println("---------------------------------")

	y = append(y, TestContent{x: "ali"})
	merkletree.AddNode(y, t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)
	fmt.Println("---------------------------------")

	y = append(y, TestContent{x: "zahra"})
	y = append(y, TestContent{x: "sara"})
	y = append(y, TestContent{x: "mm"})
	merkletree.AddNode(y, t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)

}
