package renderer

import (
	"fmt"
	"image/color"
	"io/fs"
	"log"

	"github.com/WillKirkmanM/chip-8/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// Scale factor for the display
	Scale = 10

	// Cycles per update
	CyclesPerUpdate = 10

	// RGBA components for colours
	ColourR = 0
	ColourG = 1
	ColourB = 2
	ColourA = 3

	// Colour values
	ColourOn  = 0xFF // White
	ColourOff = 0x00 // Black

	// Bytes per pixel
	BytesPerPixel = 4 // R, G, B, A
)

type Game struct {
	Chip8       *chip8.Chip8
	PixelBuffer []byte
	KeyMap      map[ebiten.Key]uint8
}

func New(c *chip8.Chip8) *Game {
	game := &Game{
		Chip8:       c,
		PixelBuffer: make([]byte, chip8.DisplayWidth*chip8.DisplayHeight*BytesPerPixel),

		// Chip8 keypad layout: 	Mapped to keyboard:
		// 1 2 3 C  				1 2 3 4
		// 4 5 6 D 					Q W E R
		// 7 8 9 E					A S D F
		// A 0 B F					Z X C V

		KeyMap: map[ebiten.Key]uint8{
			ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
			ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
			ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
			ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
		},
	}

	return game
}

func (g *Game) checkForDroppedFiles() {
	droppedFiles := ebiten.DroppedFiles()
	if droppedFiles != nil {
		entries, err := fs.ReadDir(droppedFiles, ".")
		if err != nil || len(entries) == 0 {
			return
		}

		fileName := entries[0].Name()
		file, err := fs.ReadFile(droppedFiles, fileName)
		if err != nil {
			log.Printf("Failed to read dropped file: %v", err)
			return
		}

		if err := g.Chip8.LoadROMFromBytes(file); err != nil {
			log.Printf("Failed to load dropped file as ROM: %v", err)
			return
		}

		log.Printf("Successfully loaded ROM: %s", fileName)
	}
}

func (g *Game) Update() error {
	g.checkForDroppedFiles()

	for key, value := range g.KeyMap {
		if ebiten.IsKeyPressed(key) {
			g.Chip8.Keys[value] = true
		} else {
			g.Chip8.Keys[value] = false
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("game ended by player")
	}

	for i := 0; i < CyclesPerUpdate; i++ {
		g.Chip8.EmulateCycle()
		g.Chip8.UpdateTimers()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Chip8.DrawFlag {
		g.updatePixelBuffer()
		g.Chip8.DrawFlag = false
	}

	for y := 0; y < chip8.DisplayHeight; y++ {
		var r, green, b, a byte
		for x := 0; x < chip8.DisplayWidth; x++ {
			pixelPos := (y*chip8.DisplayWidth + x) * BytesPerPixel

			r = g.PixelBuffer[pixelPos+ColourR]
			green = g.PixelBuffer[pixelPos+ColourG]
			b = g.PixelBuffer[pixelPos+ColourB]
			a = g.PixelBuffer[pixelPos+ColourA]

			screen.Set(x*Scale+Scale/2, y*Scale+Scale/2, color.RGBA{r, green, b, a})

			for dy := 0; dy < Scale; dy++ {
				for dx := 0; dx < Scale; dx++ {
					screen.Set(x*Scale+dx, y*Scale+dy, color.RGBA{r, green, b, a})
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return chip8.DisplayWidth * Scale, chip8.DisplayHeight * Scale
}

func (g *Game) updatePixelBuffer() {
	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			pixelPos := (y*chip8.DisplayWidth + x) * BytesPerPixel

			colourVal := ColourOff
			if g.Chip8.Display[y][x] {
				colourVal = ColourOn
			}

			colour := byte(colourVal)

			g.PixelBuffer[pixelPos+ColourR] = colour
			g.PixelBuffer[pixelPos+ColourG] = colour
			g.PixelBuffer[pixelPos+ColourB] = colour
			g.PixelBuffer[pixelPos+ColourA] = 0xFF
		}
	}
}
