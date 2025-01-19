package gamestats2

type GRNG struct {
	initialSeed uint32
	seed uint32
}

func NewGRNG(checksum uint32) GRNG {
	ret := GRNG{}

	ret.initialSeed = checksum | (checksum << 16)
	ret.seed = ret.initialSeed

	return ret
}

func (g *GRNG) Next() uint32 {
	g.seed = (g.seed * 0x45 + 0x1111) & 0x7FFFFFFF
	return g.seed
}

func (g *GRNG) ResetGRNG() uint32 {
	g.seed = g.initialSeed
	return g.seed
}

func (g *GRNG) GetVal() uint32 {
	return g.seed
}