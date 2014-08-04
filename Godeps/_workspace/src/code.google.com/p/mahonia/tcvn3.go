package mahonia

// Converters for TCVN3 encoding.

import (
	"sync"
)

var (
	onceTCVN3 sync.Once
	dataTCVN3 = struct {
		UnicodeToWord map[rune][2]byte
		WordToUnicode [256]struct {
			r rune
			m *[256]rune
		}
	}{}
)

func init() {
	p := new(Charset)
	p.Name = "TCVN3"
	p.NewDecoder = func() Decoder {
		onceTCVN3.Do(buildTCVN3Tables)
		return decodeTCVN3
	}
	p.NewEncoder = func() Encoder {
		onceTCVN3.Do(buildTCVN3Tables)
		return encodeTCVN3
	}
	RegisterCharset(p)
}

func decodeTCVN3(p []byte) (rune, int, Status) {
	if len(p) == 0 {
		return 0, 0, NO_ROOM
	}
	item := &dataTCVN3.WordToUnicode[p[0]]
	if item.m != nil && len(p) > 1 {
		if r := item.m[p[1]]; r != 0 {
			return r, 2, SUCCESS
		}
	}
	if item.r != 0 {
		return item.r, 1, SUCCESS
	}
	if p[0] < 0x80 {
		return rune(p[0]), 1, SUCCESS
	}
	return '?', 1, INVALID_CHAR
}

func encodeTCVN3(p []byte, c rune) (int, Status) {
	if len(p) == 0 {
		return 0, NO_ROOM
	}
	if c < rune(0x80) {
		p[0] = byte(c)
		return 1, SUCCESS
	}
	if v, ok := dataTCVN3.UnicodeToWord[c]; ok {
		if v[1] != 0 {
			if len(p) < 2 {
				return 0, NO_ROOM
			}
			p[0] = v[0]
			p[1] = v[1]
			return 2, SUCCESS
		} else {
			p[0] = v[0]
			return 1, SUCCESS
		}
	}
	p[0] = '?'
	return 1, INVALID_CHAR
}

