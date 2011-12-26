package rs232

/*
#include <stdlib.h>
#include <fcntl.h>
#include <termios.h>
*/
import "C"


import (
	"os"
	"syscall"
)

type SerialPort struct {
	port *os.File
}

func OpenSerialPort(port string) (rv SerialPort) {
	f, err := os.OpenFile(port,
		os.O_RDWR | os.O_NOCTTY | os.O_NDELAY,
		0666)
	if err != nil {
		panic("Couldn't open port.");
	}
	rv.port = f

	fd := rv.port.Fd()

	var options C.struct_termios
	if C.tcgetattr(C.int(fd), &options) < 0 {
		panic("tcgetattr failed")
	}

	if C.cfsetispeed(&options, C.B57600) < 0 {
		panic("cfsetispeed failed")
	}
	if C.cfsetospeed(&options, C.B57600) < 0 {
		panic("cfsetospeed failed")
	}
	// 8N1
	// options.c_cflag ^= C.PARENB
	// options.c_cflag ^= C.CSTOPB
	// options.c_cflag ^= C.CSIZE
	// options.c_cflag |= C.CS8
	// Local
	options.c_cflag |= (C.CLOCAL | C.CREAD)
	// no hardware flow control
	options.c_cflag ^= C.CRTSCTS

	if C.tcsetattr(C.int(fd), C.TCSANOW, &options) < 0 {
		panic("tcsetattr failed")
	}

	if syscall.SetNonblock(fd, false) != nil {
		panic("Error disabling blocking")
	}

	return
}

func (port *SerialPort) Read(p []byte) (n int, err error) {
	return port.port.Read(p)
}
