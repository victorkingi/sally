package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start boot node: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("boot node started on port :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		c := s.newClient(conn)
		epoch := 0
		s.join(c)
		go c.readInput()
		go func() {
			i := 0
			for {
				time.Sleep(1 * time.Second)
				c.table.broadcast(c, "STATE!!"+c.stateHash+"!!"+c.nick) // broadcast state to all nodes
				if i%20 == 0 {
					fmt.Println(c.conn.RemoteAddr(), "state:", c.stateHash)
				}
				i++
			}
		}()
		// Reach consensus every 10 seconds by comparing state hash
		go func() {
			j := 0
			for {
				rand.Seed(time.Now().UnixNano())
				min := 10
				max := 12
				n := rand.Intn(max - min + 1) + min
				time.Sleep(time.Duration(n) * time.Second)
				confirmations := 0 // non-faulty node
				for _, elem := range c.knownNodesState {
					if elem == c.stateHash {
						confirmations += 1
					}
				}
				// 3f + 1 rule will be used to reach consensus
				if len(c.table.members) > (3*(len(c.table.members)-confirmations)+1) || (len(c.table.members) == 3 && confirmations == 2) {
					if j % 5 == 0 {
						// slow down number of times this is printed, give user time to enter more commands
						c.msg("CONSENSUS REACHED FOR EPOCH: " + fmt.Sprint(epoch))
					}
					epoch += 1
				} else {
					c.msg("UNABLE TO TRANSITION FROM EPOCH: " + fmt.Sprint(epoch))
				}
				j++
			}
		}()
	}
}
