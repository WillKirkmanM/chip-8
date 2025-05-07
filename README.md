<p align="center">
  <img src="https://avatars.githubusercontent.com/u/138057124?s=200&v=4" width="150" />
</p>
<h1 align="center">CHIP-8</h1>

<p align="center">
CHIP-8 is an interpreted programming language created by Joseph Weisbecker in the mid-1970s. It was designed to simplify game programming on 8-bit microcomputers like the COSMAC VIP and Telmac 1800. CHIP-8 programs run on a virtual machine with the following specifications:
</p>

- **Memory:** 4KB (4096 bytes)
- **Display:** 64×32 pixels, monochrome
- **Registers:** 16 8-bit registers (V0 to VF)
- **Index Register:** 16-bit (I)
- **Program Counter:** 16-bit (PC)
- **Stack:** 16-level for return addresses
- **Timers:** 8-bit delay and sound timers
- **Input:** 16-key hexadecimal keypad (0-F)

## Build Instructions

### Prerequisites

- Go 1.16+

### Troubleshooting

**If you see "Go is not defined":**

-   Verify `wasm_exec.js` was copied correctly
-   Check browser console for errors
-   Ensure you're using a supported browser
-   Try clearing browser cache

**If WebAssembly fails to load:**

-   Check that the WASM file was built correctly
-   Verify you're using a web server (not opening HTML directly)
-   Check browser console for detailed error messages

### Web Browser Support

-   Chrome/Edge (recommended) ≥ 57
-   Firefox ≥ 52
-   Safari ≥ 11
- Web browser with WebAssembly support

### Project Structure

```
chip-8/
├── web/               # Web deployment files
│   ├── index.html    # Web interface
│   ├── wasm_exec.js  # Go's WebAssembly support
│   └── chip8.wasm    # Compiled WebAssembly binary
├── roms/             # CHIP-8 ROM files
├── main.go           # Main application entry
├── web.go            # Web-specific code
└── build_web.bat     # Build automation script
```

### Build Commands

1.  **Create necessary directories:**

    ```bash
    mkdir web
    mkdir roms
    ```
2.  **Copy WebAssembly support file:**

    ```bash
    copy "%GOROOT%\misc\wasm\wasm_exec.js" web\
    ```
3.  **Build WebAssembly binary:**

    ```bash
    set GOOS=js
    set GOARCH=wasm
    go build -o web/chip8.wasm
    ```
4.  **Start development server:**

    ```bash
    cd web
    python -m http.server 8080
    ```

## Input Mapping

| CHIP-8 Keypad | Keyboard |
| :-----------: | :------: |
|     +-+-+-+-+     | +-+-+-+-+ |
|     \|1\|2\|3\|C\|     | \|1\|2\|3\|4\| |
|     +-+-+-+-+     | +-+-+-+-+ |
|     \|4\|5\|6\|D\|     | \|Q\|W\|E\|R\| |
|     +-+-+-+-+     | +-+-+-+-+ |