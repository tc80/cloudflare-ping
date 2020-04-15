package ping

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// Wait constants based off the man page for 'ping'.
	waitFlag = "i"
	waitHelp = "Set the number of seconds to wait between sending each packet.\n" +
		"The number can be a fraction (ex. 0.1). If unset,\n" +
		"the default is a one second interval between packets.\n" +
		"This flag (-i) is incompatible with flood (-f)."
	waitInvalid      = "count must be greater than or equal to 0"
	waitDefault      = time.Second
	waitInputBitSize = 64 // float64 accepted as input
)

var (
	// error for invalid wait
	errWaitInvalid = errors.New(waitInvalid)
)

// Wait is a wrapper around a boolean and unsigned integer
// to use for command-line argument flag parsing.
type Wait struct {
	IsSet bool
	Value time.Duration
}

// Init initializes a Wait instance.
func (w *Wait) Init() {
	w.Value = waitDefault
}

// String is used to format Wait's value and is required
// to satisfy the flag.Value interface.
func (w *Wait) String() string {
	return fmt.Sprintf("set=%v, value=%v", w.IsSet, w.Value)
}

// Set will initialize Wait's value using a string, and is
// required to satisfy the flag.Value interface.
func (w *Wait) Set(val string) error {
	res, err := strconv.ParseFloat(val, waitInputBitSize)
	if err != nil {
		return err
	}
	if res < 0 {
		return errWaitInvalid
	}
	w.IsSet = true
	w.Value = time.Duration(float64(time.Second) * res)
	return nil
}

// Flag gets the command-line flag used for Wait.
func (*Wait) Flag() string {
	return waitFlag
}

// Help gets the command-line help for Wait.
func (*Wait) Help() string {
	return waitHelp
}
