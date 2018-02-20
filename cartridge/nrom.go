package cartridge

import "fmt"

type NROM struct {
	prgRamSize             uint16
	prgRom, chrRom, prgRam []byte
}

func (nrom *NROM) ReadByte(address uint16) (byte, error) {
	if address >= 0x6000 && address <= 0x7FFF {
		if nrom.prgRamSize > 0 {
			return nrom.prgRam[(address-0x6000)%nrom.prgRamSize], nil
		} else {
			return 0x00, fmt.Errorf("NROM.ReadByte(): Unmapped memory address 0x%x.", address)
		}

	} else if address >= 0x8000 && address <= 0xBFFF {
		return nrom.prgRom[address-0x8000], nil

	} else if address >= 0xC000 {
		// If only 16KB of PRG-ROM, then 0x8000-0xBFFF is mirrored by 0xC000-0xFFFF
		if len(nrom.prgRom) > 16384 {
			return nrom.prgRom[address-0x8000], nil
		} else {
			return nrom.prgRom[address-0xC000], nil
		}

	} else { // address < 0x6000
		return 0x00, fmt.Errorf("NROM.ReadByte(): Unmapped memory address 0x%x.", address)
	}
}

func (nrom *NROM) WriteByte(address uint16, data byte) error {
	if address >= 0x6000 && address <= 0x7FFF && nrom.prgRamSize > 0 {
		nrom.prgRam[(address-0x6000)%nrom.prgRamSize] = data
	}

	return nil // TODO implement
}
