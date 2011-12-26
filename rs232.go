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

func baudConversion(rate int) (flag _Ctype_speed_t) {
	if C.B9600 != 9600 {
		panic("Baud rates may not map directly.")
	}
	return _Ctype_speed_t(rate)
}

type SerConf int

const (
	S_8N1 = iota
	S_7E1
	S_7O1
)

func OpenPort(port string, baudRate int, serconf SerConf) (rv SerialPort, err error) {
	f, open_err := os.OpenFile(port,
		os.O_RDWR | os.O_NOCTTY | os.O_NDELAY,
		0666)
	if open_err != nil {
		err = open_err
		return
	}
	rv.port = f

	fd := rv.port.Fd()

	var options C.struct_termios
	if C.tcgetattr(C.int(fd), &options) < 0 {
		panic("tcgetattr failed")
	}

	if C.cfsetispeed(&options, baudConversion(baudRate)) < 0 {
		panic("cfsetispeed failed")
	}
	if C.cfsetospeed(&options, baudConversion(baudRate)) < 0 {
		panic("cfsetospeed failed")
	}
	switch serconf {
	case S_8N1: {
			options.c_cflag &^= C.PARENB
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS8
		}
	case S_7E1: {
			options.c_cflag |= C.PARENB
			options.c_cflag &^= C.PARODD
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS7
		}
	case S_7O1: {
			options.c_cflag |= C.PARENB
			options.c_cflag |= C.PARODD
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS7
		}
	}
	// Local
	options.c_cflag |= (C.CLOCAL | C.CREAD)
	// no hardware flow control
	options.c_cflag &^= C.CRTSCTS

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

func (port *SerialPort) Write(p []byte) (n int, err error) {
	return port.port.Write(p)
}
