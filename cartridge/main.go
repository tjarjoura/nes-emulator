package cartridge

import (
	"bytes"
	"fmt"
	"github.com/tjarjoura/nes-emulator/types"
	"io"
	"os"
)

func getCartridge(mapperNumber byte, prgRom []byte, chrRom []byte) (types.MappedHardware, error) {
	var cartridge types.MappedHardware

	switch mapperNumber {
	case 1:
		return new(MMC1), nil
	default:
		return cartridge, fmt.Errorf("Unsupported mapper number: %d\n", mapperNumber)
	}
}

func CartridgeFromFile(filename string) (types.MappedHardware, error) {
	var cartridge types.MappedHardware

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

	prgRom := make([]byte, prgRomSize)
	_, err = io.ReadFull(romFile, prgRom)
	if err != nil {
		return cartridge, err
	}

	chrRom := make([]byte, chrRomSize)
	_, err = io.ReadFull(romFile, chrRom)
	if err != nil {
		return cartridge, err
	}

	return getCartridge(mapperHi|mapperLo, prgRom, chrRom)
}
