// The huffman tree
package main

import "fmt"

type Node struct {
	symbol    uint8 // the symbol held by the node
	frequency int64 // the symbol frequency
	rNode     *Node // the right child node
	lNode     *Node // the left child node
}

type NodeArray struct {
	nodes []Node
}

func (na *NodeArray) sort() {
	n := len(na.nodes) - 1
	for a := 0; a < n; a++ {
		for b := 0; b < n; b++ {
			if na.nodes[b+1].frequency < na.nodes[b].frequency {
				t := na.nodes[b+1]
				na.nodes[b+1] = na.nodes[b]
				na.nodes[b] = t
			}
		}
	}
}

// helper function to get []Node from the frequency table
func getNodes(ft map[uint8]int64) []Node {
	na := []Node{}
	for k, v := range ft {
		n := Node{
			symbol:    k,
			frequency: v,
			rNode:     nil,
			lNode:     nil,
		}
		na = append(na, n)
	}
	return na
}

// joins arr[0] and arr[1] into 1 node
func (na *NodeArray) join() {
	lNode := (*na).nodes[0]
	rNode := (*na).nodes[1]

	n := Node{
		symbol:    0x00,
		frequency: lNode.frequency + rNode.frequency,
		lNode:     &lNode,
		rNode:     &rNode,
	}

	// remove arr[0] and arr[1] and add the new node
	_nodeArray := []Node{}
	for a := 2; a < len(na.nodes); a++ {
		_nodeArray = append(_nodeArray, (*na).nodes[a])
	}
	// add the new node
	_nodeArray = append(_nodeArray, n)
	// set the value of na.nodes
	na.nodes = _nodeArray
	// sort the nodes
	na.sort()
}

// Todo: codes should not equate to 0xff
// Todo: codes should have a maximum length of 16 bits
// retunrs the corresponding code for a given symbol
func generateCodes(ft map[uint8]int64) map[uint8]uint32 {
	// get the nodes
	nodes := getNodes(ft)
	// create the node array
	na := &NodeArray{nodes: nodes}
	// sort the nodes
	na.sort()
	// create the huffman tree by joining 'len - 1' times
	n := len(na.nodes) - 1
	for a := 0; a < n; a++ {
		na.join()
	}
	codes := &map[uint8]uint32{}
	// the root node is na.nodes[0]
	rootNode := &(*na).nodes[0]
	// generate all the nodes using the trav function
	trav(rootNode, uint16(0), uint16(0), codes)
	return *codes
}

func trav(n *Node, code uint16, cLength uint16, mp *map[uint8]uint32) {
	if n.lNode == nil && n.rNode == nil {
		sym := n.symbol
		// 'encode' the uint32
		d := uint32(0x00000000)
		// upper 16 bits are for the code length
		d |= uint32(cLength)
		d <<= 16
		d |= uint32(code)
		// set the code
		(*mp)[sym] = d
	} else {
		trav(n.lNode, (code<<1)|0, cLength+1, mp)
		trav(n.rNode, (code<<1)|1, cLength+1, mp)
	}
}

func codeStr(cLength uint32, code uint32) string {
	cd := fmt.Sprintf("%b", int(code))
	rem := cLength - uint32(len(cd))
	for a := uint32(0); a < rem; a++ {
		cd = "0" + cd
	}
	return cd
}

// helper function to print all the codes
func printCodes(ct map[uint8]uint32) {
	for k, v := range ct {
		// print the symbol
		fmt.Printf("0x")
		fmt.Printf("%x", (k>>4)&0xf)
		fmt.Printf("%x", (k>>0)&0xf)
		fmt.Printf(" -> ")
		// print the code
		fmt.Printf("0x")
		fmt.Printf("%x", (v>>28)&0xf)
		fmt.Printf("%x", (v>>24)&0xf)
		fmt.Printf("%x", (v>>20)&0xf)
		fmt.Printf("%x", (v>>16)&0xf)
		fmt.Printf("%x", (v>>12)&0xf)
		fmt.Printf("%x", (v>>8)&0xf)
		fmt.Printf("%x", (v>>4)&0xf)
		fmt.Printf("%x", (v>>0)&0xf)
		fmt.Printf(" -> ")
		// print the bin code
		cLength := v >> uint32(16)
		code := v & uint32(0xffff)
		fmt.Printf("(%d) ", cLength)
		fmt.Printf("%s", codeStr(cLength, code))
		fmt.Printf("\n")
	}
}
