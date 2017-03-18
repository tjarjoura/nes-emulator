package types

type MappedHardware interface {
	ReadByte(address uint16) (byte, error)
	WriteByte(address uint16, data byte) error
}
