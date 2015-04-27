# serialcat

Serialcat is a simple program that copies all of the data read from a
serial port to stdout.

    Usage: serialcat [args] /dev/tty.whatever
      -baud=57600: Baud rate
      -mode="8N1": 8N1 | 7E1 | 7O1

## Example

    serialcat -baud 115200 /dev/tty.usbserial-A900acmg | tee output

