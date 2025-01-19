package handlers

import (
	"encoding/base64"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
 * Set the log level
 */
func SetLogLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "i":
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break

	case "e":
	case "error":
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	}
}

/*
 * Set human readable logs
 */
func SetHumanLogs(useHuman bool) {
	
}

/* 
 * Log an error
 */
func LogError(err error) {
	log.Error().Err(err).Msg("")
}

/*
 * Log info
 */
func LogInfo(msg string) {
	log.Info().Msg(msg)
}

/*
 * Encode a string to "Nintendo-Base64"
 */
func EncodeNintendoB64(s string) string {
	tmp := base64.StdEncoding.EncodeToString( []byte( s ) )
	tmp = strings.ReplaceAll(tmp, "=", "*")
	tmp = strings.ReplaceAll(tmp, "/", "-")
	return strings.ReplaceAll(tmp, "+", ".")
}

/*
 * Decode a "Nintendo-Base64" string
 */
func DecodeNintendoB64(s string) ([]byte, error) {
	tmp := strings.ReplaceAll(s, "*", "=")
	tmp = strings.ReplaceAll(tmp, "-", "/")
	tmp = strings.ReplaceAll(tmp, ".", "+")
	return base64.StdEncoding.DecodeString( string(tmp) )
}


var cardsRootPath = "./"
func SetRootPath(path string) {
	cardsRootPath = path
}

func GetRootPath() string {
	return cardsRootPath
}