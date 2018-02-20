package cpu

const (
	MODE_IMPLIED int = iota
	MODE_ACCUMULATOR
	MODE_IMMEDIATE
	MODE_ZERO_PAGE
	MODE_ABSOLUTE
	MODE_RELATIVE
	MODE_INDIRECT
	MODE_INDEX_INDIRECT
	MODE_INDIRECT_INDEX
	REG_NONE
	REG_X
	REG_Y
)

type addressMode struct{ mode, reg int }

type instruction struct {
	handler func(cpu *Cpu, arg uint16, mode addressMode) error
	addressMode
	neumonic string
}

func adc(cpu *Cpu, arg uint16, mode addressMode) error {
	var result, carryVal uint16

	if cpu.carryFl {
		carryVal = 1
	} else {
		carryVal = 0
	}

	if mode.mode == MODE_IMMEDIATE {
		result = carryVal + uint16(cpu.a) + arg
		cpu.overflowFl = ((cpu.a ^ byte(result)) & (byte(arg) ^ byte(result)) & 0x80) != 0
	} else {
		result = carryVal + uint16(cpu.a) + uint16(cpu.byteAt(arg))
		cpu.overflowFl = ((cpu.a ^ byte(result)) & (cpu.byteAt(arg) ^ byte(result)) & 0x80) != 0
	}

	cpu.a = byte(result)

	cpu.carryFl = (result > 0xFF)
	cpu.zeroFl = (cpu.a == 0)
	cpu.signFl = (int8(cpu.a) < 0)

	return nil
}

func and(cpu *Cpu, arg uint16, mode addressMode) error {
	if mode.mode == MODE_IMMEDIATE {
		cpu.a &= byte(arg)
	} else {
		cpu.a &= cpu.byteAt(arg)
	}

	cpu.zeroFl = (cpu.a == 0)
	cpu.signFl = (int8(cpu.a) < 0)

	return nil
}

func asl(cpu *Cpu, arg uint16, mode addressMode) error {
	var oldCarry byte = 0x0
	if cpu.carryFl {
		oldCarry = 0x1
	}

	if mode.mode == MODE_ACCUMULATOR {
		cpu.carryFl = cpu.a&0x80 > 0
		cpu.a <<= 1
		cpu.a |= oldCarry

		return nil

	} else {
		target := cpu.byteAt(arg)
		cpu.carryFl = target&0x80 > 0
		target <<= 1
		target |= oldCarry

		err := cpu.writeByte(arg, target)
		return err
	}
}

func bcc(cpu *Cpu, arg uint16, mode addressMode) error {
	if !cpu.carryFl {
		cpu.pc = arg
	}

	return nil
}

func bcs(cpu *Cpu, arg uint16, mode addressMode) error {
	if cpu.carryFl {
		cpu.pc = arg
	}

	return nil
}

func beq(cpu *Cpu, arg uint16, mode addressMode) error {
	if cpu.zeroFl {
		cpu.pc = arg
	}

	return nil
}

func bit(cpu *Cpu, arg uint16, mode addressMode) error {
	target := cpu.byteAt(arg)

	result := cpu.a & target

	cpu.zeroFl = result == 0
	cpu.overflowFl = target&0x80 > 0
	cpu.signFl = target&0xF0 > 0

	return nil
}

func bmi(cpu *Cpu, arg uint16, mode addressMode) error {
	if cpu.signFl {
		cpu.pc = arg
	}

	return nil
}

func bne(cpu *Cpu, arg uint16, mode addressMode) error {
	if !cpu.zeroFl {
		cpu.pc = arg
	}

	return nil
}

func bpl(cpu *Cpu, arg uint16, mode addressMode) error {
	if !cpu.signFl {
		cpu.pc = arg
	}

	return nil
}

func brk(cpu *Cpu, arg uint16, mode addressMode) error {
	return nil
}

