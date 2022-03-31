package main

import (
	"bufio"
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
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
	conn         net.Conn
	nick         string
	table        *table
	commands     chan<- command
	currentState int // This will hold the currrent state of the node
	stateHash    string
}

var localLog []Event // Contains all events ever executed on this state machine

type stack struct {
	lock sync.Mutex
	s    []int
}

var receivedStates []string
var nodesReceivedFrom []string

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

func (c *client) getHash() (string, error) {
	logStr, err := json.Marshal(localLog)
	if err != nil {
		c.msg("error parsing json: "+err.Error())
		return "", err
	}
	data := []byte(fmt.Sprint(c.currentState)+string(logStr))
	hash := sha256.Sum256(data)
	c.msg("New State Hash: "+hex.EncodeToString(hash[:]))
	return hex.EncodeToString(hash[:]), nil
}

func (c *client) sendState() {
	newStateHash, _ := c.getHash()
	c.stateHash = newStateHash
	c.table.broadcast(c, "STATE!!"+newStateHash+"!!"+c.nick)
}

func (c *client) msg(msg string) {
	msgArray := strings.Split(msg, "!!")

	if msgArray[0] == "STATE" {
		for _, a := range nodesReceivedFrom {
			if a == msgArray[2] {
				return
			}
		}
		receivedStates = append(receivedStates, msgArray[1])
		if len(receivedStates) < len(c.table.members) {
			c.msg("total states received in this epoch: "+fmt.Sprint(len(receivedStates)))
			return
		}
		if len(receivedStates) == len(c.table.members) {
			matches := 0 // good nodes
			for _, a := range receivedStates {
				if a == c.stateHash {
					matches += 1
				}
			}

			//3f+1 rule i.e. if a network can tolerate 1 faulty node then it should have at least 4 other nodes
			faulty := (len(c.table.members) - 1) / 3 // safe bet
			if (len(c.table.members) - matches) > faulty {
				c.msg("RSM COMPROMISED! THE STATE TRANSITION WILL NOT BE CORRECT")
				return
			} else {
				c.msg("CONSENSUS REACHED, TRANSITIONING TO NEXT STATE")
				receivedStates = nil
				return
			}
		}

	} else if msgArray[0] == "LOG" {
		c.conn.Write([]byte(msg + "\n"))

	} else if msgArray[0] == "UPDATE" {
		json.Unmarshal([]byte(msgArray[1]), &localLog)

	} else if msgArray[0] == "EXECUTE" {
		c.msg("New Event...")
		timestamp, err := strconv.ParseInt(msgArray[3], 10, 64)

		if err != nil {
			fmt.Println("timestamp error parsing", err)
			return
		}

		if contains(localLog, timestamp) {
			c.msg("Duplicate log: " + fmt.Sprint(timestamp))
			// prevents infinite loop where same event is always broadcasted
			return
		}
		var broadcastedCode []Instruction

		json.Unmarshal([]byte(msgArray[1]), &broadcastedCode)
		newState, err := c.runCode(broadcastedCode)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second) // this adds latency to each node as in a real world example, nodes might in different locations
		if err != nil {
			c.msg("Execution failed: " + err.Error())
			return
		}
		c.msg("Execution result: " + fmt.Sprint(newState))
		c.sendState() // broadcast new state to all nodes

	} else {
		c.conn.Write([]byte("> " + msg + "\n"))
	}
}
