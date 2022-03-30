package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/rs/xid"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ACTIVE:
			s.nodes(cmd.client)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) nodes(c *client) {
	c.msg(fmt.Sprintf("Total active nodes: %d", len(c.room.members)))
	if len(c.room.members) == 0 {
		return
	}
	for addr, m := range c.room.members {
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
	}
}

func (s *server) join(c *client) {

	r, ok := s.rooms["main_table"]
	if !ok {
		r = &room{
			name:   "main_table",
			members: make(map[net.Addr]*client),
		}
		s.rooms["main_table"] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("node id: %s connected", c.nick))

	c.msg(fmt.Sprintf("Node successfully synced"))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.msg("Node has to be added to table first")
		return
	}
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("node disconnected: %s id: %s", c.conn.RemoteAddr().String(), c.nick)
	s.quitCurrentRoom(c)
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
