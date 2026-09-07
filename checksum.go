package main

import (
	"hash/crc32"
)

func calcHeaderCRC(data []byte) uint16 {
	crc := crc32.ChecksumIEEE(data)
	return uint16(crc & 0xFFFF)
}