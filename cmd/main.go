package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/songgao/water"
)

type Message struct {
	IsInitialized bool   `json:"isInitialized"`
	SourceMacAddr string `json:"sourceMacAddr"`
}

func getInterfaceInfo(name string) (string, string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return "", "", fmt.Errorf("failed to get interface: %v", err)
	}

	// Get MAC address
	macAddr := iface.HardwareAddr.String()

	// Get IP address
	addrs, err := iface.Addrs()
	if err != nil {
		return "", "", fmt.Errorf("failed to get addresses: %v", err)
	}

	var ipAddr string
	for _, addr := range addrs {
		// Check if address is IP network
		if ipnet, ok := addr.(*net.IPNet); ok {
			// Use only IPv4 address
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				ipAddr = ip4.String()
				break
			}
		}
	}

	if ipAddr == "" {
		return "", "", fmt.Errorf("no IPv4 address found")
	}

	return ipAddr, macAddr, nil
}

func main() {

	var wg sync.WaitGroup

	serverAddr := "192.168.122.1:8080" // Using default UDP port 53, adjust if needed

	// Get interface information
	_, macAddr, err := getInterfaceInfo("nvp-tap")
	if err != nil {
		fmt.Printf("Error getting interface info: %v\n", err)
		os.Exit(1)
	}

	// Create message
	msg := Message{
		IsInitialized: true,
		SourceMacAddr: macAddr,
	}

	// Marshal message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Resolve UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		os.Exit(1)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to UDP server at %s\n", serverAddr)

	// Send the message
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sent message: %s\n", string(jsonData))

	// Create a buffer for receiving data
	buffer := make([]byte, 1024)

	// Read response from server
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
		return
	}

	fmt.Printf("Received %d bytes from %v: %s\n", n, remoteAddr, string(buffer[:n]))

	// Start packet reader and UDP sender

	// Create TAP interface
	config := water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: "nvp-tap",
		},
	}
	tapInterface, err := water.New(config)
	if err != nil {
		log.Fatalf("Failed to create TAP interface: %v", err)
	}
	defer tapInterface.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 2048) // Buffer for reading packets
		for {
			n, err := tapInterface.Read(buffer)
			if err != nil {
				log.Printf("Error reading from TAP interface: %v", err)
				continue
			}

			// Send packet to UDP server
			_, err = conn.WriteToUDP(buffer[:n], udpAddr)
			if err != nil {
				log.Printf("Error sending packet to UDP server: %v", err)
			}
		}
	}()
	wg.Wait()
}
