package main

import (
	"github.com/tarm/serial"
	"sync"
	"time"
)

var (
	onceScanner sync.Once
	scanner     *Scanner
)

//Scanner scanner
type Scanner struct {
	port *serial.Port
	id   string
}

//GetScanner get scanner
func GetScanner() *Scanner {

	onceScanner.Do(func() {
		scanner = &Scanner{}
	})

	return scanner

}

//OpenPort
func (scan *Scanner) OpenPort(portName string) error {

	s, err := serial.OpenPort(&serial.Config{
		Name:        portName,
		Baud:        9600,
		ReadTimeout: time.Millisecond * 50,
	})
	if err != nil {
		return err
	}
	scan.port = s
	go func() {

		buffer := make([]byte, 128)
		frame := make([]byte, 128)
		frameLen := 0
		for {

			length, err := s.Read(buffer)
			if err != nil && err.Error() != "EOF" {
				break
			}
			for i := 0; i < length; i++ {
				frame[frameLen] = buffer[i]
				frameLen++
			}
			if length == 0 && frameLen > 0 {
				scan.processFrame(frame[0:frameLen])
				frameLen = 0
			}

			if frameLen >= len(frame) {
				frameLen = 0
			}

		}

	}()

	return nil
}

func (scan *Scanner) processFrame(frame []byte) {

	cmd1 := []byte{0x7e, 0x00, 0x00, 0x05, 0x33, 0x48, 0x30, 0x33, 0x30, 0xb2}
	cmd2 := []byte("NLS0006000;")

	cmd1Match := true
	cmd2Match := true
	for i := 0; i < len(frame); i++ {
		if i >= len(cmd2) {

			cmd2Match = false
		} else if cmd2[i] != frame[i] {
			cmd2Match = false
		}

		if i >= len(cmd1) {

			cmd1Match = false
		} else if cmd1[i] != frame[i] {
			cmd1Match = false
		}
	}

	if cmd1Match {
		scan.port.Write([]byte{0x7e, 0x00, 0x00, 0x05, '1', '2', '3', '4', '5', 0xb2})

	} else if cmd2Match {
		scan.port.Write([]byte{0x06})
	}

}

//SendCode send code
func (scan *Scanner) SendCode(code string) {
	scan.port.Write([]byte(code))
}
