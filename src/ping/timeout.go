package ping

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// Timeout constants based off the man page for 'ping'.
	timeoutFlag    = "t"
	timeoutHelp    = "Set the timeout in seconds for the program to exit, regardless of the number of packets sent/received. If unset, the program will behave normally."
	timeoutInvalid = "timeout must be greater than or equal to 0"
)

// Timeout is a wrapper around a boolean and unsigned integer
// to use for command-line argument flag parsing.
type Timeout struct {
	IsSet bool
	Value uint32
}

// Init initializes a Timeout instance.
// It has an empty body since its zeroed fields
// are sufficient.
func (t *Timeout) Init() {
}

// String is used to format Timeout's value and is required
// to satisfy the flag.Value interface.
func (t *Timeout) String() string {
	return fmt.Sprint(*t)
}

// Set will initialize Timeout's value using a string, and is
// required to satisfy the flag.Value interface.
func (t *Timeout) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res <= 0 {
		return errors.New(timeoutInvalid)
	}
	t.IsSet = true
	t.Value = uint32(res)
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

// Uint32 gets the inner integer value Timeout wraps around.
func (t *Timeout) Uint32() uint32 {
	return t.Value
}
