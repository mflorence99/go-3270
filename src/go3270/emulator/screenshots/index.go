package screenshots

// 🟧 Prefabricated outbound data streams for testing and performance measurement, captured via WireShark

var Index map[string][]byte

func init() {
	Index = make(map[string][]byte)
	Index["ge"] = GE
	Index["termtest"] = TERMTEST
}
