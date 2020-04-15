package main

import (
	"cloudflare-ping/src/ping"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	hostArgIndex = 0
	argCount     = 1
	usageExample = "sudo go run main/main.go [-c count] [-f] [-i wait] [-m ttl] [-s packetsize] [-t timeout] host"
)

// flagArg interface allows us to process the command-line
// arguments generically.
type flagArg interface {
	Init()            // initializes the value
	String() string   // value as a string
	Set(string) error // sets the value from a string
	Flag() string     // gets the flag representation
	Help() string     // gets help message
}

// print usage to stderr and quit
func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %v\n\n", usageExample)
	flag.PrintDefaults()    // print flag defaults
	fmt.Fprintln(os.Stderr) // space for readability
	os.Exit(1)              // exit program
}

// main method
func main() {
	p := parse()        // parse args
	err := p.Validate() // check if valid
	if err != nil {
		fmt.Printf("Failed to ping: %v\n", err)
		usage()
	}
	err = p.Start() // start pinging
	if err != nil {
		log.Fatalf("ping failure: %v\n", err)
	}
	fmt.Println("Ping finished successfully.")
	// maybe pass back some stats to print out?
	// shutdown hook for ping

}

// Parses the command-line arguments and flags passed
// to the program.
func parse() *ping.Ping {
	p := ping.Ping{}
	flags := []flagArg{
		&p.TTL,
		&p.Count,
		&p.Timeout,
		&p.PacketSize,
		&p.Flood,
		&p.Wait,
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
		// invalid number or order of arguments
		usage()
	}
	p.HostName = args[hostArgIndex]
	// return pointer to ping.Ping
	return &p
}
