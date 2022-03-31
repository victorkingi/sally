package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"github.com/rs/xid"
)

type server struct {
	nodeTable map[string]*table
	commands  chan command
}

type Instruction struct {
	Opcode Opcode
	Data   int
}

var GlobalStateHash string    // This represents the state hash of the whole rsm after consensus was reached

type Event struct {
	timestamp int64
	Instructions []Instruction
}

func contains(s []Event, e int64) bool {
	for _, a := range s {
		if a.timestamp == e {
			return true
		}
	}
	return false
}

func newServer() *server {
	return &server{
		nodeTable: make(map[string]*table),
		commands:  make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ACTIVE:
			s.nodes(cmd.client)
		case CMD_STATE:
			s.getState(cmd.client)
		case CMD_LOG:
			s.getLog(cmd.client)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) getState(c *client) {
	c.msg("Current State: " + fmt.Sprint(c.currentState))
}

func (s *server) getLog(c *client) {
	c.msg(fmt.Sprintf("Size of log: %d", len(c.log)))
	c.msg(fmt.Sprint("LOG", c.log))
}

func (s *server) nodes(c *client) {
	c.msg(fmt.Sprintf("Total active nodes: %d", len(c.table.members)))
	if len(c.table.members) == 0 {
		return
	}
	for addr, m := range c.table.members {
		c.msg(fmt.Sprintf("%s: %s", m.nick, addr))
	}
}

func (s *server) newClient(conn net.Conn) *client {
	id := xid.New()
	log.Printf("new node joined: %s", conn.RemoteAddr().String())
	log.Printf("unique node ID created: %s\n", id.String())
	log.Printf("syncing with node: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     string(id.String()),
		commands: s.commands,
		currentState: 0,
		stateHash: "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9", // sha256 of 0
		confirmations: 0,
		knownNodesState: make(map[string]string, 0),
	}
}

func (s *server) join(c *client) {

	r, ok := s.nodeTable["node_table"]
	if !ok {
		r = &table{
			name:    "node_table",
			members: make(map[net.Addr]*client),
		}
		s.nodeTable["node_table"] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.removeFromTable(c)
	c.table = r

	r.broadcast(c, fmt.Sprintf("node id: %s connected", c.nick))
}

func (s *server) sanitizeArg(msg []string) ([]Instruction, error) {
	var commands []Instruction

	for _, op := range msg {
		if strings.HasPrefix(op, "PUSH") {
			strOp := string(op[4:])
			intOp, err := strconv.Atoi(strOp)
			if err != nil {
				return nil, errors.New(err.Error())
			}
			opcode, data := push(intOp)
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "ADD") {
			opcode, data := add()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "SUB") {
			opcode, data := sub()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "MUL") {
			opcode, data := mul()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "DIV") {
			opcode, data := div()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "MOD") {
			opcode, data := mod()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "AND") {
			opcode, data := and()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "OR") {
			opcode, data := or()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "XOR") {
			opcode, data := xor()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "LSHIFT") {
			opcode, data := lshift()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "RSHIFT") {
			opcode, data := rshift()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else if strings.HasPrefix(op, "DUMP") {
			opcode, data := dump()
			inst := Instruction{Opcode: opcode, Data: data}
			commands = append(commands, inst)

		} else {
			return nil, errors.New("Invalid instruction, message dropped" + op)
		}
	}
	return commands, nil
}

func (s *server) msg(c *client, args []string) {
	now := time.Now()
	nowNano := now.UTC().UnixNano()

	if c.table == nil {
		c.msg("Node has to be added to table first")
		return
	}
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	// first time server sees the message
	msg := strings.Split(args[1], ";")
	code, err := s.sanitizeArg(msg)

	if err != nil {
		fmt.Println("error parsing instructions", err)
		return
	}

	out, err := json.Marshal(code)
	if err != nil {
		fmt.Println("error parsing json", err)
		return
	}
	c.table.broadcast(c, "EXECUTE!!"+string(out)+"!!"+c.nick+"!!"+fmt.Sprint(nowNano))
	newState, err := c.runCode(code)
	if err != nil {
		c.msg("Execution failed: " + err.Error())
		return
	}
	c.msg("Execution result: " + fmt.Sprint(newState))
	event := Event{timestamp: nowNano, Instructions: code}
	log := append(c.log, event)
	jsonLog, err := json.Marshal(log)
	if err != nil {
		fmt.Println("error parsing json", err)
		return
	}
	c.table.broadcast(c, "UPDATE!!"+string(jsonLog))
	c.setState()
	c.log = log
}

func (s *server) quit(c *client) {
	log.Printf("node disconnected: %s id: %s", c.conn.RemoteAddr().String(), c.nick)
	s.removeFromTable(c)
	c.conn.Close()
}

func (s *server) removeFromTable(c *client) {
	if c.table != nil {
		oldTable := s.nodeTable[c.table.name]
		delete(s.nodeTable[c.table.name].members, c.conn.RemoteAddr())
		oldTable.broadcast(c, fmt.Sprintf("%s: %s disconnected", c.nick, c.conn.RemoteAddr()))
	}
}
