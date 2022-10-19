package stoRNA

import (
	"CommitDAG/CommitDAG"
	"CommitDAG/por"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"
)

type TestContent struct {
	x *big.Int
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

func Store(fileName string) (*rsa.PublicKey, *os.File, Tau, []*big.Int, *rsa.PublicKey) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	spk, ssk := por.Keygen()
	tau, authenticators := por.St(ssk, file)
	tg := tau
	randomSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randomSource)
	keygen, err := rsa.GenerateKey(r, 256)
	if err != nil {
		panic(err)
	}
	rs := keygen.PublicKey
	return spk, &file, tg, authenticators, rs
}

func Prove(file *os.File, tag tau, rs *big.Int, depositTime int, auditFrequency int, []*big.Int, []*big.Int) {
	i := 0
	st := rs
	for et := 0; et <= depositTime; et++ {
		spk, file, tau, authenticators, rs := store(file)
		q := por.Verify_one(tag, spk)
		mu, st = por.Prove(q, authenticators, spk, file)
		i++
		var t CommitDAG.CommitDAG
		var c []byte
		var h, newProve []*big.Int
		if i == 1 { // first node to genesis DAG
			t1 := TestContent{x: st}
			t, err := CommitDAG.NewDAGGenesis(t1)
			if err != nil {
				log.Fatal(err)
			}
			c[i] = t.dagRoot
		} else {
			temp := TestContent{x: st}
			hash, t, err := CommitDAG.AddNodeToDAG(temp, t)
			if err != nil {
				log.Fatal(err)
			}
			c[i] = hash
		}
		h[i] = st
		newmu, pi[i] := por.Prove(q, authenticators, spk, file)
		st = CommitDAG.Update(newProve[i], t.Node[i-1], t)
		et = et + auditFrequency
	}
	return t.dagRoot, h, pi
}

func Verify(com []byte, t *CommitDAG, tag Tau, rs *rsa.PublicKey, spk *rsa.PublicKey, h []*big.Int, pi []*big.Int) bool {
	for i:=0; i < len(t.Nodes); i++ {
		if por.Verify_two(tag, q, mus h[i+1], pi[i+1], spk) == false {
			return false
		}
		nodeParents := t.Nodes[i].Parents
		if len(nodeParents) >= 1 {
			var parentsSum []byte
			for _, parent := range nodeParents {
				parentsSum = parentsSum + parent.Hash
			}
			parentsHash := sha256.New()
			if _, err := parentsHash.Write([]byte(parentsSum)); err != nil {
				return nil, err
			}
			if t.Nodes[i].hash != parentsHash {
				return false
			}
		} 
		if t.Nodes[i].Left != nil && t.Nodes[i].Right != nil {
			var parentsSum []byte
			parentsSum = t.Nodes[i].Left.Hash + t.Nodes[i].Right.Hash
			parentsHash := sha256.New()
			if _, err := parentsHash.Write([]byte(parentsSum)); err != nil {
				return nil, err
			}
			if t.Nodes[i].hash != parentsHash {
				return false
			}
		}
		return true
	}
}
