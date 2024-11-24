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
			log.Println("Down Forwarder: waiting for packet...")
			var frame ethernet.Frame
			frame.Resize(1500)
			n, _, err := f.UdpConn.ReadFromUDP([]byte(frame))
			if err != nil {
				log.Println("Error reading from UDP connection:", err)
			}
			frame = frame[:n]
			_, err = f.Ifce.Write([]byte(frame))
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
