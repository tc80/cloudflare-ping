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

# ping cloudflare 10 times with a 15ms wait time
ping-cloudflare-waittime:
	sudo ./main/ping -c 10 -W 15 cloudflare.com

# flood google for 1 second
ping-google-flood:
	sudo ./main/ping -t 1 -f google.com

# ping google with a large packet size
ping-google-large-packet:
	sudo ./main/ping -c 5 -s 300 google.com

# ping google and always exceed ttl
ping-google-exceed-ttl:
	sudo ./main/ping -c 5 -m 0 google.com

# ping localhost
ping-localhost:
	sudo ./main/ping localhost

# ping localhost with a bunch of options
# 	wait time = 1 ms
#	count = 1000 requests
# 	flood option = true
#	ttl = 0
#	packet size = 80 bytes
#	timeout = 2s
ping-localhost-combo:
	sudo ./main/ping -W 1 -c 1000 -f -m 0 -s 80 -t 2 localhost

# cleans the executable
clean: 
	rm main/ping