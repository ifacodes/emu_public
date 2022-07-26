package hardware

import "log"

// 8-bit Loadss

func ldR8nn(gbc *GBC, r REGISTER8, nn byte) {
	gbc.REG[r] = nn
}

func ldR8(gbc *GBC, r REGISTER8, rr REGISTER8) {
	ldR8nn(gbc, r, gbc.REG[rr])
}

func ldR16nn(gbc *GBC, r16 REGISTER16, nn byte) {
	gbc.Write(gbc.Reg16(r16), nn)
}

func ldnnR8(gbc *GBC, nn uint16, r REGISTER8) {
	gbc.Write(nn, gbc.REG[r])
}

func xF2(gbc *GBC) {
	gbc.REG[A] = gbc.Read(0xFF00 + uint16(gbc.REG[C]))
}

func xE2(gbc *GBC) {
	gbc.Write(0xFF00+uint16(gbc.REG[C]), gbc.REG[A])
}

func ldAHLDEC(gbc *GBC) {
	addr := gbc.Reg16(HL)
	gbc.REG[A] = gbc.Read(addr)
	gbc.SetReg16(HL, addr-1)
}

func ldAHLINC(gbc *GBC) {
	addr := gbc.Reg16(HL)
	gbc.REG[A] = gbc.Read(addr)
	gbc.SetReg16(HL, addr+1)
}

func xE0(gbc *GBC) {
	gbc.Write(0xFF00+uint16(gbc.Read(gbc.currPC+1)), gbc.REG[A])
	gbc.PC++
}

func xF0(gbc *GBC) {
	gbc.REG[A] = gbc.Read(0xFF00 + uint16(gbc.Read(gbc.currPC+1)))
	gbc.PC++
}

// 16-bit Loads

func ldR16u16(gbc *GBC, r16 REGISTER16) {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	gbc.SetReg16(r16, (u<<8)|l)

}

func xF9(gbc *GBC) {
	gbc.SP = gbc.Reg16(HL)
}

func xF8(gbc *GBC) {
	l, r := gbc.SP, uint16(gbc.Read(gbc.currPC+1))
	gbc.PC++
	result := l + r
	carry := l ^ r ^ result
	gbc.SetReg16(HL, result)
	gbc.setFlags(false, false, carry>>4&7 != 0, carry>>7&7 != 0)
}

func x08(gbc *GBC) {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	gbc.Write((u<<8)|l, gbc.Read(gbc.SP))
}

func pushR16(gbc *GBC, r1, r2 REGISTER8) {
	gbc.Write(gbc.SP-1, gbc.REG[r1])
	gbc.Write(gbc.SP-2, gbc.REG[r2])
	gbc.SP -= 2
}

func popR16(gbc *GBC, r2, r1 REGISTER8) {
	gbc.REG[r2] = gbc.Read(gbc.SP)
	gbc.REG[r1] = gbc.Read(gbc.SP + 1)
	gbc.SP += 2
}

func pushAF(gbc *GBC) {
	gbc.Write(gbc.SP-1, gbc.REG[A])
	gbc.Write(gbc.SP-2, gbc.REG[F]&0xF0)
	gbc.SP -= 2
}

func popAF(gbc *GBC) {
	gbc.REG[F] = (gbc.Read(gbc.SP) & 0xF0)
	gbc.REG[A] = gbc.Read(gbc.SP + 1)
	gbc.SP += 2
}

// 8-bit Arithmetic

func add(gbc *GBC, value byte) {
	l, r := gbc.REG[A], value
	result := uint16(l) + uint16(r)
	carry := uint16(l) ^ uint16(r) ^ result
	gbc.REG[A] = byte(result)
	gbc.setFlags(byte(result) == 0, false, (carry&(1<<4)) != 0, (carry&(1<<8)) != 0)
}

