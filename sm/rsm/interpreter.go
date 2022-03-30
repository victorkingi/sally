package main

import (
	"sync"
	"errors"
	"fmt"
)

type Opcode int

const (
	// Default represents the default state of the system.
	Default State = "0"

	// NoOp represents a no-op event.
	NoOp EventType = "NoOp"

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

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []int
}

func NewStack() *stack {
	return &stack{sync.Mutex{}, make([]int, 0)}
}

func (s *stack) Push(v int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func push(x int) (Opcode, int, bool) { return PUSH, x, true }
func add() (Opcode, int, bool) { return ADD, 0, true }
func sub() (Opcode, int, bool) { return SUB, 0, true }
func mul() (Opcode, int, bool) { return MUL, 0, true }
func div() (Opcode, int, bool) { return DIV, 0, true }
func mod() (Opcode, int, bool) { return MOD, 0, true }
func and() (Opcode, int, bool) { return AND, 0, true }
func or() (Opcode, int, bool)  { return OR, 0, true }
func xor() (Opcode, int, bool) { return XOR, 0, true }
func lshift() (Opcode, int, bool) { return LSHIFT, 0, true }
func rshift() (Opcode, int, bool) { return RSHIFT, 0, true }
func dump() (Opcode, int, bool) { return DUMP, 0, true }

func runCode(program []struct{Opcode; int; bool}) (int, error) {
	stack := NewStack()
	for _, op := range program {
		switch op.Opcode {
		case PUSH:
			stack.Push(op.int)
		case ADD:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a + b)
		case SUB:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a - b)
		case MUL:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a * b)
		case DIV:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a / b)
		case MOD:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a % b)
		case AND:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a & b)
		case OR:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a | b)
		case XOR:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a ^ b)
		case LSHIFT:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a << b)
		case RSHIFT:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(a >> b)
		case DUMP:
			a, err := stack.Pop()
			if err != nil {
				return 0, err
			}
			fmt.Println(a)
		default:
			panic("Unreachable case")
		}
	}
	ans, err := stack.Pop()
	if err != nil {
		return 0, err
	}
	return ans, nil
}
