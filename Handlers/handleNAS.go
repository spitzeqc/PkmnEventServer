package handlers

import (
//	"encoding/base64"
	"errors"
//	"net"
	"net/http"
	"strings"
	"time"

	tokens "cornchip.com/pkmneventserver/Tokens"
)

const loginActionString = "bG9naW4*"
const acctcreateActionString = "YWNjdGNyZWF0ZQ**"
const svclocActionString = "U1ZDTE9D"

const locatorAddress = "gamespy.com"
const svchostAddress = "dls1.nintendowifi.net"

/*
 * Handle requests to "/nas"
 */
func HandleNAS(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/nas/ac" {
		w.WriteHeader(502)
		return
	}

	if r.Method != "POST" {
		return
	}
	
	err := r.ParseForm()
	if err != nil {
		LogError(err)
		return
	}

	t := time.Now()
	w.Header().Set( "Date", t.Format("RFC1123") )
	w.Header().Set( "Server", "Cornchip" )
	w.Header().Set( "NODE", "wifiappw3" )
	w.Header().Set( "Vary", "Accept-Encoding" )

	retBuilder := strings.Builder{}

	// bG9naW4
	if r.FormValue("action") == loginActionString {
		LogInfo("/nas/ac login requested")
		token := tokens.GenerateToken( r.FormValue("macadr") )
		if token == "" {
			errorMessage := "No token generated for " + r.FormValue("macadr")
			LogError( errors.New(errorMessage) )
			return
		}

		dateVal := t.Format("20060102150405")

		retBuilder.WriteString("challenge=")
		retBuilder.WriteString( tokens.GenerateChallenge() )
		retBuilder.WriteString("&locator=")
		retBuilder.WriteString( EncodeNintendoB64(locatorAddress) )
		retBuilder.WriteString("&retry=MA**")
		retBuilder.WriteString("&returncd=MDAx")
		retBuilder.WriteString("&token=")
		retBuilder.WriteString( token )
		retBuilder.WriteString("&datetime=")
		retBuilder.WriteString( EncodeNintendoB64(dateVal) )

	// U1ZDTE9D
	} else if r.FormValue("action") == svclocActionString {
		LogInfo("/nas/ac svcloc requested")
		token, exists := tokens.GetToken( r.FormValue("macadr") )
		if !exists {
			errorMessage := "No token found for " + r.FormValue("macadr")
			LogError( errors.New(errorMessage) )
			return
		}

		dateVal := t.Format("20060102150405")

		retBuilder.WriteString("retry=MA**")
		retBuilder.WriteString("&returncd=MDA3")
		retBuilder.WriteString("&token=")
		retBuilder.WriteString( token )
		retBuilder.WriteString("&servicetoken=")
		retBuilder.WriteString( token )

		retBuilder.WriteString("&statusdata=WQ**")
		retBuilder.WriteString("&svchost=")
		retBuilder.WriteString( EncodeNintendoB64(svchostAddress) )

		retBuilder.WriteString("&datetime=")
		retBuilder.WriteString( EncodeNintendoB64(dateVal) )


	// YWNjdGNyZWF0ZQ**
	} else if r.FormValue("action") == acctcreateActionString {
		LogInfo("/nas/ac acctcreate requested")
		dateVal := t.Format("20060102150405")
		token, exists := tokens.GetToken( r.FormValue("macadr") )
		if !exists {
			LogInfo("No token found for " + r.FormValue("macadr") )
			tokens.GenerateToken( r.FormValue("macadr") )
			retBuilder.WriteString("retry=MA**")
			retBuilder.WriteString("&returncd=MDAy")
			retBuilder.WriteString("&userid=")

			retBuilder.WriteString("NzEyODk4NzYzNTgxMjcyMg**")

			retBuilder.WriteString("&datetime=")
			retBuilder.WriteString(dateVal)
		} else {
			retBuilder.WriteString("challege=")
			retBuilder.WriteString( EncodeNintendoB64("687TF0EG") )
			retBuilder.WriteString("&locator=")
			retBuilder.WriteString( EncodeNintendoB64(locatorAddress) )
			retBuilder.WriteString("&retry=MA**") // encoding for '0'
			retBuilder.WriteString("&returncd=MDAx") // encoding for '001'
			retBuilder.WriteString("&token=")
			retBuilder.WriteString( token )
			retBuilder.WriteString("&datetime=")
			retBuilder.WriteString( EncodeNintendoB64(dateVal) )
		}

	} else {
		return
	}

	retBuilder.WriteString("\r\n")
	w.Write( []byte(retBuilder.String()) )
}