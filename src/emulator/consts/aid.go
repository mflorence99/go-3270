package consts

import (
	"emulator/utils"
	"regexp"
	"strconv"
	"strings"
)

var (
	INBOUND byte = 0x88
	CLEAR   byte = 0x6D
	ENTER   byte = 0x7D
	PA1     byte = 0x6C
	PA2     byte = 0x6E
	PA3     byte = 0x6B
	PF1     byte = 0xF1
	PF2     byte = 0xF2
	PF3     byte = 0xF3
	PF4     byte = 0xF4
	PF5     byte = 0xF5
	PF6     byte = 0xF6
	PF7     byte = 0xF7
	PF8     byte = 0xF8
	PF9     byte = 0xF9
	PF10    byte = 0x7A
	PF11    byte = 0x7B
	PF12    byte = 0x7C
	PF13    byte = 0xC1
	PF14    byte = 0xC2
	PF15    byte = 0xC3
	PF16    byte = 0xC4
	PF17    byte = 0xC5
	PF18    byte = 0xC6
	PF19    byte = 0xC7
	PF20    byte = 0xC8
	PF21    byte = 0xC9
	PF22    byte = 0x4A
	PF23    byte = 0x4B
	PF24    byte = 0x4C
)

var aids = map[byte]string{
	0x88: "INBOUND",
	0x6D: "CLEAR",
	0x7D: "ENTER",
	0x6C: "PA1",
	0x6E: "PA2",
	0x6B: "PA3",
	0xF1: "PF1",
	0xF2: "PF2",
	0xF3: "PF3",
	0xF4: "PF4",
	0xF5: "PF5",
	0xF6: "PF6",
	0xF7: "PF7",
	0xF8: "PF8",
	0xF9: "PF9",
	0x7A: "PF10",
	0x7B: "PF11",
	0x7C: "PF12",
	0xC1: "PF13",
	0xC2: "PF14",
	0xC3: "PF15",
	0xC4: "PF16",
	0xC5: "PF17",
	0xC6: "PF18",
	0xC7: "PF19",
	0xC8: "PF20",
	0xC9: "PF21",
	0x4A: "PF22",
	0x4B: "PF23",
	0x4C: "PF24",
}

var aidsLookup = make(map[string]byte)

func AIDFor(aid byte) string {
	return aids[aid]
}

func AIDOf(key string, alt, ctrl, shift bool) byte {
	code := strings.ToUpper(key)
	re := regexp.MustCompile(`F([0-9]+)`)
	matches := re.FindStringSubmatch(code)
	switch {
	case code == "ENTER":
		return ENTER
	case code == "ESCAPE":
		return CLEAR
	case !alt && !ctrl && !shift && len(matches) == 2:
		num, _ := strconv.Atoi(matches[1])
		return aidsLookup["PF"+strconv.Itoa(num)]
	case !alt && !ctrl && shift && len(matches) == 2:
		num, _ := strconv.Atoi(matches[1])
		return aidsLookup["PF"+strconv.Itoa(num+12)]
	case alt && !ctrl && !shift && len(matches) == 2:
		return aidsLookup["PA"+matches[1]]
	}
	return 0
}

func PAx(aid byte) bool {
	str, ok := aids[aid]
	return ok && strings.HasPrefix(str, "PA")
}

func PFx(aid byte) bool {
	str, ok := aids[aid]
	return ok && strings.HasPrefix(str, "PF")
}

func init() {
	aidsLookup = utils.Invert(aids)
}
