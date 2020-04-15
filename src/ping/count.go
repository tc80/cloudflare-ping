package ping

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// Count constants based off the man page for 'ping'.
	countFlag = "c"
	countHelp = "Set the number of echo packets that are sent and received\n" +
		"before stopping the program. If unset, the program will loop until interrupted."
	countInvalid = "count must be greater than 0"
)

var (
	// error for invalid count
	errCountInvalid = errors.New(countInvalid)
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
func (*Count) Init() {
}

// String is used to format Count's value and is required
// to satisfy the flag.Value interface.
func (c *Count) String() string {
	return fmt.Sprintf("set=%v, value=%v", c.IsSet, c.Value)
}

// Set will initialize Count's value using a string, and is
// required to satisfy the flag.Value interface.
func (c *Count) Set(val string) error {
	res, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if res <= 0 {
		return errCountInvalid
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
