package ping

import (
	"fmt"
	"strconv"
)

const (
	// Flood constants based off the man page for 'ping'.
	// Made implementation easier by changing the flood mode from
	// the max of 100 requests/second or as fast as the requests are received
	// to simply 100 requests/second + as fast as the requests are received.
	// With respect to 'ping', the output of flood will not be suppressed
	// with periods/backspaces. It will be consistent with the other options.
	floodFlag = "f"
	floodHelp = "Set the mode to flood. In flood mode, packets are output\n" +
		"100 times per second plus as fast as they are received.\n" +
		"If unset, the program will behave normally. This flag (-f)\n" +
		"is incompatible with wait (-i)."
	floodTimesPerSecond = 100
)

// Flood is a wrapper around a boolean
// to use for command-line argument flag parsing.
type Flood bool

// Init initializes a Flood instance.
// It has an empty body since its zeroed fields
// are sufficient.
func (*Flood) Init() {
}

// String is used to format Flood's value and is required
// to satisfy the flag.Value interface.
func (f *Flood) String() string {
	return fmt.Sprintf("value=%v", *f)
}

// Set will initialize Flood's value using a string, and is
// required to satisfy the flag.Value interface.
func (f *Flood) Set(val string) error {
	res, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*f = Flood(res)
	return nil
}

// Flag gets the command-line flag used for Flood.
func (*Flood) Flag() string {
	return floodFlag
}

// Help gets the command-line help for Flood.
func (*Flood) Help() string {
	return floodHelp
}

// IsBoolFlag is used to notify that Flood is
// a boolean flag, so '-f' defaults to '-f=true' or '-f true'.
func (*Flood) IsBoolFlag() bool {
	return true
}
