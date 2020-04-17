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


