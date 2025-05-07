package chip8

import (
	"fmt"
	"math/rand"
	"os"
)

const (
	// Display resolution
	DisplayWidth  = 64
	DisplayHeight = 32

	// Memory configuration
	MemorySize       = 4096  // 4KB of memory
	ProgramStartAddr = 0x200 // Programs start at 0x200 (512 decimal)
	FontsetStartAddr = 0x000 // Fontset stored at 0x000-0x050
	FontEndAddr      = 0x050 // End of fontset in memory

	// Register configuration
	RegisterCount = 16 // 16 general-purpose registers (V0-VF)

	// Stack configuration
	StackSize = 16 // 16 levels of stack

	// Keypad configuration
	KeyCount = 16 // 16 keys for input

	// Sprite configuration
	FontSpriteHeight = 5 // Each font sprite is 5 bytes tall
	SpriteWidth      = 8 // Sprites are 8 pixels wide

	// Register indices
	VFRegister = 0xF // VF is the flag register

	// Opcode masks
	OpcodeMask    = 0xF000 // Mask for first nibble (instruction type)
	XRegisterMask = 0x0F00 // Mask for X register in opcode
	YRegisterMask = 0x00F0 // Mask for Y register in opcode
	NibbleMask    = 0x000F // Mask for lowest nibble
	AddressMask   = 0x0FFF // Mask for 12-bit address
	ByteMask      = 0x00FF // Mask for 8-bit immediate value

	// Opcode shifts
	XRegisterShift = 8 // Shift amount to get X register
	YRegisterShift = 4 // Shift amount to get Y register

	// Instruction size
	InstructionSize = 2 // Each instruction is 2 bytes

	// Math constants
	ByteMax       = 0xFF  // Maximum value of a byte (255)
	SpriteLeftBit = 0x80  // Leftmost bit in a byte (10000000)
	AddressMax    = 0xFFF // Maximum 12-bit address (4095)

	// BCD constants
	Hundreds = 100 // For BCD conversion
	Tens     = 10  // For BCD conversion
)

// Chip8 represents the emulator's state
type Chip8 struct {
	// CPU registers
	V  [RegisterCount]byte // 16 8-bit registers V0-VF
	I  uint16              // 16-bit index register
	PC uint16              // Program counter
	SP byte                // Stack pointer
	DT byte                // Delay timer
	ST byte                // Sound timer

	// Memory components
	Memory [MemorySize]byte  // 4KB memory
	Stack  [StackSize]uint16 // 16 levels of stack

	// Display
	Display [DisplayHeight][DisplayWidth]bool

	// Input
	Keys [KeyCount]bool // 16 keys for input

	// Flags
	DrawFlag bool // Indicates screen needs redrawing
}

func New() *Chip8 {
	c := &Chip8{}

	// Load fonts into memory (at the beginning of memory)
	fontSet := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	// Copy font data to memory
	copy(c.Memory[:], fontSet)

	// Initialize other components
	c.PC = ProgramStartAddr
	c.DrawFlag = true

	return c
}

func (c *Chip8) LoadROM(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load ROM: %v", err)
	}

	return c.LoadROMFromBytes(rom)
}

func (c *Chip8) LoadROMFromBytes(rom []byte) error {
	// First, clear memory except for font data (first 80 bytes typically)
	// Keep font area intact (usually 0x000-0x1FF)
	for i := FontEndAddr; i < MemorySize; i++ {
		c.Memory[i] = 0
	}

	maxRomSize := MemorySize - ProgramStartAddr
	if len(rom) > maxRomSize {
		return fmt.Errorf("ROM size exceeds available memory (max: %d bytes)", maxRomSize)
	}

	for i := 0; i < len(rom); i++ {
		c.Memory[ProgramStartAddr+i] = rom[i]
	}

	// Reset CPU state when loading a new ROM
	c.PC = ProgramStartAddr
	c.I = 0
	c.SP = 0
	c.DT = 0
	c.ST = 0
	c.DrawFlag = true

	// Clear the display
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			c.Display[y][x] = false
		}
	}

	// Reset registers
	for i := 0; i < RegisterCount; i++ {
		c.V[i] = 0
	}

	return nil
}

