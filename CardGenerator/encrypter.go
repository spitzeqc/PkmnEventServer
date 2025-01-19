/*
 * Reimplementation of ProjectPokemon's PCDWC4Converter written in Go
 * https://github.com/projectpokemon/PCDWC4Converter/tree/9c1b1f1281f2117230616b5b8e80bbebbdb34169
 * AGPL-3.0 License
 */

package wondercard

import (
	"encoding/binary"
)

// LCRng
type lcrng struct {
	seed uint32
}
func newLCRng(s uint32) lcrng {
	return lcrng{seed: s}
}

func (lc *lcrng) next() uint32 {
	lc.seed = (lc.seed * 0x41C64E6D) + 0x00006073
	return lc.seed
}
func (lc *lcrng) nextH() uint32 {
	return lc.next() >> 0x10
}
func (lc *lcrng) prev() uint32 {
	lc.seed = (lc.seed * 0xEEB9EB65) + 0xA3561A1
	return lc.seed
}
func (lc *lcrng) prevH() uint32 {
	return lc.prev() >> 0x10
}

// do not modify this
var blockPositions = [...]byte{
	0, 1, 2, 3,
	0, 1, 3, 2,
	0, 2, 1, 3,
	0, 3, 1, 2,
	0, 2, 3, 1,
	0, 3, 2, 1,
	1, 0, 2, 3,
	1, 0, 3, 2,
	2, 0, 1, 3,
	3, 0, 1, 2,
	2, 0, 3, 1,
	3, 0, 2, 1,
	1, 2, 0, 3,
	1, 3, 0, 2,
	2, 1, 0, 3,
	3, 1, 0, 2,
	2, 3, 0, 1,
	3, 2, 0, 1,
	1, 2, 3, 0,
	1, 3, 2, 0,
	2, 1, 3, 0,
	3, 1, 2, 0,
	2, 3, 1, 0,
	3, 2, 1, 0,
}
var blockPositionInvert = [...]byte{
	0, 1, 2, 4, 3, 5, 6, 7, 12, 18, 13, 19, 8, 10, 14, 20, 16, 22, 9, 11, 15, 21, 17, 23,
}

func xorCrypt(data []byte) []byte{
	pid := binary.LittleEndian.Uint32(data[0:4])
	initialSeed := binary.LittleEndian.Uint16(data[6:9]) // Go excludes last element 

	rng := newLCRng( uint32(initialSeed) )

	for i:=8; i<236; i+=2 {
		if i == 136 {
			rng = newLCRng( pid )
		}

		dataBlock := data[i:]

		value := binary.LittleEndian.Uint16(dataBlock)

		tmp := uint16( uint32(value) ^ rng.nextH() )
		data[i+1] = byte(tmp >> 8)
		data[i] = byte(tmp)
	}

	return data
}

func shuffle(data []byte, shiftValue uint32) []byte{
	originalData := make([]byte, 32*4)
	copy(originalData, data[8:(32*4)+8])

	var shuffleData []byte
	for i:=0; i<4; i++ {
		newIndex := 32 * blockPositions[uint32(i) + shiftValue * 4]

		shuffleData = append(shuffleData, originalData[newIndex:newIndex+32]...)
	}

	// stitch shuffled data back into provided data
	ret := make([]byte, len(data))
	copy(ret, data)
	for i, b := range shuffleData {
		ret[8 + i] = b
	}

	return ret
}

func DecryptData(pkmData []byte) []byte {
	pid := binary.LittleEndian.Uint32(pkmData[0:4])
	shiftValue := ((pid & 0x3E000) >> 0xD) % 24

	pkmData = xorCrypt(pkmData)
	return shuffle( pkmData, uint32(shiftValue) ) 
}

func EncryptData(pkmData []byte) []byte {
	pid := binary.LittleEndian.Uint32(pkmData[0:4])
	shiftValue := ((pid & 0x3E000) >> 0xD) % 24

	pkmData = shuffle( pkmData, uint32(blockPositionInvert[shiftValue]) )
	return xorCrypt(pkmData)
}
