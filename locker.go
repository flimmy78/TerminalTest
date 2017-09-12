package main

import (
	"github.com/tarm/serial"
	"sync"
)

//CabinetManager a manaer
type CabinetManager struct {
	portName string
	port     *serial.Port
	cabinets []Cabinet
	frames   chan []byte
}

//Cabinet cabinet
type Cabinet struct {
	cabMng     *CabinetManager
	serialName string
	address    int
	boxes      []Box

	state int
}

//Box a box
type Box struct {
	name      string
	boxNum    int
	openState int //0 open 1 close
	state     int //0 normal 1 bad
}

var (
	cabMng *CabinetManager
	once   sync.Once
)

// GetCabMng create cab instance
func GetCabMng() *CabinetManager {

	once.Do(func() {

		if cabMng == nil {
			cabMng = &CabinetManager{
				frames: make(chan []byte, 1),
			}

		}
	})
	return cabMng
}

//OpenPort open serialport
func (cabMng *CabinetManager) OpenPort() int {

	s, err := serial.OpenPort(&serial.Config{Name: cabMng.portName, Baud: 9600})
	if err != nil {
		return -1
	}
	cabMng.port = s
	go func(cabMng *CabinetManager) {
		frame := make([]byte, 256)
		frameLen := 0

		buffer := make([]byte, 128)
		for {

			lenData, err := s.Read(buffer)
			if err != nil {
				break
			}

			for i := 0; i < lenData; i++ {
				if buffer[i] == 0xaa {

					frame[0] = 0xaa
					frameLen = 1
				} else if frameLen > 0 {

					if buffer[i] == 0x55 {

						frame[frameLen] = buffer[i]
						frameLen++
						f := make([]byte, frameLen)
						copy(f, frame[:frameLen])
						cabMng.frames <- f

					} else {

						frame[frameLen] = buffer[i]
						frameLen++
					}

				}

				if frameLen >= len(frame) {
					frameLen = 0
				}

			}
		}

	}(cabMng)

	go func() {

		for {
			select {
			case frame := <-cabMng.frames:

				if len(frame) > 4 {
					addr := int(frame[3])
					if addr < len(cabMng.cabinets) {

						switch frame[4] {
						case 1:
							boxIndex := -1
							for i := 0; i < 3; i++ {
								for j := uint(0); j < 8; j++ {
									if (frame[5+i] & (1 << j)) != 0 {
										boxIndex = (2-i)*8 + int(j)
										break
									}
								}
							}
							if boxIndex >= 0 {
								cabMng.cabinets[addr].OpenBox(boxIndex)
							}
							break
						case 2:
							cabMng.cabinets[addr].QueryStatus()
							break
						}

					}
				}

			}
		}
	}()

	return 0

}

// AddCabnit add cabniet
func (cabMng *CabinetManager) AddCabnit(addr, boxNum int) {

	cab := Cabinet{
		cabMng:  cabMng,
		address: addr,
		boxes:   make([]Box, boxNum),
		state:   0,
	}
	cabMng.cabinets = append(cabMng.cabinets, cab)
}

//MarkCab state :0 OK 1:broken
func (cabMng *CabinetManager) MarkCab(addr int, state int) {

	cabMng.cabinets[addr].state = state

}

//MarkBox state 0: ok 1: broken
func (cabMng *CabinetManager) MarkBox(addr, boxNum int, state, openState int) {
	cabMng.cabinets[addr].boxes[boxNum].state = state
	cabMng.cabinets[addr].boxes[boxNum].openState = openState
}

//AckRequeest ack req
func (cabMng *CabinetManager) AckRequeest(addr, cmd byte, param []byte) {

	cnt := 8
	if param != nil {
		cnt += len(param)
	}

	frame := make([]byte, cnt)

	frame[0] = 0xaa //帧头
	cnt -= 8
	frame[1] = byte(cnt >> 8)
	frame[2] = byte(cnt & 0xff)
	frame[3] = addr
	frame[4] = cmd

	for i := 0; i < cnt; i++ {
		frame[5+i] = param[i]
	}

	//CRC
	crc := CRC16(frame, 1, 4+cnt)
	frame[5+cnt] = (byte)(crc >> 8)
	frame[6+cnt] = (byte)(crc & 0xff)

	frame[7+cnt] = 0x55 //帧尾

	cabMng.port.Write(frame)
}

//OpenBox simulate openbox
func (cab *Cabinet) OpenBox(index int) {
	if cab.state == 0 {

		if cab.boxes[index].state == 0 {
			cab.boxes[index].openState = 1
		}
		cab.cabMng.AckRequeest(byte(cab.address), 1, []byte{0})
	}
}

//QueryStatus query status
func (cab *Cabinet) QueryStatus() {

	if cab.state == 0 {

		status := make([]byte, 6)
		for i := 0; i < 24; i++ {

			if cab.boxes[i].openState == 0 {
				byteIndex := 2 - i/8
				bitIndex := uint(i % 8)
				status[byteIndex] = byte(status[byteIndex] | (1 << bitIndex))
			}

		}
		cab.cabMng.AckRequeest(byte(cab.address), 2, status)

	}

}

//SetAllBox close all box
func (cabMng *CabinetManager) SetAllBox(state int) {

	cabNum := len(cabMng.cabinets)
	for i := 0; i < cabNum; i++ {

		boxNum := len(cabMng.cabinets[i].boxes)
		for j := 0; j < boxNum; j++ {
			cabMng.cabinets[i].boxes[j].openState = state
		}
	}

}
