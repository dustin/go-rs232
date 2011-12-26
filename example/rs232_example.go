package main

import (
	"bufio"
	"flag"
	"github.com/dustin/rs232.go"
	"log"
)

func main() {
	flag.Parse()
	portString := flag.Args()[0]
	log.Printf("Opening '%s'", portString)
	port, err := rs232.OpenPort(portString, 57600, rs232.S_8N1)
	if err != nil {
		log.Fatalf("Error opening port: %s", err)
	}
	defer port.Close()

	r := bufio.NewReader(&port)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatalf("Error reading:  %s", err)
		}
		log.Printf("<: %s", line)
	}
}
