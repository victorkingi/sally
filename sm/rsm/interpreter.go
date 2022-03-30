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
	lock sync.Mutex
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
	// at this stage the stack should contain only one element since all operations
	// where successful
	ans, err := stack.Pop()
	if err != nil {
		return 0, err
	}
	return ans, nil
}
