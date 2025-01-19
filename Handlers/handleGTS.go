package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tokens "cornchip.com/pkmneventserver/Tokens"
)

func HandleGTS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(502)
		return
	}
	passedVals := r.URL.Query()

	tmp, err := strconv.ParseInt( passedVals.Get("pid"), 10, 32)
	if err != nil {
		fmt.Println("Invalid PID provided")
		return
	}

	pid := uint32( tmp )
	hash := passedVals.Get("hash")
	data := passedVals.Get("data")

	// Client is requesting a challenge
	retBuilder := strings.Builder{}
	if pid != 0 && hash == "" && data == "" {
		retBuilder.WriteString("challenge=")
		retBuilder.WriteString( tokens.GenerateChallenge() )
	}

	if pid != 0 && hash != "" && data != "" {

	}

	w.Write( []byte(retBuilder.String()) )
}