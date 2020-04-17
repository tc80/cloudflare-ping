package ping

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// handles an interrupt to the program,
// printing the ping stats before exiting
func createInterruptHandler(p *Ping) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c            // received interrupt
		p.printStats() // print stats
		os.Exit(0)     // exit the program
	}()
}

// prints statistics
func (p *Ping) printStats() {
	fmt.Printf("\n--- %v ping statistics ---\n", p.HostName)
	p.sentMux.Lock()
	defer p.sentMux.Unlock()
	transmitted := len(p.sent) // number of packets sent
	if transmitted == 0 {
		fmt.Println("<no packets sent>")
		return // no packets, so no stats to show (avoid division by 0 too)
	}
	var received, exceeded int                  // number of packets recv, wait time exceeded
	var min, sum, max, avg, stDev time.Duration // rtt stats
	min = time.Duration(math.MaxInt64)          // set to max time
	for _, packet := range p.sent {
		if !packet.received {
			continue // not received, so continue
		}
		received++
		if packet.waitTimeExceeded {
			exceeded++
		}
		rtt := packet.roundtripTime
		if rtt > max {
			max = rtt // found new max
		}
		if rtt < min {
			min = rtt // found new min
		}
		sum += rtt // increment sum
	}
	// calculate average rtt
	avg = time.Duration(sum.Nanoseconds() / int64(transmitted))
	// used for calculating stdev
	var sumSquaredDiff float64
	for _, packet := range p.sent {
		rtt := packet.roundtripTime
		sumSquaredDiff += math.Pow(float64((rtt - avg).Nanoseconds()), 2) // add squared difference from mean
	}
	stDev = time.Duration(math.Sqrt(sumSquaredDiff / float64(transmitted)))             // calculate standard deviation
	var packetLoss float64 = 100 * float64(transmitted-received) / float64(transmitted) // calculate packet loss
	packetLoss = math.Ceil(packetLoss*10) / 10                                          // round up (formatting to 1 decimal places)
	fmt.Printf("%v packets transmitted, %v packets received, %.1f%% packet loss",
		transmitted, received, packetLoss)
	if exceeded > 0 {
		fmt.Printf(", %v packets out of wait time", exceeded) // only print exceeded packets if > 0
	}
	fmt.Printf("\nround-trip min/avg/max/stddev = %v/%v/%v/%v\n", min, avg, max, stDev)
}
