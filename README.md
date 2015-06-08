# raspicam Go API
[![Build Status](https://travis-ci.org/dhowden/raspicam.svg?branch=master)](https://travis-ci.org/dhowden/raspicam)
[![GoDoc](https://godoc.org/github.com/dhowden/raspicam?status.svg)](https://godoc.org/github.com/dhowden/raspicam)

Provides a simple Go interface to the [Raspberry Pi Camera board](http://www.raspberrypi.org/archives/tag/camera-board).

For the moment we call the existing command line tools, though eventually hope to hook directly into the C API.

We aim to create an idiomatic Go API whilst mirroring the functionality of the existing command line tools to allow for the best interoperability.

## Status

Implemented:

- Interface to raspistill
- Interface to raspiyuv
- Interface to raspivid

Todo:

- More tests!

## Installation

	go get github.com/dhowden/raspicam

## Getting Started

### Write to a file

	package main
	
	import (
		"fmt"
		"log"
		"os"
	
		"github.com/dhowden/raspicam"
	)
	
	func main() {
		f, err := os.Create(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "create file: %v", err)
			return
		}
		defer f.Close()
	
		s := raspicam.NewStill()
		errCh := make(chan error)
		go func() {
			for x := range errCh {
				fmt.Fprintf(os.Stderr, "%v\n", x)
			}
		}()
		log.Println("Capturing image...")
		raspicam.Capture(s, f, errCh)
	}


### Simple Raspicam TCP Server

	package main

	import (
		"log"
		"fmt"
		"net"
		"os"
	
		"github.com/dhowden/raspicam"
	)
	
	func main() {
		listener, err := net.Listen("tcp", "0.0.0.0:6666")
		if err != nil {
			fmt.Fprintf(os.Stderr, "listen: %v", err)
			return
		}
		log.Println("Listening on 0.0.0.0:6666")

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, "accept: %v", err)
				return
			}
			log.Printf("Accepted connection from: %v\n", conn.RemoteAddr())
			go func() {
				s := raspicam.NewStill()
				errCh := make(chan error)
				go func() {
					for x := range errCh {
						fmt.Fprintf(os.Stderr, "%v\n", x)
				}
				}()
				log.Println("Capturing image...")
				raspicam.Capture(s, conn, errCh)
				log.Println("Done")
				conn.Close()
			}()
		}
	}
