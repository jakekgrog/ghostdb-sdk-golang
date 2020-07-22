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
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	NO_MORE_SERVERS_ERROR = "All nodes marked as dead: Failed to establish a connection to any servers: (Check your fleet status)"
)

const (
	HTTP = "http://"
	HTTPS = "https://"
)

type Cache struct {
	deadServers    map[string]bool
	configFilepath string
	ring           *Ring
	protocol       string
	port           string
	reviveInterval int32
}

func NewCache(configFilepath string, http bool, port string) *Cache {
	var protocol string
	if http {
		protocol = HTTP
	} else {
		protocol = HTTPS
	}

	cache := &Cache{
		deadServers: make(map[string]bool),
		configFilepath: configFilepath,
		ring: NewRing(configFilepath, 1),
		protocol: protocol,
		port: port,
		reviveInterval: 30,
	}

	go startServerRevival(cache)
	return cache
}

func (this *Cache) Get(key string) (CacheResponse, error) {
	node := this.ring.GetPoint(key)
	if node == nil {
		return CacheResponse{}, errors.New(NO_MORE_SERVERS_ERROR)
	}

	serviceRequestParams := CacheRequestParams{
		Key: key,
		Value: "",
		TTL: -1,
	}

	response, err := this.makeServiceRequest("get", node, serviceRequestParams)
	if err != nil {
		this.markDead(node)
		return this.Get(key)
	}
	return response, nil
}

func (this *Cache) NodeSize(ip string) (CacheResponse, error) {
	node := this.ring.GetPoint(ip)
	if node == nil {
		return CacheResponse{}, errors.New(NO_MORE_SERVERS_ERROR)
	}

	serviceRequestParams := CacheRequestParams{
		Key: "",
		Value: "",
		TTL: -1,
	}

	response, err := this.makeServiceRequest("getNodeSize", node, serviceRequestParams)
	if err != nil {
		this.markDead(node)
		return this.NodeSize(ip)
	}
	return response, nil
}

func (this *Cache) Add(key string, value interface{}, ttl int) (CacheResponse, error) {
	node := this.ring.GetPoint(key)
	if node == nil {
		return CacheResponse{}, errors.New(NO_MORE_SERVERS_ERROR)
	}

	serviceRequestParams := CacheRequestParams{
		Key: key,
		Value: value,
		TTL: ttl,
	}

	response, err := this.makeServiceRequest("add", node, serviceRequestParams)
	if err != nil {
		this.markDead(node)
		return this.Add(key, value, ttl)
	}
	return response, nil
}

func (this *Cache) Put(key string, value interface{}, ttl int) (CacheResponse, error) {
	node := this.ring.GetPoint(key)
	if node == nil {
		return CacheResponse{}, errors.New(NO_MORE_SERVERS_ERROR)
	}

	serviceRequestParams := CacheRequestParams{
		Key: key,
		Value: value,
		TTL: ttl,
	}

	response, err := this.makeServiceRequest("put", node, serviceRequestParams)
	if err != nil {
		this.markDead(node)
		return this.Put(key, value, ttl)
	}
	return response, nil
}

func (this *Cache) Delete(key string) (CacheResponse, error) {
	node := this.ring.GetPoint(key)
	if node == nil {
		return CacheResponse{}, errors.New(NO_MORE_SERVERS_ERROR)
	}

	serviceRequestParams := CacheRequestParams{
		Key: key,
		Value: "",
		TTL: -1,
	}

	response, err := this.makeServiceRequest("delete", node, serviceRequestParams)
	if err != nil {
		this.markDead(node)
		return this.Delete(key)
	}
	return response, nil
}

func (this *Cache) Flush() (bool, error) {
	var nodes []*VirtualPoint = this.ring.GetPoints()
	for _, vp := range nodes {
		serviceRequestParams := CacheRequestParams{
			Key: "",
			Value: "",
			TTL: -1,
		}
		pairNode := &Pair{index: vp.index, value: vp}
		_, err := this.makeServiceRequest("flush", pairNode, serviceRequestParams)
		if err != nil {
			this.markDead(pairNode)
			return this.Flush()
		}
	}
	return true, nil
}

func (this *Cache) GetSysMetrics() []*Metric {
	res := this.recGetSysMetrics([]*Metric{}, []string{})
	return res
}

