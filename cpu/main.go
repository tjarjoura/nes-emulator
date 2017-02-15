package cpu

type Cpu struct {
	a, x, y, pc, sp, p byte
	ram                [2000]byte
}
