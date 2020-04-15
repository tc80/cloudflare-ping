package main

import (
	"cloudflare-ping/src/ping"
	"flag"
	"log"
)

const (
	hostArgIndex = 0
	argCount     = 1
	usage        = "go run main/main.go [-c count] [-m ttl] [-t timeout] host"
)

// flagArg interface allows us to process the command-line
// arguments generically
type flagArg interface {
	Init()
	String() string
	Set(string) error
	Flag() string
	Help() string
}

func main() {
	p := parse()  // parse args
	ping.Start(p) // start pinging
}

// Parses the command-line arguments and flags passed
// to the program.
func parse() *ping.Ping {
	p := ping.Ping{}
	flags := []flagArg{
		&p.TTL,
		&p.Count,
		&p.Timeout,
	}
	// parse each flag, each implements flag.Value
	for _, f := range flags {
		f.Init()
		flag.Var(f, f.Flag(), f.Help())
	}
	flag.Parse()
	// parse host name argument
	args := flag.Args()
	if len(args) != argCount {
		log.Fatalf("\nInvalid number of arguments. Expected %v, received %v.\nUsage: %v\n", argCount, len(args), usage)
	}
	p.Host = args[hostArgIndex]
	// return pointer to ping.Ping
	return &p
}
