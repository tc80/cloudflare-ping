package ping

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
)

var (
	readTimeout = time.Second // timeout for reading icmp packets
)

// sends an ICMP "echo request" to a host for a particular
// sequence using the Ping request
func (p *Ping) receiver(done <-chan bool, fatal chan<- error) {
	defer p.waitGroup.Done()
	for {
		select {
		case <-done:
			return
		default:
			buffer := make([]byte, icmpPacketMaxSize)           // assuming max packet
			p.conn.SetReadDeadline(time.Now().Add(readTimeout)) // avoid blocking read (might want to clean up)
			n, _, err := p.conn.ReadFrom(buffer)                // read incoming icmp packets
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue // timed out, try to read again
			}
			if err != nil {
				fatal <- fmt.Errorf("failed to read: %v", err)
				return
			}
			p.waitGroup.Add(1)
			go p.handleEchoReply(buffer[:n])
		}
	}
}

// handles an echo reply
func (p *Ping) handleEchoReply(reply []byte) {
	defer p.waitGroup.Done()
	message, err := icmp.ParseMessage(p.proto, reply)
	if err != nil {
		panic(err)
	}

	m := message.Body.(*icmp.Echo)

	//bytes, _ := m.Marshal(p.proto)

	fmt.Printf("received: %v\n", m.Seq)
	p.sentMux.Lock()
	defer p.sentMux.Unlock()
	// only handle valid sequence numbers
	if packet, ok := p.sent[m.Seq]; ok {
		if packet.received {
			return // ignore - packet has been seen before
		}
		// only print if wait time not exceeded
		if !packet.waitTimeExceeded {

		}
	}
	p.sentMux.Unlock()
}
