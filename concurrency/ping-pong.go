package main

import (
	"fmt"
)

func pingPlayer(ping <-chan string, pong chan<- string, done chan<- bool) {
	roundCount := 0
	maxRounds := 5

	for roundCount < maxRounds {
		msg := <-ping

		fmt.Printf("Received %s for round %d\n", msg, roundCount)

		roundCount++
		if roundCount <= maxRounds {
			pong <- "pong"
		}
	}
	done <- true
}

func pongPlayer(pong <-chan string, ping chan<- string) {
	for {
		msg, ok := <-pong
		if !ok {
			return
		}
		fmt.Printf("Received %s\n", msg)
		ping <- "ping"
	}

}

func PingPong() {
	ping, pong := make(chan string), make(chan string)

	done := make(chan bool)

	go pingPlayer(ping, pong, done)
	go pongPlayer(pong, ping)

	ping <- "ping"

	<-done

	close(pong)

}
