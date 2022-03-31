package main

import (
	"net"
)

type table struct {
	name    string
	members map[net.Addr]*client
}

func (r *table) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
				m.msg(msg)
		}
	}
}
