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

func TestCache(t *testing.T) {
	cache := NewCache("simulation.conf", true, "7991")

	// TEST CORRECT RESPONSE WHEN CACHE EMPTY
	response, err := cache.Get("Test")
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, response.Message, "CACHE_MISS", "")
	}

	// TEST THERE ARE NO KEY/VALUE PAIRS IN THE CACHE
	response, err = cache.NodeSize("127.0.0.1")
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, response.Gobj.Value, float64(0), "")
	}

	// TEST ADD
	response, err = cache.Add("Ireland", "Dublin", -1)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, response.Status, int32(1), "")
	}

	// TEST ADD WORKED
	response, err = cache.Get("Ireland")
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, response.Gobj.Value, "Dublin", "")
	}

	// TEST PUT
	response, err = cache.Put("Ireland", "Not Dublin", -1)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, response.Status, int32(1), "")
	}

	// TEST PUT WORKED
	response, err = cache.Get("Ireland")
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, response.Gobj.Value, "Not Dublin", "")
	}

	// TEST GET SYS METRICS - ONLY REACHABLE SERVERS IN DATA
	data := cache.GetSysMetrics()
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, data[0].node, "127.0.0.1", "")
	}

	// TEST GET APP METRICS - ONLY REACHABLE SERVERS IN DATA
	data = cache.GetAppMetrics()
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, data[0].node, "127.0.0.1", "")
	}

	// TEST Ping - ONLY REACHABLE SERVERS IN DATA
	data = cache.Ping()
	if err != nil {
		t.Fatal(err.Error())
	} else {
		AssertEqual(t, data[0].node, "127.0.0.1", "")
	}
}