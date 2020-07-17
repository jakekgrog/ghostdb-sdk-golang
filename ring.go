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
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
)

const (
	EMPTY_CONFIG_ERR = "Cluster configuration file is empty!"
)

type Ring struct {
	replicas int
	ring     *AVLTree
}

func NewRing(clusterConfig string, replicas int) *Ring {
	var ring *Ring = &Ring{
		replicas: replicas,
		ring: NewAvlTree(),
	}
	ring.initRing(clusterConfig)
	return ring
}

func (this *Ring) Add(node string) {
	for i := 0; i < this.replicas; i++ {
		var index string = keyHash(node, i)
		var vp *VirtualPoint = NewVirtualPoint(node, index)
		this.ring.InsertNode(index, vp)
	}
}

func (this *Ring) Delete(node string) {
	var index string
	for i := 0; i < this.replicas; i++ {
		index = keyHash(node, i)
		this.ring.RemoveNode(index)
	}
}

func (this *Ring) initRing(clusterConfig string) {
	nodes, err := readFileByLine(clusterConfig)
	if err != nil {
		log.Fatalf("Failed to read from cluster configuration: %s", err.Error())
	}

	if len(nodes) != 0 {
		for _, node := range nodes {
			this.Add(node)
		}
	} else {
		log.Fatal(EMPTY_CONFIG_ERR)
	}
}

func keyHash(key string, index ...int) string {
	var keyToHash string

	if len(index) > 0 {
		keyToHash = fmt.Sprintf("%s:%d", key, index[0])
	} else {
		keyToHash = key
	}
	crc32Uint32 := crc32.ChecksumIEEE([]byte(keyToHash))
	crc32String := strconv.FormatUint(uint64(crc32Uint32), 16)
	return crc32String
}
