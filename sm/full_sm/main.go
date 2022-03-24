package main

import "fmt"
import "strings"
import "errors"
import "sync"

// ErrEventRejected is the error returned when the state machine cannot process
// an event in the state that it is in.
var ErrEventRejected = errors.New("event rejected")

const (
	// Default represents the default state of the system.
	Default StateType = ""

	// NoOp represents a no-op event.
	NoOp EventType = "NoOp"
)

// StateType represents an extensible state type in the state machine.
type StateType string

// EventType represents an extensible event type in the state machine.
type EventType string

// EventContext represents the context to be passed to the action implementation.
type EventContext interface{}

// Action represents the action to be executed in a given state.
type Action interface {
	Execute(eventCtx EventContext) EventType
}

// Events represents a mapping of events and states.
type Events map[EventType]StateType

// State binds a state with an action and a set of events it can handle.
type State struct {
	Action Action
	Events Events
}

// States represents a mapping of states and their implementations.
type States map[StateType]State

// StateMachine represents the state machine.
type StateMachine struct {
	// Previous represents the previous state.
	Previous StateType

	// Current represents the current state.
	Current StateType

	// States holds the configuration of states and events handled by the state machine.
	States States

	// mutex ensures that only 1 event is processed by the state machine at any given time.
	mutex sync.Mutex
}

// getNextState returns the next state for the event given the machine's current
// state, or an error if the event can't be handled in the given state.
func (s *StateMachine) getNextState(event EventType) (StateType, error) {
	if state, ok := s.States[s.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}
	return Default, ErrEventRejected
}

// SendEvent sends an event to the state machine.
func (s *StateMachine) SendEvent(event EventType, eventCtx EventContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {
		// Determine the next state for the event given the machine's current state.
		nextState, err := s.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}

		// Identify the state definition for the next state.
		state, ok := s.States[nextState]
		if !ok || state.Action == nil {
			// configuration error
		}

		// Transition over to the next state.
		s.Previous = s.Current
		s.Current = nextState

		// Execute the next state's action and loop over again if the event returned
		// is not a no-op.
		nextEvent := state.Action.Execute(eventCtx)
		if nextEvent == NoOp {
			return nil
		}
		event = nextEvent
	}
}

const (
	CreatingOrder     StateType = "CreatingOrder"
	OrderFailed       StateType = "OrderFailed"
	OrderPlaced       StateType = "OrderPlaced"
	ChargingCard      StateType = "ChargingCard"
	TransactionFailed StateType = "TransactionFailed"
	OrderShipped      StateType = "OrderShipped"

	CreateOrder     EventType = "CreateOrder"
	FailOrder       EventType = "FailOrder"
	PlaceOrder      EventType = "PlaceOrder"
	ChargeCard      EventType = "ChargeCard"
	FailTransaction EventType = "FailTransaction"
	ShipOrder       EventType = "ShipOrder"
)

type OrderCreationContext struct {
	items []string
	err   error
}

func (c *OrderCreationContext) String() string {
	return fmt.Sprintf("OrderCreationContext [ items: %s, err: %v ]",
		strings.Join(c.items, ","), c.err)
}

type OrderShipmentContext struct {
	cardNumber string
	address    string
	err        error
}

func (c *OrderShipmentContext) String() string {
	return fmt.Sprintf("OrderShipmentContext [ cardNumber: %s, address: %s, err: %v ]",
		c.cardNumber, c.address, c.err)
}

type CreatingOrderAction struct{}

func (a *CreatingOrderAction) Execute(eventCtx EventContext) EventType {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Validating, order:", order)
	if len(order.items) == 0 {
		order.err = errors.New("Insufficient number of items in order")
		return FailOrder
	}
	return PlaceOrder
}

type OrderFailedAction struct{}

func (a *OrderFailedAction) Execute(eventCtx EventContext) EventType {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Order failed, err:", order.err)
	return NoOp
}

type OrderPlacedAction struct{}

func (a *OrderPlacedAction) Execute(eventCtx EventContext) EventType {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Order placed, items:", order.items)
	return NoOp
}

type ChargingCardAction struct{}

func (a *ChargingCardAction) Execute(eventCtx EventContext) EventType {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Validating card, shipment:", shipment)
	if shipment.cardNumber == "" {
		shipment.err = errors.New("Card number is invalid")
		return FailTransaction
	}
	return ShipOrder
}

type TransactionFailedAction struct{}

func (a *TransactionFailedAction) Execute(eventCtx EventContext) EventType {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Transaction failed, err:", shipment.err)
	return NoOp
}

type OrderShippedAction struct{}

func (a *OrderShippedAction) Execute(eventCtx EventContext) EventType {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Order shipped, address:", shipment.address)
	return NoOp
}

func newOrderFSM() *StateMachine {
	return &StateMachine{
		States: States{
			Default: State{
				Events: Events{
					CreateOrder: CreatingOrder,
				},
			},
			CreatingOrder: State{
				Action: &CreatingOrderAction{},
				Events: Events{
					FailOrder:  OrderFailed,
					PlaceOrder: OrderPlaced,
				},
			},
			OrderFailed: State{
				Action: &OrderFailedAction{},
				Events: Events{
					CreateOrder: CreatingOrder,
				},
			},
			OrderPlaced: State{
				Action: &OrderPlacedAction{},
				Events: Events{
					ChargeCard: ChargingCard,
				},
			},
			ChargingCard: State{
				Action: &ChargingCardAction{},
				Events: Events{
					FailTransaction: TransactionFailed,
					ShipOrder:       OrderShipped,
				},
			},
			TransactionFailed: State{
				Action: &TransactionFailedAction{},
				Events: Events{
					ChargeCard: ChargingCard,
				},
			},
			OrderShipped: State{
				Action: &OrderShippedAction{},
			},
		},
	}
}
