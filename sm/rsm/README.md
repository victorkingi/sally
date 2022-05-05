# Simple Programmable Replicated State Machine

This directory contains an implementation of a replicated state machine. Note that this is not a true RSM as it relies on `server.go` to accept connections, remove and add peers to the node table, sanitize messages and propagate messages to other nodes. With this, it might seem as just another client server model but in most ways, it does behave as an RSM. For example:
1. If two nodes end up with different state hashes, they both stop transitioning to the next state, or rather, epoch.
2. It follows the byzantium fault tolerance system which is the `3f+1` rule where `f` is the number of faulty clients. Therefore, if we have 1 faulty client, we should have at least 4 nodes in the network. Furthermore, I included a special case of 3 nodes in the network where if there is 1 faulty actor, the network will still proceed given 2 safe nodes.

## Limitations
This RSM should not be treated as a production-ready implementation as it is only meant for learning purposes. Some of its limitations are:
1. Once the server is started, all nodes have to be started sequentially before the first state transition happens when 3 nodes are connected. Note that having 2 nodes will make the network stay at start state forever.
2. If a node joins the network later i.e `epoch` 10, the node cannot sync upto the latest state since nodes cannot send syncing messages yet. It will furthermore cause the network to stop transitioning if the `3f+1` rule ends up being violated by the node joining i.e. we initially had 3 nodes then 2 nodes end up joining, hence, we will need at least `3*2+1` or 7 nodes in the network for a transition to be successful.

## Overview
Our simple RSM is programmable, meaning that it can perform some computation given a message received and alter its current state. `interpreter.go` file contains all opcodes that the RSM can accept. If an opcode that is not part of the list is sent, the message returns an error and is rejected.<br />
The function `runCode()` in `client.go` file handles opcode executions. For example, typing `/msg PUSH34;PUSH12;ADD;MUL` will perform the specified operations on the stack. Contrary to normal stacks, this stack always contains the current state as the first value.
<br />
Other commands accepted include `/log` which will print out all execution messages and their timestamps in the lifetime of the program, `/state` will print out the current state which is a number, `/nodes` will print out all connected nodes with their id and remote address and `/quit` safely exits a node by deleting itself from the server node table.

## How to run
1. While in this directory, type `go mod download` to download the dependencies. If an error occurs due to `GO111MODULE` being off, turn it on using `export GO111MODULE=on` first before running `go mod download`.
2. Run `go run *.go` to start the server or rather `go build .` to generate an executable then run `./rsm`.
3. In multiple other terminals, run `telnet localhost 8888`, which will start an instance of a node.