func bvc(cpu *Cpu, arg uint16, mode addressMode) error {
	if !cpu.overflowFl {
		cpu.pc = arg
	}

	return nil
}

func bvs(cpu *Cpu, arg uint16, mode addressMode) error {
	if cpu.overflowFl {
		cpu.pc = arg
	}

	return nil
}

func clc(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.carryFl = false
	return nil
}

func cld(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.decimalFl = false
	return nil
}

func cli(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.interruptFl = false
	return nil
}

func clv(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.overflowFl = false
	return nil
}

func cmp(cpu *Cpu, arg uint16, mode addressMode) error {
	var result int8
	if mode.mode == MODE_IMMEDIATE {
		result = int8(cpu.a - byte(arg))
	} else {
		result = int8(cpu.a - cpu.byteAt(arg))
	}

	cpu.carryFl = (result >= 0)
	cpu.signFl = (result < 0)
	cpu.zeroFl = (result == 0)

	return nil
}

func cpx(cpu *Cpu, arg uint16, mode addressMode) error {
	var result int8
	if mode.mode == MODE_IMMEDIATE {
		result = int8(cpu.x - byte(arg))
	} else {
		result = int8(cpu.x - cpu.byteAt(arg))
	}

	cpu.carryFl = (result >= 0)
	cpu.signFl = (result < 0)
	cpu.zeroFl = (result == 0)

	return nil
}

func cpy(cpu *Cpu, arg uint16, mode addressMode) error {
	var result int8
	if mode.mode == MODE_IMMEDIATE {
		result = int8(cpu.y - byte(arg))
	} else {
		result = int8(cpu.y - cpu.byteAt(arg))
	}

	cpu.carryFl = (result >= 0)
	cpu.signFl = (result < 0)
	cpu.zeroFl = (result == 0)

	return nil
}

func dec(cpu *Cpu, arg uint16, mode addressMode) error {
	data := cpu.byteAt(arg) - 1
	cpu.zeroFl = data == 0
	cpu.signFl = int8(data) < 0
	cpu.writeByte(arg, data)
	return nil
}

func dex(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.x -= 1
	cpu.zeroFl = cpu.x == 0
	cpu.signFl = int8(cpu.x) < 0
	return nil
}

func dey(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.y -= 1
	cpu.zeroFl = cpu.y == 0
	cpu.signFl = int8(cpu.y) < 0
	return nil
}

func eor(cpu *Cpu, arg uint16, mode addressMode) error {
	if mode.mode == MODE_IMMEDIATE {
		cpu.a ^= byte(arg)
	} else {
		cpu.a ^= cpu.byteAt(arg)
	}

	cpu.zeroFl = cpu.a == 0
	cpu.signFl = int8(cpu.a) < 0
	return nil
}

func inc(cpu *Cpu, arg uint16, mode addressMode) error {
	data := cpu.byteAt(arg) + 1
	cpu.zeroFl = data == 0
	cpu.signFl = int8(data) < 0
	cpu.writeByte(arg, data)
	return nil
}

func inx(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.x += 1
	cpu.zeroFl = cpu.x == 0
	cpu.signFl = int8(cpu.x) < 0
	return nil
}

func iny(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.y += 1
	cpu.zeroFl = cpu.y == 0
	cpu.signFl = int8(cpu.y) < 0
	return nil
}

func jmp(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.pc = arg
	return nil
}

func jsr(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.pushWordToStack(cpu.pc)
	cpu.pc = arg
	return nil
}

func lda(cpu *Cpu, arg uint16, mode addressMode) error {
	if mode.mode == MODE_IMMEDIATE {
		cpu.a = byte(arg)
	} else {
		cpu.a = cpu.byteAt(arg)
	}

	cpu.zeroFl = cpu.a == 0
	cpu.signFl = int8(cpu.a) < 0

	return nil
}

