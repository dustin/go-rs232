// The RS232 package lets you access old serial junk from go
package rs232

//go:generate ./genc.sh

/*
#include <stdlib.h>
#include <fcntl.h>
#include <termios.h>

void initBaudRates();
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

var baudConversionMap = map[int]_Ctype_speed_t{}

// This is your serial port handle.
type SerialPort struct {
	port *os.File
}

func init() {
	C.initBaudRates()
}

//export addBaudRate
func addBaudRate(num int, val _Ctype_speed_t) {
	baudConversionMap[num] = val
}

func baudConversion(rate int) (flag _Ctype_speed_t) {
	return baudConversionMap[rate]
}

// SerConf represents the basic serial configuration to provide to OpenPort.
type SerConf int

const (
	S_8N1 SerConf = iota
	S_7E1
	S_7O1
)

// Opens and returns a non-blocking serial port.
// The device, baud rate, and SerConf is specified.
//
// Example:  rs232.OpenPort("/dev/ttyS0", 115200, rs232.S_8N1)
func OpenPort(port string, baudRate int, serconf SerConf) (*SerialPort, error) {
	rv := &SerialPort{}
	f, err := os.OpenFile(port,
		syscall.O_RDWR|syscall.O_NOCTTY, 0666)
	if err != nil {
		return nil, err
	}
	rv.port = f

	fd := rv.port.Fd()

	var options C.struct_termios
	if C.tcgetattr(C.int(fd), &options) < 0 {
		defer f.Close()
		return nil, fmt.Errorf("tcgetattr failed")
	}

	if C.cfsetispeed(&options, baudConversion(baudRate)) < 0 {
		defer f.Close()
		return nil, fmt.Errorf("cfsetispeed failed")
	}
	if C.cfsetospeed(&options, baudConversion(baudRate)) < 0 {
		defer f.Close()
		return nil, fmt.Errorf("cfsetospeed failed")
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
	// Don't EOF on a zero read, just block
	options.c_cc[C.VMIN] = 1

	if C.tcsetattr(C.int(fd), C.TCSANOW, &options) < 0 {
		defer f.Close()
		return nil, fmt.Errorf("tcsetattr failed")
	}

	return rv, nil
}

// SetInputAttr sets VMIN and VTIME for control serial reads.
//
// In non-canonical input processing mode, input is not assembled into
// lines and input processing (erase, kill, delete, etc.) does not
// occur. Two parameters control the behavior of this mode:
// c_cc[VTIME] sets the character timer, and c_cc[VMIN] sets the
// minimum number of characters to receive before satisfying the read.
//
// If MIN > 0 and TIME = 0, MIN sets the number of characters to
//receive before the read is satisfied. As TIME is zero, the timer is
//not used.
//
// If MIN = 0 and TIME > 0, TIME serves as a timeout value. The read
//will be satisfied if a single character is read, or TIME is exceeded
//(t = TIME *0.1 s). If TIME is exceeded, no character will be
//returned.
//
// If MIN > 0 and TIME > 0, TIME serves as an inter-character
//timer. The read will be satisfied if MIN characters are received, or
//the time between two characters exceeds TIME. The timer is restarted
//every time a character is received and only becomes active after the
//first character has been received.
//
// If MIN = 0 and TIME = 0, read will be satisfied immediately. The
//number of characters currently available, or the number of
//characters requested will be returned. According to Antonino (see
//contributions), you could issue a fcntl(fd, F_SETFL, FNDELAY);
//before reading to get the same result.
//
// By modifying newtio.c_cc[VTIME] and newtio.c_cc[VMIN] all modes
//described above can be tested.
//
// -- copied from http://tldp.org/HOWTO/Serial-Programming-HOWTO/x115.html
func (port *SerialPort) SetInputAttr(minBytes int, timeout time.Duration) error {
	fd := port.port.Fd()

	var options C.struct_termios
	if C.tcgetattr(C.int(fd), &options) < 0 {
		return fmt.Errorf("tcgetattr failed")
	}
	options.c_cc[C.VMIN] = _Ctype_cc_t(minBytes)
	options.c_cc[C.VTIME] = _Ctype_cc_t(timeout / (time.Second / 10))

	if C.tcsetattr(C.int(fd), C.TCSANOW, &options) < 0 {
		return fmt.Errorf("tcsetattr failed")
	}
	return nil
}

// SetNonblock enables nonblocking serial IO.
func (port *SerialPort) SetNonblock() error {
	fd := port.port.Fd()

	return syscall.SetNonblock(int(fd), true)
}

// Read from the port.
func (port *SerialPort) Read(p []byte) (int, error) {
	return port.port.Read(p)
}

// Write to the port.
func (port *SerialPort) Write(p []byte) (int, error) {
	return port.port.Write(p)
}

// Close the port.
func (port *SerialPort) Close() error {
	return port.port.Close()
}
