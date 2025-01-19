package genIV

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"strconv"
	"unicode"
)

// Charmap doesnt have [,], use them to delimit special chars?
var genIVEncoding = [][]string {
	{"", "　", "ぁ", "あ", "ぃ", "い", "ぅ", "う", "ぇ", "え", "ぉ", "お", "か", "が", "き", "ぎ"},
	{"く", "ぐ", "け", "げ", "こ", "ご", "さ", "ざ", "し", "じ", "す", "ず", "せ", "ぜ", "そ", "ぞ"},
	{"た", "だ", "ち", "ぢ", "っ", "つ", "づ", "て", "で", "と", "ど", "な", "に", "ぬ", "ね", "の"},
	{"は", "ば", "ぱ", "ひ", "び", "ぴ", "ふ", "ぶ", "ぷ", "へ", "べ", "ぺ", "ほ", "ぼ", "ぽ", "ま"},
	{"み", "む", "め", "も", "ゃ", "や", "ゅ", "ゆ", "ょ", "よ", "ら", "り", "る", "れ", "ろ", "わ"},
	{"を", "ん", "ァ", "ア", "ィ", "イ", "ゥ", "ウ", "ェ", "エ", "ォ", "オ", "カ", "ガ", "キ", "ギ"},
	{"ク", "グ", "ケ", "ゲ", "コ", "ゴ", "サ", "ザ", "シ", "ジ", "ス", "ズ", "セ", "ゼ", "ソ", "ゾ"},
	{"タ", "ダ", "チ", "ヂ", "ッ", "ツ", "ヅ", "テ", "デ", "ト", "ド", "ナ", "ニ", "ヌ", "ネ", "ノ"},
	{"ハ", "バ", "パ", "ヒ", "ビ", "ピ", "フ", "ブ", "プ", "ヘ", "ベ", "ペ", "ホ", "ボ", "ポ", "マ"},
	{"ミ", "ム", "メ", "モ", "ャ", "ヤ", "ュ", "ユ", "ョ", "ヨ", "ラ", "リ", "ル", "レ", "ロ", "ワ"},
	{"ヲ", "ン", "０", "１", "２", "３", "４", "５", "６", "７", "８", "９", "Ａ", "Ｂ", "Ｃ", "Ｃ"},
	{"Ｅ", "Ｆ", "Ｇ", "Ｈ", "Ｉ", "Ｊ", "Ｋ", "Ｌ", "Ｍ", "Ｎ", "Ｏ", "Ｐ", "Ｑ", "Ｒ", "Ｓ", "Ｔ"},
	{"Ｕ", "Ｖ", "Ｗ", "Ｘ", "Ｙ", "Ｚ", "ａ", "ｂ", "ｃ", "ｄ", "ｅ", "ｆ", "ｇ", "ｈ", "ｉ", "ｊ"},
	{"ｋ", "ｌ", "ｍ", "ｎ", "ｏ", "ｐ", "ｑ", "ｒ", "ｓ", "ｔ", "ｕ", "ｖ", "ｗ", "ｘ", "ｙ", "ｚ"},
	{"",  "！", "？", "、", "。", "…", "・", "／", "「", "」", "『", "』", "（", "）", "[f♂]", "[f♀]"},
	{"＋", "ー", "×", "÷", "＝", "～", "：", "；", "．", "，", "[fspade]", "[fclub]", "[fheart]", "[fdiamond]", "[fstar]", "◎"},
	{"[fcircle]", "[fsquare]", "[fdelta]", "[fhollow diamond]", "＠", "[fmusic]", "％", "[fsun]", "[fcloud]", "[fumbrella]", "[fsnow man]", "[fsmile]", "[flaugh]", "[fyell]", "[ffrown]", "[fpoint up]"},
	{"[fpoint down]", "[fzz]", "円", "[items]", "[key]", "[tm]", "[mail]", "[medicine]", "[berries]", "[balls]", "[battle]", "←", "↑", "↓", "→", "►"},
	{"＆", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E"},
	{"F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U"},
	{"V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"},
	{"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "À"},
	{"Á", "Â", "Ã", "Ä", "Å", "Æ", "Ç", "È", "É", "Ê", "Ë", "Ì", "Í", "Î", "Ï", "Ð"},
	{"Ñ", "Ò", "Ó", "Ô", "Õ", "Ö", "[x]", "Ø", "Ù", "Ú", "Û", "Ü", "Ý", "Þ", "ß", "à"},
	{"á", "â", "ã", "ä", "å", "æ", "ç", "è", "é", "ê", "ë", "ì", "í", "î", "ï", "ð"},
	{"ñ", "ò", "ó", "ô", "õ", "ö", "÷", "ø", "ù", "ú", "û", "ü", "ý", "þ", "ÿ", "Œ"},
	{"œ", "Ş", "ş", "ª", "º", "[er]", "[re]", "[r]", "$", "¡", "¿", "!", "?", ",", ".", "…"},
	{"･", "/", "‘", "'", "“", "”", "„", "«", "»", "(", ")", "♂", "♀", "+", "-", "*"},
	{"#", "=", "&", "~", ":", ";", "[spade]", "[club]", "[heart]", "[diamond]", "[star]", "◎", "[circle]", "[square]", "[delta]", "[hollow diamond]"},
	{"@", "[music]", "%", "[sun]", "[cloud]", "[umbrella]", "[snow man]", "[smile]", "[laugh]", "[yell]", "[frown]", "[point up]", "[point down]", "[zz]", " ", "[e]"},
	{"[pk]", "[mn]", "[figure space]", "[1px space]", "[2px space]", "[4px space]", "[8px space]", "[16px space]", "°", "_", "＿", "․", "‥"},
}

func GetGenIVEncoding(s string) ([]byte, error){
	var ret bytes.Buffer
	scratchBytes := make([]byte, 2)

	// Loop over every character in the string
	for i := 0; i < len( []rune(s) ); i++ {
		var row, col uint16
		var c strings.Builder // current "character" we are looking at
		c.WriteRune( []rune(s)[i] )

		// We are dealing with a special "character"
		if []rune(s)[i] == '[' {
			var tmp int
			for tmp = 1;; {
				// Out of characters, throw error
				if (tmp + i) >= len(s) {
					errorMessage := "unclosed '[' at index " + strconv.Itoa(i)
					return nil, errors.New(errorMessage)
				}

				c.WriteRune(unicode.ToLower( []rune(s)[i+tmp] ))

				if ( []rune(s)[tmp + i] ) == ']' { break }

				tmp++
			}

			i += tmp // increment by 'tmp' to skip over special character data
		}

		// Loop over every possible character encoding until we find our indecies
		breaking := false
		for row = 0; row < uint16( len(genIVEncoding) ); row++ {
			for col = 0; col < uint16( len(genIVEncoding[row]) ); col++ {
				if genIVEncoding[row][col] == c.String() {
					breaking = true
					break
				}
			}
			if breaking { break }
		}
		// if we have not set 'breaking', weve run out of characters to look at in our table
		if !breaking {
			errorMessage := "Unable to locate character '" + c.String() + "'"
			return nil, errors.New(errorMessage)
		}

		row <<= 4
		row |= col
		binary.LittleEndian.PutUint16(scratchBytes, row)
		ret.Write(scratchBytes)
	}


	return ret.Bytes(), nil;
}