package cartridge

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Cartridge struct {
	PrgRom, chrRom []byte
}

func CartridgeFromFile(filename string) (Cartridge, error) {
	var cartridge Cartridge

	romFile, err := os.Open(filename)
	if err != nil {
		return cartridge, err
	}
	defer romFile.Close()

	header := make([]byte, 16)
	romFile.Read(header)

	inesMagic := []byte{'N', 'E', 'S', 0x1A}
	if !bytes.Equal(header[0:4], inesMagic) {
		return cartridge, fmt.Errorf("%s: Unrecognized file format", filename)
	}

	prgRomSize := uint(header[4]) * 16384
	fmt.Printf("PRG ROM: %d bytes\n", prgRomSize)
	chrRomSize := uint(header[5]) * 16384
	fmt.Printf("CHR ROM: %d bytes\n", chrRomSize)

	mapperLo := (header[6] & 0xF0) >> 4
	mapperHi := header[7] & 0xF0
	fmt.Printf("Mapper number: %d\n", mapperHi|mapperLo)

	if prgRomSize > 0 {
		cartridge.PrgRom = make([]byte, prgRomSize)
		_, err := io.ReadFull(romFile, cartridge.PrgRom)
		if err != nil {
			return cartridge, err
		}
	}

	if chrRomSize > 0 {
		cartridge.chrRom = make([]byte, chrRomSize)
		_, err := io.ReadFull(romFile, cartridge.chrRom)
		if err != nil {
			return cartridge, err
		}
	}

	return cartridge, nil
}
