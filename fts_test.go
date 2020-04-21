package fts

import (
	"testing"
	// "encoding/hex"
	"fmt"
	// "math/big"
	"math/rand"
	// "bytes"
	// "strconv"
	// "strings"
	// "github.com/fts/rlp"
	// "golang.org/x/crypto/sha3"
)

func TestRand(t *testing.T)  {
	r1,r2 := rand.New(rand.NewSource(100)),rand.New(rand.NewSource(100))
	for i:=0;i<10;i++ {
		fmt.Println(r1.Intn(30))
	}
	for i:=0;i<10;i++ {
		fmt.Println(r2.Intn(30))
	}
	fmt.Println("finish")
}
func TestFts(t *testing.T) {
	stakeholders := make([]*Stakeholder,0,0)
	c := 20
	for i := 0; i<8; i++ {
		stakeholders = append(stakeholders,&Stakeholder{
			name: 	fmt.Sprintf("Stakeholder %d",i),
			coins:	c,
		})
		if c % 2 == 0 {
			c = c / 2
		} else {
			c = c * 3 + 1
		}
	}
	tree := CreateMerkleTree(stakeholders)
	fmt.Println("Doing follow-the-satoshi in the stake tree")
	res := FtsTree(tree,rand.New(rand.NewSource(25)))
	fmt.Println("res:",res.toString())
	fmt.Println("Verifying the result.")
	FtsVerify(tree[1].getMerkleHash(),res,rand.New(rand.NewSource(25)))
	fmt.Println("finish")
}