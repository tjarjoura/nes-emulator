package cpu

import (
	"fmt"
	"github.com/tjarjoura/nes-emulator/types"
)

type Cpu struct {
	a, x, y, p                   byte
	carryFl, zeroFl, interruptFl bool
	brkFl, overflowFl, signFl    bool
	pc, sp                       uint16
	ram                          [2048]byte
	cartridge, ppu, apu          types.MappedHardware
}

func (cpu *Cpu) LoadProgram(cartridge types.MappedHardware) {
	cpu.cartridge = cartridge
}

func (cpu *Cpu) byteAt(address uint16) byte {
	if address < 0x2000 {
		return cpu.ram[address%0x800]
	} else if address >= 0x4020 {
		data, err := cpu.cartridge.ReadByte(address - 0x4020)

		if err != nil {
			fmt.Printf("Error reading byte at address %x: %s. Returning 0x00\n", address, err)
			return 0x00
		}

		return data
	} else {
		fmt.Printf("Unmapped memory address %x. Returning 0x00\n", address)
		return 0x00
	}
}

func (cpu *Cpu) wordAt(address uint16) uint16 {
	dataHi := cpu.byteAt(address + 1)
	dataLo := cpu.byteAt(address)

	return uint16(dataHi)<<8 | uint16(dataLo)
}

func (cpu *Cpu) writeByte(address uint16, data uint8) error {
	if address < 0x2000 {
		cpu.ram[address%0x800] = data
		return nil
	} else if address >= 0x4020 {
		return cpu.cartridge.WriteByte(address-0x4020, data)
	} else {
		return fmt.Errorf("Cpu.writeByte(): Unmapped memory address %x.\n", address)
	}
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
