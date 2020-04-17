package ping

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// Timeout constants based off the man page for 'ping'.
	timeoutFlag = "t"
	timeoutHelp = "Set the timeout in seconds for the program to exit, regardless\n" +
		"of the number of packets sent/received. If unset, the\n" +
		"program will behave normally."
	timeoutInvalid = "timeout must be greater than or equal to 0"
)

var (
	// error for invalid timeout
	errTimeoutInvalid = errors.New(timeoutInvalid)
)

// Timeout is a wrapper around a boolean and a time.Duration
// to use for command-line argument flag parsing.
type Timeout struct {
	IsSet bool
	Value time.Duration
}

// Init initializes a Timeout instance.
// It has an empty body since its zeroed fields
// are sufficient.
func (*Timeout) Init() {
}

// String is used to format Timeout's value and is required
// to satisfy the flag.Value interface.
func (t *Timeout) String() string {
	return fmt.Sprintf("set=%v, value=%v", t.IsSet, t.Value)
}

// Set will initialize Timeout's value using a string, and is
// required to satisfy the flag.Value interface.
func (t *Timeout) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res < 0 {
		return errTimeoutInvalid
	}
	t.IsSet = true
	t.Value = time.Second * time.Duration(res)
	return nil
}

// Flag gets the command-line flag used for Timeout.
func (*Timeout) Flag() string {
	return timeoutFlag
}

// Help gets the command-line help for Timeout.
func (*Timeout) Help() string {
	return timeoutHelp
}