func adc(gbc *GBC, value byte) {
	var c uint8
	if gbc.getFlag(CARRY) {
		c = 1
	} else {
		c = 0
	}
	l, r := gbc.REG[A], value
	result := l + r + c
	halfcarry, carry := (l&0xF)+(r&0xF)+c, uint16(l)+uint16(r)+uint16(c)
	gbc.REG[A] = byte(result)
	gbc.setFlags(result == 0, false, halfcarry&(1<<4) != 0, carry&(1<<8) != 0)
}

func sub(gbc *GBC, value byte) {
	l, r := gbc.REG[A], value
	result := uint16(l) - uint16(r)
	carry := uint16(l) ^ uint16(r) ^ result
	gbc.REG[A] = byte(result)
	gbc.setFlags(result == 0, true, carry>>4&7 != 0, l < byte(result))
}

func sbc(gbc *GBC, value byte) {
	var c uint16
	if gbc.getFlag(CARRY) {
		c = 1
	} else {
		c = 0
	}
	l, r := gbc.REG[A], value
	result := uint16(l) - (uint16(r) + c)
	carry := uint16(l) ^ (uint16(r) + c) ^ result
	gbc.REG[A] = byte(result)
	gbc.setFlags(result == 0, true, carry>>4&7 != 0, l < byte(result))
}

func and(gbc *GBC, value byte) {
	gbc.REG[A] &= value
	gbc.setFlags(gbc.REG[A] == 0, false, true, false)
}

func or(gbc *GBC, value byte) {
	gbc.REG[A] |= value
	gbc.setFlags(gbc.REG[A] == 0, false, false, false)
}

func xor(gbc *GBC, value byte) {
	gbc.REG[A] ^= value
	gbc.setFlags(gbc.REG[A] == 0, false, false, false)
}

func cp(gbc *GBC, value byte) {
	l, r := gbc.REG[A], value
	result := l - r
	carry := l ^ r ^ result
	gbc.setFlags(result == 0, true, carry&(1<<4) != 0, l < byte(result))
}

func incR8(gbc *GBC, r REGISTER8) {
	old := gbc.REG[r]
	gbc.REG[r]++
	carry := old ^ gbc.REG[r] ^ 1
	gbc.setZNH(gbc.REG[r] == 0, false, carry>>4&7 != 0)
}

func decR8(gbc *GBC, r REGISTER8) {
	old := gbc.REG[r]
	gbc.REG[r]--
	carry := old ^ gbc.REG[r] ^ 1
	gbc.setZNH(gbc.REG[r] == 0, true, carry>>4&7 != 0)
}

// 16-bit Arithmetic

func addHLR16(gbc *GBC, rr REGISTER16) {
	l, r := gbc.Reg16(HL), gbc.Reg16(rr)
	result := uint32(l) + uint32(r)
	carry := uint32(l) ^ uint32(r) ^ result
	gbc.SetReg16(HL, uint16(result))
	gbc.setFlags(byte(result) == 0, false, carry>>4&7 != 0, carry>>7&7 != 0)
}

func incR16(gbc *GBC, r REGISTER16) {
	gbc.SetReg16(r, gbc.Reg16(r)+1)
}

func decR16(gbc *GBC, r REGISTER16) {
	gbc.SetReg16(r, gbc.Reg16(r)-1)
}

// Jumps

func jri8(gbc *GBC) uint64 {
	d := int8(gbc.Read(gbc.currPC + 1))
	gbc.PC++
	gbc.PC = uint16(int32(gbc.PC) + int32(d))
	return 12
}

func jrcc(gbc *GBC, f FLAG) uint64 {
	d := int8(gbc.Read(gbc.currPC + 1))
	gbc.PC++
	if gbc.getFlag(f) {
		gbc.PC = uint16(int32(gbc.PC) + int32(d))
		return 12
	} else {
		return 8
	}
}

func jrncc(gbc *GBC, f FLAG) uint64 {
	d := int8(gbc.Read(gbc.currPC + 1))
	gbc.PC++
	if !gbc.getFlag(f) {
		gbc.PC = uint16(int32(gbc.PC) + int32(d))
		return 12
	} else {
		return 8
	}
}

