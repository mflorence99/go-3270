package types

import (
	"emulator/utils"
	"regexp"
	"strconv"
	"strings"
)

// 🟧 3270 AIDs (attention identifiers)

type AID byte

// 🟦 Lookup tables

const (
	INBOUND AID = 0x88
	CLEAR   AID = 0x6D
	ENTER   AID = 0x7D
	PA1     AID = 0x6C
	PA2     AID = 0x6E
	PA3     AID = 0x6B
	PF1     AID = 0xf1
	PF2     AID = 0xf2
	PF3     AID = 0xf3
	PF4     AID = 0xf4
	PF5     AID = 0xf5
	PF6     AID = 0xf6
	PF7     AID = 0xf7
	PF8     AID = 0xf8
	PF9     AID = 0xf9
	PF10    AID = 0x7A
	PF11    AID = 0x7B
	PF12    AID = 0x7C
	PF13    AID = 0xc1
	PF14    AID = 0xc2
	PF15    AID = 0xc3
	PF16    AID = 0xc4
	PF17    AID = 0xc5
	PF18    AID = 0xc6
	PF19    AID = 0xc7
	PF20    AID = 0xc8
	PF21    AID = 0xc9
	PF22    AID = 0x4A
	PF23    AID = 0x4B
	PF24    AID = 0x4C
)

var aids = map[AID]string{
	0x88: "INBOUND",
	0x6D: "CLEAR",
	0x7D: "ENTER",
	0x6C: "PA1",
	0x6E: "PA2",
	0x6B: "PA3",
	0xf1: "PF1",
	0xf2: "PF2",
	0xf3: "PF3",
	0xf4: "PF4",
	0xf5: "PF5",
	0xf6: "PF6",
	0xf7: "PF7",
	0xf8: "PF8",
	0xf9: "PF9",
	0x7A: "PF10",
	0x7B: "PF11",
	0x7C: "PF12",
	0xc1: "PF13",
	0xc2: "PF14",
	0xc3: "PF15",
	0xc4: "PF16",
	0xc5: "PF17",
	0xc6: "PF18",
	0xc7: "PF19",
	0xc8: "PF20",
	0xc9: "PF21",
	0x4A: "PF22",
	0x4B: "PF23",
	0x4C: "PF24",
}

var aidsLookup = make(map[string]AID)

func init() {
	aidsLookup = utils.Invert(aids)
}

// 🟦 Constructor

func AIDOf(key string, alt, ctrl, shift bool) AID {
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

// 🟦 Public functions

func (a AID) PAx() bool {
	str, ok := aids[a]
	return ok && strings.HasPrefix(str, "PA")
}

func (a AID) PFx() bool {
	str, ok := aids[a]
	return ok && strings.HasPrefix(str, "PF")
}

func (a AID) ShortRead() bool {
	return a == CLEAR || a.PAx()
}

// 🟦 Stringer implementation

func AIDFor(a AID) string {
	return aids[a]
}

func (a AID) String() string {
	return AIDFor(a)
}