func ldx(cpu *Cpu, arg uint16, mode addressMode) error {
	if mode.mode == MODE_IMMEDIATE {
		cpu.x = byte(arg)
	} else {
		cpu.x = cpu.byteAt(arg)
	}

	cpu.zeroFl = cpu.x == 0
	cpu.signFl = int8(cpu.x) < 0

	return nil
}

func ldy(cpu *Cpu, arg uint16, mode addressMode) error {
	if mode.mode == MODE_IMMEDIATE {
		cpu.y = byte(arg)
	} else {
		cpu.y = cpu.byteAt(arg)
	}

	cpu.zeroFl = cpu.y == 0
	cpu.signFl = int8(cpu.y) < 0

	return nil
}

func lsr(cpu *Cpu, arg uint16, mode addressMode) error {
	var result byte
	if mode.mode == MODE_ACCUMULATOR {
		cpu.carryFl = cpu.a&0x1 > 0
		result = cpu.a >> 1
		cpu.a = result
	} else {
		data := cpu.byteAt(arg)
		cpu.carryFl = data&0x1 > 0
		result = data >> 1
		cpu.writeByte(arg, result)
	}

	cpu.zeroFl = result == 0
	cpu.signFl = result < 0
	return nil
}

func nop(cpu *Cpu, arg uint16, mode addressMode) error {
	// This instruction intentionally does nothing
	return nil
}

func ora(cpu *Cpu, arg uint16, mode addressMode) error {
	var operand byte

	if mode.mode == MODE_IMMEDIATE {
		operand = byte(arg)
	} else {
		operand = cpu.byteAt(arg)
	}

	cpu.a |= operand
	cpu.zeroFl = (cpu.a == 0)
	cpu.signFl = (int8(cpu.a) < 0)

	return nil
}

func pha(cpu *Cpu, arg uint16, mode addressMode) error {
	err := cpu.pushByteToStack(cpu.a)
	return err
}

func php(cpu *Cpu, arg uint16, mode addressMode) error {
	statusFlagsByte := cpu.getStatusFlagsByte()
	statusFlagsByte &= 0x10 // Set B flag
	err := cpu.pushByteToStack(statusFlagsByte)
	return err
}

func pla(cpu *Cpu, arg uint16, mode addressMode) error {
	var err error
	cpu.a, err = cpu.pullByteFromStack()
	return err
}

func plp(cpu *Cpu, arg uint16, mode addressMode) error {
	statusFlagsByte, err := cpu.pullByteFromStack()
	if err != nil {
		return err
	}

	cpu.restoreStatusFlags(statusFlagsByte)
	return nil
}

func rol(cpu *Cpu, arg uint16, mode addressMode) error {
	var oldCarry byte = 0x0
	if cpu.carryFl {
		oldCarry = 0x1
	}

	if mode.mode == MODE_ACCUMULATOR {
		cpu.carryFl = cpu.a&0x80 > 0
		cpu.a <<= 1
		cpu.a |= oldCarry
		cpu.zeroFl = cpu.a == 0
		cpu.signFl = int8(cpu.a) < 0

	} else {
		target := cpu.byteAt(arg)
		cpu.carryFl = target&0x80 > 0
		target <<= 1
		target |= oldCarry
		cpu.zeroFl = target == 0
		cpu.signFl = int8(target) < 0
		cpu.writeByte(arg, target)
	}

	return nil
}

func ror(cpu *Cpu, arg uint16, mode addressMode) error {
	var oldCarry byte = 0x0
	if cpu.carryFl {
		oldCarry = 0x80
	}

	if mode.mode == MODE_ACCUMULATOR {
		cpu.carryFl = cpu.a&0x01 > 0
		cpu.a >>= 1
		cpu.a |= oldCarry
		cpu.zeroFl = cpu.a == 0
		cpu.signFl = int8(cpu.a) < 0

	} else {
		target := cpu.byteAt(arg)
		cpu.carryFl = target&0x01 > 0
		target >>= 1
		target |= oldCarry
		cpu.zeroFl = target == 0
		cpu.signFl = int8(target) < 0
		cpu.writeByte(arg, target)
	}

	return nil
}

