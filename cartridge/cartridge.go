package cartridge

type Cartridge struct {
	prgRom, chrRom []byte
}

func (cartridge *Cartridge) ReadByte(address uint16) (byte, error) {
	return 0, nil // TODO implement
}

func (cartridge *Cartridge) WriteByte(address uint16, data byte) error {
	return nil // TODO implement
}
