package ping

import "fmt"

// Ping ...
type Ping struct {
	TTL     TimeToLive
	Count   Count
	Timeout Timeout
	Host    string
}

// Start ...
func Start(p *Ping) {
	fmt.Printf("received: %v\n", p)
}
