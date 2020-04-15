package ping

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"syscall"
)

const (
	// TTL constants based off the man page for 'ping'.
	ttlFlag    = "m"
	ttlSysVar  = "net.inet.ip.ttl"
	ttlHelp    = "Set the time to live (ttl) for outgoing packets as an integer. If unset, the default ttl is the system value sysctl " + ttlSysVar + "."
	ttlInvalid = "time to live (ttl) must be greater than or equal to 0"
)

// TimeToLive is a wrapper around an unsigned integer
// to use for command-line argument flag parsing.
type TimeToLive uint32

// Init initializes a TimeToLive instance by setting its
// value to the Management Information Base (MIB) variable
// for 'net.inet.ip.ttl'.
func (t *TimeToLive) Init() {
	ttlDefault, err := syscall.SysctlUint32(ttlSysVar)
	if err != nil {
		log.Fatalf("Failed to query sysctl for %v: %v\n", ttlSysVar, err)
	}
	*t = TimeToLive(ttlDefault)
}

// String is used to format TimeToLive's value and is required
// to satisfy the flag.Value interface.
func (t *TimeToLive) String() string {
	return fmt.Sprint(*t)
}

// Set will initialize TimeToLive's value using a string, and is
// required to satisfy the flag.Value interface.
func (t *TimeToLive) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res < 0 {
		return errors.New(ttlInvalid)
	}
	*t = TimeToLive(res)
	return nil
}

// Flag gets the command-line flag used for TimeToLive.
func (*TimeToLive) Flag() string {
	return ttlFlag
}

// Help gets the command-line help for TimeToLive.
func (*TimeToLive) Help() string {
	return ttlHelp
}

// Uint32 gets the inner integer value TimeToLive wraps around.
func (t *TimeToLive) Uint32() uint32 {
	return uint32(*t)
}
