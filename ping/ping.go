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
	TTL          TimeToLive          // time to live (uint32)
	PacketSize   PacketSize          // packet size (uint16)
	Count        Count               // if set, number of echo response packets sent and received
	Timeout      Timeout             // if set, time before program exits
	Flood        Flood               // flood mode
	Wait         Wait                // wait time between sending pings
	WaitTime     WaitTime            // max round-trip time for outputting response
	HostName     string              // host name as a string
	hostAddr     *net.IPAddr         // host as an address
	isIPv4       bool                // if the host is IPv4
	proto        int                 // iana protocol
	conn         *icmp.PacketConn    // connection for sending/receiving
	id           int                 // id for requests/responses
	requestType  icmp.Type           // ICMP request type
	replyType    icmp.Type           // ICMP response type
	sent         map[int]*icmpPacket // sent sequences (seq -> sent packet)
	sentMux      sync.Mutex          // mutex for sent map
	floodRecv    int                 // how many packets received in the last second
	floodRecvMux sync.Mutex          // mutex for flood recv
	waitGroup    sync.WaitGroup      // wait group to wait for all helper goroutines to finish
}

// Validate checks if the Ping request is valid,
// returning a non-nil error if invalid.
// Should be called before calling Start().
// Requirements:
// 		Count > 0
//		Host must be valid IPv4 or IPv6 address
//		Cannot have both wait flag (-i) and flood flag (-f) at a time
func (p *Ping) Validate() error {
	if p.Count.IsSet && p.Count.Value == 0 {
		return errCountInvalid
	}
	if _, _, err := ResolveHost(p.HostName); err != nil {
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
	addr, IPv4, err := ResolveHost(p.HostName)
	if err != nil {
		panic(err)
	}
	p.hostAddr = addr
	p.isIPv4 = IPv4
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
	p.sent = make(map[int]*icmpPacket)
	p.sentMux = sync.Mutex{}
	p.floodRecvMux = sync.Mutex{}
	// create wait group
	p.waitGroup = sync.WaitGroup{}
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
	// print stats if program interrupted
	createInterruptHandler(p)
	fmt.Printf("PING %v (%v): %v data bytes\n", p.HostName, p.hostAddr.String(), p.PacketSize)
	// make channel for timeout
	timeout := make(chan bool)
	defer close(timeout)
	if p.Timeout.IsSet {
		go func() {
			<-time.After(p.Timeout.Value) // wait
			timeout <- true               // notify timeout
		}()
	}
	// make channels for done, errors
	// done will notify threads to clean up and stop
	// errors will be used to pass errors from threads
	done := make(chan bool)
	errors := make(chan error)
	defer close(errors)
	defer close(done)
	p.waitGroup.Add(1)
	// start receiving
	go p.receiver(done, errors)
	p.waitGroup.Add(1)
	// start sending
	if bool(p.Flood) {
		go p.floodSender(done, errors)
	} else {
		go p.sender(done, errors)
	}
	// clean up function that notifies sender/receiver to stop
	cleanUp := func() {
		done <- true
		done <- true
	}
	err = nil
	// wait for an error to occur or a timeout
	// note: an error can be nil, indicating a successful
	// sender termination if we are sending finite packets
	select {
	case <-timeout:
	case err = <-errors:
	}
	go cleanUp()
	// wait for all threads to clean up
	p.waitGroup.Wait()
	// print stats
	p.printStats()
	return err
}