func rti(cpu *Cpu, arg uint16, mode addressMode) error {
	statusFlagsByte, err := cpu.pullByteFromStack()
	if err != nil {
		return err
	}

	cpu.restoreStatusFlags(statusFlagsByte)
	savedPC, err := cpu.pullWordFromStack()
	if err != nil {
		return err
	}

	cpu.pc = savedPC
	return nil
}

func rts(cpu *Cpu, arg uint16, mode addressMode) error {
	savedPC, err := cpu.pullWordFromStack()
	if err != nil {
		return err
	}

	cpu.pc = savedPC
	return nil
}

func sbc(cpu *Cpu, arg uint16, mode addressMode) error {
	var result, carryVal uint16

	if cpu.carryFl {
		carryVal = 1
	} else {
		carryVal = 0
	}

	if mode.mode == MODE_IMMEDIATE {
		result = uint16(cpu.a) - arg - (1 - carryVal)
		cpu.overflowFl = ((cpu.a ^ byte(result)) & (byte(arg) ^ byte(result)) & 0x80) != 0
	} else {
		result = uint16(cpu.a) - uint16(cpu.byteAt(arg)) - (1 - carryVal)
		cpu.overflowFl = ((cpu.a ^ byte(result)) & (cpu.byteAt(arg) ^ byte(result)) & 0x80) != 0
	}

	cpu.a = byte(result)

	cpu.carryFl = (result <= 0xFF)
	cpu.zeroFl = (cpu.a == 0)
	cpu.signFl = (int8(cpu.a) < 0)

	return nil
}

func sec(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.carryFl = true
	return nil
}

func sed(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.decimalFl = true
	return nil
}

func sei(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.interruptFl = true
	return nil
}

func sta(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.writeByte(arg, cpu.a)
	return nil
}

func stx(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.writeByte(arg, cpu.x)
	return nil
}

func sty(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.writeByte(arg, cpu.y)
	return nil
}

func tax(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.x = cpu.a
	cpu.zeroFl = cpu.x == 0
	cpu.signFl = int8(cpu.x) < 0
	return nil
}

func tay(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.y = cpu.a
	cpu.zeroFl = cpu.y == 0
	cpu.signFl = int8(cpu.y) < 0
	return nil
}

func tsx(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.x = cpu.sp
	cpu.zeroFl = cpu.x == 0
	cpu.signFl = int8(cpu.x) < 0
	return nil
}

func txa(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.a = cpu.x
	cpu.zeroFl = cpu.a == 0
	cpu.signFl = int8(cpu.a) < 0
	return nil
}

func txs(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.sp = cpu.x
	return nil
}

func tya(cpu *Cpu, arg uint16, mode addressMode) error {
	cpu.a = cpu.y
	cpu.zeroFl = cpu.a == 0
	cpu.signFl = int8(cpu.a) < 0
	return nil
}

