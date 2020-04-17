package ping

import (
	"errors"
	"net"
)

const (
	// constants for resolving addresses as ipv4, ipv6
	ipv4Network          = "ip4"
	ipv6Network          = "ip6"
	ipv4ICMPNetwork      = "ip4:icmp"
	ipv6ICMPNetwork      = "ip6:ipv6-icmp"
	ipv4BindAddress      = "0.0.0.0" // capture all ipv4 addresses
	ipv6BindAddress      = "::"      // capture all ipv6 addresses
	hostInvalid          = "invalid IPv4 or IPv6 address"
	ianaProtocolIPv4ICMP = 1
	ianaProtocolIPv6ICMP = 58
)

var (
	// error for invalid host
	errHostInvalid = errors.New(hostInvalid)
)

// ResolveHost attempts to resolve a string hostname
// into an IPv4 or IPv6 address.
// Returns the ip addr pointer, a boolean 'true' if
// the address is IPv4, and an error if anything went wrong.
func ResolveHost(host string) (*net.IPAddr, bool, error) {
	// try to resolve as IPv4
	ipAddr, err := net.ResolveIPAddr(ipv4Network, host)
	if err == nil {
		return ipAddr, true, nil
	}
	ipAddr, err = net.ResolveIPAddr(ipv6Network, host)
	if err == nil {
		return ipAddr, false, nil
	}
	return ipAddr, false, errHostInvalid
}
