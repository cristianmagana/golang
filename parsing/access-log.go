package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

func readWholeFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file))
}

func readFileLineByLine(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}
}

func readFileInBufferStreamChunks(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buffer[:n]))
	}
}

func readFileWithChannels(path string) (<-chan string, <-chan error) {
	lines := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(lines)
		file, err := os.Open(path)
		if err != nil {
			errc <- err
			close(errc)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err == scanner.Err() && err != nil {
			errc <- err
		}
		close(errc)
	}()

	return lines, errc
}

type Log struct {
	IP      string `json:"ip"`
	Level   string `json:"level"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (l *Log) parseJson(path string) {
	lines, errc := readFileWithChannels(path)

	for line := range lines {
		var entry Log
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			if entry.Status == "500" || entry.Level == "error" {
				fmt.Println("Error log:", entry)
			}
		} else {
			fmt.Println("Skipping invalid JSON:", line)
		}
	}

	if err := <-errc; err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func parseWebLogs(path string) {
	lines, errc := readFileWithChannels(path)

	for line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 8 {
			status := parts[8]
			if status == "500" {
				fmt.Println("500 error line:", line)
			}
		}
	}
	if err := <-errc; err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func parseWebLogsWithRegex(path string) {
	lines, errc := readFileWithChannels(path)

	var serverErrorRegex = regexp.MustCompile((`"\s5\d{2}\s`))
	for line := range lines {
		if serverErrorRegex.MatchString(line) {
			fmt.Println("5xx error line:", line)
		}

	}
	if err := <-errc; err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func AccessLogs() {
	fmt.Printf("Staring by parsing in buffer chunks\n")
	fmt.Println("")
	start := time.Now()
	path := "./access.log"
	readFileInBufferStreamChunks(path)
	fmt.Printf("\nTime elapsed: %d\n", time.Since(start).Microseconds())

	fmt.Printf("Staring by parsing in channels\n")
	fmt.Println("")
	start2 := time.Now()
	lines, errc := readFileWithChannels(path)

	for line := range lines {
		fmt.Println(line)
	}
	if err := <-errc; err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("\nTime elapsed: %d\n", time.Since(start2).Microseconds())

	fmt.Printf("Staring by parsing in json struct in channels\n")
	fmt.Println("")
	start3 := time.Now()
	var log Log
	log.parseJson("./access-logs.json")

	fmt.Printf("\nTime elapsed: %d\n", time.Since(start3).Microseconds())

	fmt.Printf("Staring by parsing in json struct in channels\n")
	fmt.Println("")
	start4 := time.Now()
	parseWebLogs(path)

	fmt.Printf("\nTime elapsed: %d\n", time.Since(start4).Microseconds())
}
