package hardware

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type Processor struct {
	cycles uint64
}

func NewProcessor() *Processor               { return &Processor{} }
func (p *Processor) Cycle() uint64           { return p.cycles }
func (p *Processor) AddCycles(amount uint64) { p.cycles += amount }

type GBC struct {
	cycles uint64
	currOP byte
	currPC uint16
	Register
	MMU
	debug_compare bufio.Scanner
	setPendingIME bool
	halted        bool
	stoped        bool
}

func NewGBC(file []byte, compare_file io.Reader) *GBC {
	return &GBC{
		Register: Register{
			REG: [8]byte{0x00, 0x13, 0x00, 0xD8, 0x01, 0x4D, 0xB0, 0x01},
			SP:  0xFFFE,
			PC:  0x0100,
			IME: false,
		},
		MMU:           NewMMU(file),
		debug_compare: *bufio.NewScanner(compare_file),
	}
}

func (gbc *GBC) DebugStep() {
	currentMEM := fmt.Sprintf("A: %02X F: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X SP: %04X PC: 00:%04X (%02X %02X %02X %02X)", gbc.REG[A], gbc.REG[F], gbc.REG[B], gbc.REG[C], gbc.REG[D], gbc.REG[E], gbc.REG[H], gbc.REG[L], gbc.SP, gbc.PC, gbc.Read(gbc.PC), gbc.Read(gbc.PC+1), gbc.Read(gbc.PC+2), gbc.Read(gbc.PC+3))
	log.Println(currentMEM)
	test := gbc.debug_compare.Scan()
	if !test && gbc.debug_compare.Err() == nil {
		return
	} else {
		fmt.Println(gbc.debug_compare.Text())
		logAtStep := gbc.debug_compare.Text()
		if logAtStep != currentMEM && test {
			panic("Current Memory Does not Match!")
		}
	}
}

func (gbc *GBC) Step() {
	// timer stuff

	gbc.DebugStep()

	if gbc.halted {
	} else {
		if gbc.stoped {
			gbc.cycles += 4
		} else {
			//gbc.HandleInterrupts()
			gbc.currPC = gbc.PC
			gbc.currOP = gbc.Read(gbc.currPC)
			gbc.PC++

			inst := instructions[gbc.currOP]
			log.Println(inst.label)
			inst.f(gbc)
		}
	}
	outputCheck := gbc.Read(0xFF02)
	output := gbc.Read(0xFF01)
	if outputCheck == 0x81 {
		fmt.Printf("blargg output: %q\n", output)
		gbc.Write(0xFF02, 0x00)
	}
}

func (gbc *GBC) HandleInterrupts() {
	IME := gbc.IME
	if gbc.setPendingIME {
		gbc.Register.IME = !gbc.Register.IME
		gbc.setPendingIME = false
	}

	IE := gbc.Read(IE)
	IF := gbc.Read(0xFF0F)

	if IME && IE&IF&0x1F != 0 {
		// do the interrupt stuff
		gbc.halted = false
		gbc.Write(gbc.SP-1, byte(gbc.PC>>8))
		gbc.Write(gbc.SP-2, byte(gbc.PC&0x00FF))
		gbc.SP -= 2
		gbc.cycles += 12

	}
}

/*
	if toggle next cycle
		toggle ime
	run instruction
	read current instruction and move to the next address
	if ime
		handle interrupt
*/
