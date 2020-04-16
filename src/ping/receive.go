package ping

import (
	"fmt"

	"golang.org/x/net/icmp"
)

// sends an ICMP "echo request" to a host for a particular
// sequence using the Ping request
func (p *Ping) receive() error {
	for {
		buffer := make([]byte, packetSizeMax)
		if n, _, err := p.conn.ReadFrom(buffer); err != nil {
			panic(err)
		} else {
			go p.handleEchoReply(buffer[:n])
			message, err := icmp.ParseMessage(p.proto, buffer[:n])
			if err != nil {
				panic(err)
			}

			m := message.Body.(*icmp.Echo)

			//bytes, _ := m.Marshal(p.proto)

			fmt.Printf("received: %v\n", m.Data)
		}

	}
	return nil
}

func (p *Ping) handleEchoReply(reply []byte) {

}
