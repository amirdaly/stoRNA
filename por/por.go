package por

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/binary"
	"encoding/gob"
	"math"
	"math/big"
	"os"
)

func Keygen() (*rsa.PublicKey, *rsa.PrivateKey) {
	ssk, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	return &ssk.PublicKey, ssk
}

type Tau_zero struct {
	name []byte
	n    int64
	U    []big.Int
}

type Tau struct {
	Tau_zero
	signature []byte
}

func Split(file *os.File) (M [][]byte, S int64, N int64) {
	file.Seek(0, 0)
	s := int64(1)

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	size := fileInfo.Size()
	n := int64(math.Ceil(float64(size / s)))
	// matrix is indexed as m_ij, so the first dimension has n items and the second has s.
	matrix := make([][]byte, n)
	for i := int64(0); i < n; i++ {
		piece := make([]byte, s)
		_, err := file.Read(piece)
		if err != nil {
			panic(err)
		}
		matrix[i] = piece
	}
	return matrix, s, n
}

func hashNameI(name []byte, i int64) *big.Int {
	i_bytes := make([]byte, 4)
	binary.PutVarint(i_bytes, i)
	hashArgument := append(name, i_bytes...)
	hash_array := sha512.Sum512(hashArgument)
	return new(big.Int).SetBytes(hash_array[:])
}

func GenerateAuthenticator(i int64, s int64, tau_zero Tau_zero, piece []byte, ssk *rsa.PrivateKey) *big.Int {
	hash_bigint := hashNameI(tau_zero.name, i+1)

	productory := big.NewInt(1)
	for j := int64(0); j < s; j++ {
		piece_bigint := new(big.Int).SetBytes([]byte{piece[j]})
		productory.Mul(productory, new(big.Int).Exp(&tau_zero.U[j], piece_bigint, nil))
	}

	innerProduct := new(big.Int).Mul(hash_bigint, productory)
	return new(big.Int).Exp(innerProduct, ssk.D, ssk.PublicKey.N)
}

func St(ssk *rsa.PrivateKey, file *os.File) (_tau Tau, _sigma []*big.Int) {
	matrix, s, n := Split(file)
	tau_zero := Tau_zero{n: n}

	tau_zero.name = make([]byte, 512)
	_, err := rand.Read(tau_zero.name)
	if err != nil {
		panic(err)
	}

	tau_zero.U = make([]big.Int, s)
	for i := int64(0); i < s; i++ {
		result, err := rand.Int(rand.Reader, ssk.PublicKey.N)
		if err != nil {
			panic(err)
		}
		tau_zero.U[i] = *result
	}

	var tau_zero_bytes bytes.Buffer
	enc := gob.NewEncoder(&tau_zero_bytes)
	err = enc.Encode(tau_zero)
	if err != nil {
		panic(err)
	}

	hashed_t_0 := sha512.Sum512(tau_zero_bytes.Bytes())
	t_0_signature, err := rsa.SignPKCS1v15(nil, ssk, crypto.SHA512, hashed_t_0[:])
	if err != nil {
		panic(err)
	}
	tau := Tau{Tau_zero: tau_zero, signature: t_0_signature}

	sigmas := make([]*big.Int, n)
	// http://www.golangpatterns.info/concurrency/parallel-for-loop
	sem := make(chan byte, n)
	for i := int64(0); i < n; i++ {
		go func(i int64) {
			sigmas[i] = GenerateAuthenticator(i, s, tau_zero, matrix[i], ssk)
			sem <- 0
		}(i)
	}
	for i := int64(0); i < n; i++ {
		<-sem
	}
	return tau, sigmas
}

type QElement struct {
	I int64
	V int64
}

func Verify_one(tau Tau, spk *rsa.PublicKey) []QElement {
	// Verifica che tau sia corretto
	var tau_zero_bytes bytes.Buffer
	enc := gob.NewEncoder(&tau_zero_bytes)
	err := enc.Encode(tau.Tau_zero)
	if err != nil {
		panic(err)
	}

	hashed_t_0 := sha512.Sum512(tau_zero_bytes.Bytes())
	err = rsa.VerifyPKCS1v15(spk, crypto.SHA512, hashed_t_0[:], tau.signature)
	if err != nil {
		panic(err)
	}

	// l := tau.Tau_zero.n / 2
	l := int64(2)
	n_bigint := big.NewInt(tau.Tau_zero.n)
	ret := make([]QElement, l)
	for i := int64(0); i < l; i++ {
		I_bignum := new(big.Int)
		for {
			I_bignum, err = rand.Int(rand.Reader, n_bigint)
			if err != nil {
				panic(err)
			}
			if I_bignum.Cmp(big.NewInt(0)) == +1 {
				break
			}
		}
		ret[i].I = I_bignum.Int64()

		Q_bignum := new(big.Int)
		for {
			Q_bytes := make([]byte, 4)
			_, err = rand.Read(Q_bytes)
			if err != nil {
				panic(err)
			}
			Q_bignum = new(big.Int).SetBytes(Q_bytes)
			if Q_bignum.Cmp(big.NewInt(0)) == +1 {
				break
			}
		}
		ret[i].V = Q_bignum.Int64()
	}
	return ret
}

func Prove(q []QElement, authenticators []*big.Int, spk *rsa.PublicKey, file *os.File) (_Mu []*big.Int, _Sigma *big.Int) {
	matrix, s, _ := Split(file)

	mu := make([]*big.Int, s)
	for j := int64(0); j < s; j++ {
		mu_j := big.NewInt(0)
		for _, qelem := range q {
			char := new(big.Int).SetBytes([]byte{matrix[qelem.I-1][j]})
			product := new(big.Int).Mul(new(big.Int).SetInt64(qelem.V), char)
			mu_j.Add(mu_j, product)
		}
		mu[j] = mu_j
	}

	sigma := new(big.Int).SetInt64(1)
	for _, qelem := range q {
		sigma.Mul(sigma, new(big.Int).Exp(authenticators[qelem.I-1], new(big.Int).SetInt64(qelem.V), spk.N))
	}
	sigma.Mod(sigma, spk.N)
	return mu, sigma
}

func Verify_two(tau Tau, q []QElement, mus []*big.Int, sigma *big.Int, spk *rsa.PublicKey) bool {
	// Todo: check that the values are in range
	first := new(big.Int).SetInt64(1)
	for _, qelem := range q {
		hash := hashNameI(tau.Tau_zero.name, qelem.I)
		hash.Exp(hash, new(big.Int).SetInt64(qelem.V), spk.N)
		first.Mul(first, hash)
	}
	first.Mod(first, spk.N)

	second := new(big.Int).SetInt64(1)
	s := int64(1)
	for j := int64(0); j < s; j++ {
		second.Mul(second, new(big.Int).Exp(&tau.Tau_zero.U[j], mus[j], spk.N))
	}
	second.Mod(second, spk.N)

	return new(big.Int).Mod(new(big.Int).Mul(first, second), spk.N).Cmp(new(big.Int).Exp(sigma, new(big.Int).SetInt64(int64(spk.E)), spk.N)) == 0
}
