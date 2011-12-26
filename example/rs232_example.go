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
	port := rs232.OpenSerialPort(portString)

	r := bufio.NewReader(&port)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatalf("Error reading:  %s", err)
		}
		log.Printf("<: %s", line)
	}
}
