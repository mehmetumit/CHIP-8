# CHIP-8
* CHIP-8 emulator written in go
## References
* https://en.wikipedia.org/wiki/CHIP-8
* https://www.cs.columbia.edu/~sedwards/classes/2016/4840-spring/designs/Chip8.pdf
* https://github.com/corax89/chip8-test-rom
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
## Custom Keymap
```
Default     Custom
  '1'   ->   '1',
  '2'   ->   '2',
  '3'   ->   '3',
  'C'   ->   '4',
  '4'   ->   'Q',
  '5'   ->   'W',
  '6'   ->   'E',
  'D'   ->   'R',
  '7'   ->   'A',
  '8'   ->   'S',
  '9'   ->   'D',
  'E'   ->   'F',
  'A'   ->   'Z',
  '0'   ->   'X',
  'B'   ->   'C',
  'F'   ->   'V'
```
