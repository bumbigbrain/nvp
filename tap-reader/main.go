package main

import (
	"fmt"
	"log"

	"github.com/songgao/packets/ethernet"
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

	var frame ethernet.Frame

	// Continuously read packets

	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Fatal(err)
		}

		frame = frame[:n]
		// Process the packet
		fmt.Printf("Received packet of length %d bytes\n", n)

		// Print packet details in hex format
		log.Printf("Entire Frame: %v\n", frame)

		log.Printf("Dst: %s\n", frame.Destination())
		log.Printf("Src: %s\n", frame.Source())
		log.Printf("Ethertype: % x\n", frame.Ethertype())
		log.Printf("Payload: % x\n", frame.Payload())

		// You can add more packet processing logic here
		// For example, parsing Ethernet frames, IP headers, etc.
	}
}