func jp(gbc *GBC) uint64 {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	gbc.PC = (u << 8) | l
	return 16
}

func jpcc(gbc *GBC, f FLAG) uint64 {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	if gbc.getFlag(f) {
		gbc.PC = (u << 8) | l
		return 16
	}
	return 12
}

func jpncc(gbc *GBC, f FLAG) uint64 {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	if !gbc.getFlag(f) {
		gbc.PC = (u << 8) | l
		return 16
	}
	return 12
}

// Calls

func call(gbc *GBC) {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	uu, ll := byte(gbc.PC>>8), byte(gbc.PC&0x00FF)
	gbc.Write(gbc.SP-1, uu)
	gbc.Write(gbc.SP-2, ll)
	gbc.SP -= 2
	gbc.PC = (u << 8) | l
}

func callcc(gbc *GBC, f FLAG) uint64 {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	if gbc.getFlag(f) {

		uu, ll := byte(gbc.PC>>8), byte(gbc.PC&0x00FF)
		gbc.Write(gbc.SP-1, uu)
		gbc.Write(gbc.SP-2, ll)
		gbc.SP -= 2
		gbc.PC = (u << 8) | l
		return 24
	}
	return 12
}

func callncc(gbc *GBC, f FLAG) uint64 {
	l, u := uint16(gbc.Read(gbc.currPC+1)), uint16(gbc.Read(gbc.currPC+2))
	gbc.PC += 2
	if !gbc.getFlag(f) {

		uu, ll := byte(gbc.PC>>8), byte(gbc.PC&0x00FF)
		gbc.Write(gbc.SP-1, uu)
		gbc.Write(gbc.SP-2, ll)
		gbc.SP -= 2
		gbc.PC = (u << 8) | l
		return 24
	}
	return 12
}

// Restart

func rst(gbc *GBC, addr uint16) {
	u, l := byte(gbc.PC>>8), byte(gbc.PC)
	gbc.Write(gbc.SP-1, u)
	gbc.Write(gbc.SP-2, l)
	gbc.SP -= 2
	gbc.PC = addr
}

// Returns

func ret(gbc *GBC) uint64 {
	l, u := uint16(gbc.Read(gbc.SP)), uint16(gbc.Read(gbc.SP+1))
	gbc.SP += 2
	gbc.PC = (u << 8) | l
	return 16
}

func reti(gbc *GBC) uint64 {
	l, u := uint16(gbc.Read(gbc.SP)), uint16(gbc.Read(gbc.SP+1))
	gbc.SP += 2
	gbc.PC = (u << 8) | l
	gbc.setPendingIME = true
	return 16
}

func retcc(gbc *GBC, f FLAG) uint64 {
	if gbc.getFlag(f) {
		l, u := uint16(gbc.Read(gbc.SP)), uint16(gbc.Read(gbc.SP+1))
		gbc.SP += 2
		gbc.PC = (u << 8) | l
		return 20
	}
	return 8
}

func retncc(gbc *GBC, f FLAG) uint64 {
	if !gbc.getFlag(f) {
		l, u := uint16(gbc.Read(gbc.SP)), uint16(gbc.Read(gbc.SP+1))
		gbc.SP += 2
		gbc.PC = (u << 8) | l
		return 20
	}
	return 8
}

// Rotates and Shifts

func rlca(gbc *GBC) {
	a := gbc.REG[A] >> 7
	gbc.REG[A] = gbc.REG[A] << 1
	gbc.REG[A] |= a
	gbc.setFlags(gbc.REG[A] == 0, false, false, a != 0)
}

func rla(gbc *GBC) {
	a := gbc.REG[A] >> 7
	gbc.REG[A] = gbc.REG[A] << 1
	if gbc.getFlag(CARRY) {
		gbc.REG[A] |= 1
	}
	gbc.setFlags(gbc.REG[A] == 0, false, false, a != 0)
}

