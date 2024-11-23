package main

import (
	"fmt"
	"log"

	"github.com/songgao/water"
)

func main() {
	// Configure TAP interface
	config := water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: "nvp-tap",
		},
	}

	// Create a new TAP interface
	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// Close the interface when the program exits
	defer ifce.Close()

	fmt.Printf("Interface %s created. Starting packet capture...\n", ifce.Name())

	// Buffer for reading packets
	packet := make([]byte, 2048)

	// Continuously read packets
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}

		// Process the packet
		fmt.Printf("Received packet of length %d bytes\n", n)

		// Print packet details in hex format
		fmt.Printf("Packet content: %x\n", packet[:n])

		// You can add more packet processing logic here
		// For example, parsing Ethernet frames, IP headers, etc.
	}
}
