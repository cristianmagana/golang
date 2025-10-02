package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

/*
{"timestamp":"2025-09-30T10:45:12.123Z",
"level":"INFO","service":"api-gateway",
"client_ip":"10.5.5.50",
"user_id":"user123",
"endpoint":"/api/v1/orders",
"method":"GET",
"status":200,
"response_time_ms":45}
*/

type AccessLogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	ClientIP  string `json:"client_ip"`
	Status    int    `json:"status"`
}

func parseJsonEntry(scanner *bufio.Scanner) *AccessLogEntry {
	var entry AccessLogEntry
	line := scanner.Bytes()
	err := json.Unmarshal(line, &entry)
	if err != nil {
		fmt.Printf("error reading from json %v\n", err)
	}

	return &entry
}

func calculateMostIps(logMap map[string]int) {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range logMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	for _, kv := range ss {
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
	}
}

func ParsingPractice() {
	logChan := make(chan []byte, 10)
	errChan := make(chan error, 1)
	path := "./access.log"
	file, err := os.Open(path)
	if err != nil {
		//channels
		errChan <- err
		// fmt.Printf("error reading from file: %s\n", path)
	}
	defer file.Close()

	logMap := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry AccessLogEntry
		line := scanner.Bytes()
		logChan <- line
		// entry := parseJsonEntry(scanner)
		// logMap[entry.ClientIP]++
	}

	calculateMostIps(logMap)
}
