package forwarder

import (
	"log"
	"net"
	"sync"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

type DownForwarder struct {
	Ifce     *water.Interface
	UdpConn  *net.UDPConn
	wgGlobal *sync.WaitGroup
}

func (f *DownForwarder) Run() {

	f.wgGlobal.Add(1)

	go func() {
		defer f.wgGlobal.Done()
		for {
			buffer := make([]byte, 1500)
			var frame ethernet.Frame
			frame.Resize(1500)
			n, _, err := f.UdpConn.ReadFromUDP(buffer)
			if err != nil {
				log.Println("Error reading from UDP connection:", err)
			}

			frame = ethernet.Frame(buffer[:n])
			if n > 0 {
				log.Printf("UDP -> TAP: packet length=%d, src=%s, dst=%s, type=0x%x\n",
					n, frame.Source(), frame.Destination(), frame.Ethertype())
			}

			_, err = f.Ifce.Write(buffer[:n])
			if err != nil {
				log.Println("Error writing to TAP interface:", err)
			}
		}
	}()

}

func NewDownForwarder(
	ifce *water.Interface,
	udpConn *net.UDPConn,
	wgGlobal *sync.WaitGroup,
) *DownForwarder {
	return &DownForwarder{
		Ifce:     ifce,
		UdpConn:  udpConn,
		wgGlobal: wgGlobal,
	}
}
