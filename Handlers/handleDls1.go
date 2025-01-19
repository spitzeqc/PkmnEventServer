package handlers

import (
//	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"

	wondercard "cornchip.com/pkmneventserver/CardGenerator"
	genIV "cornchip.com/pkmneventserver/CardGenerator/GenIV"
)

/*
 * Generate WonderCard data that no game will accept (return on error)
 */
func generateErrorWonderCard() []byte {
	tmp := make([]byte, 80, 80)
	card := wondercard.GetEmptyWondercard()
	ret := append(tmp, card.ToBytes()[:]...)
	ret[73] = 0x20 // Set "supported games" to something invalid to prevent any game from reading
	return ret
}

/*
 * Get WonderCard data from 'GenIVCardPath'
 */
func getWonderCardData(gameId string) ([]byte, error) {
	cardPath := path.Join( GetRootPath(), "geniv" )
	f, err := os.Open( cardPath )
	// folder doesnt exist, make it
	if errors.Is(err, os.ErrNotExist) {
		errorMessage := "Directory " + cardPath + " does not exist"
		return nil, errors.New(errorMessage)
	}

	//TODO: filter by extension
	cardList, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	if len(cardList) == 0 {
		errorMessage := "No wondercard files found in " + cardPath
		return nil, errors.New(errorMessage)
	}

	
	startIndex := rand.Intn( len(cardList) )
	var file *os.File
	for i := startIndex + 1; i != startIndex; i++{
		if i >= len(cardList) { i = 0 } // Loop back to 0 if needed

		if filepath.Ext(cardList[i]) != ".wc4" && filepath.Ext(cardList[i]) != "pcd" { continue } // Skip invalid files

		file, err = os.Open( path.Join(cardPath, cardList[i]) )
		if err == nil {
			tmp := make( []byte, 2 )
			file.Seek( 332, 0 ) // Seek to supportedgames offset
			file.Read(tmp)

			var gameVal uint16 = uint16(tmp[0]) << 8 | uint16(tmp[1])
			gameFlag := genIV.GenIVGames[gameId] // "invalid" game id will return 0 and mask off all game bits

			//TODO: Update to check based on connected game
			if (gameVal & gameFlag) != 0 { break }
		}

		file.Close()
		file = nil
	}
	if file == nil {
		errorMessage := "Could not find valid wondercard for game version " + gameId
		return nil, errors.New(errorMessage)
	}

	cardData := make([]byte, 0x358, 0x358)
	file.Seek(0x00, 0)
	file.Read(cardData)

	// Encrypt WC data if file is wc4
	stat, err := file.Stat()
	if err != nil { return nil, err }
	if filepath.Ext( stat.Name() ) == ".wc4" {
		encData := wondercard.EncryptData( cardData[8:] ) // Program calls it with top 8 bytes removed
		decData := cardData[0:8]
		cardData = append(decData, encData[:]...)
	}

	// Description shown when receiving wondercard
	contentDesc := [72]byte{}
	for i:=0; i<72; i++ {
		contentDesc[i] = 0xFF
	}
	miscFlags := [8]byte{}

	// Set supported games to match the WC file
	miscFlags[0] = cardData[332]
	miscFlags[1] = cardData[333]

	// Set card id to match the WC file
	miscFlags[4] = cardData[336] 
	miscFlags[5] = cardData[337]


	miscFlags[6] = cardData[ 338 ] // Save card flag? 0x0D = save card, 0x00 = no card?


	enc := cardData[260:332]
	for i, c := range enc {
		contentDesc[i] = c
	}

	var ret []byte
	ret = append(ret, contentDesc[:]...)
	ret = append(ret, miscFlags[:]...)
	ret = append(ret, cardData[:]...)
	
	return ret, nil
}


const countActionString = "Y291bnQ*"
const listActionString = "bGlzdA**"
const contentsActionString = "Y29udGVudHM*"

/*
 * Handle requests to "/dls1" 
 */
func HandleDls1(w http.ResponseWriter, r *http.Request) {
	LogInfo("/dls1 requested")
	if r.URL.String() != "/dls1/download" {
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

	// Return count of items on server
	if r.FormValue( "action" ) == countActionString {
		LogInfo("/dls1/download count requested")
		w.Write( []byte("1") )
		return
	} else if r.FormValue( "action" ) == listActionString {
		LogInfo("/dls1/download list requested")
		w.Write( []byte("21dppUS.myg\t\t\t\t\t936\r\n") )

	// Return raw pct data
	} else if r.FormValue( "action" ) == contentsActionString {
		LogInfo("/dls1/download count requested")
		w.Header().Set("Content-Disposition", "attachment; filename=21dppUS.myg")

		tmp, err := DecodeNintendoB64(r.FormValue("gamecd"))
		if err != nil {
			errorMessage := "Could not decode Nintendo-B64 value " + r.FormValue("gamecd")
			LogError( errors.New(errorMessage) )
			w.Write( generateErrorWonderCard() )
		}

		card, err := getWonderCardData( string(tmp)[:3])
		if err != nil {
			LogError(err)
			card = generateErrorWonderCard()
		}

		w.Write( card )
	}

}