package hardware

import "log"

type Instruction struct {
	label string
	f     func(g *GBC) uint64
}

var (
	instructions = [256]Instruction{
		{
			"0x00; NOP",
			func(g *GBC) uint64 {
				return 4
			},
		},
		{
			"0x01; LD BC, u16",
			func(gbc *GBC) uint64 {
				ldR16u16(gbc, BC)
				return 12
			},
		},
		{
			"0x02; LD (BC), A",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(BC), A)
				return 8
			},
		},
		{
			"0x03; INC BC",
			func(gbc *GBC) uint64 {
				incR16(gbc, BC)
				return 8
			},
		},
		{
			"0x04; INC B",
			func(gbc *GBC) uint64 {
				incR8(gbc, B)
				return 4
			},
		},
		{
			"0x05; DEC B",
			func(gbc *GBC) uint64 {
				decR8(gbc, B)
				return 4
			},
		},
		{
			"0x06; LD B, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, B)
				return 8
			},
		},
		{
			"0x07; RLCA",
			func(gbc *GBC) uint64 {
				rlca(gbc)
				return 4
			},
		},
		{
			"0x08; LD (u16), SP",
			func(gbc *GBC) uint64 {
				l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
				gbc.PC += 2
				addr := (u << 8) | l
				uu, ll := byte(gbc.SP>>8), byte(gbc.SP)
				gbc.Write(addr, ll)
				gbc.Write(addr+1, uu)
				return 20
			},
		},
		{
			"0x09; ADD HL, BC",
			func(gbc *GBC) uint64 {
				addR16(gbc, HL, BC)
				return 8
			},
		},
		{
			"0x0A; LD A, (BC)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, A, gbc.Read(gbc.Reg16(BC)))
				return 8
			},
		},
		{
			"0x0B; DEC BC",
			func(gbc *GBC) uint64 {
				decR16(gbc, BC)
				return 8
			},
		},
		{
			"0x0C; INC C",
			func(gbc *GBC) uint64 {
				incR8(gbc, C)
				return 4
			},
		},
		{
			"0x0D; DEC C",
			func(gbc *GBC) uint64 {
				decR8(gbc, C)
				return 4
			},
		},
		{
			"0x0E; LD C, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, C)
				return 8
			},
		},
		{
			"0x0F; RRCA",
			func(gbc *GBC) uint64 {
				rrca(gbc)
				return 4
			},
		},
		{
			"0x10; STOP",
			func(gbc *GBC) uint64 {
				stop(gbc)
				return 4
			},
		},
		{
			"0x11; LD DE, u16",
			func(gbc *GBC) uint64 {
				ldR16u16(gbc, DE)
				return 12
			},
		},
		{
			"0x12; LD (DE), A",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(DE), A)
				return 8
			},
		},
		{
			"0x13; INC DE",
			func(gbc *GBC) uint64 {
				incR16(gbc, DE)
				return 8
			},
		},
		{
			"0x14; INC D",
			func(gbc *GBC) uint64 {
				incR8(gbc, D)
				return 4
			},
		},
		{
			"0x15; DEC D",
			func(gbc *GBC) uint64 {
				decR8(gbc, D)
				return 4
			},
		},
		{
			"0x16; LD D, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, D)
				return 8
			},
		},
		{
			"0x17; RLA",
			func(gbc *GBC) uint64 {
				rla(gbc)
				return 4
			},
		},
		{
			"0x18; JR i8",
			func(gbc *GBC) uint64 {
				jri8(gbc)
				return 12
			},
		},
		{
			"0x19; ADD HL, DE",
			func(gbc *GBC) uint64 {
				addR16(gbc, HL, DE)
				return 8
			},
		},
		{
			"0x1A; LD A, (DE)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, A, gbc.Read(gbc.Reg16(DE)))
				return 8
			},
		},
		{
			"0x1B; DEC DE",
			func(gbc *GBC) uint64 {
				decR16(gbc, DE)
				return 8
			},
		},
		{
			"0x1C; INC E",
			func(gbc *GBC) uint64 {
				incR8(gbc, E)
				return 4
			},
		},
		{
			"0x1D; DEC E",
			func(gbc *GBC) uint64 {
				decR8(gbc, E)
				return 4
			},
		},
		{
			"0x1E; LD E, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, E)
				return 8
			},
		},
		{
			"0x1F; RRA",
			func(gbc *GBC) uint64 {
				rra(gbc)
				return 4
			},
		},
		{
			"0x20; JR NZ, i8",
			func(gbc *GBC) uint64 {
				return jrncc(gbc, ZERO)
			},
		},
		{
			"0x21; LD HL, u16",
			func(gbc *GBC) uint64 {
				ldR16u16(gbc, HL)
				return 12
			},
		},
		{
			"0x22; LD (HL+), A",
			func(gbc *GBC) uint64 {
				ldHLINCA(gbc)
				return 8
			},
		},
		{
			"0x23; INC HL",
			func(gbc *GBC) uint64 {
				incR16(gbc, HL)
				return 8
			},
		},
		{
			"0x24; INC H",
			func(gbc *GBC) uint64 {
				incR8(gbc, H)
				return 4
			},
		},
		{
			"0x25; DEC H",
			func(gbc *GBC) uint64 {
				decR8(gbc, H)
				return 4
			},
		},
		{
			"0x26; LD H, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, H)
				return 8
			},
		},
		{
			"0x27; DAA",
			func(gbc *GBC) uint64 {
				daa(gbc)
				return 4
			},
		},
		{
			"0x28; JR Z, i8",
			func(gbc *GBC) uint64 {
				return jrcc(gbc, ZERO)
			},
		},
		{
			"0x29; ADD HL, HL",
			func(gbc *GBC) uint64 {
				addR16(gbc, HL, HL)
				return 8
			},
		},
		{
			"0x2A; LD A, (HL+)",
			func(gbc *GBC) uint64 {
				ldAHLINC(gbc)
				return 8
			},
		},
		{
			"0x2B; DEC HL",
			func(gbc *GBC) uint64 {
				decR16(gbc, HL)
				return 8
			},
		},
		{
			"0x2C; INC L",
			func(gbc *GBC) uint64 {
				incR8(gbc, L)
				return 4
			},
		},
		{
			"0x2D; DEC L",
			func(gbc *GBC) uint64 {
				decR8(gbc, L)
				return 4
			},
		},
		{
			"0x2E; LD L, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, L)
				return 8
			},
		},
		{
			"0x2F; CPL",
			func(gbc *GBC) uint64 {
				cpl(gbc)
				return 4
			},
		},
		{
			"0x30; JR NC, i8",
			func(gbc *GBC) uint64 {
				return jrncc(gbc, CARRY)
			},
		},
		{
			"0x31; LD SP, u16",
			func(gbc *GBC) uint64 {
				ldR16u16(gbc, SP)
				return 12
			},
		},
		{
			"0x32; LD (HL-), A",

			func(gbc *GBC) uint64 {
				ldHLDECA(gbc)
				return 8
			},
		},
		{
			"0x33; INC SP",
			func(gbc *GBC) uint64 {
				gbc.SP++
				return 8
			},
		},
		{
			"0x34; INC (HL)",
			func(gbc *GBC) uint64 {
				indirectIncHL(gbc)
				return 12
			},
		},
		{
			"0x35; DEC (HL)",
			func(gbc *GBC) uint64 {
				indirectDecHL(gbc)
				return 12
			},
		},
		{
			"0x36; LD (HL), u8",
			func(gbc *GBC) uint64 {
				indirectLDR16u8(gbc, HL)
				return 12
			},
		},
		{
			"0x37; SCF",
			func(gbc *GBC) uint64 {
				scf(gbc)
				return 4
			},
		},
		{
			"0x38; JR C, i8",
			func(gbc *GBC) uint64 {
				return jrcc(gbc, CARRY)
			},
		},
		{
			"0x39; ADD HL, SP",
			func(gbc *GBC) uint64 {
				addR16(gbc, HL, SP)
				return 8
			},
		},
		{
			"0x3A; LD A, (HL-)",
			func(gbc *GBC) uint64 {
				ldAHLDEC(gbc)
				return 8
			},
		},
		{
			"0x3B; DEC SP",
			func(gbc *GBC) uint64 {
				gbc.SP--
				return 8
			},
		},
		{
			"0x3C; INC A",
			func(gbc *GBC) uint64 {
				incR8(gbc, A)
				return 4
			},
		},
		{
			"0x3D; DEC A",
			func(gbc *GBC) uint64 {
				decR8(gbc, A)
				return 4
			},
		},
		{
			"0x3E; LD A, u8",
			func(gbc *GBC) uint64 {
				ldR8u8(gbc, A)
				return 8
			},
		},
		{
			"0x3F; CCF",
			func(gbc *GBC) uint64 {
				ccf(gbc)
				return 4
			},
		},
		{
			"0x40; LD B, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, B)
				return 4
			},
		},
		{
			"0x41; LD B, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, C)
				return 4
			},
		},
		{
			"0x42; LD B, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, D)
				return 4
			},
		},
		{
			"0x43; LD B, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, E)
				return 4
			},
		},
		{
			"0x44; LD B, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, H)
				return 4
			},
		},
		{
			"0x45; LD B, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, L)
				return 4
			},
		},
		{
			"0x46; LD B, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, B, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x47; LD B, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, B, A)
				return 4
			},
		},
		{
			"0x48; LD C, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, B)
				return 4
			},
		},
		{
			"0x49; LD C, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, C)
				return 4
			},
		},
		{
			"0x4A; LD C, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, D)
				return 4
			},
		},
		{
			"0x4B; LD C, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, E)
				return 4
			},
		},
		{
			"0x4C; LD C, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, H)
				return 4
			},
		},
		{
			"0x4D; LD C, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, L)
				return 4
			},
		},
		{
			"0x4E; LD C, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, C, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x4F; LD C, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, C, A)
				return 4
			},
		},
		{
			"0x50; LD D, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, B)
				return 4
			},
		},
		{
			"0x51; LD D, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, C)
				return 4
			},
		},
		{
			"0x52; LD D, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, D)
				return 4
			},
		},
		{
			"0x53; LD D, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, E)
				return 4
			},
		},
		{
			"0x54; LD D, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, H)
				return 4
			},
		},
		{
			"0x55; LD D, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, L)
				return 4
			},
		},
		{
			"0x56; LD D, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, D, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x57; LD D, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, D, A)
				return 4
			},
		},
		{
			"0x58; LD E, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, B)
				return 4
			},
		},
		{
			"0x59; LD E, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, C)
				return 4
			},
		},
		{
			"0x5A; LD E, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, D)
				return 4
			},
		},
		{
			"0x5B; LD E, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, F)
				return 4
			},
		},
		{
			"0x5C; LD E, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, H)
				return 4
			},
		},
		{
			"0x5D; LD E, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, L)
				return 4
			},
		},
		{
			"0x5E; LD E, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, E, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x5F; LD E, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, E, A)
				return 4
			},
		},
		{
			"0x60; LD H, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, B)
				return 4
			},
		},
		{
			"0x61; LD H, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, C)
				return 4
			},
		},
		{
			"0x62; LD H, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, D)
				return 4
			},
		},
		{
			"0x63; LD H, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, E)
				return 4
			},
		},
		{
			"0x64; LD H, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, H)
				return 4
			},
		},
		{
			"0x65; LD H, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, L)
				return 4
			},
		},
		{
			"0x66; LD H, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, H, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x67; LD H, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, H, A)
				return 4
			},
		},
		{
			"0x68; LD L, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, B)
				return 4
			},
		},
		{
			"0x69; LD L, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, C)
				return 4
			},
		},
		{
			"0x6A; LD L, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, D)
				return 4
			},
		},
		{
			"0x6B; LD L, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, E)
				return 4
			},
		},
		{
			"0x6C; LD L, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, H)
				return 4
			},
		},
		{
			"0x6D; LD L, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, L)
				return 4
			},
		},
		{
			"0x6E; LD L, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, L, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x6F; LD L, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, L, A)
				return 4
			},
		},
		{
			"0x70; LD (HL), B",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), B)
				return 8
			},
		},
		{
			"0x71; LD (HL), C",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), C)
				return 8
			},
		},
		{
			"0x72; LD (HL), D",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), D)
				return 8
			},
		},
		{
			"0x73; LD (HL), E",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), E)
				return 8
			},
		},
		{
			"0x74; LD (HL), H",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), H)
				return 8
			},
		},
		{
			"0x75; LD (HL), L",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), L)
				return 8
			},
		},
		{
			"0x76; HALT",
			func(gbc *GBC) uint64 {
				halt(gbc)
				return 4
			},
		},
		{
			"0x77; LD (HL), A",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, gbc.Reg16(HL), A)
				return 8
			},
		},
		{
			"0x78; LD A, B",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, B)
				return 4
			},
		},
		{
			"0x79; LD A, C",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, C)
				return 4
			},
		},
		{
			"0x7A; LD A, D",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, D)
				return 4
			},
		},
		{
			"0x7B; LD A, E",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, E)
				return 4
			},
		},
		{
			"0x7C; LD A, H",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, H)
				return 4
			},
		},
		{
			"0x7D; LD A, L",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, L)
				return 4
			},
		},
		{
			"0x7E; LD A, (HL)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, A, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x7F; LD A, A",
			func(gbc *GBC) uint64 {
				ldR8(gbc, A, A)
				return 4
			},
		},
		{
			"0x80; ADD A, B",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0x81; ADD A, C",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0x82; ADD A, D",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0x83; ADD A, E",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0x84; ADD A, H",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0x85; ADD A, L",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0x86; ADD A, (HL)",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x87; ADD A, A",
			func(gbc *GBC) uint64 {
				add(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0x88; ADC A, B",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0x89; ADC A, C",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0x8A; ADC A, D",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0x8B; ADC A, E",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0x8C; ADC A, H",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0x8D; ADC A, L",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0x8E; ADC A, (HL)",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x8F; ADC A, A",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0x90; SUB A, B",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0x91; SUB A, C",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0x92; SUB A, D",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0x93; SUB A, E",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0x94; SUB A, H",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0x95; SUB A, L",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0x96; SUB A, (HL)",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x97; SUB A, A",
			func(gbc *GBC) uint64 {
				sub(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0x98; SBC A, B",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0x99; SBC A, C",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0x9A; SBC A, D",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0x9B; SBC A, E",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0x9C; SBC A, H",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0x9D; SBC A, L",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0x9E; SBC A, (HL)",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0x9F; SBC A, A",
			func(gbc *GBC) uint64 {
				sbc(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0xA0; AND A, B",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0xA1; AND A, C",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0xA2; AND A, D",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0xA3; AND A, E",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0xA4; AND A, H",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0xA5; AND A, L",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0xA6; AND A, (HL)",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0xA7; AND A, A",
			func(gbc *GBC) uint64 {
				and(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0xA8; XOR A, B",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0xA9; XOR A, C",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0xAA; XOR A, D",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0xAB; XOR A, E",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0xAC; XOR A, H",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0xAD; XOR A, L",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0xAE; XOR A, (HL)",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0xAF; XOR A, A",
			func(gbc *GBC) uint64 {
				xor(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0xB0; OR A, B",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0xB1; OR A, C",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0xB2; OR A, D",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0xB3; OR A, E",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0xB4; OR A, H",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0xB5; OR A, L",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0xB6; OR A, (HL)",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0xB7; OR A, A",
			func(gbc *GBC) uint64 {
				or(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0xB8; CP A, B",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[B])
				return 4
			},
		},
		{
			"0xB9; CP A, C",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[C])
				return 4
			},
		},
		{
			"0xBA; CP A, D",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[D])
				return 4
			},
		},
		{
			"0xBB; CP A, E",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[E])
				return 4
			},
		},
		{
			"0xBC; CP A,H",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[H])
				return 4
			},
		},
		{
			"0xBD; CP A, L",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[L])
				return 4
			},
		},
		{
			"0xBE; CP A, (HL)",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.Read(gbc.Reg16(HL)))
				return 8
			},
		},
		{
			"0xBF; CP A, A",
			func(gbc *GBC) uint64 {
				cp(gbc, gbc.REG[A])
				return 4
			},
		},
		{
			"0xC0; RET NZ",
			func(gbc *GBC) uint64 {
				return retncc(gbc, ZERO)
			},
		},
		{
			"0xC1; POP BC",
			func(gbc *GBC) uint64 {
				popR16(gbc, C, B)
				return 12
			},
		},
		{
			"0xC2; JP NZ, u16",
			func(gbc *GBC) uint64 {
				return jpncc(gbc, ZERO)
			},
		},
		{
			"0xC3; JP u16",
			func(gbc *GBC) uint64 {
				return jp(gbc)
			},
		},
		{
			"0xC4; CALL NZ, u16",
			func(gbc *GBC) uint64 {
				return callncc(gbc, ZERO)
			},
		},
		{
			"0xC5; PUSH BC",
			func(gbc *GBC) uint64 {
				pushR16(gbc, B, C)
				return 16
			},
		},
		{
			"0xC6; ADD A, u8",
			func(gbc *GBC) uint64 {
				add(gbc, byte(gbc.Read(gbc.currPC+1)))
				gbc.PC++
				return 8
			},
		},
		{
			"0xC7; RST 00h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0000)
				return 16
			},
		},
		{
			"0xC8; RET Z",
			func(gbc *GBC) uint64 {
				return retcc(gbc, ZERO)
			},
		},
		{
			"0xC9; RET",
			func(gbc *GBC) uint64 {
				return ret(gbc)
			},
		},
		{
			"0xCA; JP Z, u16",
			func(gbc *GBC) uint64 {
				return jpcc(gbc, ZERO)
			},
		},
		{
			"0xCB; PREFIX CB",
			func(gbc *GBC) uint64 {
				cbop := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				inst := cb_instructions[cbop]
				log.Println(inst.label)
				return 4 + inst.f(gbc)
			},
		},
		{
			"0xCC; CALL Z, u16",
			func(gbc *GBC) uint64 {
				return callcc(gbc, ZERO)
			},
		},
		{
			"0xCD; CALL, u16",
			func(gbc *GBC) uint64 {
				call(gbc)
				return 24
			},
		},
		{
			"0xCE; ADC a, u8",
			func(gbc *GBC) uint64 {
				adc(gbc, gbc.Read(gbc.currPC+1))
				gbc.PC++
				return 8
			},
		},
		{
			"0xCF; RST 08h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0008)
				return 16
			},
		},
		{
			"0xD0; RET NC",
			func(gbc *GBC) uint64 {
				return retncc(gbc, CARRY)
			},
		},
		{
			"0xD1; POP DE",
			func(gbc *GBC) uint64 {
				popR16(gbc, E, D)
				return 12
			},
		},
		{
			"0xD2; JP NC, u16",
			func(gbc *GBC) uint64 {
				return jpncc(gbc, CARRY)
			},
		},
		{
			"0xD3; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xD4; CALL NC, u16",
			func(gbc *GBC) uint64 {
				return callncc(gbc, CARRY)
			},
		},
		{
			"0xD5; PUSH DE",
			func(gbc *GBC) uint64 {
				pushR16(gbc, D, E)
				return 16
			},
		},
		{
			"0xD6; SUB A, u8",
			func(gbc *GBC) uint64 {
				u8 := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				sub(gbc, u8)
				return 8
			},
		},
		{
			"0xD7; RST 10h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0010)
				return 16
			},
		},
		{
			"0xD8; RET C",
			func(gbc *GBC) uint64 {
				return retcc(gbc, CARRY)
			},
		},
		{
			"0xD9; RETI",
			func(gbc *GBC) uint64 {
				return reti(gbc)
			},
		},
		{
			"0xDA; JP C, u16",
			func(gbc *GBC) uint64 {
				return jpcc(gbc, CARRY)
			},
		},
		{
			"0xDB; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xDC; CALL C, u16",
			func(gbc *GBC) uint64 {
				return callcc(gbc, CARRY)
			},
		},
		{
			"0xDD; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xDE; SBC A, u8",
			func(gbc *GBC) uint64 {
				u8 := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				sbc(gbc, u8)
				return 8
			},
		},
		{
			"0xDF; RST 18h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0018)
				return 16
			},
		},
		{
			"0xE0; LD (FF00+u8), A",
			func(gbc *GBC) uint64 {
				l := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				ldnnR8(gbc, 0xFF00+uint16(l), A)
				return 12
			},
		},
		{
			"0xE1; POP HL",
			func(gbc *GBC) uint64 {
				popR16(gbc, L, H)
				return 12
			},
		},
		{
			"0xE2; LD (FF00+C), A",
			func(gbc *GBC) uint64 {
				ldnnR8(gbc, 0xFF00+uint16(gbc.REG[C]), A)
				return 8
			},
		},
		{
			"0xE3; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xE4; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xE5; PUSH HL",
			func(gbc *GBC) uint64 {
				pushR16(gbc, H, L)
				return 16
			},
		},
		{
			"0xE6; AND A, u8",
			func(gbc *GBC) uint64 {
				u8 := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				and(gbc, u8)
				return 8
			},
		},
		{
			"0xE7; RST 20h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0020)
				return 16
			},
		},
		{
			"0xE8; ADD SP, i8",
			func(gbc *GBC) uint64 {
				i8 := int8(gbc.Read(gbc.currPC + 1))
				gbc.PC++
				result := int32(gbc.SP) + int32(i8)
				carry := uint32(gbc.SP) ^ uint32(i8) ^ uint32(result)
				gbc.SP = uint16(result)
				gbc.setFlags(false, false, carry&(1<<4) != 0, carry&(1<<8) != 0)
				return 16
			},
		},
		{
			"0xE9; JP HL",
			func(gbc *GBC) uint64 {
				gbc.PC = gbc.Reg16(HL)
				return 4
			},
		},
		{
			"0xEA; LD (u16), A",
			func(gbc *GBC) uint64 {
				l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
				gbc.PC += 2
				ldnnR8(gbc, (u<<8)|l, A)
				return 0
			},
		},
		{
			"0xEB; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xEC; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xED; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xEE; XOR A, u8",
			func(gbc *GBC) uint64 {
				u8 := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				xor(gbc, u8)
				return 8
			},
		},
		{
			"0xEF; RST 28h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0028)
				return 16
			},
		},
		{
			"0xF0; LD A, (FF00+u8)",
			func(gbc *GBC) uint64 {
				l := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				log.Printf("%04X: %02X", 0xFF00+uint16(l), gbc.Read(0xFF00+uint16(l)))
				ldR8nn(gbc, A, gbc.Read(0xFF00+uint16(l)))
				return 12
			},
		},
		{
			"0xF1; POP AF",
			func(gbc *GBC) uint64 {
				popAF(gbc)
				return 12
			},
		},
		{
			"0xF2; LD A, (FF00+C)",
			func(gbc *GBC) uint64 {
				ldR8nn(gbc, A, gbc.Read(0xFF00+uint16(gbc.REG[C])))
				return 8
			},
		},
		{
			"0xF3; DI",
			func(gbc *GBC) uint64 {
				di(gbc)
				return 4
			},
		},
		{
			"0xF4; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xF5; PUSH AF",
			func(gbc *GBC) uint64 {
				pushAF(gbc)
				return 16
			},
		},
		{
			"0xF6; OR A, u8",
			func(gbc *GBC) uint64 {
				u8 := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				or(gbc, u8)
				return 8
			},
		},
		{
			"0xF7; RST 30h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0030)
				return 16
			},
		},
		{
			"0xF8; LD HL, SP+i8",
			func(gbc *GBC) uint64 {
				l, r := gbc.SP, int8(gbc.Read(gbc.currPC+1))
				gbc.PC++
				result := int32(l) + int32(r)
				carry := uint32(l) ^ uint32(r) ^ uint32(result)
				gbc.SetReg16(HL, uint16(result))
				gbc.setFlags(false, false, carry&(1<<4) != 0, carry&(1<<8) != 0)
				return 12
			},
		},
		{
			"0xF9; LD SP, HL",
			func(gbc *GBC) uint64 {
				gbc.SP = gbc.Reg16(HL)
				return 8
			},
		},
		{
			"0xFA; LD A, u16",
			func(gbc *GBC) uint64 {
				l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
				gbc.PC += 2
				ldR8nn(gbc, A, gbc.Read((u<<8)|l))
				return 16
			},
		},
		{
			"0xFB; EI",
			func(gbc *GBC) uint64 {
				ei(gbc)
				return 4
			},
		},
		{
			"0xFC; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xFD; INVALID",
			func(gbc *GBC) uint64 {
				return 0
			},
		},
		{
			"0xFE; CP A, u8",
			func(gbc *GBC) uint64 {
				u8 := gbc.Read(gbc.currPC + 1)
				gbc.PC++
				cp(gbc, u8)
				return 8
			},
		},
		{
			"0xFF; RST 38h",
			func(gbc *GBC) uint64 {
				rst(gbc, 0x0038)
				return 16
			},
		},
	}

	cb_instructions = [256]Instruction{{
		"CBx00; RLC B",
		func(gbc *GBC) uint64 {
			rlc(gbc, B)
			return 8
		},
	},
		{
			"CBx01; RLC C",
			func(gbc *GBC) uint64 {
				rlc(gbc, C)
				return 8
			},
		},
		{
			"CBx02; RLC D",
			func(gbc *GBC) uint64 {
				rlc(gbc, D)
				return 8
			},
		},
		{
			"CBx03; RLC E",
			func(gbc *GBC) uint64 {
				rlc(gbc, E)
				return 8
			},
		},
		{
			"CBx04; RLC H",
			func(gbc *GBC) uint64 {
				rlc(gbc, H)
				return 8
			},
		},
		{
			"CBx05; RLC L",
			func(gbc *GBC) uint64 {
				rlc(gbc, L)
				return 8
			},
		},
		{
			"CBx06; RLC (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				_, bit7 := _rl(gbc, value)
				new := (value << 1) | (bit7)
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit7 != 0)
				return 16
			},
		},
		{
			"CBx07; RLC A",
			func(gbc *GBC) uint64 {
				rlc(gbc, A)
				return 8
			},
		},
		{
			"CBx08; RRC B",
			func(gbc *GBC) uint64 {
				rrc(gbc, B)
				return 8
			},
		},
		{
			"CBx09; RRC C",
			func(gbc *GBC) uint64 {
				rrc(gbc, C)
				return 8
			},
		},
		{
			"CBx0A; RRC D",
			func(gbc *GBC) uint64 {
				rrc(gbc, D)
				return 8
			},
		},
		{
			"CBx0B; RRC E",
			func(gbc *GBC) uint64 {
				rrc(gbc, E)
				return 8
			},
		},
		{
			"CBx0C; RRC H",
			func(gbc *GBC) uint64 {
				rrc(gbc, H)
				return 8
			},
		},
		{
			"CBx0D; RRC L",
			func(gbc *GBC) uint64 {
				rrc(gbc, L)
				return 8
			},
		},
		{
			"CBx0E; RRC (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				_, bit0 := _rr(gbc, value)
				new := (value >> 1) | (bit0 << 7)
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit0 != 0)
				return 16
			},
		},
		{
			"CBx0F; RRC A",
			func(gbc *GBC) uint64 {
				rrc(gbc, A)
				return 8
			},
		},
		{
			"CBx10; RL B",
			func(gbc *GBC) uint64 {
				rl(gbc, B)
				return 8
			},
		},
		{
			"CBx11; RL C",
			func(gbc *GBC) uint64 {
				rl(gbc, C)
				return 8
			},
		},
		{
			"CBx12; RL D",
			func(gbc *GBC) uint64 {
				rl(gbc, D)
				return 8
			},
		},
		{
			"CBx13; RL E",
			func(gbc *GBC) uint64 {
				rl(gbc, E)
				return 8
			},
		},
		{
			"CBx14; RL H",
			func(gbc *GBC) uint64 {
				rl(gbc, H)
				return 8
			},
		},
		{
			"CBx15; RL L",
			func(gbc *GBC) uint64 {
				rl(gbc, L)
				return 8
			},
		},
		{
			"CBx16; RL (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				carry, bit7 := _rl(gbc, value)
				new := (value << 1) | (carry)
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit7 != 0)
				return 16
			},
		},
		{
			"CBx17; RL A",
			func(gbc *GBC) uint64 {
				rl(gbc, A)
				return 8
			},
		},
		{
			"CBx18; RR B",
			func(gbc *GBC) uint64 {
				rr(gbc, B)
				return 8
			},
		},
		{
			"CBx19; RR C",
			func(gbc *GBC) uint64 {
				rr(gbc, C)
				return 8
			},
		},
		{
			"CBx1A; RR D",
			func(gbc *GBC) uint64 {
				rr(gbc, D)
				return 8
			},
		},
		{
			"CBx1B; RR E",
			func(gbc *GBC) uint64 {
				rr(gbc, E)
				return 8
			},
		},
		{
			"CBx1C; RR H",
			func(gbc *GBC) uint64 {
				rr(gbc, H)
				return 8
			},
		},
		{
			"CBx1D; RR L",
			func(gbc *GBC) uint64 {
				rr(gbc, L)
				return 8
			},
		},
		{
			"CBx1E; RR (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				carry, bit0 := _rr(gbc, value)
				new := (value << 1) | (carry << 7)
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit0 != 0)
				return 16
			},
		},
		{
			"CBx1F; RR A",
			func(gbc *GBC) uint64 {
				rr(gbc, A)
				return 8
			},
		},
		{
			"CBx20; SLA B",
			func(gbc *GBC) uint64 {
				sla(gbc, B)
				return 8
			},
		},
		{
			"CBx21; SLA C",
			func(gbc *GBC) uint64 {
				sla(gbc, C)
				return 8
			},
		},
		{
			"CBx22; SLA D",
			func(gbc *GBC) uint64 {
				sla(gbc, D)
				return 8
			},
		},
		{
			"CBx23; SLA E",
			func(gbc *GBC) uint64 {
				sla(gbc, E)
				return 8
			},
		},
		{
			"CBx24; SLA H",
			func(gbc *GBC) uint64 {
				sla(gbc, H)
				return 8
			},
		},
		{
			"CBx25; SLA L",
			func(gbc *GBC) uint64 {
				sla(gbc, L)
				return 8
			},
		},
		{
			"CBx26; SLA (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				_, bit0 := _rl(gbc, value)
				new := (value << 1)
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit0 != 0)
				return 16
			},
		},
		{
			"CBx27; SLA A",
			func(gbc *GBC) uint64 {
				sla(gbc, A)
				return 8
			},
		},
		{
			"CBx28; SRA B",
			func(gbc *GBC) uint64 {
				sra(gbc, B)
				return 8
			},
		},
		{
			"CBx29; SRA C",
			func(gbc *GBC) uint64 {
				sra(gbc, C)
				return 8
			},
		},
		{
			"CBx2A; SRA D",
			func(gbc *GBC) uint64 {
				sra(gbc, D)
				return 8
			},
		},
		{
			"CBx2B; SRA E",
			func(gbc *GBC) uint64 {
				sra(gbc, E)
				return 8
			},
		},
		{
			"CBx2C; SRA H",
			func(gbc *GBC) uint64 {
				sra(gbc, H)
				return 8
			},
		},
		{
			"CBx2D; SRA L",
			func(gbc *GBC) uint64 {
				sra(gbc, L)
				return 8
			},
		},
		{
			"CBx2E; SRA (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				_, bit7 := _rr(gbc, value)
				new := (value >> 1)
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit7 != 0)
				return 16
			},
		},
		{
			"CBx2F; SRA A",
			func(gbc *GBC) uint64 {
				sra(gbc, A)
				return 8
			},
		},
		{
			"CBx30; SWAP B",
			func(gbc *GBC) uint64 {
				swap(gbc, B)
				return 8
			},
		},
		{
			"CBx31; SWAP C",
			func(gbc *GBC) uint64 {
				swap(gbc, C)
				return 8
			},
		},
		{
			"CBx32; SWAP D",
			func(gbc *GBC) uint64 {
				swap(gbc, D)
				return 8
			},
		},
		{
			"CBx33; SWAP E",
			func(gbc *GBC) uint64 {
				swap(gbc, E)
				return 8
			},
		},
		{
			"CBx34; SWAP H",
			func(gbc *GBC) uint64 {
				swap(gbc, H)
				return 8
			},
		},
		{
			"CBx35; SWAP L",
			func(gbc *GBC) uint64 {
				swap(gbc, L)
				return 8
			},
		},
		{
			"CBx36; SWAP (HL)",
			func(gbc *GBC) uint64 {
				_swap(gbc, HL)
				return 16
			},
		},
		{
			"CBx37; SWAP A",
			func(gbc *GBC) uint64 {
				swap(gbc, A)
				return 8
			},
		},
		{
			"CBx38; SRL B",
			func(gbc *GBC) uint64 {
				srl(gbc, B)
				return 8
			},
		},
		{
			"CBx39; SRL C",
			func(gbc *GBC) uint64 {
				srl(gbc, C)
				return 8
			},
		},
		{
			"CBx3A; SRL D",
			func(gbc *GBC) uint64 {
				srl(gbc, D)
				return 8
			},
		},
		{
			"CBx3B; SRL E",
			func(gbc *GBC) uint64 {
				srl(gbc, E)
				return 8
			},
		},
		{
			"CBx3C; SRL H",
			func(gbc *GBC) uint64 {
				srl(gbc, H)
				return 8
			},
		},
		{
			"CBx3D; SRL L",
			func(gbc *GBC) uint64 {
				srl(gbc, L)
				return 8
			},
		},
		{
			"CBx3E; SRL (HL)",
			func(gbc *GBC) uint64 {
				value := gbc.Read(gbc.Reg16(HL))
				_, bit7 := _rr(gbc, value)
				new := (value >> 1) &^ 0x80
				gbc.Write(gbc.Reg16(HL), new)
				gbc.setFlags(new == 0, false, false, bit7 != 0)
				return 16
			},
		},
		{
			"CBx3F; SRL A",
			func(gbc *GBC) uint64 {
				srl(gbc, A)
				return 8
			},
		},
		{
			"CBx40; BIT 0, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, B)
				return 8
			},
		},
		{
			"CBx41; BIT 0, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, C)
				return 8
			},
		},
		{
			"CBx42; BIT 0, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, D)
				return 8
			},
		},
		{
			"CBx43; BIT 0, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, E)
				return 8
			},
		},
		{
			"CBx44; BIT 0, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, H)
				return 8
			},
		},
		{
			"CBx45; BIT 0, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, L)
				return 8
			},
		},
		{
			"CBx46; BIT 0, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 0, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx47; BIT 0, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 0, A)
				return 8
			},
		},
		{
			"CBx48; BIT 1, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, B)
				return 8
			},
		},
		{
			"CBx49; BIT 1, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, C)
				return 8
			},
		},
		{
			"CBx4A; BIT 1, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, D)
				return 8
			},
		},
		{
			"CBx4B; BIT 1, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, E)
				return 8
			},
		},
		{
			"CBx4C; BIT 1, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, H)
				return 8
			},
		},
		{
			"CBx4D; BIT 1, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, L)
				return 8
			},
		},
		{
			"CBx4E; BIT 1, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 1, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx4F; BIT 1, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 1, A)
				return 8
			},
		},
		{
			"CBx50; BIT 2, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, B)
				return 8
			},
		},
		{
			"CBx51; BIT 2, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, C)
				return 8
			},
		},
		{
			"CBx52; BIT 2, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, D)
				return 8
			},
		},
		{
			"CBx53; BIT 2, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, E)
				return 8
			},
		},
		{
			"CBx54; BIT 2, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, H)
				return 8
			},
		},
		{
			"CBx55; BIT 2, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, L)
				return 8
			},
		},
		{
			"CBx56; BIT 2, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 2, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx57; BIT 2, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 2, A)
				return 8
			},
		},
		{
			"CBx58; BIT 3, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, B)
				return 8
			},
		},
		{
			"CBx59; BIT 3, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, C)
				return 8
			},
		},
		{
			"CBx5A; BIT 3, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, D)
				return 8
			},
		},
		{
			"CBx5B; BIT 3, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, E)
				return 8
			},
		},
		{
			"CBx5C; BIT 3, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, H)
				return 8
			},
		},
		{
			"CBx5D; BIT 3, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, L)
				return 8
			},
		},
		{
			"CBx5E; BIT 3, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 3, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx5F; BIT 3, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 3, A)
				return 8
			},
		},
		{
			"CBx60; BIT 4, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, B)
				return 8
			},
		},
		{
			"CBx61; BIT 4, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, C)
				return 8
			},
		},
		{
			"CBx62; BIT 4, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, D)
				return 8
			},
		},
		{
			"CBx63; BIT 4, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, E)
				return 8
			},
		},
		{
			"CBx64; BIT 4, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, H)
				return 8
			},
		},
		{
			"CBx65; BIT 4, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, L)
				return 8
			},
		},
		{
			"CBx66; BIT 4, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 4, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx67; BIT 4, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 4, A)
				return 8
			},
		},
		{
			"CBx68; BIT 5, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, B)
				return 8
			},
		},
		{
			"CBx69; BIT 5, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, C)
				return 8
			},
		},
		{
			"CBx6A; BIT 5, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, D)
				return 8
			},
		},
		{
			"CBx6B; BIT 5, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, E)
				return 8
			},
		},
		{
			"CBx6C; BIT 5, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, H)
				return 8
			},
		},
		{
			"CBx6D; BIT 5, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, L)
				return 8
			},
		},
		{
			"CBx6E; BIT 5, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 5, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx6F; BIT 5, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 5, A)
				return 8
			},
		},
		{
			"CBx70; BIT 6, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, B)
				return 8
			},
		},
		{
			"CBx71; BIT 6, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, C)
				return 8
			},
		},
		{
			"CBx72; BIT 6, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, D)
				return 8
			},
		},
		{
			"CBx73; BIT 6, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, E)
				return 8
			},
		},
		{
			"CBx74; BIT 6, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, H)
				return 8
			},
		},
		{
			"CBx75; BIT 6, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, L)
				return 8
			},
		},
		{
			"CBx76; BIT 6, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 6, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx77; BIT 6, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 6, A)
				return 8
			},
		},
		{
			"CBx78; BIT 7, B",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, B)
				return 8
			},
		},
		{
			"CBx79; BIT 7, C",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, C)
				return 8
			},
		},
		{
			"CBx7A; BIT 7, D",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, D)
				return 8
			},
		},
		{
			"CBx7B; BIT 7, E",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, E)
				return 8
			},
		},
		{
			"CBx7C; BIT 7, H",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, H)
				return 8
			},
		},
		{
			"CBx7D; BIT 7, L",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, L)
				return 8
			},
		},
		{
			"CBx7E; BIT 7, (HL)",
			func(gbc *GBC) uint64 {
				_bit(gbc, 7, gbc.Read(gbc.Reg16(HL)))
				return 16
			},
		},
		{
			"CBx7F; BIT 7, A",
			func(gbc *GBC) uint64 {
				bit(gbc, 7, A)
				return 8
			},
		},
		{
			"CBx80; RES 0, B",
			func(gbc *GBC) uint64 {
				res(gbc, 0, B)
				return 8
			},
		},
		{
			"CBx81; RES 0, C",
			func(gbc *GBC) uint64 {
				res(gbc, 0, C)
				return 8
			},
		},
		{
			"CBx82; RES 0, D",
			func(gbc *GBC) uint64 {
				res(gbc, 0, D)
				return 8
			},
		},
		{
			"CBx83; RES 0, E",
			func(gbc *GBC) uint64 {
				res(gbc, 0, E)
				return 8
			},
		},
		{
			"CBx84; RES 0, H",
			func(gbc *GBC) uint64 {
				res(gbc, 0, H)
				return 8
			},
		},
		{
			"CBx85; RES 0, L",
			func(gbc *GBC) uint64 {
				res(gbc, 0, L)
				return 8
			},
		},
		{
			"CBx86; RES 0, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 0, HL)
				return 16
			},
		},
		{
			"CBx87; RES 0, A",
			func(gbc *GBC) uint64 {
				res(gbc, 0, A)
				return 8
			},
		},
		{
			"CBx88; RES 1, B",
			func(gbc *GBC) uint64 {
				res(gbc, 1, B)
				return 8
			},
		},
		{
			"CBx89; RES 1, C",
			func(gbc *GBC) uint64 {
				res(gbc, 1, C)
				return 8
			},
		},
		{
			"CBx8A; RES 1, D",
			func(gbc *GBC) uint64 {
				res(gbc, 1, D)
				return 8
			},
		},
		{
			"CBx8B; RES 1, E",
			func(gbc *GBC) uint64 {
				res(gbc, 1, E)
				return 8
			},
		},
		{
			"CBx8C; RES 1, H",
			func(gbc *GBC) uint64 {
				res(gbc, 1, H)
				return 8
			},
		},
		{
			"CBx8D; RES 1, L",
			func(gbc *GBC) uint64 {
				res(gbc, 1, L)
				return 8
			},
		},
		{
			"CBx8E; RES 1, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 1, HL)
				return 16
			},
		},
		{
			"CBx8F; RES 1, A",
			func(gbc *GBC) uint64 {
				res(gbc, 1, A)
				return 8
			},
		},
		{
			"CBx90; RES 2, B",
			func(gbc *GBC) uint64 {
				res(gbc, 2, B)
				return 8
			},
		},
		{
			"CBx91; RES 2, C",
			func(gbc *GBC) uint64 {
				res(gbc, 2, C)
				return 8
			},
		},
		{
			"CBx92; RES 2, D",
			func(gbc *GBC) uint64 {
				res(gbc, 2, D)
				return 8
			},
		},
		{
			"CBx93; RES 2, E",
			func(gbc *GBC) uint64 {
				res(gbc, 2, E)
				return 8
			},
		},
		{
			"CBx94; RES 2, H",
			func(gbc *GBC) uint64 {
				res(gbc, 2, H)
				return 8
			},
		},
		{
			"CBx95; RES 2, L",
			func(gbc *GBC) uint64 {
				res(gbc, 2, L)
				return 8
			},
		},
		{
			"CBx96; RES 2, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 2, HL)
				return 16
			},
		},
		{
			"CBx97; RES 2, A",
			func(gbc *GBC) uint64 {
				res(gbc, 2, A)
				return 8
			},
		},
		{
			"CBx98; RES 3, B",
			func(gbc *GBC) uint64 {
				res(gbc, 3, B)
				return 8
			},
		},
		{
			"CBx99; RES 3, C",
			func(gbc *GBC) uint64 {
				res(gbc, 3, C)
				return 8
			},
		},
		{
			"CBx9A; RES 3, D",
			func(gbc *GBC) uint64 {
				res(gbc, 3, D)
				return 8
			},
		},
		{
			"CBx9B; RES 3, E",
			func(gbc *GBC) uint64 {
				res(gbc, 3, E)
				return 8
			},
		},
		{
			"CBx9C; RES 3, H",
			func(gbc *GBC) uint64 {
				res(gbc, 3, H)
				return 8
			},
		},
		{
			"CBx9D; RES 3, L",
			func(gbc *GBC) uint64 {
				res(gbc, 3, L)
				return 8
			},
		},
		{
			"CBx9E; RES 3, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 3, HL)
				return 16
			},
		},
		{
			"CBx9F; RES 3, A",
			func(gbc *GBC) uint64 {
				res(gbc, 3, A)
				return 8
			},
		},
		{
			"CBxA0; RES 4, B",
			func(gbc *GBC) uint64 {
				res(gbc, 4, B)
				return 8
			},
		},
		{
			"CBxA1; RES 4, C",
			func(gbc *GBC) uint64 {
				res(gbc, 4, C)
				return 8
			},
		},
		{
			"CBxA2; RES 4, D",
			func(gbc *GBC) uint64 {
				res(gbc, 4, D)
				return 8
			},
		},
		{
			"CBxA3; RES 4, E",
			func(gbc *GBC) uint64 {
				res(gbc, 4, E)
				return 8
			},
		},
		{
			"CBxA4; RES 4, H",
			func(gbc *GBC) uint64 {
				res(gbc, 4, H)
				return 8
			},
		},
		{
			"CBxA5; RES 4, L",
			func(gbc *GBC) uint64 {
				res(gbc, 4, L)
				return 8
			},
		},
		{
			"CBxA6; RES 4, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 4, HL)
				return 16
			},
		},
		{
			"CBxA7; RES 4, A",
			func(gbc *GBC) uint64 {
				res(gbc, 4, A)
				return 8
			},
		},
		{
			"CBxA8; RES 5, B",
			func(gbc *GBC) uint64 {
				res(gbc, 5, B)
				return 8
			},
		},
		{
			"CBxA9; RES 5, C",
			func(gbc *GBC) uint64 {
				res(gbc, 5, C)
				return 8
			},
		},
		{
			"CBxAA; RES 5, D",
			func(gbc *GBC) uint64 {
				res(gbc, 5, D)
				return 8
			},
		},
		{
			"CBxAB; RES 5, E",
			func(gbc *GBC) uint64 {
				res(gbc, 5, E)
				return 8
			},
		},
		{
			"CBxAC; RES 5, H",
			func(gbc *GBC) uint64 {
				res(gbc, 5, H)
				return 8
			},
		},
		{
			"CBxAD; RES 5, L",
			func(gbc *GBC) uint64 {
				res(gbc, 5, L)
				return 8
			},
		},
		{
			"CBxAE; RES 5, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 5, HL)
				return 16
			},
		},
		{
			"CBxAF; RES 5, A",
			func(gbc *GBC) uint64 {
				res(gbc, 5, A)
				return 8
			},
		},
		{
			"CBxB0; RES 6, B",
			func(gbc *GBC) uint64 {
				res(gbc, 6, B)
				return 8
			},
		},
		{
			"CBxB1; RES 6, C",
			func(gbc *GBC) uint64 {
				res(gbc, 6, C)
				return 8
			},
		},
		{
			"CBxB2; RES 6, D",
			func(gbc *GBC) uint64 {
				res(gbc, 6, D)
				return 8
			},
		},
		{
			"CBxB3; RES 6, E",
			func(gbc *GBC) uint64 {
				res(gbc, 6, E)
				return 8
			},
		},
		{
			"CBxB4; RES 6, H",
			func(gbc *GBC) uint64 {
				res(gbc, 6, H)
				return 8
			},
		},
		{
			"CBxB5; RES 6, L",
			func(gbc *GBC) uint64 {
				res(gbc, 6, L)
				return 8
			},
		},
		{
			"CBxB6; RES 6, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 6, HL)
				return 16
			},
		},
		{
			"CBxB7; RES 6, A",
			func(gbc *GBC) uint64 {
				res(gbc, 6, A)
				return 8
			},
		},
		{
			"CBxB8; RES 7, B",
			func(gbc *GBC) uint64 {
				res(gbc, 7, B)
				return 8
			},
		},
		{
			"CBxB9; RES 7, C",
			func(gbc *GBC) uint64 {
				res(gbc, 7, C)
				return 8
			},
		},
		{
			"CBxBA; RES 7, D",
			func(gbc *GBC) uint64 {
				res(gbc, 7, D)
				return 8
			},
		},
		{
			"CBxBB; RES 7, E",
			func(gbc *GBC) uint64 {
				res(gbc, 7, E)
				return 8
			},
		},
		{
			"CBxBC; RES 7, H",
			func(gbc *GBC) uint64 {
				res(gbc, 7, H)
				return 8
			},
		},
		{
			"CBxBD; RES 7, L",
			func(gbc *GBC) uint64 {
				res(gbc, 7, L)
				return 8
			},
		},
		{
			"CBxBE; RES 7, (HL)",
			func(gbc *GBC) uint64 {
				_res(gbc, 7, HL)
				return 16
			},
		},
		{
			"CBxBF; RES 7, A",
			func(gbc *GBC) uint64 {
				res(gbc, 7, A)
				return 8
			},
		},
		{
			"CBxC0; SET 0, B",
			func(gbc *GBC) uint64 {
				set(gbc, 0, B)
				return 8
			},
		},
		{
			"CBxC1; SET 0, C",
			func(gbc *GBC) uint64 {
				set(gbc, 0, C)
				return 8
			},
		},
		{
			"CBxC2; SET 0, D",
			func(gbc *GBC) uint64 {
				set(gbc, 0, D)
				return 8
			},
		},
		{
			"CBxC3; SET 0, E",
			func(gbc *GBC) uint64 {
				set(gbc, 0, E)
				return 8
			},
		},
		{
			"CBxC4; SET 0, H",
			func(gbc *GBC) uint64 {
				set(gbc, 0, H)
				return 8
			},
		},
		{
			"CBxC5; SET 0, L",
			func(gbc *GBC) uint64 {
				set(gbc, 0, L)
				return 8
			},
		},
		{
			"CBxC6; SET 0, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 0, HL)
				return 16
			},
		},
		{
			"CBxC7; SET 0, A",
			func(gbc *GBC) uint64 {
				set(gbc, 0, A)
				return 8
			},
		},
		{
			"CBxC8; SET 1, B",
			func(gbc *GBC) uint64 {
				set(gbc, 1, B)
				return 8
			},
		},
		{
			"CBxC9; SET 1, C",
			func(gbc *GBC) uint64 {
				set(gbc, 1, C)
				return 8
			},
		},
		{
			"CBxCA; SET 1, D",
			func(gbc *GBC) uint64 {
				set(gbc, 1, D)
				return 8
			},
		},
		{
			"CBxCB; SET 1, E",
			func(gbc *GBC) uint64 {
				set(gbc, 1, E)
				return 8
			},
		},
		{
			"CBxCC; SET 1, H",
			func(gbc *GBC) uint64 {
				set(gbc, 1, H)
				return 8
			},
		},
		{
			"CBxCD; SET 1, L",
			func(gbc *GBC) uint64 {
				set(gbc, 1, L)
				return 8
			},
		},
		{
			"CBxCE; SET 1, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 1, HL)
				return 16
			},
		},
		{
			"CBxCF; SET 1, A",
			func(gbc *GBC) uint64 {
				set(gbc, 1, A)
				return 8
			},
		},
		{
			"CBxD0; SET 2, B",
			func(gbc *GBC) uint64 {
				set(gbc, 2, B)
				return 8
			},
		},
		{
			"CBxD1; SET 2, C",
			func(gbc *GBC) uint64 {
				set(gbc, 2, C)
				return 8
			},
		},
		{
			"CBxD2; SET 2, D",
			func(gbc *GBC) uint64 {
				set(gbc, 2, D)
				return 8
			},
		},
		{
			"CBxD3; SET 2, E",
			func(gbc *GBC) uint64 {
				set(gbc, 2, E)
				return 8
			},
		},
		{
			"CBxD4; SET 2, H",
			func(gbc *GBC) uint64 {
				set(gbc, 2, H)
				return 8
			},
		},
		{
			"CBxD5; SET 2, L",
			func(gbc *GBC) uint64 {
				set(gbc, 2, L)
				return 8
			},
		},
		{
			"CBxD6; SET 2, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 2, HL)
				return 16
			},
		},
		{
			"CBxD7; SET 2, A",
			func(gbc *GBC) uint64 {
				set(gbc, 2, A)
				return 8
			},
		},
		{
			"CBxD8; SET 3, B",
			func(gbc *GBC) uint64 {
				set(gbc, 3, B)
				return 8
			},
		},
		{
			"CBxD9; SET 3, C",
			func(gbc *GBC) uint64 {
				set(gbc, 3, C)
				return 8
			},
		},
		{
			"CBxDA; SET 3, D",
			func(gbc *GBC) uint64 {
				set(gbc, 3, D)
				return 8
			},
		},
		{
			"CBxDB; SET 3, E",
			func(gbc *GBC) uint64 {
				set(gbc, 3, E)
				return 8
			},
		},
		{
			"CBxDC; SET 3, H",
			func(gbc *GBC) uint64 {
				set(gbc, 3, H)
				return 8
			},
		},
		{
			"CBxDD; SET 3, L",
			func(gbc *GBC) uint64 {
				set(gbc, 3, L)
				return 8
			},
		},
		{
			"CBxDE; SET 3, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 3, HL)
				return 16
			},
		},
		{
			"CBxDF; SET 3, A",
			func(gbc *GBC) uint64 {
				set(gbc, 3, A)
				return 8
			},
		},
		{
			"CBxE0; SET 4, B",
			func(gbc *GBC) uint64 {
				set(gbc, 4, B)
				return 8
			},
		},
		{
			"CBxE1; SET 4, C",
			func(gbc *GBC) uint64 {
				set(gbc, 4, C)
				return 8
			},
		},
		{
			"CBxE2; SET 4, D",
			func(gbc *GBC) uint64 {
				set(gbc, 4, D)
				return 8
			},
		},
		{
			"CBxE3; SET 4, E",
			func(gbc *GBC) uint64 {
				set(gbc, 4, E)
				return 8
			},
		},
		{
			"CBxE4; SET 4, H",
			func(gbc *GBC) uint64 {
				set(gbc, 4, H)
				return 8
			},
		},
		{
			"CBxE5; SET 4, L",
			func(gbc *GBC) uint64 {
				set(gbc, 4, L)
				return 8
			},
		},
		{
			"CBxE6; SET 4, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 4, HL)
				return 16
			},
		},
		{
			"CBxE7; SET 4, A",
			func(gbc *GBC) uint64 {
				set(gbc, 4, A)
				return 8
			},
		},
		{
			"CBxE8; SET 5, B",
			func(gbc *GBC) uint64 {
				set(gbc, 5, B)
				return 8
			},
		},
		{
			"CBxE9; SET 5, C",
			func(gbc *GBC) uint64 {
				set(gbc, 5, C)
				return 8
			},
		},
		{
			"CBxEA; SET 5, D",
			func(gbc *GBC) uint64 {
				set(gbc, 5, D)
				return 8
			},
		},
		{
			"CBxEB; SET 5, E",
			func(gbc *GBC) uint64 {
				set(gbc, 5, E)
				return 8
			},
		},
		{
			"CBxEC; SET 5, H",
			func(gbc *GBC) uint64 {
				set(gbc, 5, H)
				return 8
			},
		},
		{
			"CBxED; SET 5, L",
			func(gbc *GBC) uint64 {
				set(gbc, 5, L)
				return 8
			},
		},
		{
			"CBxEE; SET 5, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 5, HL)
				return 16
			},
		},
		{
			"CBxEF; SET 5, A",
			func(gbc *GBC) uint64 {
				set(gbc, 5, A)
				return 8
			},
		},
		{
			"CBxF0; SET 6, B",
			func(gbc *GBC) uint64 {
				set(gbc, 6, B)
				return 8
			},
		},
		{
			"CBxF1; SET 6, C",
			func(gbc *GBC) uint64 {
				set(gbc, 6, C)
				return 8
			},
		},
		{
			"CBxF2; SET 6, D",
			func(gbc *GBC) uint64 {
				set(gbc, 6, D)
				return 8
			},
		},
		{
			"CBxF3; SET 6, E",
			func(gbc *GBC) uint64 {
				set(gbc, 6, E)
				return 8
			},
		},
		{
			"CBxF4; SET 6, H",
			func(gbc *GBC) uint64 {
				set(gbc, 6, H)
				return 8
			},
		},
		{
			"CBxF5; SET 6, L",
			func(gbc *GBC) uint64 {
				set(gbc, 6, L)
				return 8
			},
		},
		{
			"CBxF6; SET 6, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 6, HL)
				return 16
			},
		},
		{
			"CBxF7; SET 6, A",
			func(gbc *GBC) uint64 {
				set(gbc, 6, A)
				return 8
			},
		},
		{
			"CBxF8; SET 7, B",
			func(gbc *GBC) uint64 {
				set(gbc, 7, B)
				return 8
			},
		},
		{
			"CBxF9; SET 7, C",
			func(gbc *GBC) uint64 {
				set(gbc, 7, C)
				return 8
			},
		},
		{
			"CBxFA; SET 7, D",
			func(gbc *GBC) uint64 {
				set(gbc, 7, D)
				return 8
			},
		},
		{
			"CBxFB; SET 7, E",
			func(gbc *GBC) uint64 {
				set(gbc, 7, E)
				return 8
			},
		},
		{
			"CBxFC; SET 7, H",
			func(gbc *GBC) uint64 {
				set(gbc, 7, H)
				return 8
			},
		},
		{
			"CBxFD; SET 7, L",
			func(gbc *GBC) uint64 {
				set(gbc, 7, L)
				return 8
			},
		},
		{
			"CBxFE; SET 7, (HL)",
			func(gbc *GBC) uint64 {
				_set(gbc, 7, HL)
				return 16
			},
		},
		{
			"CBxFF; SET 7, A",
			func(gbc *GBC) uint64 {
				set(gbc, 7, A)
				return 8
			},
		},
	}
)
