package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/dustin/go-rs232"
)

func parseMode(s string) (rv rs232.SerConf) {
	switch s {
	case "8N1":
		rv = rs232.S_8N1
	case "7E1":
		rv = rs232.S_7E1
	case "7O1":
		rv = rs232.S_7O1
	default:
		log.Fatalf("Invalid mode: %v", s)
	}
	return
}

func main() {
	baudRate := flag.Int("baud", 57600, "Baud rate")
	mode := flag.String("mode", "8N1", "8N1 | 7E1 | 7O1")
	flag.Parse()

	port, err := rs232.OpenPort(flag.Arg(0), *baudRate,
		parseMode(*mode))
	if err != nil {
		log.Fatalf("Error opening port: %s", err)
	}

	io.Copy(os.Stdout, port)
}
