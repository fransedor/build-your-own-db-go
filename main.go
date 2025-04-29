package main

import (
	"fmt"
	"log"

	"github.com/fransedor/build-your-own-db-go/bnode"
)

const (
	BNODE_INTERNAL = 1 // internal nodes with pointers
	BNODE_LEAF     = 2 // leaf nodes with values
)

type Node struct {
	children []*Node
	keys     [][]byte
	vals     [][]byte
}

//func Encode(node *Node) []byte
//func Decode(page []byte) (*Node, error)

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init() {
	node1max := 4 + 1*8 + 1*2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	if node1max > BTREE_PAGE_SIZE {
		log.Fatal("Node max size exceeded")
	}

}

func main() {

	node := bnode.BNode(make([]byte, BTREE_PAGE_SIZE))
	node.SetHeader(BNODE_LEAF, 2)
	node.AppendKV(0, 0, []byte("k1"), []byte("hi"))
	// ^ 1st KV
	node.AppendKV(1, 0, []byte("k3"), []byte("hello"))
	//                 ^ 2nd KV
	firstKey := string(node.GetKey(0))
	secondKey := string(node.GetKey(1))

	firstVal := string(node.GetVal(0))
	secondVal := string(node.GetVal(1))

	fmt.Printf("First key and val: %v: %v\n", firstKey, firstVal)
	fmt.Printf("Second key and val: %v: %v\n", secondKey, secondVal)
}
