package cartridge

type MMC1 struct {
	Cartridge
}

func (mmc1 *MMC1) ReadByte(address uint16) (byte, error) {
	read_byte := mmc1.prgRom[address]
	return read_byte, nil // TODO implement
}

func (mmc1 *MMC1) WriteByte(address uint16, data byte) error {
	return nil // TODO implement
}
