package main

import (
	"log"
	"sync"

	"github.com/bumbigbrain/nvp/internal/forwarder"
	"github.com/bumbigbrain/nvp/internal/tap"
	"github.com/bumbigbrain/nvp/internal/udp"
)

func main() {
	var wg sync.WaitGroup

	conn, err := udp.Connect("192.168.121.1:8080")
	if err != nil {
		log.Println("Error connecting to UDP server:", err)
	}

	ifce, err := tap.Setup("nvp-tap")
	if err != nil {
		log.Println("Error creating TAP interface:", err)
	}

	downFowarder := forwarder.NewDownForwarder(
		ifce,
		conn,
		&wg,
	)

	upFowarder := forwarder.NewUpForwarder(
		ifce,
		conn,
		&wg,
	)

	downFowarder.Run()
	upFowarder.Run()
	wg.Wait()

}
