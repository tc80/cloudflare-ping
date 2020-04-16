package ping

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Ping is used to represent a request to
// send ICMP "echo requests" to a particular host.
type Ping struct {
	TTL                 TimeToLive            // time to live (uint32)
	PacketSize          PacketSize            // packet size (uint16)
	Count               Count                 // if set, number of echo response packets sent and received
	Timeout             Timeout               // if set, time before program exits
	Flood               Flood                 // flood mode
	Wait                Wait                  // wait time between sending pings
	WaitTime            WaitTime              // max round-trip time for outputting response
	HostName            string                // host name as a string
	hostAddr            *net.IPAddr           // host as an address
	isIPv4              bool                  // if the host is IPv4
	proto               int                   // iana protocol
	conn                *icmp.PacketConn      // connection for sending/receiving
	id                  int                   // id for requests/responses
	requestType         icmp.Type             // ICMP request type
	replyType           icmp.Type             // ICMP response type
	waitTimeExceeded    map[int]bool          // response sequences that exceeded waittime (seq -> true)
	waitTimeExceededMux sync.Mutex            // mutex for waitTimeExceeded map
	seen                map[int]time.Duration // response sequences seen (seq -> round-trip time)
	seenMux             sync.Mutex            // mutex for seen map
	sent                map[int][]byte        // sent sequences (seq -> payload)
	sentMux             sync.Mutex            // mutex for sent map
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
	var icmpNetwork, bindAddress string
	if p.isIPv4 {
		icmpNetwork = ipv4ICMPNetwork
		bindAddress = ipv4BindAddress
		p.requestType = ipv4.ICMPTypeEcho
		p.replyType = ipv4.ICMPTypeEchoReply
		p.proto = ianaProtocolIPv4ICMP
	} else {
		icmpNetwork = ipv6ICMPNetwork
		bindAddress = ipv6BindAddress
		p.requestType = ipv6.ICMPTypeEchoRequest
		p.replyType = ipv6.ICMPTypeEchoReply
		p.proto = ianaProtocolIPv6ICMP
	}
	conn, err := icmp.ListenPacket(icmpNetwork, bindAddress)
	if err != nil {
		return fmt.Errorf("failed to get packet conn: %v", err)
	}
	// set ttl (ipv4) / hop limit (ipv6)
	if p.isIPv4 {
		conn.IPv4PacketConn().SetTTL(int(p.TTL))
	} else {
		conn.IPv6PacketConn().SetHopLimit(int(p.TTL))
	}
	// set packet connection
	p.conn = conn
	// set id based on process id
	p.id = os.Getpid()
	// initialize maps and mutexes
	p.waitTimeExceeded = make(map[int]bool)
	p.seen = make(map[int]time.Duration)
	p.sent = make(map[int][]byte)
	p.waitTimeExceededMux = sync.Mutex{}
	p.seenMux = sync.Mutex{}
	p.sentMux = sync.Mutex{}
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
	go p.receive()
	// printing same message as 'ping' command on start
	fmt.Printf("PING %v (%v): %v data bytes\n", p.HostName, p.hostAddr.String(), p.PacketSize)
	p.send(0)
	<-time.After(time.Second)
	p.send(1)
	<-time.After(time.Second)
	p.send(2)
	// errors := make(chan error) // channel for receiving errors
	// // gracefully kill ping
	// select {
	// case _ = <-errors:
	// 	break
	// case <-time.After(p.Timeout.Value):
	// 	fmt.Printf("Timed out!, so print the stats")
	// }
	// fmt.Printf("received: %v\n", p)
	return nil
}
