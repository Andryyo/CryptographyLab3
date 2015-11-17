// CryptographyLab3 project main.go
package main

import (
	"fmt"
	"time"
	"math/big"
	"github.com/Andryyo/CryptographyLab2/CryptographyLab2"
)


type Gen struct {
	S []uint8
	D []uint8
	I []uint8
	aes CryprographyLab2.AES
}

func NewGen() *Gen {
	gen := &Gen{}
	key := [4][4]byte{}
	key[0] = [4]byte{0x54, 0x68, 0x61, 0x74}
	key[1] = [4]byte{0x73, 0x20, 0x6D, 0x79}
	key[2] = [4]byte{0x20, 0x4B, 0x75, 0x6E}
	key[3] = [4]byte{0x67, 0x20, 0x46, 0x75}
	gen.aes = *CryprographyLab2.NewAES(key)
	var i uint8
	now := time.Now().Unix()
	gen.D = make([]uint8, 16)
	for i = 0; i < 8; i++ {
		gen.D[i] = (byte)(now >> (i * 8)) 
		gen.D[i + 7] = (byte)(now >> (i * 8)) 
	}
	bigInt := big.NewInt(1)
	fmt.Println("Initializing S...")
	for len(bigInt.Bytes()) < 16 {
		before := time.Now().Nanosecond()
		fmt.Scanln()
		after := time.Now().Nanosecond()
		fmt.Println(after - before)
		bigInt = bigInt.Mul(bigInt, big.NewInt(int64(after - before)))
	}
	fmt.Println("Success")
	gen.S = bigInt.Bytes()[len(bigInt.Bytes()) - 16:]
	gen.I = gen.aes.Encode(gen.D)
	return gen 			
}

func (gen *Gen) GetNext() uint8 {
	R := make([]uint8, 16)
	for i := 0; i < 16; i++ {
		R[i] = gen.I[i] ^ gen.S[i]
	}
	X := gen.aes.Encode(R)
	T := make([]uint8, 16)
	for i := 0; i < 16; i++ {
		T[i] = X[i] ^ gen.I[i]
	}
	gen.S = gen.aes.Encode(T)
	return X[15] & 0x01
}

func (gen *Gen) Test() {
	fmt.Println("Starting test")
	successes := 0
	fails := 0
	nothing := 0
	prediction := 0
	M := [16][2]int{}
	results := [2000000]uint8{};
	for i := 0; i < 4; i++ {
		results[i] = gen.GetNext()
	}
	for i := 4; i < len(results); i++ {
		results[i] = gen.GetNext()
		if prediction == -1 {
			nothing++
		} else if prediction == int(results[i]) {
			successes++
		} else {
			fails++
		}
		if i % 10000 == 0 {
			fmt.Printf("%v: Successes: %v%%\nFails: %v%%\nNothing: %v%%\n", 
				i,
				float32(successes)/float32(i)*100, 
				float32(fails)/float32(i)*100, 
				float32(nothing)/float32(i)*100)
		}
		Iold := 8*results[i-4] + 4*results[i-3] + 2*results[i-2] + results[i-1]
		M[Iold][results[i]]++
		Inew := 8*results[i-3] + 4*results[i-2] + 2*results[i-1] + results[i]
		if M[Inew][0] < M[Inew][1] {
			prediction = 1
		} else if M[Inew][0] > M[Inew][1] {
			prediction = 0
		} else {
			prediction = -1
		}
	}
	//fmt.Printf("Successes: %v\nFails: %v\nNothing: %v\n", successes, fails, nothing)
}

func main() {
	fmt.Println("Hello World!")
	gen := NewGen()
	gen.Test()
}
