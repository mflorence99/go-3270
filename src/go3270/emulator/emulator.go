package emulator

import (
	"go3270/emulator/buffer"
	"go3270/emulator/bus"
	"go3270/emulator/keystroke"
	"image"

	"golang.org/x/image/font"
)

type Config struct {
	BgColor      string
	BoldFace     *font.Face
	CLUT         map[int][2]string
	Color        [2]string
	Cols         int
	FontHeight   float64
	FontSize     float64
	FontWidth    float64
	NormalFace   *font.Face
	PaddedHeight float64
	PaddedWidth  float64
	Rgba         *image.RGBA
	Rows         int
}

type Emulator struct {
	bus *bus.Bus
	buf *buffer.Buffer
	key *keystroke.Keystroke
}

func NewEmulator(bus *bus.Bus) *Emulator {
	e := new(Emulator)
	e.bus = bus
	e.bus.Subscribe("close", e.close)
	e.bus.Subscribe("config", e.configure)
	return e
}

func (e *Emulator) close() {}

func (e *Emulator) configure(cfg *Config) {
	e.buf = buffer.NewBuffer(cfg.Rows * cfg.Cols)
	e.key = keystroke.NewKeystroke(e.bus)
}
