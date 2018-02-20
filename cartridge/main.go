package cartridge

import (
	"bytes"
	"fmt"
	"github.com/tjarjoura/nes-emulator/types"
	"io"
	"os"
)

func getCartridge(mapperNumber byte, prgRom []byte, chrRom []byte, prgRamSize uint16) (types.MappedHardware, error) {
	var cartridge types.MappedHardware

	prgRam := make([]byte, prgRamSize)

	switch mapperNumber {
	case 0:
		return &NROM{prgRamSize, prgRom, chrRom, prgRam}, nil
	default:
		return cartridge, fmt.Errorf("Unsupported mapper number: %d", mapperNumber)
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
	chrRomSize := uint(header[5]) * 8192

	fmt.Printf("PRG ROM Size: %d\tCHR ROM Rize: %d\n", prgRomSize, chrRomSize)

	var prgRamSize uint16 = 0
	hasPrgRam := (header[6] & 0x02) > 0
	if hasPrgRam {
		if header[8] > 0 {
			prgRamSize = uint16(header[8]) * 8192
		} else {
			prgRamSize = 8192
		}
	}

	fmt.Printf("PRG RAM Size: %x, header[8]: %x\n", prgRamSize, header[8])

	mapperLo := (header[6] & 0xF0) >> 4
	mapperHi := header[7] & 0xF0
	fmt.Printf("Mapper number: %d\n", mapperHi|mapperLo)

	var trainer []byte
	hasTrainer := (header[6] & 0x04) > 0
	if hasTrainer {
		trainer = make([]byte, 512)
		_, err = io.ReadFull(romFile, trainer)
		if err != nil {
			return cartridge, err
		}
	}

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

	return getCartridge(mapperHi|mapperLo, prgRom, chrRom, prgRamSize)
}
