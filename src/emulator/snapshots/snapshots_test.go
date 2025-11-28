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

	"github.com/google/go-cmp/cmp"
)

// 👁️ .vscode/settings.json
// this test ONLY runs manually from VSCode, as it rebuilds all the snapshots
// used in other tests

func TestNewSnapshots(t *testing.T) {
	// 🔥 to be ABSOLUTELY sure you only run this when you have to
	//    recreate the snapshots, change below to "true"
	if os.Getenv("VSCODE") == "true" {

		// 👇 create snapshots in THIS directory
		_, file, _, _ := runtime.Caller(0)
		dir := filepath.Dir(file)
		var perm os.FileMode = 0777 // 👈 seem to need this to work

		for nm, stream := range Index {
			t.Run(fmt.Sprintf("create %s snapshot", nm), func(t *testing.T) {

				// 👇 a RW directory for each snapshot
				os.MkdirAll(filepath.Join(dir, nm), perm)

				// 👇 run each snapshot through the emulator
				emu := core.MockEmulator(32, 80)
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
	} else {
		t.Skip("🔥 snapshot creation disabled")
	}
}

// 🟦 this test compares the snapshots with what's actually being
//    produced now

func TestOldSnapshots(t *testing.T) {
	// 👇 snapshots reside in THIS directory
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	for nm, stream := range Index {
		t.Run(fmt.Sprintf("validate %s snapshot", nm), func(t *testing.T) {

			// 👇 run each snapshot through the emulator
			emu := core.MockEmulator(32, 80)
			emu.Initialize()
			emu.Bus.PubOutbound(stream)

			// 👇 what is expectedFlds was recorded on disk
			var expectedFlds []core.Flds
			raw, _ := os.ReadFile(filepath.Join(dir, nm, "flds.json"))
			json.Unmarshal(raw, &expectedFlds)

			// 👇 un/marshal the actualFlds Flds to wipe unexported fields
			var actualFlds []core.Flds
			flds, _ := json.Marshal(emu.Flds)
			json.Unmarshal(flds, &actualFlds)

			// 👇 compare expected vs actual Flds
			if diff := cmp.Diff(expectedFlds, actualFlds); diff != "" {
				t.Log(diff)
				t.Fail()
			}

			// 👇 now do the same for the screen image
			expectedImg, _ := os.ReadFile(filepath.Join(dir, nm, "screen.png"))
			img := emu.Cfg.RGBA
			var actualImg bytes.Buffer
			png.Encode(&actualImg, img)
			if !bytes.Equal(expectedImg, actualImg.Bytes()) {
				t.Logf("🔥 %s/screen.png differs from snapshot", nm)
				t.Fail()
			}

		})
	}
}
