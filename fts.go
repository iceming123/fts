package fts

import (
	// "errors"
	"encoding/hex"
	"fmt"
	// "math/big"
	"math/rand"
	"bytes"
	"strconv"
	"strings"
	"github.com/fts/rlp"
	"golang.org/x/crypto/sha3"
)

type Hash [32]byte
func (h Hash) Hex() string { return hex.EncodeToString(h[:]) }

func RlpHash(x interface{}) (h Hash) {
	hw := sha3.New256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type Stakeholder struct {
	name string;
	coins int;
}
func (s *Stakeholder) getName() string {
	return s.name
}
func (s *Stakeholder) getCoins() int {
	return s.coins
}
func (s *Stakeholder) toBytes() []byte {
	ss := fmt.Sprintf("%s%d",s.name,s.coins)
	return []byte(ss)
}
func (s *Stakeholder) toString() string {
	return s.name
}
func (s *Stakeholder) clone() *Stakeholder {
	return &Stakeholder{
		name:	s.name,
		coins: 	s.coins,
	}
}
////////////////////////////////////////////////////////////////
type node struct {
	left 		*node
	right 		*node
	sholder 	*Stakeholder
	hash 		Hash
}
func (n *node) isLeaf() bool {
	return n.sholder != nil
}
func (n *node) getStakeholder() *Stakeholder {
	return n.sholder
}
func (n *node) getLeftNode() *node {
	return n.left
}
func (n *node) getRightNode() *node {
	return n.right
}
func (n *node) getMerkleHash() Hash {
	return n.hash
}
func (n *node) getCoins() int {
	if n.isLeaf() {
		return n.sholder.getCoins()
	}
	return n.left.getCoins() + n.right.getCoins()
}
func newNodeFromSHolder(s *Stakeholder) *node {
	return &node {
		left:		nil,
		right: 		nil,
		sholder:	s.clone(),
		hash:		RlpHash(s.toBytes()),
	}
}
func newNode1(left,right *node,hash Hash) *node {
	return &node {
		left:		left,
		right: 		right,
		sholder:	nil,
		hash:		hash,
	}
}
////////////////////////////////////////////////////////////////
type ProofEntry struct {
	hash 		Hash
	x1			int
	x2 			int
}
func (p *ProofEntry) getLeftBound() int {
	return p.x1
}
func (p *ProofEntry) getRightBound() int {
	return p.x2
}
func (p *ProofEntry) getMerkleHash() Hash {
	return p.hash
}
func (p *ProofEntry) toString() string {
	return fmt.Sprintf("%s,%d,%d",p.hash.Hex(),p.x1,p.x2)
}
func newProofEntry(hash Hash,x1,x2 int) *ProofEntry {
	return &ProofEntry{
		hash: 	hash,
		x1:		x1,
		x2:		x2,
	}
}
////////////////////////////////////////////////////////////////
type ftsResult struct {
	sholder 		*Stakeholder
	merkleProof 	[]*ProofEntry
}
func (p *ftsResult) getStakeholder() *Stakeholder {
	return p.sholder
}
func (p *ftsResult) getMerkleProof() []*ProofEntry {
	return p.merkleProof
}
func (p *ftsResult) toString() string {
	proofs := make([]string, len(p.merkleProof))
	for i,v := range p.merkleProof {
		proofs[i] = v.toString()
	}
	return fmt.Sprintf("merkleProof {\n %s }\n stakeholder \n {%s} \n",
	strings.Join(proofs, "\n"),p.sholder.toString())
}
func newFtsResult(sholder *Stakeholder,proofs []*ProofEntry) *ftsResult {
	return &ftsResult{
		sholder: 	sholder,
		merkleProof: proofs,
	}
}
///////////////////////////////////////////////////////////////////////////////////
func makeNodeHash(left,right,leftValue,rightValue []byte) Hash {
	b := make([]byte,0,0)
	b = append(b,left...)
	b = append(b,right...)
	b = append(b,leftValue...)
	b = append(b,rightValue...)
	return RlpHash(b)
}
func nextInt(max int,rnd *rand.Rand) int {
	return rnd.Intn(max)
}
func CreateMerkleTree(stakeholders []*Stakeholder) []*node {
	tree := make([]*node,len(stakeholders)*2)
	fmt.Println("Creating Merkle tree with ",len(tree)-1,"nodes.")

	for i,v := range stakeholders {
		tree[len(stakeholders) + i] = newNodeFromSHolder(v)
	}
	for i:=len(stakeholders)-1;i>0;i-- {
		left,right := tree[i*2],tree[i*2 + 1]
		b1,b2 := left.getMerkleHash(),right.getMerkleHash()
		hash := makeNodeHash(b1[:],b2[:],
		[]byte(strconv.Itoa(left.getCoins())),
		[]byte(strconv.Itoa(right.getCoins())))
		tree[i] = newNode1(left,right,hash)
	}
	for i,v := range tree {
		if v != nil {
			fmt.Println("HASH:",v.getMerkleHash().Hex()," Index:",i)
		}
	}
	return tree;
}
func FtsTree(tree []*node,rnd *rand.Rand) *ftsResult {
	merkleProof := make([]*ProofEntry,0,0)
	i := 1
	for {
		if tree[i].isLeaf() {
			return newFtsResult(tree[i].getStakeholder(),merkleProof)
		}
		x1,x2 := tree[i].getLeftNode().getCoins(),tree[i].getRightNode().getCoins()
		fmt.Println("left subtree coins:",x1," right subtree coins:",x2)
		r := nextInt(x1 + x2,rnd) + 1
		fmt.Println("Picking coin number:",r)
		if r <= x1 {
			fmt.Println("Choosing left subtree...")
			i *= 2
			merkleProof = append(merkleProof,newProofEntry(tree[i + 1].getMerkleHash(),x1,x2))
		} else {
			fmt.Println("Choosing right subtree...")
			i = 2*i + 1
			merkleProof = append(merkleProof,newProofEntry(tree[i - 1].getMerkleHash(),x1,x2))
		}
	}
}
func FtsVerify(merkleRootHash Hash, res *ftsResult,rnd *rand.Rand) bool {
	resPath := make([]byte,0,0)
	for _,v := range res.getMerkleProof() {
		x1,x2 := v.getLeftBound(),v.getRightBound()
		r := nextInt(x1 + x2,rnd) + 1
		if r <= x1 {
			fmt.Println("0 ")
			resPath = append(resPath,'0')
		} else {
			fmt.Println("1 ")
			resPath = append(resPath,'1')
		}
	}
	fmt.Println("OK")
	hx := RlpHash(res.getStakeholder().toBytes())
	for i:=len(res.getMerkleProof()) -1; i >= 0; i-- {
		proof := res.getMerkleProof()[i]
		x1 := []byte(strconv.Itoa(proof.getLeftBound()))
		x2 := []byte(strconv.Itoa(proof.getRightBound()))
		hy := proof.getMerkleHash()
		if resPath[i] == '0' {
			hx = makeNodeHash(hx[:],hy[:],x1,x2)
		} else {
			hx = makeNodeHash(hy[:],hx[:],x1,x2)
		}
		fmt.Println("Next hash:",hx.Hex())
	}
	b := bytes.Equal(hx[:],merkleRootHash[:])
	if b {
		fmt.Println("Root hash matches!")
	} else {
		fmt.Println("Invalid Merkle proof.")
	}
	return b
}