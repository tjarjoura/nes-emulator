package cartridge

type MMC1 struct {
	prgRom, chrRom []byte
}

func (mmc1 *MMC1) ReadByte(address uint16) (byte, error) {
}

func (mmc1 *MMC1) WriteByte(address uint16, data byte) error {
}
