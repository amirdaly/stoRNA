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
	t10 := TestContent{x: "j"}
	merkletree.AddNodeToTree(t10, t, 3)
	t11 := TestContent{x: "k"}
	merkletree.AddNodeToTree(t11, t, 3)
	t12 := TestContent{x: "l"}
	merkletree.AddNodeToTree(t12, t, 3)
	t13 := TestContent{x: "m"}
	merkletree.AddNodeToTree(t13, t, 3)
	t14 := TestContent{x: "n"}
	merkletree.AddNodeToTree(t14, t, 3)
	t15 := TestContent{x: "o"}
	merkletree.AddNodeToTree(t15, t, 4)
	t16 := TestContent{x: "p"}
	merkletree.AddNodeToTree(t16, t, 4)
	t17 := TestContent{x: "q"}
	merkletree.AddNodeToTree(t17, t, 4)
	t18 := TestContent{x: "r"}
	merkletree.AddNodeToTree(t18, t, 4)

	fmt.Println(t)

	/*  POR example run

		func main() {
		fmt.Printf("Generating RSA keys...\n")
		spk, ssk := Keygen()
		fmt.Printf("Generated!\n")

		fmt.Printf("Signing file...\n")
		file, err := os.Open("./example.txt")
		if err != nil {
			panic(err)
		}
		tau, authenticators := St(ssk, file)
		fmt.Printf("Signed!\n")

		fmt.Printf("Generating challenge...\n")
		q := Verify_one(tau, spk)
		fmt.Printf("Generated!\n")

		fmt.Printf("Issuing proof...\n")
		mu, sigma := Prove(q, authenticators, spk, file)
		fmt.Printf("Issued!\n")

		fmt.Printf("Verifying proof...\n")
		yes := Verify_two(tau, q, mu, sigma, spk)
		fmt.Printf("Result: %t!\n", yes)
		if yes {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}


	*/
}