func rrca(gbc *GBC) {
	a := gbc.REG[A] & 1
	gbc.REG[A] = gbc.REG[A] >> 1
	gbc.REG[A] |= (a << 7)
	gbc.setFlags(gbc.REG[A] == 0, false, false, a != 0)
}

func rra(gbc *GBC) {
	a := gbc.REG[A] & 1
	gbc.REG[A] >>= 1
	if gbc.getFlag(CARRY) {
		gbc.REG[A] |= 0x80
	}
	log.Printf("%0b\n", gbc.REG[A])
	gbc.setFlags(false, false, false, a != 0)
}

// Miscellaneous

func daa(gbc *GBC) {
	if !gbc.getFlag(NEGATIVE) {
		if gbc.getFlag(CARRY) || gbc.REG[A] > 0x99 {
			gbc.REG[A] += 0x60
			gbc.setFlag(CARRY, true)
		}
		if gbc.getFlag(HALF_CARRY) || (gbc.REG[A]&0x0F) > 0x09 {
			gbc.REG[A] += 0x06
		}
	} else {
		if gbc.getFlag(CARRY) {
			gbc.REG[A] -= 0x60
		}
		if gbc.getFlag(HALF_CARRY) {
			gbc.REG[A] -= 0x06
		}
	}
	gbc.setFlag(ZERO, gbc.REG[A] == 0)
	gbc.setFlag(HALF_CARRY, false)
}

func cpl(gbc *GBC) {
	gbc.REG[A] = ^gbc.REG[A]
	gbc.setNH(true, true)
}

func ccf(gbc *GBC) {
	gbc.setNHC(false, false, !gbc.getFlag(CARRY))
}

func scf(gbc *GBC) {
	gbc.setNHC(false, false, true)
}

func nop(gbc *GBC) uint64 {
	return 4
}

func halt(gbc *GBC) uint64 {
	return 4
}

func stop(gbc *GBC) uint64 {
	return 4
}

func di(gbc *GBC) uint64 {
	gbc.IME = false
	return 4
}

func ei(gbc *GBC) uint64 {
	gbc.setPendingIME = true
	return 4
}

// LD R, u8 8-bit
func ldR8u8(gbc *GBC, r REGISTER8) {
	gbc.Register.REG[r] = gbc.Read(uint16(gbc.currPC + 1))
	gbc.PC++
}

func indirectLDR16u8(gbc *GBC, r REGISTER16) {
	gbc.Write(gbc.Reg16(r), gbc.Read(gbc.currPC+1))
	gbc.PC++
}

func ldHLINCA(gbc *GBC) {
	addr := gbc.Reg16(HL)
	gbc.Write(addr, gbc.REG[A])
	gbc.SetReg16(HL, addr+1)
}

func ldHLDECA(gbc *GBC) {
	addr := gbc.Reg16(HL)
	gbc.Write(addr, gbc.REG[A])
	gbc.SetReg16(HL, addr-1)
}

func addR16(gbc *GBC, r1 REGISTER16, r2 REGISTER16) {
	l, r := gbc.Reg16(r1), gbc.Reg16(r2)
	result := uint32(l) + uint32(r)
	carry := uint32(l) ^ uint32(r) ^ result
	gbc.SetReg16(r1, uint16(result))
	gbc.setNHC(false, carry&(1<<12) != 0, carry&(1<<16) != 0)
}

func indirectIncHL(gbc *GBC) {
	value := gbc.Read(gbc.Reg16(HL))
	result := value + 1
	carry := value ^ result ^ 1
	gbc.Write(gbc.Reg16(HL), result)
	gbc.setZNH(result == 0, false, carry>>4&7 != 0)
}

func indirectDecHL(gbc *GBC) {
	value := gbc.Read(gbc.Reg16(HL))
	result := value - 1
	carry := value ^ result ^ 1
	gbc.Write(gbc.Reg16(HL), result)
	gbc.setZNH(result == 0, true, carry>>4&7 != 0)
}
