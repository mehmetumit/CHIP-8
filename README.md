# CHIP-8
* CHIP-8 emulator written in go
## Demo
![IBMLogo](/assets/IBMLogo.png)
![Test](/assets/Test.png)
![Pong](/assets/Pong.png)
![Tetris](/assets/Tetris.png)
## Usage
```
Usage of ./CHIP-8:
  -path string
        The file path of rom (default "./roms/Instruction-Test.ch8")
  -scale int
        The display scale (default 12)
  -speed uint
        The emulation speed (default 3)
```
### Run
```
# Without creating executable in current folder
# It can take some time on first run because of the sdl2 package
$ go run . -path <./roms/Pong.ch8> -speed <3> -scale <12>
# Using executable file which is created after build operation
$ ./CHIP-8 -path <./roms/Pong.ch8> -speed <3> -scale <12>
```
### Build
```
# Print the build process using flags
# It can take some time on first build because of the sdl2 package
$ go build -x -v
```
## Dependencies
* `Go`
* `SDL2(go-sdl2)`
* `cgo`
## Opcodes
```
00E0 00EE (0NNN)-> not necessary for most roms
1NNN
2NNN
3XNN
4XNN
5XY0
6XNN
7XNN
8XY0 8XY1 8XY2 8XY3 8XY4 8XY5 8XY6 8XY7 8XYE
9XY0
ANNN
BNNN
CXNN
DXYN
EX9E EXA1
FX07 FX0A FX15 FX18 FX1E FX29 FX33 FX55 FX65
```
## Keymaps
```
Default        Custom
1 2 3 C        1 2 3 4
4 5 6 D   ->   Q W E R
7 8 9 E        A S D F
A 0 B F        Z X C V
```
## References
* https://en.wikipedia.org/wiki/CHIP-8
* https://www.cs.columbia.edu/~sedwards/classes/2016/4840-spring/designs/Chip8.pdf
* https://github.com/corax89/chip8-test-rom
