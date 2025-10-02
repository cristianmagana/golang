package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*
192.168.1.100 - - [30/Sep/2025:10:23:45 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
10.0.0.15 - - [30/Sep/2025:10:23:46 +0000] "POST /api/login HTTP/1.1" 401 512 "-" "curl/7.68.0"
192.168.1.100 - - [30/Sep/2025:10:23:47 +0000] "GET /api/data HTTP/1.1" 200 5678 "-" "Mozilla/5.0"
172.16.0.5 - admin [30/Sep/2025:10:23:48 +0000] "GET /admin/dashboard HTTP/1.1" 200 9876 "-" "Chrome/91.0"
*/

type AccessEntry struct {
	IP     string
	Status string
}

func parseLog(line string) *AccessEntry {
	parts := strings.Split(line, " ")
	//fmt.Printf("IP: %s, Status: %s\n", parts[0], parts[8])

	return &AccessEntry{
		IP:     parts[0],
		Status: parts[8],
	}
}

func calculateIps(accessLogs map[string]int) {
	type kv struct {
		Key   string
		Value int
	}

	var sortedMap []kv

	for k, v := range accessLogs {
		sortedMap = append(sortedMap, kv{k, v})
	}

	sort.Slice(sortedMap, func(i, j int) bool {
		return sortedMap[i].Value < sortedMap[j].Value
	})

	for _, kv := range sortedMap {
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
	}

}

func CalculateLegacy() {
	path := "./access.log"
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	accessLogs := make(map[string]int)

	for scanner.Scan() {
		entry := parseLog(scanner.Text())
		value, ok := accessLogs[entry.IP]
		if ok {
			accessLogs[entry.IP] = value + 1
		} else {
			accessLogs[entry.IP] = 1
		}
	}

	calculateIps(accessLogs)
}
