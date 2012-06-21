// The RS232 package lets you access old serial junk from go
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

var baudConversionMap = map[int]_Ctype_speed_t{}

// This is your serial port handle.
type SerialPort struct {
	port *os.File
}

func init() {
	baudConversionMap[0] = C.B0
	baudConversionMap[50] = C.B50
	baudConversionMap[75] = C.B75
	baudConversionMap[110] = C.B110
	baudConversionMap[134] = C.B134
	baudConversionMap[150] = C.B150
	baudConversionMap[200] = C.B200
	baudConversionMap[300] = C.B300
	baudConversionMap[600] = C.B600
	baudConversionMap[1200] = C.B1200
	baudConversionMap[1800] = C.B1800
	baudConversionMap[2400] = C.B2400
	baudConversionMap[4800] = C.B4800
	baudConversionMap[9600] = C.B9600
	baudConversionMap[19200] = C.B19200
	baudConversionMap[38400] = C.B38400
	// baudConversionMap[7200] = C.B7200
	// baudConversionMap[14400] = C.B14400
	// baudConversionMap[28800] = C.B28800
	baudConversionMap[57600] = C.B57600
	// baudConversionMap[76800] = C.B76800
	baudConversionMap[115200] = C.B115200
	baudConversionMap[230400] = C.B230400
}

func baudConversion(rate int) (flag _Ctype_speed_t) {
	return baudConversionMap[rate]
}

// SerConf represents the basic serial configuration to provide to OpenPort.
type SerConf int

const (
	S_8N1 = iota
	S_7E1
	S_7O1
)

// Opens and returns a non-blocking serial port.
// The device, baud rate, and SerConf is specified.
//
// Example:  rs232.OpenPort("/dev/ttyS0", 115200, rs232.S_8N1)
func OpenPort(port string, baudRate int, serconf SerConf) (rv SerialPort, err error) {
	f, open_err := os.OpenFile(port,
		syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY,
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
	case S_8N1:
		{
			options.c_cflag &^= C.PARENB
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS8
		}
	case S_7E1:
		{
			options.c_cflag |= C.PARENB
			options.c_cflag &^= C.PARODD
			options.c_cflag &^= C.CSTOPB
			options.c_cflag &^= C.CSIZE
			options.c_cflag |= C.CS7
		}
	case S_7O1:
		{
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

	if syscall.SetNonblock(int(fd), false) != nil {
		panic("Error disabling blocking")
	}

	return
}

// Read from the port.
func (port *SerialPort) Read(p []byte) (n int, err error) {
	return port.port.Read(p)
}

// Write to the port.
func (port *SerialPort) Write(p []byte) (n int, err error) {
	return port.port.Write(p)
}

// Close the port.
func (port *SerialPort) Close() error {
	return port.port.Close()
}
