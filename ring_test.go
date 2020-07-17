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

func TestKeyHash(t *testing.T) {
	var key, hash, expectedHash string
	
	key = "10.23.20.2"
	hash = keyHash(key)
	expectedHash = "d80ceccd"
	AssertEqual(t, hash, expectedHash, "")

	key = "10.23.34.4"
	hash = keyHash(key)
	expectedHash = "8eda8641"
	AssertEqual(t, hash, expectedHash, "")
}

func TestRingAddNode(t *testing.T) {
	ring := NewRing("", 1)

	ring.Add("10.23.20.2") // Hash Key - 0xd80ceccd
	ring.Add("10.23.34.4") // Hash Key - 0x8eda8641
	key1 := "TEST_KEY"     // Hash Key - 0x2269b0e
	key2 := "ANOTHER_KEY"  // Hash Key - 0xd3918bd2

	node := ring.GetPoint(key1)
	AssertEqual(t, node.value.ip, "10.23.34.4", "")

	node = ring.GetPoint(key2)
	AssertEqual(t, node.value.ip, "10.23.20.2", "")
}

func TestRingDeleteNode(t *testing.T) {
	ring := NewRing("", 1)

	ring.Add("10.23.20.2") // Hash Key - 0xd80ceccd
	ring.Add("10.23.34.4") // Hash Key - 0x8eda8641
	key1 := "TEST_KEY"     // Hash Key - 0x2269b0e
	key2 := "ANOTHER_KEY"  // Hash Key - 0xd3918bd2

	ring.Delete("10.23.20.2")
	node := ring.GetPoint(key1)
	AssertEqual(t, node.value.ip, "10.23.34.4", "")

	node = ring.GetPoint(key2)
	AssertEqual(t, node.value.ip, "10.23.34.4", "")
}

func TestRingInitFromConfig(t *testing.T) {
	ring := NewRing("./testconfig.conf", 1)
	nodes := ring.GetPoints()
	AssertEqual(t, nodes[0].index, "95412376", "")
	AssertEqual(t, nodes[1].index, "af102aa1", "")
}