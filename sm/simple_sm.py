# Credit to Gary Explains, https://www.youtube.com/watch?v=2OiWs-h_M3A
# Simple state machine implementation of a traffic light system in python

def state_start_handler():
    print("-> START -> ", end="")
    return "RED"

def state_red_handler():
    print("RED -> ", end="")
    return "RED_AMBER"

def state_red_amber_handler():
    print("RED & AMBER -> ", end="")
    return "GREEN"

def state_green_handler():
    print("GREEN -> ", end="")
    return "AMBER"

def state_amber_handler():
    print("AMBER -> ", end="")
    return "END"


class simpleFSM:
    def __init__(self) -> None:
        self.transitions = {}
        
    def add_state(self, state, transitionFunc):
        self.transitions[state] = transitionFunc
        
    def run(self, startState):
        transitionState = self.transitions[startState]
        while True:
            newState = transitionState() # State transition function
            if newState == "END":
                print("END")
                break
            
            # Update the state transition function to reflect what transitions does the new state support
            transitionState = self.transitions[newState]


fsm = simpleFSM()
fsm.add_state("START", state_start_handler) # Start state

fsm.add_state("RED", state_red_handler)
fsm.add_state("RED_AMBER", state_red_amber_handler)
fsm.add_state("GREEN", state_green_handler)
fsm.add_state("AMBER", state_amber_handler)

# All possible state transitions added, we can now run the State Machine
fsm.run("START")
