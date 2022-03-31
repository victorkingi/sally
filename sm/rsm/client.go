package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"errors"
	"sync"
)

type client struct {
	conn     net.Conn
	nick     string
	table    *table
	commands chan<- command
	currentState int             // This will hold the currrent state of the node
}

var localLog []Instructions      // Contains all events ever executed on this state machine

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
		return 0, errors.New("empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (c *client) runCode(program []Instruction) (int, error) {
	stack := NewStack()
	//push current state, the stack always has the current state as the first entry
	stack.Push(c.currentState)
	for _, op := range program {
		switch op.Opcode {
		case PUSH:
			stack.Push(op.Data)
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
	c.currentState = ans
	return c.currentState, nil
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		case "/state":
			c.commands <- command{
				id:     CMD_STATE,
				client: c,
			}
		case "/nodes":
			c.commands <- command{
				id:     CMD_ACTIVE,
				client: c,
			}
		case "/log":
			c.commands <- command{
				id:     CMD_LOG,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	msgArray := strings.Split(msg, "!!")
	if msgArray[0] == "UPDATE" {
		json.Unmarshal([]byte(msgArray[1]), &localLog)

	} else if msgArray[0] == "EXECUTE" {
		c.msg("New Event...")
		timestamp, err := strconv.ParseInt(msgArray[3], 10, 64)

		if err != nil {
			fmt.Println("timestamp error parsing", err)
			return
		}

		if contains(localLog, timestamp) {
			c.msg("Duplicate log: "+string(timestamp))
			// prevents infinite loop where same event is always broadcasted
			return
		}
		var broadcastedCode []Instruction

		json.Unmarshal([]byte(msgArray[1]), &broadcastedCode)
		newState, err := c.runCode(broadcastedCode)
		if err != nil {
			c.msg("Execution failed: "+err.Error())
			return
		}
		c.msg("Execution result: "+fmt.Sprint(newState))
		fmt.Println("before", pendingMessages)
		findAndDelete(pendingMessages, timestamp)
		fmt.Println("after", pendingMessages)

	} else {
		c.conn.Write([]byte("> " + msg + "\n"))
	}
}
