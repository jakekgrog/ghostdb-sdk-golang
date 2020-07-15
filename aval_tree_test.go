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
	"testing"
)

func TestAvlTree(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	
	// Test insert correctly
	tree.InsertNode("2", vp2)
	AssertEqual(t, tree.node.index, "2", "")
	AssertEqual(t, tree.node.vp.ip, "127.0.0.2", "")
	AssertEqual(t, tree.node.vp.index, "2",  "")
	AssertEqual(t, tree.node.left.node == nil, true, "")
	AssertEqual(t, tree.node.right.node == nil, true, "")

	tree.InsertNode("1", vp1)
	AssertEqual(t, tree.node.left.node.index, "1", "")
}

func TestAvlTreeRebalanceLeftRotation(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	var vp3 *VirtualPoint = NewVirtualPoint("127.0.0.3", "3")

	// When initially inserting, "3" should be the root node
	// After inserting all three entries, "2" should be the root node
	tree.InsertNode("3", vp3)
	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)

	// Assert the root node is "2"
	AssertEqual(t, tree.node.index, "2", "")
}

func TestAvlTreeRebalanceRightRotation(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	var vp3 *VirtualPoint = NewVirtualPoint("127.0.0.3", "3")
	var vp4 *VirtualPoint = NewVirtualPoint("127.0.0.4", "4")
	var vp5 *VirtualPoint = NewVirtualPoint("127.0.0.5", "5")

	tree.InsertNode("3", vp3)
	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.InsertNode("4", vp4)
	tree.InsertNode("5", vp5)

	AssertEqual(t, tree.node.right.node.index, "4", "")
}

func TestAvlTreeRemoveNode(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")

	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.RemoveNode("1")

	AssertEqual(t, tree.node.index, "2", "")
}

func TestAvlTreeRemoveRootNode(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")

	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.RemoveNode("2")

	AssertEqual(t, tree.node.index, "1", "")
}

func TestAvlTreeRemoveNodeWithTwoChildren(t *testing.T) {
    //        2               3
	//       / \             / \
	//      1   4   ---->   1   4  
	//         /
	//        3
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	var vp3 *VirtualPoint = NewVirtualPoint("127.0.0.3", "3")
	var vp4 *VirtualPoint = NewVirtualPoint("127.0.0.4", "4")

	tree.InsertNode("4", vp4)
	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.InsertNode("3", vp3)
	
	tree.RemoveNode("2")

	AssertEqual(t, tree.node.index, "3", "")
}

func TestAvlTreeInOrderTraverse(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	var vp3 *VirtualPoint = NewVirtualPoint("127.0.0.3", "3")
	var vp4 *VirtualPoint = NewVirtualPoint("127.0.0.4", "4")

	tree.InsertNode("4", vp4)
	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.InsertNode("3", vp3)
	
	var output []string = tree.InOrderTraverse()
	var expectedOutput []string = []string{"1", "2", "3", "4"}
	AssertEqual(t, expectedOutput[0], output[0], "")
	AssertEqual(t, expectedOutput[1], output[1], "")
	AssertEqual(t, expectedOutput[2], output[2], "")
	AssertEqual(t, expectedOutput[3], output[3], "")
}

func TestAvlTreePreOrderTraverse(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	var vp3 *VirtualPoint = NewVirtualPoint("127.0.0.3", "3")
	var vp4 *VirtualPoint = NewVirtualPoint("127.0.0.4", "4")

	tree.InsertNode("4", vp4)
	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.InsertNode("3", vp3)
	
	var output []string = tree.PreOrderTraverse()
	var expectedOutput []string = []string{"2", "1", "3", "4"}
	AssertEqual(t, expectedOutput[0], output[0], "")
}

func TestAvlTreePostOrderTraverse(t *testing.T) {
	var tree *AVLTree = NewAvlTree()
	var vp1 *VirtualPoint = NewVirtualPoint("127.0.0.1", "1")
	var vp2 *VirtualPoint = NewVirtualPoint("127.0.0.2", "2")
	var vp3 *VirtualPoint = NewVirtualPoint("127.0.0.3", "3")
	var vp4 *VirtualPoint = NewVirtualPoint("127.0.0.4", "4")

	tree.InsertNode("4", vp4)
	tree.InsertNode("2", vp2)
	tree.InsertNode("1", vp1)
	tree.InsertNode("3", vp3)
	
	var output []string = tree.PostOrderTraverse()
	var expectedOutput []string = []string{"1", "3", "4", "2"}
	AssertEqual(t, expectedOutput[0], output[0], "")
}