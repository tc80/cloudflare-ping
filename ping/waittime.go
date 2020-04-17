package ping

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// WaitTime constants based off the man page for 'ping'.
	waitTimeFlag = "W"
	waitTimeHelp = "Set the time in milliseconds to wait for a reply with\n" +
		"each packet sent. If a reply arrives after the interval,\n" +
		"it is not printed, but counted as a replied packet\n" +
		"for the statistics. If unset, waittime is 4 seconds."
	waitTimeInvalid       = "waittime must be greater than or equal to 0"
	waitTimeDefaultMillis = 4000
)

var (
	// error for invalid waittime
	errWaitTimeInvalid = errors.New(waitTimeInvalid)
)

// WaitTime is a wrapper around a time.Duration
// to use for command-line argument flag parsing.
type WaitTime time.Duration

// Init initializes a WaitTime instance by setting its
// default value.
func (w *WaitTime) Init() {
	*w = WaitTime(time.Millisecond * waitTimeDefaultMillis)
}

// String is used to format WaitTime's value and is required
// to satisfy the flag.Value interface.
func (w *WaitTime) String() string {
	return fmt.Sprintf("value=%v", time.Duration(*w))
}

// Set will initialize WaitTime's value using a string, and is
// required to satisfy the flag.Value interface.
func (w *WaitTime) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res < 0 {
		return errWaitTimeInvalid
	}
	*w = WaitTime(time.Millisecond * time.Duration(res))
	return nil
}

// Flag gets the command-line flag used for WaitTime.
func (*WaitTime) Flag() string {
	return waitTimeFlag
}

// Help gets the command-line help for WaitTime.
func (*WaitTime) Help() string {
	return waitTimeHelp
}
