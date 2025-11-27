//go:build dev

package snapshots

import (
	"bytes"
	"emulator/core"
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// 👁️ .vscode/settings.json
// this test ONLY runs manually from VSCode, as it rebuilds all the snapshots
// used in other tests

func TestNewSnapshots(t *testing.T) {
	if os.Getenv("VSCODE") == "true" {

		// 👇 create snapshots in THIS directory
		_, file, _, _ := runtime.Caller(0)
		dir := filepath.Dir(file)
		var perm os.FileMode = 0777 // 👈 seem to need this to work

		for nm, stream := range Index {
			t.Run(fmt.Sprintf("%s snapshot", nm), func(t *testing.T) {

				// 👇 a RW directory for each snapshot
				os.MkdirAll(filepath.Join(dir, nm), perm)

				// 👇 run each snapshot through the emulator
				emu := core.MockEmulator(24, 80)
				emu.Initialize()
				emu.Bus.PubOutbound(stream)

				// 👇 now the Flds and the RGBA we were passed should be complete
				flds, _ := json.Marshal(emu.Flds)
				img := emu.Cfg.RGBA
				var buf bytes.Buffer
				png.Encode(&buf, img)

				// 👇 emit the snapshot
				os.WriteFile(filepath.Join(dir, nm, "flds.json"), []byte(flds), perm)
				os.WriteFile(filepath.Join(dir, nm, "screen.png"), buf.Bytes(), perm)
			})
		}
	}
}
