package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"
)

var validTokens map[string]string

func encodeBytesNintendoB64(b []byte) string {
	tmp := base64.StdEncoding.EncodeToString( b )
	tmp = strings.ReplaceAll(tmp, "=", "*")
	tmp = strings.ReplaceAll(tmp, "/", "-")
	return strings.ReplaceAll(tmp, "+", ".")
}

func GenerateToken(mac string) string {
	if validTokens == nil {
		validTokens = map[string]string{}
	}

	if mac == "" {
		return ""
	}

	buf := make([]byte, 64)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	ret := encodeBytesNintendoB64(buf)
	validTokens[mac] = ret
	return ret
}

func GetToken(mac string) (string, bool) {
	ret, ret2 := validTokens[mac]
	return ret, ret2
}

func RemoveToken(mac string) {
	if _, b := validTokens[mac]; b {
		delete(validTokens, mac)
	}
}

func GenerateChallenge() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	ret := encodeBytesNintendoB64(buf)
	return ret
}