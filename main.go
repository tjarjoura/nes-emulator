package main

import (
	"github.com/tjarjoura/nes-emulator/cartridge"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s FILENAME\n", os.Args[0])
	}

	filename := os.Args[1]
	_, err := cartridge.CartridgeFromFile(filename)

	if err != nil {
		log.Fatal(err)
	}
}
