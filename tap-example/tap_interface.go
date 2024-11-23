package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/songgao/water"
)

func main() {
	// Create TAP interface
	config := water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: "nvp-tap",
		},
	}

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer ifce.Close()

	fmt.Printf("Interface created: %s\n", ifce.Name())

	// Create a WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Create a channel to handle interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create a channel to signal goroutines to stop
	done := make(chan bool)

	// Start writer goroutine
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	// Sample packet (a simple ethernet frame)
	// 	packet := []byte{
	// 		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // Destination MAC (broadcast)
	// 		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Source MAC
	// 		0x08, 0x00, // Ethertype (IPv4)
	// 		// Payload
	// 		0x45, 0x00, 0x00, 0x1c, // IPv4 header
	// 		0x00, 0x00, 0x40, 0x00,
	// 		0x40, 0x11, 0x00, 0x00,
	// 		0x0a, 0x00, 0x00, 0x01, // Source IP (10.0.0.1)
	// 		0x0a, 0x00, 0x00, 0x02, // Destination IP (10.0.0.2)
	// 	}

	// 	for {
	// 		select {
	// 		case <-done:
	// 			fmt.Println("Writer stopping...")
	// 			return
	// 		default:
	// 			n, err := ifce.Write(packet)
	// 			if err != nil {
	// 				log.Printf("Error writing to interface: %v\n", err)
	// 				continue
	// 			}
	// 			fmt.Printf("Wrote %d bytes to %s\n", n, ifce.Name())
	// 			time.Sleep(time.Second) // Add delay to avoid flooding
	// 		}
	// 	}
	// }()

	// Start reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		packet := make([]byte, 2048)
		for {
			select {
			case <-done:
				fmt.Println("Reader stopping...")
				return
			default:
				n, err := ifce.Read(packet)
				if err != nil {
					log.Printf("Error reading from interface: %v\n", err)
					continue
				}

				fmt.Printf("\nReceived %d bytes from %s:\n", n, ifce.Name())
				fmt.Printf("Raw packet (hex):\n%s\n", hex.Dump(packet[:n]))

				// Basic Ethernet frame parsing
				if n >= 14 { // Minimum Ethernet frame size
					dst := packet[0:6]
					src := packet[6:12]
					etherType := packet[12:14]

					fmt.Printf("Ethernet Frame:\n")
					fmt.Printf("  Destination MAC: %02x:%02x:%02x:%02x:%02x:%02x\n",
						dst[0], dst[1], dst[2], dst[3], dst[4], dst[5])
					fmt.Printf("  Source MAC: %02x:%02x:%02x:%02x:%02x:%02x\n",
						src[0], src[1], src[2], src[3], src[4], src[5])
					fmt.Printf("  EtherType: 0x%02x%02x\n", etherType[0], etherType[1])
				}
				fmt.Println("----------------------------------------")
			}
		}
	}()

	// Wait for interrupt signal
	<-interrupt
	fmt.Println("\nReceived interrupt signal. Shutting down...")

	// Signal goroutines to stop
	close(done)

	// Wait for goroutines to finish
	wg.Wait()
	fmt.Println("Program terminated gracefully")
}
