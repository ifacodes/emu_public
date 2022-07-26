package hardware

type MODEL byte

type MBC byte

const (
	DMG MODEL = iota
	CGB
	PGB
)

const (
	MBC0 MBC = iota
	MBC1
	MBC2
	MBC3
	MBC4
	MBC5
	MBC6
	MBC7
	HuC1
)

type Cartridge struct {
	title  string
	mode   MODEL
	header [50]byte
	mbc    MBC
}
