package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hex"
	"strings"
	"sync"
)

// This state machine proccesses events, which consist of bitwise and arithmetic
// operations in batch on a number which is the state. Causing a state transition
/**
	accepted instructions:
	PUSH(x)
	ADD
	SUB
	MUL
	DIV
	MOD
	EXP
	AND
	OR
	NOT
**/

// ErrEventRejected is the error returned when the state machine cannot process
// an event in the state that it is in.
var ErrEventRejected = errors.New("event rejected")

type Hash string

type State Hash

// EventType represents an extensible event type in the state machine.
type EventType string

// EventContext represents the context to be passed to the action implementation.
type EventContext interface{}

// Action represents the action to be executed in a given state.
type Action interface {
	Execute(eventCtx EventContext) EventType
}

// An event is a struct containing, id of node sending it, code to be executed and hash of final state the node has
type Event struct {
	Hash        Hash
	Id          string
	code        int
	LatestState Hash
	Nonce       int
}

type EventLog struct {
	Events []Event
	Hash   Hash
}

// StateMachine represents the state machine.
type StateMachine struct {
	// Previous represents the previous state.
	Previous Hash

	// Current represents the current state.
	Current Hash

	// Log holds all events that have ever occured in order.
	Log EventLog

	// mutex ensures that only 1 event is processed by the state machine at any given time.
	mutex sync.Mutex
}

// getNextState returns the next state for the event given the machine's current
// state, or an error if the event can't be handled in the given state.
func (s *StateMachine) getNextState(event EventType) (Hash, error) {
	log := s.Log
	newState := string(log.Hash) + string(s.Previous)
	newState = s.getHash(newState)
	return Hash(newState), ErrEventRejected
}

func (s *StateMachine) getHash(str string) string {
	data := []byte(str)
	hash := sha256.Sum256(data)
	fmt.Printf("New State Hash: %x", hash[:])
	return hex.EncodeToString(hash[:])
}

// SendEvent sends an event to the state machine.
func (s *StateMachine) SendEvent(event EventType, eventCtx EventContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Determine the next state for the event given the machine's current state.
	nextState, err := s.getNextState(event)
	if err != nil {
			return ErrEventRejected
	}

	// Transition over to the next state.
	s.Previous = s.Current
	s.Current = nextState
	return nil
}

func newOrderFSM() *StateMachine {
	return nil
}
