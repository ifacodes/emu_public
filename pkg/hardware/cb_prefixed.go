package hardware

import "log"

func swap(gbc *GBC, r REGISTER8) {
	old := gbc.REG[r]
	gbc.REG[r] = (old << 4) | (old >> 4)
	gbc.setFlags(gbc.REG[r] == 0, false, false, false)
}

func rlc(gbc *GBC, r REGISTER8) {
	_, bit7 := _rl(gbc, gbc.REG[r])
	gbc.REG[r] = (gbc.REG[r] << 1) | bit7
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit7 != 0)
}

func rl(gbc *GBC, r REGISTER8) {
	carry, bit7 := _rl(gbc, gbc.REG[r])
	gbc.REG[r] = (gbc.REG[r] << 1) | carry
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit7 != 0)
}

func rrc(gbc *GBC, r REGISTER8) {
	_, bit0 := _rr(gbc, gbc.REG[r])
	gbc.REG[r] = (gbc.REG[r] >> 1) | (bit0 << 7)
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit0 != 0)
}

func rr(gbc *GBC, r REGISTER8) {
	carry, bit0 := _rl(gbc, gbc.REG[r])
	log.Printf("current bit0: %0b", bit0)
	gbc.REG[r] = (gbc.REG[r] >> 1) | (carry << 7)
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit0 != 0)
}

func sla(gbc *GBC, r REGISTER8) {
	_, bit7 := _rl(gbc, gbc.REG[r])
	gbc.REG[r] = (gbc.REG[r] << 1)
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit7 != 0)
}

func sra(gbc *GBC, r REGISTER8) {
	_, bit0 := _rr(gbc, gbc.REG[r])
	_, bit7 := _rl(gbc, gbc.REG[r])
	gbc.REG[r] = (gbc.REG[r] >> 1) | (bit7 << 7)
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit0 != 0)
}

func srl(gbc *GBC, r REGISTER8) {
	bit0 := gbc.REG[r] & 1
	gbc.REG[r] = (gbc.REG[r] >> 1)
	gbc.setFlags(gbc.REG[r] == 0, false, false, bit0 != 0)
}

func bit(gbc *GBC, bit int, r REGISTER8) {
	_bit(gbc, bit, gbc.REG[r])
}

func set(gbc *GBC, bit int, r REGISTER8) {
	gbc.REG[r] |= (1 << bit)
}

func res(gbc *GBC, bit int, r REGISTER8) {
	gbc.REG[r] &^= (1 << bit)
}

func _swap(gbc *GBC, r REGISTER16) {
	old := gbc.Read(gbc.Reg16(r))
	new := (old << 4) | (old >> 4)
	gbc.Write(gbc.Reg16(HL), new)
	gbc.setFlags(new == 0, false, false, false)
}

func _rl(gbc *GBC, mm byte) (uint8, byte) {
	var c uint8
	if gbc.getFlag(CARRY) {
		c = 1
	} else {
		c = 0
	}
	bit7 := (mm & 1)
	return c, bit7
}

func _rr(gbc *GBC, mm byte) (uint8, byte) {
	var c uint8
	if gbc.getFlag(CARRY) {
		c = 1
	} else {
		c = 0
	}
	bit0 := (mm >> 7) & 1
	return c, bit0
}

func _bit(gbc *GBC, bit int, mm byte) {
	gbc.setZNH(mm>>bit == 0, false, true)
}

func _set(gbc *GBC, bit int, r REGISTER16) {
	old := gbc.Read(gbc.Reg16(HL))
	gbc.Write(gbc.Reg16(HL), old|(1<<bit))
}

func _res(gbc *GBC, bit int, r REGISTER16) {
	old := gbc.Read(gbc.Reg16(HL))
	gbc.Write(gbc.Reg16(HL), old&^(1<<bit))
}