func (this *Cache) recGetSysMetrics(metrics []*Metric, visitedNodes []string) ([]*Metric) {
	var nodes []*VirtualPoint = this.ring.GetPoints()
	for _, vp := range nodes {
		if ok := exists(visitedNodes, vp.ip); !ok {
			serviceRequestParams := CacheRequestParams{
				Key: "",
				Value: "",
				TTL: -1,
			}
			pairNode := &Pair{index: vp.index, value: vp}
			resp, err := this.makeServiceRequest("getSysMetrics", pairNode, serviceRequestParams)
			if err != nil {
				this.markDead(pairNode)
				return this.recGetSysMetrics(metrics, visitedNodes)
			}
			metrics = append(metrics, &Metric{node: pairNode.value.ip, metrics: resp})
			visitedNodes = append(visitedNodes, pairNode.value.ip)
		}
	}
	return metrics
}

func (this *Cache) GetAppMetrics() []*Metric {
	res := this.recGetAppMetrics([]*Metric{}, []string{})
	return res
}

func (this *Cache) recGetAppMetrics(metrics []*Metric, visitedNodes []string) ([]*Metric) {
	var nodes []*VirtualPoint = this.ring.GetPoints()
	for _, vp := range nodes {
		if ok := exists(visitedNodes, vp.ip); !ok {
			serviceRequestParams := CacheRequestParams{
				Key: "",
				Value: "",
				TTL: -1,
			}
			pairNode := &Pair{index: vp.index, value: vp}
			resp, err := this.makeServiceRequest("getAppMetrics", pairNode, serviceRequestParams)
			if err != nil {
				this.markDead(pairNode)
				return this.recGetAppMetrics(metrics, visitedNodes)
			}
			metrics = append(metrics, &Metric{node: pairNode.value.ip, metrics: resp})
			visitedNodes = append(visitedNodes, pairNode.value.ip)
		}
	}
	return metrics
}

func (this *Cache) Ping() []*Metric {
	res := this.recPing([]*Metric{}, []string{})
	return res
}

func (this *Cache) recPing(metrics []*Metric, visitedNodes []string) ([]*Metric) {
	var nodes []*VirtualPoint = this.ring.GetPoints()
	for _, vp := range nodes {
		if ok := exists(visitedNodes, vp.ip); !ok {
			serviceRequestParams := CacheRequestParams{
				Key: "",
				Value: "",
				TTL: -1,
			}
			pairNode := &Pair{index: vp.index, value: vp}
			resp, err := this.makeServiceRequest("ping", pairNode, serviceRequestParams)
			if err != nil {
				this.markDead(pairNode)
				return this.recPing(metrics, visitedNodes)
			}
			metrics = append(metrics, &Metric{node: pairNode.value.ip, metrics: resp})
			visitedNodes = append(visitedNodes, pairNode.value.ip)
		}
	}
	return metrics
}

func startServerRevival(cache *Cache) {
	interval := time.Duration(cache.reviveInterval) * time.Second
	ticker := time.NewTicker(interval)

	for {
		select {
			case <- ticker.C:
				go attemptRevive(cache)
		}
	}
}

func attemptRevive(cache *Cache) {
	serviceRequestParams := CacheRequestParams{
		Key: "",
		Value: "",
		TTL: -1,
	}
	requestObj := NewCacheRequest(serviceRequestParams)
	requestBody, _ := json.Marshal(requestObj)
	for server, _ := range cache.deadServers {
		url := cache.protocol + server + ":" + cache.port + "/ping"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 {
			cache.ring.Add(server)
			delete(cache.deadServers, server)
		}
	}
}

func exists(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (this *Cache) makeServiceRequest(requestType string, server *Pair, params CacheRequestParams) (CacheResponse, error) {
	requestObj := NewCacheRequest(params)
	requestBody, err := json.Marshal(requestObj)

	url := this.protocol + server.value.ip + ":" + this.port + getRequestType(requestType)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return CacheResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return CacheResponse{}, err
	}

	var responseObj CacheResponse
	err = json.Unmarshal(body, &responseObj) 
	if err != nil {
		return CacheResponse{}, err
	}
	return responseObj, nil
}

func (this *Cache) markDead(server *Pair) {
	this.deadServers[server.value.ip] = true
	this.ring.Delete(server.value.ip)
}

func getRequestType(requestType string) string {
	endpoints := map[string]string {
		"ping": "/ping",
		"put": "/put",
		"get": "/get",
		"add": "/add",
		"delete": "/delete",
		"flush": "/flush",
		"getSysMetrics": "/getSysMetrics",
		"getAppMetrics": "/getAppMetrics",
		"getNodeSize":  "/nodeSize",
	}

	return endpoints[requestType]
}