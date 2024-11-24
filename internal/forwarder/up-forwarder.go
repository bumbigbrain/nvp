package forwarder

import (
	"log"
	"net"
	"sync"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

type UpForwarder struct {
	Ifce     *water.Interface
	UdpConn  *net.UDPConn
	wgGlobal *sync.WaitGroup
}

func (f *UpForwarder) Run() {
	f.wgGlobal.Add(1)

	go func() {
		defer f.wgGlobal.Done()
		var frame ethernet.Frame
		buffer := make([]byte, 1500)

		for {
			n, err := f.Ifce.Read(buffer)
			if err != nil {
				log.Printf("Error reading from TAP interface: %v\n", err)
				continue
			}

			frame = ethernet.Frame(buffer[:n])
			if n > 0 {
				log.Printf("TAP -> UDP: packet length=%d, src=%s, dst=%s, type=0x%x\n",
					n, frame.Source(), frame.Destination(), frame.Ethertype())
			}

			_, err = f.UdpConn.Write(buffer[:n])
			if err != nil {
				log.Printf("Error writing to UDP server: %v\n", err)
				continue
			}
		}
	}()
}

func NewUpForwarder(
	ifce *water.Interface,
	udpConn *net.UDPConn,
	wgGlobal *sync.WaitGroup,
) *UpForwarder {
	return &UpForwarder{
		Ifce:     ifce,
		UdpConn:  udpConn,
		wgGlobal: wgGlobal,
	}
}
