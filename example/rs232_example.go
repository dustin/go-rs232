package main

import (
	"log"
	"flag"
	"bufio"
 	"rs232"
)

func main() {
	flag.Parse()
	portString := flag.Args()[0]
	log.Printf("Opening '%s'", portString)
	port, err := rs232.OpenPort(portString, 57600, rs232.S_8N1)
	if err != nil {
		log.Fatalf("Error opening port: %s", err)
	}

	r := bufio.NewReader(&port)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatalf("Error reading:  %s", err)
		}
		log.Printf("<: %s", line)
	}
}
