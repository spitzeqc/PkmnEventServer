package gamestats2

import (
	"crypto/sha1"
)

func calcGeneralHash(token string, hashConst string) []byte {
	hasher := sha1.New()
	hasher.Write( []byte(hashConst + token) )
	return hasher.Sum(nil)
}

func GenIVHash(token string) []byte {
	return calcGeneralHash(token, "sAdeqWo3voLeC5r16DYv")
}

func GenVHash(token string) []byte {
	return calcGeneralHash(token, "HZEdGCzcGGLvguqUEKQN")
}

func DecryptPayload(keystream []byte, g *GRNG) []byte {
	g.ResetGRNG() // reset so we know we are in sync with keystream

	for i, b := range keystream {
		byteKey := byte( (g.GetVal() >> 16) & 0xFF )
		keystream[i] = b ^ byteKey
		g.Next()
	}

	return keystream
}

// Encryption and decryption process are the same, function exists for readability reasons
func EncryptPayload(plainstream []byte, g *GRNG) []byte {
	return DecryptPayload(plainstream, g)
}

func calcGeneralChecksum(pid uint32, payload []byte, xorConst uint32) uint32 {
	ret := pid
	for _, b := range payload {
		ret += uint32(b)
	}

	return ret ^ xorConst
}

func GenIVChecksum(pid uint32, payload []byte) uint32 {
	return calcGeneralChecksum(pid, payload, 0x4a3b2c1d)
}

func GenVChecksum(pid uint32, payload []byte) uint32 {
	return calcGeneralChecksum(pid, payload, 0x2db842b2)
}