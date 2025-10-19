package consts

type QCode byte

const (
	ALPHANUMERIC_PARTITIONS QCode = 0x84
	CHARACTER_SETS          QCode = 0x85
	COLOR_SUPPORT           QCode = 0x86
	DDM                     QCode = 0x95
	HIGHLIGHTING            QCode = 0x87
	IMPLICIT_PARTITION      QCode = 0xA6
	REPLY_MODES             QCode = 0x88
	RPQ_NAMES               QCode = 0xA1
	SUMMARY                 QCode = 0x80
	USABLE_AREA             QCode = 0x81
)

var qcodes = map[QCode]string{
	0x84: "ALPHANUMERIC_PARTITIONS",
	0x85: "CHARACTER_SETS",
	0x86: "COLOR_SUPPORT",
	0x95: "DDM",
	0x87: "HIGHLIGHTING",
	0xA6: "IMPLICIT_PARTITION",
	0x88: "REPLY_MODES",
	0xA1: "RPQ_NAMES",
	0x80: "SUMMARY",
	0x81: "USABLE_AREA",
}

func QCodeFor(s QCode) string {
	return qcodes[s]
}

func (s QCode) String() string {
	return QCodeFor(s)
}
