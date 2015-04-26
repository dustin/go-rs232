package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dustin/go-rs232"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [args] /dev/tty.whatever\n", os.Args[0])
		flag.PrintDefaults()
	}
}

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

	path := flag.Arg(0)
	if path == "" {
		flag.Usage()
		os.Exit(64)
	}

	port, err := rs232.OpenPort(path, *baudRate, parseMode(*mode))
	if err != nil {
		log.Fatalf("Error opening port %q: %s", port, err)
	}

	io.Copy(os.Stdout, port)
}
