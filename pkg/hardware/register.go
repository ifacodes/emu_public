package hardware

type REGISTER8 int
type REGISTER16 int

const (
	B REGISTER8 = iota
	C
	D
	E
	H
	L
	F
	A
	AF REGISTER16 = iota
	BC
	DE
	HL
	SP
	PC
)

const (
	REG_B = 0
	REG_C = 1
	REG_D = 2
	REG_E = 3
	REG_H = 4
	REG_L = 5
	REG_F = 6
	REG_A = 7
)

type FLAG int

const (
	ZERO       FLAG = 0x80
	NEGATIVE   FLAG = 0x40
	HALF_CARRY FLAG = 0x20
	CARRY      FLAG = 0x10
)

type Register struct {
	REG [8]byte // order: [b,c,d,e,h,l,f,a]
	SP  uint16  // Stack Pointer
	PC  uint16  // Program Counter
	IME bool
}

func (r *Register) Reg16(rg REGISTER16) uint16 {
	var i1, i2 REGISTER8
	switch rg {
	case AF:
		i1, i2 = REG_A, REG_F
	case BC:
		i1, i2 = REG_B, REG_C
	case DE:
		i1, i2 = REG_D, REG_E
	case HL:
		i1, i2 = REG_H, REG_L
	case SP:
		return r.SP
	default:
		panic("Invalid 16bit Register!")
	}
	//log.Printf("%x", (uint16(r.REG[i1])<<8)|uint16(r.REG[i2]))
	return (uint16(r.REG[i1]) << 8) | uint16(r.REG[i2])
}

func (r *Register) SetReg16(rg REGISTER16, value uint16) {
	var i1, i2 REGISTER8
	switch rg {
	case AF:
		i1, i2 = REG_A, REG_F
	case BC:
		i1, i2 = REG_B, REG_C
	case DE:
		i1, i2 = REG_D, REG_E
	case HL:
		i1, i2 = REG_H, REG_L
	case SP:
		r.SP = value
		return
	default:
		panic("Invalid 16bit Register!")
	}
	r.REG[i1], r.REG[i2] = byte(value>>8), byte(value)
}

func (r *Register) getFlag(flag FLAG) bool {
	//log.Print(r.REG[F]&byte(flag), r.REG[F]&byte(flag) != 0)
	return r.REG[F]&byte(flag) != 0
}

func (r *Register) setFlag(f FLAG, set bool) {
	if set {
		r.REG[F] |= byte(f)
		return
	}
	r.REG[F] &= ^byte(f)
}

func (r *Register) setZNH(z, n, h bool) {
	r.setFlag(ZERO, z)
	r.setNH(n, h)
}

func (r *Register) setNH(n, h bool) {
	r.setFlag(NEGATIVE, n)
	r.setFlag(HALF_CARRY, h)
}

func (r *Register) setNHC(n, h, c bool) {
	r.setNH(n, h)
	r.setFlag(CARRY, c)
}

func (r *Register) setFlags(z, n, h, c bool) {
	r.setZNH(z, n, h)
	r.setFlag(CARRY, c)
}
