package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"os"
)

func CreateArchive(outputFilename string, entries []FileEntry) error {
	var archive bytes.Buffer

	// 1. Marker Block
	markerBlock := []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}
	archive.Write(markerBlock)

	// 2. Main Header
	mainHeaderWithoutCRC := []byte{
		0x73, 0x01, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	mainCRC := calcHeaderCRC(mainHeaderWithoutCRC)
	archiveHeader := make([]byte, 2)
	binary.LittleEndian.PutUint16(archiveHeader, mainCRC)
	archive.Write(archiveHeader)
	archive.Write(mainHeaderWithoutCRC)

	// 3. File Blocks
	for _, entry := range entries {
		content, err := os.ReadFile(entry.DiskPath)
		if err != nil {
			return err
		}

		dataLen := uint32(len(content))
		fileCRC := crc32.ChecksumIEEE(content)
		method := byte(0x30) // Store

		archiveNameBytes := []byte(entry.ArchiveName)
		nameSize := uint16(len(archiveNameBytes))
		fileHeadSize := uint16(32 + len(archiveNameBytes))

		fileHeadFields := new(bytes.Buffer)
		fileHeadFields.WriteByte(0x74)                                      // HEAD_TYPE: FILE_HEAD
		binary.Write(fileHeadFields, binary.LittleEndian, uint16(0x8000))   // HEAD_FLAGS
		binary.Write(fileHeadFields, binary.LittleEndian, fileHeadSize)     // HEAD_SIZE
		binary.Write(fileHeadFields, binary.LittleEndian, dataLen)          // PACK_SIZE
		binary.Write(fileHeadFields, binary.LittleEndian, dataLen)          // UNP_SIZE
		fileHeadFields.WriteByte(0x02)                                      // HOST_OS (Win32)
		binary.Write(fileHeadFields, binary.LittleEndian, fileCRC)          // FILE_CRC
		binary.Write(fileHeadFields, binary.LittleEndian, uint32(0x50000000)) // FTIME
		fileHeadFields.WriteByte(29)                                        // UNP_VER
		fileHeadFields.WriteByte(method)                                    // METHOD
		binary.Write(fileHeadFields, binary.LittleEndian, nameSize)         // NAME_SIZE
		binary.Write(fileHeadFields, binary.LittleEndian, uint32(0x20))     // ATTR
		fileHeadFields.Write(archiveNameBytes)                              // FILE_NAME

		fileHeadBytes := fileHeadFields.Bytes()
		fileBlockCRC := calcHeaderCRC(fileHeadBytes)

		fileHeader := make([]byte, 2)
		binary.LittleEndian.PutUint16(fileHeader, fileBlockCRC)

		archive.Write(fileHeader)
		archive.Write(fileHeadBytes)
		archive.Write(content)
	}

	// 4. Terminator Block
	termPayload := []byte{0x7B, 0x00, 0x40, 0x07, 0x00}
	termCrc := calcHeaderCRC(termPayload)
	terminator := make([]byte, 2)
	binary.LittleEndian.PutUint16(terminator, termCrc)
	archive.Write(terminator)
	archive.Write(termPayload)

	return os.WriteFile(outputFilename, archive.Bytes(), 0644)
}