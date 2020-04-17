# Makefile for building and running the ICMP ping app.
# @author tc80, April 2020
#
# Note, project was built with:
# 	Go version: go1.13.5 darwin/amd64
#	macOS Catalina 10.15.4

# build the project
all: 
	cd main/; go build -o ping

# ping google's IPv4 DNS 5 times
ping-google-dns-ipv4:
	sudo ./main/ping -c 5 8.8.8.8

# ping google's IPv6 DNS 5 times
ping-google-dns-ipv6:
	sudo ./main/ping -c 5 2001:4860:4860::8888

# ping cloudflare forever
ping-cloudflare-forever:
	sudo ./main/ping cloudflare.com

# ping cloudflare quickly (req/0.01 sec) for 2 seconds
ping-cloudflare-quickly:
	sudo ./main/ping -i 0.01 -t 2 cloudflare.com

# flood google for 1 second
ping-google-flood:
	sudo ./main/ping -t 1 -f google.com

# ping google with a large packet size
ping-google-large-packet:
	sudo ./main/ping -c 5 -s 1000 google.com

# add some more tests

# cleans the executable
clean: 
	rm main/ping