package main

import (
	"log"
	"os"
	"runtime"

	"github.com/WillKirkmanM/chip-8/chip8"
	"github.com/WillKirkmanM/chip-8/renderer"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
    emulator := chip8.New()

    if runtime.GOOS != "js" {
        romPath := "./roms/pong.rom"
        if len(os.Args) > 1 {
            romPath = os.Args[1]
        }

        if err := emulator.LoadROM(romPath); err != nil {
            log.Printf("Failed to load ROM %s: %v", romPath, err)
            log.Println("Continuing with no ROM. You can drag and drop a ROM file onto the window.")
        }
        } else {
            // On web: register JavaScript callbacks for loading ROMs
            registerCallbacks(emulator)
        }
    
    game := renderer.New(emulator)

    ebiten.SetWindowSize(chip8.DisplayWidth*renderer.Scale, chip8.DisplayHeight*renderer.Scale)
    ebiten.SetWindowTitle("Chip-8 Emulator")
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
