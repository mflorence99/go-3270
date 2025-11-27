package fonts

import (
	_ "embed"
)

// 🔥 Hack alert! we must use extension {js, wasm}
//    and we can't use symlinks, so this file is a copy of the font renamed

var (
	//go:embed JuliaMono-Regular.ttf.wasm
	NormalFontEmbed []byte
	//go:embed JuliaMono-Bold.ttf.wasm
	BoldFontEmbed []byte
)
