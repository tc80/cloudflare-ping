package ping

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	// PacketSize constants based off the man page for 'ping'.
	packetSizeDefault    = 56
	packetPayloadSizeMax = 65507 // max payload size
	icmpPacketMaxSize    = 65535 // includes headers
	packetSizeFlag       = "s"
	packetSizeHelp       = "Set the number of data bytes sent. If unset, 56 bytes\n" +
		"will be sent, which becomes 64 ICMP data bytes when included\n" +
		"with the ICMP header data (8 bytes). Note that due to the Go\n" +
		"ipv4/ipv6 library, small packet sizes may not work. For example,\n" +
		"the minimum header size in the ipv4 library is 20."
	packetSizeInvalid  = "packet size must be greater than or equal to 0"
	packetSizeTooLarge = "packet size too large"
)

var (
	// error for invalid packet size
	errPacketSizeInvalid = errors.New(packetSizeInvalid)
)

// represents a sent ICMP packet
type icmpPacket struct {
	sendTime         time.Time     // time sent
	receiveTime      time.Time     // time received
	roundtripTime    time.Duration // rtt time
	receivedTTL      int           // ttl when received
	received         bool          // if the packet has been received
	waitTimeExceeded bool          // if the packet exceeded its wait time
	payload          []byte        // payload
}

// PacketSize is a wrapper around an unsigned integer
// to use for command-line argument flag parsing.
type PacketSize uint16

// Init initializes a PacketSize instance.
// It has an empty body since its zeroed fields
// are sufficient.
func (p *PacketSize) Init() {
	*p = PacketSize(packetSizeDefault)
}

// String is used to format PacketSize's value and is required
// to satisfy the flag.Value interface.
func (p *PacketSize) String() string {
	return fmt.Sprintf("value=%v", *p)
}

// Set will initialize PacketSize's value using a string, and is
// required to satisfy the flag.Value interface.
func (p *PacketSize) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res < 0 {
		return errPacketSizeInvalid
	}
	if res > packetPayloadSizeMax {
		return fmt.Errorf("%v: %v > %v", packetSizeTooLarge, res, packetPayloadSizeMax)
	}
	*p = PacketSize(res)
	return nil
}

// Flag gets the command-line flag used for PacketSize.
func (*PacketSize) Flag() string {
	return packetSizeFlag
}

// Help gets the command-line help for PacketSize.
func (*PacketSize) Help() string {
	return packetSizeHelp
}

// GeneratePayload makes a random byte array for the packet size.
func (p *PacketSize) GeneratePayload() []byte {
	// init PRNG with time as seed
	rand.Seed(time.Now().UnixNano())
	// fill buffer with random bytes
	buffer := make([]byte, *p)
	rand.Read(buffer)
	return buffer
}