func (c *Chip8) UpdateTimers() {
	// Update delay timer
	if c.DT > 0 {
		c.DT--
	}

	// Update sound timer
	if c.ST > 0 {
		c.ST--
		// If we had sound, we would play it when ST > 0
	}
}

func (c *Chip8) EmulateCycle() {
	// Fetch opcode (2 bytes)
	opcode := uint16(c.Memory[c.PC]) << 8 | uint16(c.Memory[c.PC+1])

	// Increment PC before execution
	c.PC += InstructionSize

	// Decode and execute opcode
	switch opcode & OpcodeMask {
	case 0x0000:
		switch opcode & ByteMask {
		case 0x00E0: // CLS: Clear the display
			for y := 0; y < DisplayHeight; y++ {
				for x := 0; x < DisplayWidth; x++ {
					c.Display[y][x] = false
				}
			}
			c.DrawFlag = true
		case 0x00EE: // RET: Return from subroutine
			c.SP--
			c.PC = c.Stack[c.SP]
		}
	case 0x1000: // JP addr: Jump to location nnn
		c.PC = opcode & AddressMask
	case 0x2000: // CALL addr: Call subroutine at nnn
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = opcode & AddressMask
	case 0x3000: // SE Vx, byte: Skip next instruction if Vx = kk
		x := (opcode & XRegisterMask) >> XRegisterShift
		if c.V[x] == byte(opcode&ByteMask) {
			c.PC += InstructionSize
		}
	case 0x4000: // SNE Vx, byte: Skip next instruction if Vx != kk
		x := (opcode & XRegisterMask) >> XRegisterShift
		if c.V[x] != byte(opcode&ByteMask) {
			c.PC += InstructionSize
		}
	case 0x5000: // SE Vx, Vy: Skip next instruction if Vx = Vy
		x := (opcode & XRegisterMask) >> XRegisterShift
		y := (opcode & YRegisterMask) >> YRegisterShift
		if c.V[x] == c.V[y] {
			c.PC += InstructionSize
		}
	case 0x6000: // LD Vx, byte: Set Vx = kk
		x := (opcode & XRegisterMask) >> XRegisterShift
		c.V[x] = byte(opcode & ByteMask)
	case 0x7000: // ADD Vx, byte: Set Vx = Vx + kk
		x := (opcode & XRegisterMask) >> XRegisterShift
		c.V[x] += byte(opcode & ByteMask)
	case 0x8000:
		x := (opcode & XRegisterMask) >> XRegisterShift
		y := (opcode & YRegisterMask) >> YRegisterShift
		switch opcode & NibbleMask {
		case 0x0000: // LD Vx, Vy: Set Vx = Vy
			c.V[x] = c.V[y]
		case 0x0001: // OR Vx, Vy: Set Vx = Vx OR Vy
			c.V[x] |= c.V[y]
		case 0x0002: // AND Vx, Vy: Set Vx = Vx AND Vy
			c.V[x] &= c.V[y]
		case 0x0003: // XOR Vx, Vy: Set Vx = Vx XOR Vy
			c.V[x] ^= c.V[y]
		case 0x0004: // ADD Vx, Vy: Set Vx = Vx + Vy, set VF = carry
			sum := uint16(c.V[x]) + uint16(c.V[y])
			if sum > ByteMax {
				c.V[VFRegister] = 1
			} else {
				c.V[VFRegister] = 0
			}
			c.V[x] = byte(sum)
		case 0x0005: // SUB Vx, Vy: Set Vx = Vx - Vy, set VF = NOT borrow
			if c.V[x] > c.V[y] {
				c.V[VFRegister] = 1
			} else {
				c.V[VFRegister] = 0
			}
			c.V[x] -= c.V[y]
		case 0x0006: // SHR Vx: Set Vx = Vx SHR 1
			c.V[VFRegister] = c.V[x] & 0x1
			c.V[x] >>= 1
		case 0x0007: // SUBN Vx, Vy: Set Vx = Vy - Vx, set VF = NOT borrow
			if c.V[y] > c.V[x] {
				c.V[VFRegister] = 1
			} else {
				c.V[VFRegister] = 0
			}
			c.V[x] = c.V[y] - c.V[x]
		case 0x000E: // SHL Vx: Set Vx = Vx SHL 1
			c.V[VFRegister] = c.V[x] >> 7
			c.V[x] <<= 1
		}
	case 0x9000: // SNE Vx, Vy: Skip next instruction if Vx != Vy
		x := (opcode & XRegisterMask) >> XRegisterShift
		y := (opcode & YRegisterMask) >> YRegisterShift
		if c.V[x] != c.V[y] {
			c.PC += InstructionSize
		}
	case 0xA000: // LD I, addr: Set I = nnn
		c.I = opcode & AddressMask
	case 0xB000: // JP V0, addr: Jump to location nnn + V0
		c.PC = (opcode & AddressMask) + uint16(c.V[0])
	case 0xC000: // RND Vx, byte: Set Vx = random byte AND kk
		x := (opcode & XRegisterMask) >> XRegisterShift
		c.V[x] = byte(rand.Intn(256)) & byte(opcode&ByteMask)
	case 0xD000: // DRW Vx, Vy, nibble: Display n-byte sprite at (Vx, Vy)
		x := uint16(c.V[(opcode&XRegisterMask)>>XRegisterShift]) % DisplayWidth
		y := uint16(c.V[(opcode&YRegisterMask)>>YRegisterShift]) % DisplayHeight
		height := opcode & NibbleMask

		c.V[VFRegister] = 0

		for row := uint16(0); row < height; row++ {
			if int(y+row) >= DisplayHeight {
				break
			}

			sprite := c.Memory[c.I+row]

			for col := uint16(0); col < SpriteWidth; col++ {
				if int(x+col) >= DisplayWidth {
					break
				}

				// Check if current pixel in sprite is set
				if (sprite & (SpriteLeftBit >> col)) != 0 {
					// Check if we're going to flip the pixel on screen
					if c.Display[y+row][x+col] {
						c.V[VFRegister] = 1
					}
					c.Display[y+row][x+col] = !c.Display[y+row][x+col]
				}
			}
		}
		c.DrawFlag = true
	case 0xE000:
		x := (opcode & XRegisterMask) >> XRegisterShift
		switch opcode & ByteMask {
		case 0x009E: // SKP Vx: Skip next instruction if key with value Vx is pressed
			if c.Keys[c.V[x]] {
				c.PC += InstructionSize
			}
		case 0x00A1: // SKNP Vx: Skip next instruction if key with value Vx is not pressed
			if !c.Keys[c.V[x]] {
				c.PC += InstructionSize
			}
		}
	case 0xF000:
		x := (opcode & XRegisterMask) >> XRegisterShift
		switch opcode & ByteMask {
		case 0x0007: // LD Vx, DT: Set Vx = delay timer value
			c.V[x] = c.DT
		case 0x000A: // LD Vx, K: Wait for a key press, store key value in Vx
			keyPressed := false
			for i := 0; i < KeyCount; i++ {
				if c.Keys[i] {
					c.V[x] = byte(i)
					keyPressed = true
					break
				}
			}
			// If no key is pressed, repeat the instruction
			if !keyPressed {
				c.PC -= InstructionSize
			}
		case 0x0015: // LD DT, Vx: Set delay timer = Vx
			c.DT = c.V[x]
		case 0x0018: // LD ST, Vx: Set sound timer = Vx
			c.ST = c.V[x]
		case 0x001E: // ADD I, Vx: Set I = I + Vx
			if c.I+uint16(c.V[x]) > AddressMax {
				c.V[VFRegister] = 1
			} else {
				c.V[VFRegister] = 0
			}
			c.I += uint16(c.V[x])
		case 0x0029: // LD F, Vx: Set I = location of sprite for digit Vx
			c.I = uint16(c.V[x] * FontSpriteHeight) // Each font sprite is 5 bytes tall
		case 0x0033: // LD B, Vx: Store BCD representation of Vx in memory[I, I+1, I+2]
			c.Memory[c.I] = c.V[x] / Hundreds
			c.Memory[c.I+1] = (c.V[x] / Tens) % Tens
			c.Memory[c.I+2] = c.V[x] % Tens
		case 0x0055: // LD [I], Vx: Store registers V0 through Vx in memory starting at I
			for i := uint16(0); i <= x; i++ {
				c.Memory[c.I+i] = c.V[i]
			}
		case 0x0065: // LD Vx, [I]: Read registers V0 through Vx from memory starting at I
			for i := uint16(0); i <= x; i++ {
				c.V[i] = c.Memory[c.I+i]
			}
		}
	}
}
