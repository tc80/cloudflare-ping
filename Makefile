# Makefile for building and running the ICMP ping app.
# @author tc80, April 2020
#
# Note, project was built with:
# 	Go version: go1.13.5 darwin/amd64
#	macOS Catalina 10.15.4

# build the project
all: 
	cd main/; go build -o ping

# cleans the executable
clean: 
	rm main/ping