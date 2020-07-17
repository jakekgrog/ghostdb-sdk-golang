package main

type Params struct {
	node   *Node
	output []string
}

type GetVpParams struct {
	node   *Node
	output []*VirtualPoint
}

type Pair struct {
	index string
	value *VirtualPoint
}