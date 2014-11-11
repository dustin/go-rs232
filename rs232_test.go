package rs232

import "io"

var _ = io.ReadWriteCloser(&SerialPort{})
