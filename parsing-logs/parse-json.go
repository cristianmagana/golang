package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

/*
{"timestamp":"2025-09-30T10:45:12.123Z","level":"INFO","service":"api-gateway","client_ip":"10.5.5.50","user_id":"user123","endpoint":"/api/v1/orders","method":"GET","status":200,"response_time_ms":45}
{"timestamp":"2025-09-30T10:45:13.456Z","level":"WARN","service":"api-gateway","client_ip":"192.168.50.100","user_id":"user456","endpoint":"/api/v1/users","method":"POST","status":429,"response_time_ms":12,"error":"rate_limit_exceeded"}
{"timestamp":"2025-09-30T10:45:14.789Z","level":"INFO","service":"api-gateway","client_ip":"10.5.5.50","user_id":"user123","endpoint":"/api/v1/products","method":"GET","status":200,"response_time_ms":67}
{"timestamp":"2025-09-30T10:45:15.012Z","level":"ERROR","service":"api-gateway","client_ip":"172.20.10.75","user_id":"user789","endpoint":"/api/v1/checkout","method":"POST","status":500,"response_time_ms":5000,"error":"database_timeout"}
{"timestamp":"2025-09-30T10:45:16.345Z","level":"INFO","service":"api-gateway","client_ip":"192.168.50.100","user_id":"user456","endpoint":"/api/v1/users","method":"POST","status":201,"response_time_ms":89}
{"timestamp":"2025-09-30T10:45:17.678Z","level":"INFO","service":"api-gateway","client_ip":"10.5.5.50","user_id":"user123","endpoint":"/api/v1/cart","method":"PUT","status":200,"response_time_ms":34}
*/

type JsonEntry struct {
	Timestamp    string `json:"timestamp"`
	Level        string `json:"level"`
	Service      string `json:"service"`
	ClientIP     string `json:"client_ip"`
	UserID       string `json:"user_id"`
	Endpoint     string `json:"endpoint"`
	Method       string `json:"method"`
	Status       int    `json:"status"`
	ResponseTime int    `json:"response_time_ms"`
}

func ParseJson() {
	path := "./access.log"

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error reading the file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	logs := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Bytes()
		var log JsonEntry
		err := json.Unmarshal(line, &log)
		if err != nil {
			fmt.Printf("error reading from json %v\n", err)
		}
		v, ok := logs[log.ClientIP]
		if ok {
			logs[log.ClientIP] = v + 1
		} else {
			logs[log.ClientIP] = 1
		}
	}

	// calculate ips
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv

	for k, v := range logs {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	for _, kv := range ss {
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
	}

}
