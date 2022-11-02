package assignment01

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"crypto/sha1"
	"encoding/hex"
	"strings"
)

const difficulty = 4

type block struct { //block node
	Index    int
	Hash     []byte
	Data     string
	PrevHash []byte
	root     Node
}
type Blockchain struct { //block chain
	B_chain []*block
}

func AddBlock(b_chain *Blockchain, b *block) {

	b_chain.B_chain = append(b_chain.B_chain, b)

}

// data is of type Block.Data, and prev is Block.PrevHash
func CalculateHash(data string, prev []byte) []byte {
	// merge data and prev as bytes using bytes.Join

	head := bytes.Join([][]byte{prev, []byte(data)}, []byte{})

	// create a sha256 hash from this merge
	h32 := sha256.Sum256(head)

	fmt.Printf("Header hash: %x\n", h32)

	// sha256.Sum256() returns a [32]byte value, we will use it as a []byte
	// value in the next part of this article, thus the [:] trick
	return h32[:]
}

// data is of type Block.Data, and prev is Block.PrevHash
func DisplayBlock(b *block) {
	fmt.Printf("Current Block Id: %d\nCurrent Block Name: %s\n Current Block Hash: %x\nPrevious Block Hash: %x\n",
		b.Index,
		b.Data,
		b.Hash,
		b.PrevHash,
	)
	fmt.Printf("Curent Block Merkel Tree: \n")
	printTree(b.root)

}

func MineBlock(hash []byte) []byte { //function to mine the block

	target := big.NewInt(1) //for target

	target = target.Lsh(target, uint(256-difficulty)) // perform left shift to adjuest the difficulty level

	fmt.Printf("target: %x\n", target) //print the target space

	// this is the value that will be incremented and added to the header hash
	var nonce int64

	for nonce = 0; nonce < math.MaxInt64; nonce++ { // run to the size of max int

		testNum := big.NewInt(0) // creating a test number for mining

		testNum.Add(testNum.SetBytes(hash), big.NewInt(nonce)) // adding the nounce to test number

		testHash := sha256.Sum256(testNum.Bytes()) //creating the hash of test number

		fmt.Printf("\rhash calculated: %x (nonce used: %d)", testHash, nonce) //printing hash and nounce

		if target.Cmp(testNum.SetBytes(testHash[:])) > 0 { //comparing the test number hash if lies in the target space

			fmt.Println("\n<========> Congratulations!----Found <========>")

			return testHash[:] //return the hash
		}
	}

	return []byte{}
}

func NewBlock(id int, data []string, prev []byte, blockname string) *block {
	str1 := ""
	for i := 0; i < len(data); i++ {
		str1 += "  " + data[i]
	}
	root := buildTree([]Hashable{Block(data[0]), Block(data[1]), Block(data[2]), Block(data[3])})[0].(Node)

	return &block{
		// block Index
		id,
		// first compute a hash with block's header, then mine it
		MineBlock(CalculateHash(str1, prev)),
		// actual data
		blockname,
		// reference to previous block
		prev,
		root,
	}
}

type Node struct {
	left  Hashable
	right Hashable
}

func buildTree(parts []Hashable) []Hashable {
	var nodes []Hashable
	var i int
	for i = 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			nodes = append(nodes, Node{left: parts[i], right: parts[i+1]})
		} else {
			nodes = append(nodes, Node{left: parts[i], right: EmptyBlock{}})
		}
	}
	if len(nodes) == 1 {
		return nodes
	} else if len(nodes) > 1 {
		return buildTree(nodes)
	} else {
		panic("huh?!")
	}
}

type Hashable interface {
	hash() Hash
}

type Hash [20]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

type Block string

func (b Block) hash() Hash {
	return hash([]byte(b)[:])
}

type EmptyBlock struct {
}

func (_ EmptyBlock) hash() Hash {
	return [20]byte{}
}

func (n Node) hash() Hash {
	var l, r [sha1.Size]byte
	l = n.left.hash()
	r = n.right.hash()
	return hash(append(l[:], r[:]...))
}

func hash(data []byte) Hash {
	return sha1.Sum(data)
}

func printTree(node Node) {
	printNode(node, 0)
}

func printNode(node Node, level int) {
	fmt.Printf("(level: %d)  transaction hash:  %s %s\n", level, strings.Repeat(" ", level), node.hash())
	if l, ok := node.left.(Node); ok {
		printNode(l, level+1)
	} else if l, ok := node.left.(Block); ok {
		fmt.Printf("(level: %d)  transaction hash: %s %s (transactions: %s)\n", level+1, strings.Repeat(" ", level+1), l.hash(), l)
	}
	if r, ok := node.right.(Node); ok {
		printNode(r, level+1)
	} else if r, ok := node.right.(Block); ok {
		fmt.Printf("(level: %d) transaction hash: %s %s (transactions: %s)\n", level+1, strings.Repeat(" ", level+1), r.hash(), r)
	}
}
