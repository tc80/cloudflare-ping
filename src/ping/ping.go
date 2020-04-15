package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	id  = os.Getpid() // id used for echo request/reply
	seq = 0
)

// Ping is used to represent a request to
// send ICMP "echo requests" to a particular host.
type Ping struct {
	TTL         TimeToLive       // time to live (uint32)
	PacketSize  PacketSize       // packet size (uint16)
	Count       Count            // if set, number of echo response packets sent and received
	Timeout     Timeout          // if set, time before program exits
	Flood       Flood            // flood mode
	Wait        Wait             // wait time between pings
	HostName    string           // host name as a string
	hostAddr    *net.IPAddr      // host as an address
	isIPv4      bool             // if the host is IPv4
	conn        *icmp.PacketConn // connection for sending/receiving
	requestType icmp.Type        // ICMP request type
	replyType   icmp.Type        //  ICMP response type

	// do we need to check payload on receiving side?
	//payload    []byte      // payload
}

// Validate checks if the Ping request is valid,
// returning a non-nil error if invalid.
// Should be called before calling Start().
// Requirements: Count > 0, Host must be valid IPv4 or IPv6 address
func (p *Ping) Validate() error {
	if p.Count.IsSet && p.Count.Value == 0 {
		return errCountInvalid
	}
	if _, err := ResolveHost(p.HostName); err != nil {
		return err
	}
	if p.Wait.IsSet && bool(p.Flood) {
		return fmt.Errorf("incompatible flags: -%v and -%v", waitFlag, floodFlag)
	}
	return nil
}

// initializes the Ping's private fields
// for a packet connection and request/reply ICMP types
func (p *Ping) init() error {
	// resolve host
	addr, err := ResolveHost(p.HostName)
	if err != nil {
		panic(err)
	}
	p.hostAddr = addr
	// determine if IPv4 or IPv6
	p.isIPv4 = isIPv4(p.hostAddr)
	// initialize packet connection, req/resp types
	var icmpNetwork string
	if p.isIPv4 {
		icmpNetwork = ipv4ICMPNetwork
		p.requestType = ipv4.ICMPTypeEcho
		p.replyType = ipv4.ICMPTypeEchoReply
	} else {
		icmpNetwork = ipv6ICMPNetwork
		p.requestType = ipv6.ICMPTypeEchoRequest
		p.replyType = ipv6.ICMPTypeEchoReply
	}
	conn, err := icmp.ListenPacket(icmpNetwork, p.hostAddr.String())
	if err != nil {
		return fmt.Errorf("failed to get packet conn: %v", err)
	}
	p.conn = conn
	return nil
}

// sends an ICMP "echo request" to a host for a particular
// sequence using the Ping request
func (p *Ping) send(seq int) error {

	// create echo request
	message := icmp.Message{
		Type: p.requestType,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: p.PacketSize.GeneratePayload(),
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
	fmt.Println("sent")
	return nil
}

// Start begins the ICMP "echo requests"
// using the Ping request.
// Will panic if the Ping request has invalid arguments
// determined by Validate().
func (p *Ping) Start() error {
	err := p.Validate()
	if err != nil {
		panic("invalid ping: " + err.Error())
	}
	err = p.init()
	if err != nil {
		return fmt.Errorf("failed to initialize ping: %v", err)
	}
	// printing same message as 'ping' command on start
	fmt.Printf("PING %v (%v): %v data bytes\n", p.HostName, p.hostAddr.String(), p.PacketSize)
	p.send(0)
	<-time.After(time.Second)
	p.send(1)
	<-time.After(time.Second)
	p.send(2)
	return nil
}

// Start begins the ICMP "echo requests"
// using a given Ping request.
func Start(p *Ping) {
	errors := make(chan error) // channel for receiving errors
	// gracefully kill ping
	select {
	case _ = <-errors:
		break
	case <-time.After(time.Duration(p.Timeout.Value) * time.Second):
		fmt.Printf("Timed out!, so print the stats")
	}
	fmt.Printf("received: %v\n", p)
}
