package genIV

import (
	"errors"
	"strconv"
	"time"
)

// First 3 digits of serial indicate game, 4th digit is only language/region
var GenIVGames = map[string]uint16 {
	"ADA": 0x0004, //"Diamond"
	"APA": 0x0008, //"Pearl"
	"CPU": 0x0010, //"Platinum"
	"IPK": 0x8000, //"HeartGold"
	"IPG": 0x0001, //"SoulSilver"
}

type GenIVWonderCard struct {
	pgt				[260]byte // PCD uses encrypted PGT, WC4 uses decrypted PGT
	cardName		[ 72]byte // Uses custom encoding, unused bytes are 0xFF
	supportedGames	[  2]byte // DP = 0x000C, Pt = 0x0010, HGSS = 0x8001
	unk1 			[  2]byte // Always 0's?
	cardID 			[  2]byte // any number between 0x0000 - 0xFFFF. Must be unique for each WC
	saveCard		[  2]byte // Always 0x0D, 0x00? Save card flag?
	description		[500]byte // Uses custom encoding, unused bytes are 0xFF
	shareCount		[  1]byte // Number of times card can be shared
	unk3			[  3]byte // Always 0's?
	icon			[  2]byte // Index of pokemon to use, caps at 0x1ED for DPPt
	unk4			[  6]byte // Always 0's?
	receiveDate		[  2]byte // Days since Jan. 01, 2000 (caps at Dec. 31, 2099)
	unk5			[  2]byte // Always 0's?
}

// Generate a new Gen IV Wondercard
func NewGenIVWondercard()(* GenIVWonderCard) {
	ret := GenIVWonderCard{}

	for i := 0; i < 72; i++ {
		ret.cardName[i] = 0xFF
	}
	for i := 0; i < 500; i++ {
		ret.description[i] = 0xFF
	}

	ret.SetSaveCard() // most? (all?) official wondercards have this set, set by default

	return &ret
}

// Set the PGT data
func (w *GenIVWonderCard) SetPGT(pgtData []byte) error {
	if len(pgtData) != 260 {
		errorMessage := "incorrect PGT length, expected 260 but got " + strconv.Itoa( len(pgtData) )
		return errors.New(errorMessage)
	}

	for i, b := range pgtData {
		w.pgt[i] = b
	}

	return nil
}

// Set the Card Name
func (w *GenIVWonderCard) SetCardName(name string) error {
	enc, err := GetGenIVEncoding(name)
	if err != nil { return err }

	if len(enc) > 72 {
		enc = enc[:72]
	}
	for i, c := range enc {
		w.cardName[i] = c
	}

	return nil
}

// Set/clear supported game flags
func (w *GenIVWonderCard) SetDPFlag() {
	w.supportedGames[1] |= 0x0C
}
func (w *GenIVWonderCard) ClearDPFlag() {
	w.supportedGames[1] &= ^byte(0x0C)
}
func (w *GenIVWonderCard) SetPtFlag() {
	w.supportedGames[1] |= 0x10
}
func (w *GenIVWonderCard) ClearPtFlag() {
	w.supportedGames[1] &= ^byte(0x10)
}
func (w *GenIVWonderCard) SetHGSSFlag() {
	w.supportedGames[0] |= 0x80
	w.supportedGames[1] |= 0x01
}
func (w *GenIVWonderCard) ClearHGSSFlag() {
	w.supportedGames[0] &= ^byte(0x80)
	w.supportedGames[1] &= ^byte(0x01)
}

func (w *GenIVWonderCard) SetSaveCard() {
	w.saveCard[0] = 0x0D
}
func (w *GenIVWonderCard) ClearSaveCard() {
	w.saveCard[0] = 0x00
}
func (w *GenIVWonderCard) GetSaveCard() bool {
	return w.saveCard[0] != 0 
}

// Set the Card ID
func (w *GenIVWonderCard) SetCardID(id uint16) {
	w.cardID[0] = byte( (id & 0xFF00) >> 8 )
	w.cardID[1] = byte( id & 0x00FF )
}

// Set the Card Description
func (w *GenIVWonderCard) SetCardDescription(desc string) error {
	enc, err := GetGenIVEncoding(desc)
	if err != nil { return err }

	if len(enc) > 500 {
		enc = enc[:500]
	}
	for i, c := range enc {
		w.description[i] = c
	}

	return nil
}

func (w *GenIVWonderCard) SetShareCount(count uint8) {
	w.shareCount[0] = count
}

// Set Pokemon icon
func (w *GenIVWonderCard) SetIcon(pkmn uint16) error {
	if pkmn > 0x1ED {
		errorMessage := "index " + strconv.FormatUint(uint64(pkmn), 10) + " is out of bounds (max is 0x01ED)"
		return errors.New(errorMessage)
	}

	w.icon[0] = byte(pkmn & 0x00FF)
	w.icon[1] = byte( (pkmn & 0xFF00) >> 8 )

	return nil
}

// Set Received Date
func (w *GenIVWonderCard) SetDate(date time.Time) error {
	startDate := time.Date( 2000,  1,  1, 0, 0, 0, 0, date.Location() )
	endDate   := time.Date( 2099, 12, 31, 0, 0, 0, 0, date.Location() )
	truncDate := time.Date( date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location() )

	if truncDate.Before(startDate) {
		errorMessage := "date is before Jan. 01 2000"
		return errors.New(errorMessage)
	}
	if truncDate.After(endDate) {
		errorMessage := "date is after Dec. 31 2099"
		return errors.New(errorMessage)
	}


	offset := uint16( (truncDate.Sub( startDate ).Hours()) / 24 )

	w.receiveDate[1] = byte( (offset & 0xFF00) >> 8 )
	w.receiveDate[0] = byte( offset & 0x00FF )

	return nil
}


func (w *GenIVWonderCard) GetGameFlags() []byte {
	return w.supportedGames[:]
}

func (w *GenIVWonderCard) GetCardID() []byte {
	return w.cardID[:]
}

func (w *GenIVWonderCard) ToBytes() []byte {
	ret := make([]byte, 0, 0x358)

	ret = append(ret, w.pgt[:]...)
	ret = append(ret, w.cardName[:]...)
	ret = append(ret, w.supportedGames[:]...)
	ret = append(ret, w.unk1[:]...)
	ret = append(ret, w.cardID[:]...)
	ret = append(ret, w.saveCard[:]...)
	ret = append(ret, w.description[:]...)
	ret = append(ret, w.shareCount[:]...)
	ret = append(ret, w.unk3[:]...)
	ret = append(ret, w.icon[:]...)
	ret = append(ret, w.unk4[:]...)
	ret = append(ret, w.receiveDate[:]...)
	ret = append(ret, w.unk5[:]...)

	return ret
}