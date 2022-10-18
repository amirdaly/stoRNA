package main

import (
	"CommitDAG/CommitDAG"
	"CommitDAG/por"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
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
func (t TestContent) Equals(other CommitDAG.Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}
func (t TestContent) GetData() string {
	return t.x
}

func CalculateHash(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := sha256.New()
	defer duration(track("SHA256 Calculation Runtime"))
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))

}

func SortFileSizeDescend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size() > files[j].Size()
	})
}
func SortFileSizeAscend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size() < files[j].Size()
	})
}

func sha256RunTest(directoryPath string) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}
	SortFileSizeAscend(files)
	count := 1
	for _, v := range files {
		fmt.Println(count)
		path := directoryPath + v.Name()
		inf, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		fs := float64(inf.Size())

		hash, err := CalculateHash(path)
		if err != nil {
			panic(err)
		}
		fmt.Println(path)
		fmt.Printf("File size: %v B\n", fs)
		fmt.Printf("SHA256 Hash: %x\n", hash)
		count++
	}
}

func porTestRun(directoryPath string) {
	/*  POR example run  */

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}
	SortFileSizeAscend(files)
	count := 1
	for _, v := range files {
		fmt.Println(count)
		path := directoryPath + v.Name()
		inf, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		fs := float64(inf.Size())

		file, err := os.Open(path)
		if err != nil {
			return
		}
		fmt.Printf("Generating RSA keys...\n")
		spk, ssk := por.Keygen()
		fmt.Printf("Generated!\n")
		fmt.Printf("Signing file...\n")
		tau, authenticators := por.St(ssk, file)
		fmt.Printf("Signed!\n")
		fmt.Printf("Generating challenge...\n")
		q := por.Verify_one(tau, spk)
		fmt.Printf("Generated!\n")

		fmt.Printf("Issuing proof for file ..\n")

		mu, sigma := por.Prove(q, authenticators, spk, file)
		fmt.Printf("Issued!\n")

		fmt.Printf("Verifying proof of file: ")
		fmt.Println(path)
		fmt.Printf("File size: %v KB\n", fs/1024)
		yes := por.Verify_two(tau, q, mu, sigma, spk)
		fmt.Printf("Result: %t!\n", yes)
		if yes {
			file.Close()
			count++
			continue
		} else {
			file.Close()
			os.Exit(1)
		}
	}
}

func main() {
	// sha256RunTest()

	t1 := TestContent{x: "a"}
	t, err := CommitDAG.NewDAGGenesis(t1)
	if err != nil {
		log.Fatal(err)
	}
	t2 := TestContent{x: "b"}
	CommitDAG.AddNodeToDAG(t2, t)
	t3 := TestContent{x: "c"}
	CommitDAG.AddNodeToDAG(t3, t)
	t4 := TestContent{x: "d"}
	CommitDAG.AddNodeToDAG(t4, t)
	t5 := TestContent{x: "e"}
	CommitDAG.AddNodeToDAG(t5, t)
	t6 := TestContent{x: "f"}
	CommitDAG.AddNodeToDAG(t6, t)
	t7 := TestContent{x: "g"}
	CommitDAG.AddNodeToDAG(t7, t)
	t8 := TestContent{x: "h"}
	CommitDAG.AddNodeToDAG(t8, t)
	t9 := TestContent{x: "i"}
	CommitDAG.AddNodeToDAG(t9, t)
	t10 := TestContent{x: "j"}
	CommitDAG.AddNodeToDAG(t10, t)
	t11 := TestContent{x: "k"}
	CommitDAG.AddNodeToDAG(t11, t)
	t12 := TestContent{x: "l"}
	CommitDAG.AddNodeToDAG(t12, t)
	t13 := TestContent{x: "m"}
	CommitDAG.AddNodeToDAG(t13, t)
	t14 := TestContent{x: "n"}
	CommitDAG.AddNodeToDAG(t14, t)
	t15 := TestContent{x: "o"}
	CommitDAG.AddNodeToDAG(t15, t)
	t16 := TestContent{x: "p"}
	CommitDAG.AddNodeToDAG(t16, t)
	t17 := TestContent{x: "q"}
	CommitDAG.AddNodeToDAG(t17, t)
	t18 := TestContent{x: "r"}
	CommitDAG.AddNodeToDAG(t18, t)

}
