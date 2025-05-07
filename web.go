//go:build js && wasm

package main

import (
	"log"
	"syscall/js"

	"github.com/WillKirkmanM/chip-8/chip8"
)

func registerCallbacks(emulator *chip8.Chip8) {
    log.Println("Registering web callbacks")

    loadROMFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        if len(args) != 1 {
            return "Invalid number of arguments"
        }

        jsArray := args[0]
        length := jsArray.Get("length").Int()
        romData := make([]byte, length)
        
        js.CopyBytesToGo(romData, jsArray)
        
        if err := emulator.LoadROMFromBytes(romData); err != nil {
            return err.Error()
        }
        
        return "ROM loaded successfully"
    })

    js.Global().Set("loadROM", loadROMFunc)
}