// The following code was auto generated based on this table:
// http://www.thealmightyguru.com/Games/Hacking/Wiki/index.php/6502_Opcodes
var instructions = map[byte]instruction{
	0x00: instruction{brk, addressMode{MODE_IMPLIED, REG_NONE}, "BRK"},
	0x01: instruction{ora, addressMode{MODE_INDEX_INDIRECT, REG_X}, "ORA"},
	0x05: instruction{ora, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ORA"},
	0x06: instruction{asl, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ASL"},
	0x08: instruction{php, addressMode{MODE_IMPLIED, REG_NONE}, "PHP"},
	0x09: instruction{ora, addressMode{MODE_IMMEDIATE, REG_NONE}, "ORA"},
	0x0A: instruction{asl, addressMode{MODE_ACCUMULATOR, REG_NONE}, "ASL"},
	0x0D: instruction{ora, addressMode{MODE_ABSOLUTE, REG_NONE}, "ORA"},
	0x0E: instruction{asl, addressMode{MODE_ABSOLUTE, REG_NONE}, "ASL"},
	0x10: instruction{bpl, addressMode{MODE_IMPLIED, REG_NONE}, "BPL"},
	0x11: instruction{ora, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "ORA"},
	0x15: instruction{ora, addressMode{MODE_ZERO_PAGE, REG_X}, "ORA"},
	0x16: instruction{asl, addressMode{MODE_ZERO_PAGE, REG_X}, "ASL"},
	0x18: instruction{clc, addressMode{MODE_IMPLIED, REG_NONE}, "CLC"},
	0x19: instruction{ora, addressMode{MODE_ABSOLUTE, REG_Y}, "ORA"},
	0x1D: instruction{ora, addressMode{MODE_ABSOLUTE, REG_X}, "ORA"},
	0x1E: instruction{asl, addressMode{MODE_ABSOLUTE, REG_X}, "ASL"},
	0x20: instruction{jsr, addressMode{MODE_IMPLIED, REG_NONE}, "JSR"},
	0x21: instruction{and, addressMode{MODE_INDEX_INDIRECT, REG_X}, "AND"},
	0x24: instruction{bit, addressMode{MODE_ZERO_PAGE, REG_NONE}, "BIT"},
	0x25: instruction{and, addressMode{MODE_ZERO_PAGE, REG_NONE}, "AND"},
	0x26: instruction{rol, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ROL"},
	0x28: instruction{plp, addressMode{MODE_IMPLIED, REG_NONE}, "PLP"},
	0x29: instruction{and, addressMode{MODE_IMMEDIATE, REG_NONE}, "AND"},
	0x2A: instruction{rol, addressMode{MODE_ACCUMULATOR, REG_NONE}, "ROL"},
	0x2C: instruction{bit, addressMode{MODE_ABSOLUTE, REG_NONE}, "BIT"},
	0x2D: instruction{and, addressMode{MODE_ABSOLUTE, REG_NONE}, "AND"},
	0x2E: instruction{rol, addressMode{MODE_ABSOLUTE, REG_NONE}, "ROL"},
	0x30: instruction{bmi, addressMode{MODE_IMPLIED, REG_NONE}, "BMI"},
	0x31: instruction{and, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "AND"},
	0x35: instruction{and, addressMode{MODE_ZERO_PAGE, REG_X}, "AND"},
	0x36: instruction{rol, addressMode{MODE_ZERO_PAGE, REG_X}, "ROL"},
	0x38: instruction{sec, addressMode{MODE_IMPLIED, REG_NONE}, "SEC"},
	0x39: instruction{and, addressMode{MODE_ABSOLUTE, REG_Y}, "AND"},
	0x3D: instruction{and, addressMode{MODE_ABSOLUTE, REG_X}, "AND"},
	0x3E: instruction{rol, addressMode{MODE_ABSOLUTE, REG_X}, "ROL"},
	0x40: instruction{rti, addressMode{MODE_IMPLIED, REG_NONE}, "RTI"},
	0x41: instruction{eor, addressMode{MODE_INDEX_INDIRECT, REG_X}, "EOR"},
	0x45: instruction{eor, addressMode{MODE_ZERO_PAGE, REG_NONE}, "EOR"},
	0x46: instruction{lsr, addressMode{MODE_ZERO_PAGE, REG_NONE}, "LSR"},
	0x48: instruction{pha, addressMode{MODE_IMPLIED, REG_NONE}, "PHA"},
	0x49: instruction{eor, addressMode{MODE_IMMEDIATE, REG_NONE}, "EOR"},
	0x4A: instruction{lsr, addressMode{MODE_ACCUMULATOR, REG_NONE}, "LSR"},
	0x4C: instruction{jmp, addressMode{MODE_ABSOLUTE, REG_NONE}, "JMP"},
	0x4D: instruction{eor, addressMode{MODE_ABSOLUTE, REG_NONE}, "EOR"},
	0x4E: instruction{lsr, addressMode{MODE_ABSOLUTE, REG_NONE}, "LSR"},
	0x50: instruction{bvc, addressMode{MODE_IMPLIED, REG_NONE}, "BVC"},
	0x51: instruction{eor, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "EOR"},
	0x55: instruction{eor, addressMode{MODE_ZERO_PAGE, REG_X}, "EOR"},
	0x56: instruction{lsr, addressMode{MODE_ZERO_PAGE, REG_X}, "LSR"},
	0x58: instruction{cli, addressMode{MODE_IMPLIED, REG_NONE}, "CLI"},
	0x59: instruction{eor, addressMode{MODE_ABSOLUTE, REG_Y}, "EOR"},
	0x5D: instruction{eor, addressMode{MODE_ABSOLUTE, REG_X}, "EOR"},
	0x5E: instruction{lsr, addressMode{MODE_ABSOLUTE, REG_X}, "LSR"},
	0x60: instruction{rts, addressMode{MODE_IMPLIED, REG_NONE}, "RTS"},
	0x61: instruction{adc, addressMode{MODE_INDEX_INDIRECT, REG_X}, "ADC"},
	0x65: instruction{adc, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ADC"},
	0x66: instruction{ror, addressMode{MODE_ZERO_PAGE, REG_NONE}, "ROR"},
	0x68: instruction{pla, addressMode{MODE_IMPLIED, REG_NONE}, "PLA"},
	0x69: instruction{adc, addressMode{MODE_IMMEDIATE, REG_NONE}, "ADC"},
	0x6A: instruction{ror, addressMode{MODE_ACCUMULATOR, REG_NONE}, "ROR"},
	0x6C: instruction{jmp, addressMode{MODE_INDIRECT, REG_NONE}, "JMP"},
	0x6D: instruction{adc, addressMode{MODE_ABSOLUTE, REG_NONE}, "ADC"},
	0x6E: instruction{ror, addressMode{MODE_ABSOLUTE, REG_NONE}, "ROR"},
	0x70: instruction{bvs, addressMode{MODE_IMPLIED, REG_NONE}, "BVS"},
	0x71: instruction{adc, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "ADC"},
	0x75: instruction{adc, addressMode{MODE_ZERO_PAGE, REG_X}, "ADC"},
	0x76: instruction{ror, addressMode{MODE_ZERO_PAGE, REG_X}, "ROR"},
	0x78: instruction{sei, addressMode{MODE_IMPLIED, REG_NONE}, "SEI"},
	0x79: instruction{adc, addressMode{MODE_ABSOLUTE, REG_Y}, "ADC"},
	0x7D: instruction{adc, addressMode{MODE_ABSOLUTE, REG_X}, "ADC"},
	0x7E: instruction{ror, addressMode{MODE_ABSOLUTE, REG_X}, "ROR"},
	0x81: instruction{sta, addressMode{MODE_INDEX_INDIRECT, REG_X}, "STA"},
	0x84: instruction{sty, addressMode{MODE_ZERO_PAGE, REG_NONE}, "STY"},
	0x85: instruction{sta, addressMode{MODE_ZERO_PAGE, REG_NONE}, "STA"},
	0x86: instruction{stx, addressMode{MODE_ZERO_PAGE, REG_NONE}, "STX"},
	0x88: instruction{dey, addressMode{MODE_IMPLIED, REG_NONE}, "DEY"},
	0x8A: instruction{txa, addressMode{MODE_IMPLIED, REG_NONE}, "TXA"},
	0x8C: instruction{sty, addressMode{MODE_ABSOLUTE, REG_NONE}, "STY"},
	0x8D: instruction{sta, addressMode{MODE_ABSOLUTE, REG_NONE}, "STA"},
	0x8E: instruction{stx, addressMode{MODE_ABSOLUTE, REG_NONE}, "STX"},
	0x90: instruction{bcc, addressMode{MODE_IMPLIED, REG_NONE}, "BCC"},
	0x91: instruction{sta, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "STA"},
	0x94: instruction{sty, addressMode{MODE_ZERO_PAGE, REG_X}, "STY"},
	0x95: instruction{sta, addressMode{MODE_ZERO_PAGE, REG_X}, "STA"},
	0x96: instruction{stx, addressMode{MODE_ZERO_PAGE, REG_Y}, "STX"},
	0x98: instruction{tya, addressMode{MODE_IMPLIED, REG_NONE}, "TYA"},
	0x99: instruction{sta, addressMode{MODE_ABSOLUTE, REG_Y}, "STA"},
	0x9A: instruction{txs, addressMode{MODE_IMPLIED, REG_NONE}, "TXS"},
	0x9D: instruction{sta, addressMode{MODE_ABSOLUTE, REG_X}, "STA"},
	0xA0: instruction{ldy, addressMode{MODE_IMMEDIATE, REG_NONE}, "LDY"},
	0xA1: instruction{lda, addressMode{MODE_INDEX_INDIRECT, REG_X}, "LDA"},
	0xA2: instruction{ldx, addressMode{MODE_IMMEDIATE, REG_NONE}, "LDX"},
	0xA4: instruction{ldy, addressMode{MODE_ZERO_PAGE, REG_NONE}, "LDY"},
	0xA5: instruction{lda, addressMode{MODE_ZERO_PAGE, REG_NONE}, "LDA"},
	0xA6: instruction{ldx, addressMode{MODE_ZERO_PAGE, REG_NONE}, "LDX"},
	0xA8: instruction{tay, addressMode{MODE_IMPLIED, REG_NONE}, "TAY"},
	0xA9: instruction{lda, addressMode{MODE_IMMEDIATE, REG_NONE}, "LDA"},
	0xAA: instruction{tax, addressMode{MODE_IMPLIED, REG_NONE}, "TAX"},
	0xAC: instruction{ldy, addressMode{MODE_ABSOLUTE, REG_NONE}, "LDY"},
	0xAD: instruction{lda, addressMode{MODE_ABSOLUTE, REG_NONE}, "LDA"},
	0xAE: instruction{ldx, addressMode{MODE_ABSOLUTE, REG_NONE}, "LDX"},
	0xB0: instruction{bcs, addressMode{MODE_IMPLIED, REG_NONE}, "BCS"},
	0xB1: instruction{lda, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "LDA"},
	0xB4: instruction{ldy, addressMode{MODE_ZERO_PAGE, REG_X}, "LDY"},
	0xB5: instruction{lda, addressMode{MODE_ZERO_PAGE, REG_X}, "LDA"},
	0xB6: instruction{ldx, addressMode{MODE_ZERO_PAGE, REG_Y}, "LDX"},
	0xB8: instruction{clv, addressMode{MODE_IMPLIED, REG_NONE}, "CLV"},
	0xB9: instruction{lda, addressMode{MODE_ABSOLUTE, REG_Y}, "LDA"},
	0xBA: instruction{tsx, addressMode{MODE_IMPLIED, REG_NONE}, "TSX"},
	0xBC: instruction{ldy, addressMode{MODE_ABSOLUTE, REG_X}, "LDY"},
	0xBD: instruction{lda, addressMode{MODE_ABSOLUTE, REG_X}, "LDA"},
	0xBE: instruction{ldx, addressMode{MODE_ABSOLUTE, REG_Y}, "LDX"},
	0xC0: instruction{cpy, addressMode{MODE_IMMEDIATE, REG_NONE}, "CPY"},
	0xC1: instruction{cmp, addressMode{MODE_INDEX_INDIRECT, REG_X}, "CMP"},
	0xC4: instruction{cpy, addressMode{MODE_ZERO_PAGE, REG_NONE}, "CPY"},
	0xC5: instruction{cmp, addressMode{MODE_ZERO_PAGE, REG_NONE}, "CMP"},
	0xC6: instruction{dec, addressMode{MODE_ZERO_PAGE, REG_NONE}, "DEC"},
	0xC8: instruction{iny, addressMode{MODE_IMPLIED, REG_NONE}, "INY"},
	0xC9: instruction{cmp, addressMode{MODE_IMMEDIATE, REG_NONE}, "CMP"},
	0xCA: instruction{dex, addressMode{MODE_IMPLIED, REG_NONE}, "DEX"},
	0xCC: instruction{cpy, addressMode{MODE_ABSOLUTE, REG_NONE}, "CPY"},
	0xCD: instruction{cmp, addressMode{MODE_ABSOLUTE, REG_NONE}, "CMP"},
	0xCE: instruction{dec, addressMode{MODE_ABSOLUTE, REG_NONE}, "DEC"},
	0xD0: instruction{bne, addressMode{MODE_IMPLIED, REG_NONE}, "BNE"},
	0xD1: instruction{cmp, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "CMP"},
	0xD5: instruction{cmp, addressMode{MODE_ZERO_PAGE, REG_X}, "CMP"},
	0xD6: instruction{dec, addressMode{MODE_ZERO_PAGE, REG_X}, "DEC"},
	0xD8: instruction{cld, addressMode{MODE_IMPLIED, REG_NONE}, "CLD"},
	0xD9: instruction{cmp, addressMode{MODE_ABSOLUTE, REG_Y}, "CMP"},
	0xDD: instruction{cmp, addressMode{MODE_ABSOLUTE, REG_X}, "CMP"},
	0xDE: instruction{dec, addressMode{MODE_ABSOLUTE, REG_X}, "DEC"},
	0xE0: instruction{cpx, addressMode{MODE_IMMEDIATE, REG_NONE}, "CPX"},
	0xE1: instruction{sbc, addressMode{MODE_INDEX_INDIRECT, REG_X}, "SBC"},
	0xE4: instruction{cpx, addressMode{MODE_ZERO_PAGE, REG_NONE}, "CPX"},
	0xE5: instruction{sbc, addressMode{MODE_ZERO_PAGE, REG_NONE}, "SBC"},
	0xE6: instruction{inc, addressMode{MODE_ZERO_PAGE, REG_NONE}, "INC"},
	0xE8: instruction{inx, addressMode{MODE_IMPLIED, REG_NONE}, "INX"},
	0xE9: instruction{sbc, addressMode{MODE_IMMEDIATE, REG_NONE}, "SBC"},
	0xEA: instruction{nop, addressMode{MODE_IMPLIED, REG_NONE}, "NOP"},
	0xEC: instruction{cpx, addressMode{MODE_ABSOLUTE, REG_NONE}, "CPX"},
	0xED: instruction{sbc, addressMode{MODE_ABSOLUTE, REG_NONE}, "SBC"},
	0xEE: instruction{inc, addressMode{MODE_ABSOLUTE, REG_NONE}, "INC"},
	0xF0: instruction{beq, addressMode{MODE_IMPLIED, REG_NONE}, "BEQ"},
	0xF1: instruction{sbc, addressMode{MODE_INDIRECT_INDEX, REG_Y}, "SBC"},
	0xF5: instruction{sbc, addressMode{MODE_ZERO_PAGE, REG_X}, "SBC"},
	0xF6: instruction{inc, addressMode{MODE_ZERO_PAGE, REG_X}, "INC"},
	0xF8: instruction{sed, addressMode{MODE_IMPLIED, REG_NONE}, "SED"},
	0xF9: instruction{sbc, addressMode{MODE_ABSOLUTE, REG_Y}, "SBC"},
	0xFD: instruction{sbc, addressMode{MODE_ABSOLUTE, REG_X}, "SBC"},
	0xFE: instruction{inc, addressMode{MODE_ABSOLUTE, REG_X}, "INC"},
}
