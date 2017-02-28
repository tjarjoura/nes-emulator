package cpu

import (
	"fmt"
)

type Cpu struct {
	a, x, y, p                   byte
	carryFl, zeroFl, interruptFl bool
	brkFl, overflowFl, signFl    bool
	pc, sp                       uint16
	ram                          [2048]byte
	prgRom                       []byte
}

func (cpu *Cpu) String() string {
	return fmt.Sprintf("%d", len(cpu.prgRom))
}

func (cpu *Cpu) LoadProgram(prgRom []byte) {
	cpu.prgRom = make([]byte, len(prgRom))
	copy(cpu.prgRom, prgRom)
	cpu.pc = 0x4020
	fmt.Printf("len(cpu.prgRom) = %d\n", len(cpu.prgRom))
}

func (cpu *Cpu) byteAt(address uint16) byte {
	if address < 0x2000 {
		return cpu.ram[address%0x800]
	} else if address >= 0x4020 {
		return cpu.prgRom[address-0x4020]
	} else {
		fmt.Printf("Unmapped memory address %x. Reading 0x00", address)
		return 0x00
	}
}

func (cpu *Cpu) wordAt(address uint16) uint16 {
	return (uint16(cpu.prgRom[address+1]) << 8) | uint16(cpu.prgRom[address])
}

func (cpu *Cpu) writeByte(address uint16, data uint8) error {
	cpu.prgRom[address] = data
	return nil
}

func (cpu *Cpu) getArgument(mode addressMode) (uint16, uint16) {
	switch mode.mode {
	case MODE_IMMEDIATE:
		arg := cpu.byteAt(cpu.pc + 1)
		return uint16(arg), 2

	case MODE_ZERO_PAGE:
		address := uint16(cpu.byteAt(cpu.pc + 1))

		if mode.reg == REG_X {
			address += uint16(cpu.x)
		} else if mode.reg == REG_Y {
			address += uint16(cpu.y)
		}

		return uint16(cpu.byteAt(address)), 2

	case MODE_ABSOLUTE:
		address_hi := uint16(cpu.byteAt(cpu.pc + 2))
		address_lo := uint16(cpu.byteAt(cpu.pc + 1))
		address := address_hi | address_lo

		if mode.reg == REG_X {
			address += uint16(cpu.x)
		} else if mode.reg == REG_Y {
			address += uint16(cpu.y)
		}

		return address, 3

	case MODE_RELATIVE:
		offset := int16(cpu.byteAt(cpu.pc + 1))

		if offset < 0 {
			return cpu.pc - uint16(-1*offset), 2
		} else {
			return cpu.pc + uint16(offset), 2
		}

	case MODE_INDEX_INDIRECT:
		address := uint16(cpu.byteAt(cpu.pc+1)) + uint16(cpu.x)
		return cpu.wordAt(address), 2

	case MODE_INDIRECT_INDEX:
		address := cpu.wordAt(uint16(cpu.byteAt(cpu.pc + 1)))
		return cpu.wordAt(address + uint16(cpu.y)), 2
	}

	return 0, 1 // No argument, use 0 as dummy value
}

func (cpu *Cpu) disassemble(instruction instruction, incr uint16) {
	if incr == 0 {
		fmt.Printf(instruction.neumonic)
	} else if incr == 1 {
		fmt.Printf(instruction.neumonic, cpu.byteAt(cpu.pc+1))
	} else {
		fmt.Printf(instruction.neumonic, cpu.wordAt(cpu.pc+1))
	}

}

func (cpu *Cpu) Run(disassemble bool) error {
	for {
		opcode := cpu.byteAt(cpu.pc)
		instruction := instructions[opcode]

		if instruction.handler == nil {
			return fmt.Errorf("Unrecognized opcode: %x\n", opcode)
		}

		arg, incr := cpu.getArgument(instruction.addressMode)

		if disassemble {
			cpu.disassemble(instruction, incr)
		}

		err := instruction.handler(cpu, arg)
		if err != nil {
			return err
		}

		cpu.pc += incr
	}

	return nil
}
