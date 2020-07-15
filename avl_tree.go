/*
 * Copyright (c) 2020, Jake Grogan
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 *  * Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from
 *    this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"math"
)

type Params struct {
	node   *Node
	output []string
}

type AVLTree struct {
	node *Node
	height int
	balance int
}

func NewAvlTree() *AVLTree {
	return &AVLTree{
		node: nil,
		height: -1,
		balance: 0,
	}
}

func (this *AVLTree) InsertNode(index string, vp *VirtualPoint) {
	var node *Node = NewNode(index, vp)

	if (this.node == nil) {
		this.node = node
		this.node.left = NewAvlTree()
		this.node.right = NewAvlTree()
	} else if (index < this.node.index) {
		this.node.left.InsertNode(index, vp)
	} else if (index > this.node.index) {
		this.node.right.InsertNode(index, vp)
	}
	this.rebalance()
}

func (this *AVLTree) RemoveNode(index string) {
	if (this.node != nil) {
		if (index == this.node.index) {
			if (this.node.left.node == nil && this.node.right.node == nil) {
				this.node = nil
			} else if (this.node.left.node == nil) {
				this.node = this.node.right.node
			} else if (this.node.right.node == nil) {
				this.node = this.node.left.node
			} else {
				var successor *Node = this.node.right.node
				for successor != nil && successor.left.node != nil {
					successor = successor.left.node
				} 
				if successor != nil {
					this.node.index = successor.index
					this.node.right.RemoveNode(successor.index)
				}
			}
		} else if (index < this.node.index) {
			this.node.left.RemoveNode(index)
		} else if (index > this.node.index) {
			this.node.right.RemoveNode(index)
		}
		this.rebalance()
	}
}

func (this *AVLTree) InOrderTraverse() []string {
	var root *Node = this.node
	var output []string = this.inOrder(&Params{node: root, output: nil})
	return output
}

func (this *AVLTree) inOrder(params *Params) []string {
	if (params.node != nil) {
		var old *Node = params.node

		params.node = old.left.node
		this.inOrder(params)

		params.output = append(params.output, old.index)

		params.node = old.right.node
		this.inOrder(params)
	}
	return params.output
}

func (this *AVLTree) PreOrderTraverse() []string {
	var root *Node = this.node
	var output []string = this.preOrder(&Params{node: root, output: nil})
	return output
}

func (this *AVLTree) preOrder(params *Params) []string {
	if (params.node != nil) {
		var old *Node = params.node

		params.output = append(params.output, old.index)
		
		params.node = old.left.node
		this.preOrder(params)

		params.node = old.right.node
		this.preOrder(params)
	}
	return params.output
}

func (this *AVLTree) PostOrderTraverse() []string {
	var root *Node = this.node
	var output []string = this.postOrder(&Params{node: root, output: nil})
	return output
}

func (this *AVLTree) postOrder(params *Params) []string {
	if (params.node != nil) {
		var old *Node = params.node
		
		params.node = old.left.node
		this.postOrder(params)

		params.node = old.right.node
		this.postOrder(params)

		params.output = append(params.output, old.index)
	}
	return params.output
}

func (this *AVLTree) rebalance() {
	this.updateHeights()
	this.updateBalances()

	for this.balance < -1 || this.balance > 1 {
		if (this.balance > 1) {
			if (this.node.left.balance < 0) {
				this.node.left.rotateLeft()
				this.updateHeights()
				this.updateBalances()
			}
			this.rotateRight()
			this.updateHeights()
			this.updateBalances()
		}
		
		if (this.balance < -1) {
			if (this.node.right.balance > 0) {
				this.node.right.rotateRight()
				this.updateHeights()
				this.updateBalances()
			}
			this.rotateLeft()
			this.updateHeights()
			this.updateBalances()
		}
	}
}

func (this *AVLTree) updateHeights() {
	if (this.node != nil) {
		if (this.node.left != nil) {
			this.node.left.updateHeights()
		}
		if (this.node.right != nil) {
			this.node.right.updateHeights()
		}
		this.height = 1 + int(math.Max(float64(this.node.left.height), float64(this.node.right.height)))
	} else {
		this.height = -1
	}
}

func (this *AVLTree) updateBalances() {
	if (this.node != nil) {
		if (this.node.left != nil) {
			this.node.left.updateBalances()
		}
		if (this.node.right != nil) {
			this.node.right.updateBalances()
		}
		this.balance = this.node.left.height - this.node.right.height
	} else {
		this.balance = 0
	}
}

func (this *AVLTree) rotateRight() {
	var newRoot *Node = this.node.left.node
	var newLeftSub *Node = newRoot.right.node
	var oldRoot *Node = this.node

	this.node = newRoot
	oldRoot.left.node = newLeftSub
	newRoot.right.node = oldRoot
}

func (this *AVLTree) rotateLeft() {
	var newRoot *Node = this.node.right.node
	var newRightSub *Node = newRoot.left.node
	var oldRoot *Node = this.node;

	this.node = newRoot
	oldRoot.right.node = newRightSub
	newRoot.left.node = oldRoot
}