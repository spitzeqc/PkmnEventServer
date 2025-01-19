package wondercard

import (
	"time"
	genIV "cornchip.com/pkmneventserver/CardGenerator/GenIV"
)

const pkmnUserName = "CORNCHP"
const pkmnTID = 30636
const pkmnSID = 11419


func GetEmptyWondercard() *genIV.GenIVWonderCard {
	ret := genIV.NewGenIVWondercard()

	tmp := make( []byte, 260 )

	ret.SetPGT(tmp)
	ret.SetCardID(0x00)
	ret.SetCardName("Empty Card")
	ret.SetCardDescription("You should not be able to read this")
	ret.SetIcon(1)
	
	return ret
}


func GenerateWondercard() *genIV.GenIVWonderCard {
	ret := genIV.NewGenIVWondercard()

	tmp := make([]byte, 260)
	tmp[0] = 0x0a

	ret.SetPGT(tmp)
	ret.SetCardID(0x0421)
	ret.SetCardName("Test Card")
	ret.SetCardDescription("Test desc")
	ret.SetIcon(1)

	ret.SetDate( time.Now() )

	ret.SetPtFlag()
	ret.SetDPFlag()

	return ret
}