func buildTCVN3Tables() {
	dataTCVN3.UnicodeToWord = map[rune][2]byte{
		// one byte
		0x00C2: {0xA2, 0x00},
		0x00CA: {0xA3, 0x00},
		0x00D4: {0xA4, 0x00},
		0x00E0: {0xB5, 0x00},
		0x00E1: {0xB8, 0x00},
		0x00E2: {0xA9, 0x00},
		0x00E3: {0xB7, 0x00},
		0x00E8: {0xCC, 0x00},
		0x00E9: {0xD0, 0x00},
		0x00EA: {0xAA, 0x00},
		0x00EC: {0xD7, 0x00},
		0x00ED: {0xDD, 0x00},
		0x00F2: {0xDF, 0x00},
		0x00F3: {0xE3, 0x00},
		0x00F4: {0xAB, 0x00},
		0x00F5: {0xE2, 0x00},
		0x00F9: {0xEF, 0x00},
		0x00FA: {0xF3, 0x00},
		0x00FD: {0xFD, 0x00},
		0x0102: {0xA1, 0x00},
		0x0103: {0xA8, 0x00},
		0x0110: {0xA7, 0x00},
		0x0111: {0xAE, 0x00},
		0x0129: {0xDC, 0x00},
		0x0169: {0xF2, 0x00},
		0x01A0: {0xA5, 0x00},
		0x01A1: {0xAC, 0x00},
		0x01AF: {0xA6, 0x00},
		0x01B0: {0xAD, 0x00},
		0x1EA1: {0xB9, 0x00},
		0x1EA3: {0xB6, 0x00},
		0x1EA5: {0xCA, 0x00},
		0x1EA7: {0xC7, 0x00},
		0x1EA9: {0xC8, 0x00},
		0x1EAB: {0xC9, 0x00},
		0x1EAD: {0xCB, 0x00},
		0x1EAF: {0xBE, 0x00},
		0x1EB1: {0xBB, 0x00},
		0x1EB3: {0xBC, 0x00},
		0x1EB5: {0xBD, 0x00},
		0x1EB7: {0xC6, 0x00},
		0x1EB9: {0xD1, 0x00},
		0x1EBB: {0xCE, 0x00},
		0x1EBD: {0xCF, 0x00},
		0x1EBF: {0xD5, 0x00},
		0x1EC1: {0xD2, 0x00},
		0x1EC3: {0xD3, 0x00},
		0x1EC5: {0xD4, 0x00},
		0x1EC7: {0xD6, 0x00},
		0x1EC9: {0xD8, 0x00},
		0x1ECB: {0xDE, 0x00},
		0x1ECD: {0xE4, 0x00},
		0x1ECF: {0xE1, 0x00},
		0x1ED1: {0xE8, 0x00},
		0x1ED3: {0xE5, 0x00},
		0x1ED5: {0xE6, 0x00},
		0x1ED7: {0xE7, 0x00},
		0x1ED9: {0xE9, 0x00},
		0x1EDB: {0xED, 0x00},
		0x1EDD: {0xEA, 0x00},
		0x1EDF: {0xEB, 0x00},
		0x1EE1: {0xEC, 0x00},
		0x1EE3: {0xEE, 0x00},
		0x1EE5: {0xF4, 0x00},
		0x1EE7: {0xF1, 0x00},
		0x1EE9: {0xF8, 0x00},
		0x1EEB: {0xF5, 0x00},
		0x1EED: {0xF6, 0x00},
		0x1EEF: {0xF7, 0x00},
		0x1EF1: {0xF9, 0x00},
		0x1EF3: {0xFA, 0x00},
		0x1EF5: {0xFE, 0x00},
		0x1EF7: {0xFB, 0x00},
		0x1EF9: {0xFC, 0x00},
		// two bytes
		0x00C0: {0x41, 0xB5},
		0x00C1: {0x41, 0xB8},
		0x00C3: {0x41, 0xB7},
		0x00C8: {0x45, 0xCC},
		0x00C9: {0x45, 0xD0},
		0x00CC: {0x49, 0xD7},
		0x00CD: {0x49, 0xDD},
		0x00D2: {0x4F, 0xDF},
		0x00D3: {0x4F, 0xE3},
		0x00D5: {0x4F, 0xE2},
		0x00D9: {0x55, 0xEF},
		0x00DA: {0x55, 0xF3},
		0x00DD: {0x59, 0xFD},
		0x0128: {0x49, 0xDC},
		0x0168: {0x55, 0xF2},
		0x1EA0: {0x41, 0xB9},
		0x1EA2: {0x41, 0xB6},
		0x1EA4: {0xA2, 0xCA},
		0x1EA6: {0xA2, 0xC7},
		0x1EA8: {0xA2, 0xC8},
		0x1EAA: {0xA2, 0xC9},
		0x1EAC: {0xA2, 0xCB},
		0x1EAE: {0xA1, 0xBE},
		0x1EB0: {0xA1, 0xBB},
		0x1EB2: {0xA1, 0xBC},
		0x1EB4: {0xA1, 0xBD},
		0x1EB6: {0xA1, 0xC6},
		0x1EB8: {0x45, 0xD1},
		0x1EBA: {0x45, 0xCE},
		0x1EBC: {0x45, 0xCF},
		0x1EBE: {0xA3, 0xD5},
		0x1EC0: {0xA3, 0xD2},
		0x1EC2: {0xA3, 0xD3},
		0x1EC4: {0xA3, 0xD4},
		0x1EC6: {0xA3, 0xD6},
		0x1EC8: {0x49, 0xD8},
		0x1ECA: {0x49, 0xDE},
		0x1ECC: {0x4F, 0xE4},
		0x1ECE: {0x4F, 0xE1},
		0x1ED0: {0xA4, 0xE8},
		0x1ED2: {0xA4, 0xE5},
		0x1ED4: {0xA4, 0xE6},
		0x1ED6: {0xA4, 0xE7},
		0x1ED8: {0xA4, 0xE9},
		0x1EDA: {0xA5, 0xED},
		0x1EDC: {0xA5, 0xEA},
		0x1EDE: {0xA5, 0xEB},
		0x1EE0: {0xA5, 0xEC},
		0x1EE2: {0xA5, 0xEE},
		0x1EE4: {0x55, 0xF4},
		0x1EE6: {0x55, 0xF1},
		0x1EE8: {0xA6, 0xF8},
		0x1EEA: {0xA6, 0xF5},
		0x1EEC: {0xA6, 0xF6},
		0x1EEE: {0xA6, 0xF7},
		0x1EF0: {0xA6, 0xF9},
		0x1EF2: {0x59, 0xFA},
		0x1EF4: {0x59, 0xFE},
		0x1EF6: {0x59, 0xFB},
		0x1EF8: {0x59, 0xFC},
	}
	for r, b := range dataTCVN3.UnicodeToWord {
		item := &dataTCVN3.WordToUnicode[b[0]]
		if b[1] == 0 {
			item.r = r
		} else {
			if item.m == nil {
				item.m = new([256]rune)
			}
			item.m[b[1]] = r
		}
	}
}
