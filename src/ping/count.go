package ping

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// Count constants based off the man page for 'ping'.
	countFlag    = "c"
	countHelp    = "Set the number of ECHO_RESPONSE packets that are sent and received before stopping the program. If unset, the program will loop until interrupted."
	countInvalid = "count must be greater than 0"
)

// Count is a wrapper around a boolean and unsigned integer
// to use for command-line argument flag parsing.
type Count struct {
	IsSet bool
	Value uint32
}

// Init initializes a Count instance.
// It has an empty body since its zeroed fields
// are sufficient.
func (c *Count) Init() {
}

// String is used to format Count's value and is required
// to satisfy the flag.Value interface.
func (c *Count) String() string {
	return fmt.Sprint(*c)
}

// Set will initialize Count's value using a string, and is
// required to satisfy the flag.Value interface.
func (c *Count) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res <= 0 {
		return errors.New(countInvalid)
	}
	c.IsSet = true
	c.Value = uint32(res)
	return nil
}

// Flag gets the command-line flag used for Count.
func (*Count) Flag() string {
	return countFlag
}

// Help gets the command-line help for Count.
func (*Count) Help() string {
	return countHelp
}

// Uint32 gets the inner integer value Count wraps around.
func (c *Count) Uint32() uint32 {
	return c.Value
}
