package consts

import "fmt"

type SFID byte

type SFld struct {
	ID   SFID
	Info []byte
}

func (s SFld) String() string {
	return fmt.Sprintf("{ID %#02x, Info % #x}", byte(s.ID), s.Info)
}

const (
	QUERY_REPLY    SFID = 0x81
	READ_PARTITION SFID = 0x01
)

var sfids = map[SFID]string{
	0x81: "QUERY_REPLY",
	0x01: "READ_PARTITION",
}

func SFIDFor(s SFID) string {
	return sfids[s]
}

func (s SFID) String() string {
	return SFIDFor(s)
}
