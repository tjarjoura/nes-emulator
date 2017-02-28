package cpu

const (
	MODE_IMPLIED int = iota
	MODE_ACCUMULATOR
	MODE_IMMEDIATE
	MODE_ZERO_PAGE
	MODE_ABSOLUTE
	MODE_RELATIVE
	MODE_INDEX_INDIRECT
	MODE_INDIRECT_INDEX
	REG_NONE
	REG_X
	REG_Y
)

type addressMode struct{ mode, reg int }

type instruction struct {
	handler func(cpu *Cpu, arg uint16) error
	addressMode
	neumonic string
}

func asl(cpu *Cpu, arg uint16) error {
	target := cpu.byteAt(arg)
	cpu.carryFl = target&0x80 > 0

	err := cpu.writeByte(arg, target<<8)
	return err
}

func brk(cpu *Cpu, arg uint16) error {
	return nil
}

func ora(cpu *Cpu, arg uint16) error {
	cpu.a |= cpu.byteAt(arg)
	cpu.zeroFl = (cpu.a == 0)
	cpu.signFl = (int8(cpu.a) < 0)

	return nil
}

var instructions = map[byte]instruction{
	0x00: instruction{brk, addressMode{MODE_IMPLIED, REG_NONE}, "BRK\n"},
	0x01: instruction{ora, addressMode{MODE_INDEX_INDIRECT, REG_X}, "ORA ($%x, X)\n"},
	0x05: instruction{ora, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ORA $%x\n"},
	0x06: instruction{asl, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ASL $%x\n"},
}
