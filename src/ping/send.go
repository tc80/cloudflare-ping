package ping

import (
	"fmt"

	"golang.org/x/net/icmp"
)

// sends an ICMP "echo request" to a host for a particular
// sequence using the Ping request
func (p *Ping) send(seq int) error {
	// create echo request
	payload := p.PacketSize.GeneratePayload()
	message := icmp.Message{
		Type: p.requestType,
		Body: &icmp.Echo{
			ID:   p.id,
			Seq:  seq,
			Data: payload,
		},
	}
	// marshal echo request into bytes
	bytes, err := message.Marshal(nil)
	if err != nil {
		fmt.Printf("FAILED marshal: %v\n", err)
	}
	// write bytes
	_, err = p.conn.WriteTo(bytes, p.hostAddr)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	}
	fmt.Printf("sent: %v\n", payload)
	p.sent[seq] = payload
	// reminders <- new thing
	// thread that manages the reminders
	// received reminder -> check
	return nil
}
