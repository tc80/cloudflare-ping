# Simple Ping App for MacOS

A small ping CLI application using ICMP echo requests, based off the ping man page. 

It was developed with Go version 1.13.5 and macOS Catalina 10.15.4.

## Features

- [x] IPv4 and IPv6 Support
- [x] Packets Reported
    - [x] TTL, RTT
    - [x] Support for Time Limit Exceeded
    - [x] Support for Destination Unreachable
- [x] Configurable Flags
    - [x] Count
    - [x] Flood
    - [x] Wait
    - [x] TTL
    - [x] Packet Size
    - [x] Timeout
    - [x] Wait Time
- [x] Statistics Reported
    - [x] Packets Transmitted
    - [x] Packets Received
    - [x] Packet Loss
    - [x] Packets Out of Wait Time
    - [x] RTT Min/Avg/Max/Stddev

## Build

To build, use the `Makefile` provided.

Simply run `make` to build and `make clean` to remove the executable.

To run the program once built:

`sudo ./main/ping [-W waittime] [-c count] [-f] [-i wait] [-m ttl] [-s packetsize] [-t timeout] host`

The usage will be printed in the case of any errors. For instance, the flags `-i` and `-f` are mutually exclusive. Note that `host` is any valid hostname or IPv4/IPv6 address.

Make sure that this repository is located in your computer's `GOPATH` in the top-level `src` directory. Otherwise, you may need to modify the import statements for the program to build. 

## Note

This application is strongly built off of the ping man page with respect to command-line flags and output statements for packets and statistics.

Similarly to the ping command, the statistics are printed when the number of sent ICMP echo requests is satisfied (if set), or when the program is interrupted, whether by the timeout argument or manually with an interrupt signal.

To make the flood implementation slightly easier, it was altered from the ping man page to send 100 requests/second plus as fast as the packets are received. Originally, this was the maximum of the two.

Finally, when testing with IPv6 addresses, make sure IPv6 is enabled on your router.

## Tests

There are several tests/examples of running the application in the `Makefile`. For example: `make run ping-google-dns-ipv6` pings Google's IPv6 DNS 5 times and outputs the statistics. Remember to build before running.







