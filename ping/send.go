package ping

import (
	"fmt"
	"time"

	"golang.org/x/net/icmp"
)

// sends an ICMP "echo request" to a host for a particular
// sequence using the Ping request
func (p *Ping) sender(done <-chan bool, errors chan<- error) {
	defer p.waitGroup.Done()
	// keep sending forever unless count is set
	for i := 0; !p.Count.IsSet || i < int(p.Count.Value); i++ {
		select {
		case <-done:
			return // stop sending
		default:
			// send sequence i
			err := p.send(i)
			if err != nil {
				errors <- err
				return
			}
			<-time.After(time.Duration(p.Wait.Value))
		}
	}
}

// sends an ICMP "echo request" in flood mode to a host for a particular
// sequence using the Ping request
func (p *Ping) floodSender(done <-chan bool, errors chan<- error) {
	defer p.waitGroup.Done()
	// keep sending forever unless count is set
	for i := 0; !p.Count.IsSet || i < int(p.Count.Value); i++ {
		select {
		case <-done:
			return // stop sending
		default:
			// send sequence i
			err := p.send(i)
			if err != nil {
				errors <- err
				return
			}
			<-time.After(time.Duration(p.Wait.Value))
		}
	}
}

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
		return fmt.Errorf("failed to marshal echo request: %v", err)
	}
	// add sent entry
	p.sentMux.Lock()
	sendTime := time.Now()
	p.sent[seq] = &icmpPacket{
		sendTime: sendTime,
		payload:  payload,
	}
	p.sentMux.Unlock()
	// send echo request
	_, err = p.conn.WriteTo(bytes, p.hostAddr)
	if err != nil {
		return fmt.Errorf("failed to send echo request: %v", err)
	}
	// spawn wait time check
	go func() {
		<-time.After(time.Duration(p.WaitTime))
		p.sentMux.Lock()
		packet := p.sent[seq]
		if !packet.received {
			// has not been seen yet, so it is late
			packet.waitTimeExceeded = true
		}
		p.sentMux.Unlock()
	}()
	return nil
}
