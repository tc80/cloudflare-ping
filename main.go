package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"syscall"
)

// followed ping man page

const (
	ttlFlag    = "m"
	ttlSysVar  = "net.inet.ip.ttl"
	ttlHelp    = "Set the time to live (ttl) for outgoing packets as an integer. If unset, the default ttl is the system value sysctl " + ttlSysVar + "."
	ttlInvalid = "time to live (ttl) must be greater than 0"
)

type Flag interface {
	Init()
	String() string
	Set(string) error
	Flag() string
	Help() string
}

type TimeToLive uint32

func (t *TimeToLive) Init() {
	ttlDefault, err := syscall.SysctlUint32(ttlSysVar)
	if err != nil {
		log.Fatalf("Failed to query sysctl for %v: %v\n", ttlSysVar, err)
	}
	*t = TimeToLive(ttlDefault)
}

func (t *TimeToLive) String() string {
	return fmt.Sprint(*t)
}

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

func (t *TimeToLive) Flag() string {
	return ttlFlag
}

func (t *TimeToLive) Help() string {
	return ttlHelp
}

func (t *TimeToLive) Uint32() uint32 {
	return uint32(*t)
}

func main() {
	flags := []Flag{}
	t := new(TimeToLive)
	flags = append(flags, t)
	for _, f := range flags {
		f.Init()
		flag.Var(f, f.Flag(), f.Help())
	}
	flag.Parse()
	fmt.Println(t)
	// flag.Value{}
	// flag.Value{

	// }
	// //ttlDefault, err :=
	// str, err := syscall.SysctlUint32("net.inet.ip.ttl")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%v\n", str)
	// var ip = flag.Int(timeToLiveFlag, 1234, "help message for flagname")
	// flag.Parse()
	// fmt.Printf("Int pointer: %v\n", *ip)
}
