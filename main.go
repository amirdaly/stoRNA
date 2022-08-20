package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"merkletree/merkletree"
	"strconv"
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

func strPad(input string, padLength int) string {
	inputLength := len(input)
	if inputLength >= padLength {
		return input
	}
	return strings.Repeat("0", padLength-inputLength) + input
}

func integerToBinaryString(num int, length int) string {
	binaryString := strconv.FormatInt(int64(num), 2)
	if len(binaryString) >= length {
		return binaryString
	}
	return strings.Repeat("0", length-len(binaryString)) + binaryString
}

func main() {
	/*
		var list []merkletree.Content
		list = append(list, TestContent{x: "amir"})
		list = append(list, TestContent{x: "ali"})
		list = append(list, TestContent{x: "zahra"})
		list = append(list, TestContent{x: "sara"})

		t, err := merkletree.NewTree(list)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(t.Root)
	*/
	// n := 3 //depth
	//var index [n]string
	// for i := 0; i < int(math.Pow(2, float64(n))); i++ {
	// 	fmt.Print(i, "\t ")
	// 	fmt.Println(integerToBinaryString(i, n))
	// }

	var z []merkletree.Content
	z = append(z, TestContent{x: "amir"})

	t, err := merkletree.NewTreeGenesis(z, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)
}
