package main

type Opcode int

type Intepreter struct {}

const (
	PUSH Opcode = iota
	ADD
	SUB
	MUL
	DIV
	MOD
	AND
	OR
	XOR
	LSHIFT
	RSHIFT
	DUMP
)

func push(x int) (Opcode, int) { return PUSH, x }
func add() (Opcode, int) { return ADD, 0 }
func sub() (Opcode, int) { return SUB, 0 }
func mul() (Opcode, int) { return MUL, 0 }
func div() (Opcode, int) { return DIV, 0 }
func mod() (Opcode, int) { return MOD, 0 }
func and() (Opcode, int) { return AND, 0 }
func or() (Opcode, int)  { return OR, 0 }
func xor() (Opcode, int) { return XOR, 0 }
func lshift() (Opcode, int) { return LSHIFT, 0 }
func rshift() (Opcode, int) { return RSHIFT, 0 }
func dump() (Opcode, int) { return DUMP, 0 }
