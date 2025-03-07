package ping

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	readTimeout = time.Second // timeout for reading icmp packets
)

// sends an ICMP "echo request" to a host for a particular
// sequence using the Ping request
func (p *Ping) receiver(done <-chan bool, errors chan<- error) {
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
				go func() { errors <- fmt.Errorf("failed to read: %v", err) }()
				return
			}
			// handle reply
			recvTime := time.Now()
			p.waitGroup.Add(1)
			go p.handleReply(buffer[:n], recvTime)
			if bool(p.Flood) {
				// update count of recv packages since last reset
				go func() {
					p.floodRecvMux.Lock()
					defer p.floodRecvMux.Unlock()
					p.floodRecv++
				}()
			}
		}
	}
}

// handles the reply depending on its type
func (p *Ping) handleReply(reply []byte, recvTime time.Time) {
	defer p.waitGroup.Done()
	// attempt to parse message
	message, err := icmp.ParseMessage(p.proto, reply)
	if err != nil {
		return // failed to parse message, so ignore it
	}
	var header interface{} // message header
	// classify message
	switch message.Type {
	case ipv4.ICMPTypeTimeExceeded, ipv6.ICMPTypeTimeExceeded:
		body, ok := message.Body.(*icmp.TimeExceeded)
		if !ok || body == nil {
			return // failed to parse body, ignore
		}
		if p.isIPv4 {
			header, err = ipv4.ParseHeader(body.Data)
			if header == nil || err != nil {
				return // failed to parse header, ignore
			}
		} else {
			header, err = ipv6.ParseHeader(body.Data)
			if header == nil || err != nil {
				return // failed to parse header, ignore
			}
		}
		p.handleEchoTimeExceeded(reply, recvTime, header, body)
	case ipv4.ICMPTypeDestinationUnreachable, ipv6.ICMPTypeDestinationUnreachable:
		body, ok := message.Body.(*icmp.DstUnreach)
		if !ok || body == nil {
			return // failed to parse body, ignore
		}
		if p.isIPv4 {
			header, err = ipv4.ParseHeader(body.Data)
			if header == nil || err != nil {
				return // failed to parse header, ignore
			}
		} else {
			header, err = ipv6.ParseHeader(body.Data)
			if header == nil || err != nil {
				return // failed to parse header, ignore
			}
		}
		p.handleEchoDstUnreachable(reply, recvTime, header, body)
	case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
		body, ok := message.Body.(*icmp.Echo)
		if !ok || body == nil {
			return // failed to parse body, ignore
		}
		if p.isIPv4 {
			header, err = ipv4.ParseHeader(body.Data)
			if header == nil || err != nil {
				return // failed to parse header, ignore
			}
		} else {
			header, err = ipv6.ParseHeader(body.Data)
			if header == nil || err != nil {
				return // failed to parse header, ignore
			}
		}
		p.handleEchoReply(reply, recvTime, header, body)
	default:
		return // unknown or unhandled type, so ignoring
	}
}

// handles an IPv4 or IPv6 echo host time exceeded reply
// header interface argument is either
// a non-nil *ipv4.Header or non-nil *ipv6.Header
// for now, there is no validation to check
// for associated sequence / if we sent a request
func (p *Ping) handleEchoTimeExceeded(reply []byte, recvTime time.Time, header interface{}, body *icmp.TimeExceeded) {
	fmt.Printf("%v bytes from %v: Time to live exceeded\n%v\n",
		len(reply), p.hostAddr.String(), header)
}

// handles an IPv4 or IPv6 echo host unreachable reply
// header interface argument is either
// a non-nil *ipv4.Header or non-nil *ipv6.Header
// for now, there is no validation to check
// for associated sequence / if we sent a request
func (p *Ping) handleEchoDstUnreachable(reply []byte, recvTime time.Time, header interface{}, body *icmp.DstUnreach) {
	fmt.Printf("%v bytes from %v: Destination unreachable\n%v\n",
		len(reply), p.hostAddr.String(), header)
}

// handles an IPv4 or IPv6 echo reply, where the
// header interface argument is either
// a non-nil *ipv4.Header or non-nil *ipv6.Header
func (p *Ping) handleEchoReply(reply []byte, recvTime time.Time, header interface{}, body *icmp.Echo) {
	// validate
	if body.ID != p.id {
		return // echo request not sent by our client, so ignore response
	}
	p.sentMux.Lock()
	defer p.sentMux.Unlock()
	// only handle new valid sequence numbers
	if packet, ok := p.sent[body.Seq]; ok && !packet.received {
		packet.received = true
		packet.receiveTime = recvTime
		packet.roundtripTime = recvTime.Sub(packet.sendTime)
		// get ttl from header
		if p.isIPv4 {
			packet.receivedTTL = header.(*ipv4.Header).TTL
		} else {
			packet.receivedTTL = header.(*ipv6.Header).HopLimit
		}
		// only print if wait time not exceeded
		// note: currently not checking integrity of payload
		if !packet.waitTimeExceeded {
			fmt.Printf("%v bytes from %v: icmp_seq=%v ttl=%v time=%v\n",
				len(reply), p.hostAddr.String(), body.Seq, packet.receivedTTL, packet.roundtripTime)
		}
	}
}
