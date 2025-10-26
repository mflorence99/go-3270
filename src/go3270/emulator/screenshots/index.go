package screenshots

var Index map[string][]byte

func init() {
	Index = make(map[string][]byte)
	Index["termtest"] = TERMTEST
}
