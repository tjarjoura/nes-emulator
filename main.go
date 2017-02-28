package main

import (
	"fmt"
	"github.com/tjarjoura/nes-emulator/cartridge"
	"github.com/tjarjoura/nes-emulator/cpu"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s FILENAME\n", os.Args[0])
	}

	filename := os.Args[1]
	cartridge, err := cartridge.CartridgeFromFile(filename)

	if err != nil {
		log.Fatalf("CartridgeFromFile(): %s\n", err)
	}

	fmt.Printf("len(cartridge.PrgRom) = %d\n", len(cartridge.PrgRom))
	cpu := new(cpu.Cpu)
	cpu.LoadProgram(cartridge.PrgRom)
	fmt.Printf("%s\n", cpu.String())

	err = cpu.Run(true)
	if err != nil {
		log.Fatalf("cpu.Run(): %s\n", err)
	}
}
