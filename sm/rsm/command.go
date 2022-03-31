package main

type commandID int

const (
	CMD_MSG commandID = iota
	CMD_ACTIVE
	CMD_LOG
	CMD_STATE
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}