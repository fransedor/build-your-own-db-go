package main

import (
	"encoding/binary"
	"fmt"
	"log"
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
	arr := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	right := binary.LittleEndian.Uint16(arr[:4])
	fmt.Printf("%d \n %d", arr[:4], right)

	fmt.Print("Hello world")
}
