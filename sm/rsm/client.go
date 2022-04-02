package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type client struct {
	conn            net.Conn
	nick            string
	table           *table
	commands        chan<- command
	currentState    int     // This will hold the currrent state of the node
	log             []Event // Contains all events ever executed on this state machine
	stateHash       string
	confirmations   int
	knownNodesState map[string]string // a least of all connected nodes state
}

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
			stack.Push(b + a)
		case SUB:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b - a)
		case MUL:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b * a)
		case DIV:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b / a)
		case MOD:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b % a)
		case AND:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b & a)
		case OR:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b | a)
		case XOR:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b ^ a)
		case LSHIFT:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b << a)
		case RSHIFT:
			a, err := stack.Pop()
			b, err1 := stack.Pop()
			if err != nil {
				return 0, err
			}
			if err1 != nil {
				return 0, err1
			}
			stack.Push(b >> a)
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
	// if the stack has more than one value at the end,
	// i.e. instructions where only PUSH and PUSH
	// the last push represents the new state
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

func (c *client) getHash() (string, error) {
	logStr, err := json.Marshal(c.log)
	if err != nil {
		c.msg("error parsing json: " + err.Error())
		return "", err
	}
	data := []byte(fmt.Sprint(c.currentState) + string(logStr))
	hash := sha256.Sum256(data)
	c.msg("New State Hash: " + hex.EncodeToString(hash[:]))
	return hex.EncodeToString(hash[:]), nil
}

func (c *client) setState() {
	newStateHash, _ := c.getHash()
	c.stateHash = newStateHash
}

func (c *client) msg(msg string) {
	msgArray := strings.Split(msg, "!!")
	
	if msgArray[0] == "STATE" {
		c.knownNodesState[msgArray[2]] = msgArray[1]

	} else if msgArray[0] == "LOG" {
		c.conn.Write([]byte(msg + "\n"))

	} else if msgArray[0] == "UPDATE" {
		json.Unmarshal([]byte(msgArray[1]), &(c.log))

	} else if msgArray[0] == "EXECUTE" {
		c.msg("New Event...")
		timestamp, err := strconv.ParseInt(msgArray[3], 10, 64)

		if err != nil {
			fmt.Println("timestamp error parsing", err)
			return
		}

		if contains(c.log, timestamp) {
			c.msg("Duplicate log: " + fmt.Sprint(timestamp))
			// prevents infinite loop where same event is always broadcasted
			return
		}
		var broadcastedCode []Instruction

		json.Unmarshal([]byte(msgArray[1]), &broadcastedCode)
		newState, err := c.runCode(broadcastedCode)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(3)
		time.Sleep(time.Duration(n) * time.Second) // this adds latency to each node as in a real world example, nodes might be in different locations
		if err != nil {
			c.msg("Execution failed: " + err.Error())
			return
		}
		c.msg("Execution result: " + fmt.Sprint(newState))

		c.setState()

	} else {
		c.conn.Write([]byte("> " + msg + "\n"))
	}
}
