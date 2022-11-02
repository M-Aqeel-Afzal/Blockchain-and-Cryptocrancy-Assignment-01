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

const difficulty_level = 6 //Determines the difficulty to mine block
type Block string
type Hash [20]byte
type block struct { //block node
	Index     int
	Hash      []byte
	Data      string
	Prev_Hash []byte
	root      Node
}
type Blockchain struct { //block chain
	B_chain []*block
}

type EmptyBlock struct {
}

type Hashable interface {
	hash() Hash
}

// Creating the left and right node of Merkel Tree
type Node struct {
	left_node  Hashable
	right_node Hashable
}

// Appends the Block in the BlockChain
func AddBlock(b_chain *Blockchain, b *block) {
	b_chain.B_chain = append(b_chain.B_chain, b)
}

// Calculates the Hash of the Block
// Data type is Block.Data
// prev contains the hash of previous block
func CalculateHash(data string, prev []byte) []byte {
	// merging data and prev hash
	head := bytes.Join([][]byte{prev, []byte(data)}, []byte{})
	// creating sha256 hash
	hash32 := sha256.Sum256(head)
	// sha256 returns [32]byte
	fmt.Printf("Header hash: %x\n", hash32)
	return hash32[:]
}

// Displays the Data of Block
func DisplayBlock(b *block) {
	fmt.Printf("Current Block Id: %d\nCurrent Block Name: %s\n Current Block Hash: %x\nPrevious Block Hash: %x\n",
		b.Index,
		b.Data,
		b.Hash,
		b.Prev_Hash,
	)
	fmt.Printf("Curent Block Merkel Tree: \n")
	DisplayMerkelTree(b.root) // Calling Display MerkleTree
}

func MineBlock(hash []byte) []byte { //function to mine the block
	target := big.NewInt(1)                                 //for target
	target = target.Lsh(target, uint(256-difficulty_level)) // perform left shift to adjuest the difficulty level
	fmt.Printf("target: %x\n", target)                      //print the target space
	// this is the value that will be incremented and added to the header hash
	var nonce int64
	for nonce = 0; nonce < math.MaxInt64; nonce++ { // run to the size of max int
		testNum := big.NewInt(0)                                              // creating a test number for mining
		testNum.Add(testNum.SetBytes(hash), big.NewInt(nonce))                // adding the nounce to test number
		testHash := sha256.Sum256(testNum.Bytes())                            //creating the hash of test number
		fmt.Printf("\rhash calculated: %x (nonce used: %d)", testHash, nonce) //printing hash and nounce
		if target.Cmp(testNum.SetBytes(testHash[:])) > 0 {                    //comparing the test number hash if lies in the target space
			fmt.Println("\n<========> Congratulations!----Found <========>")
			return testHash[:] //return the hash
		}
	}

	return []byte{}
}

// Creates the New Block
func NewBlock(id int, data []string, prev []byte, blockname string) *block {
	str1 := ""
	for i := 0; i < len(data); i++ {
		str1 += "  " + data[i]
	}
	root := MerkleTree([]Hashable{Block(data[0]), Block(data[1]), Block(data[2]), Block(data[3])})[0].(Node)
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

// Creates the Merkle Tree
func MerkleTree(parts []Hashable) []Hashable {
	var nodes []Hashable
	var i int
	for i = 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			nodes = append(nodes, Node{left_node: parts[i], right_node: parts[i+1]})
		} else {
			nodes = append(nodes, Node{left_node: parts[i], right_node: EmptyBlock{}})
		}
	}
	if len(nodes) == 1 {
		return nodes
	} else if len(nodes) > 1 {
		return MerkleTree(nodes)
	}
	return nodes
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// Hash of the BlockChain
func (b Block) hash() Hash {
	return hash([]byte(b)[:])
}

// Hash Function Genisis Block
func (_ EmptyBlock) hash() Hash {
	return [20]byte{}
}

// Hash Function for Nodes
func (n Node) hash() Hash {
	var left, right [sha1.Size]byte
	left = n.left_node.hash()
	right = n.right_node.hash()
	return hash(append(left[:], right[:]...))
}

// Hash Function for the Data
func hash(data []byte) Hash {
	return sha1.Sum(data)
}

// Display Merkel Tree
func DisplayMerkelTree(node Node) {
	printNode(node, 0)
}

// Recursive Function to print the Merkle Tree of current Block
func printNode(node Node, level int) {
	fmt.Printf("(level: %d)  transaction hash:  %s %s\n", level, strings.Repeat(" ", level), node.hash())
	if left, check := node.left_node.(Node); check {
		printNode(left, level+1)
	} else if left, check := node.left_node.(Block); check {
		fmt.Printf("(level: %d)  transaction hash: %s %s (transactions: %s)\n", level+1, strings.Repeat(" ", level+1), left.hash(), left)
	}
	if right, check := node.right_node.(Node); check {
		printNode(right, level+1)
	} else if right, check := node.right_node.(Block); check {
		fmt.Printf("(level: %d) transaction hash: %s %s (transactions: %s)\n", level+1, strings.Repeat(" ", level+1), right.hash(), right)
	}
}
