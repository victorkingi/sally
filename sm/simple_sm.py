# Simple state machine implementation of how a traffic light system works

def state_start_handler(iterations_to_run):
    print("-> START -> ", end="")
    return ("RED", iterations_to_run)

def state_red_handler(iterations_to_run):
    print("RED -> ", end="")
    return ("RED_AMBER", iterations_to_run)

def state_red_amber_handler(iterations_to_run):
    print("RED & AMBER -> ", end="")
    return ("GREEN", iterations_to_run)

def state_green_handler(iterations_to_run):
    print("GREEN -> ", end="")
    return ("AMBER", iterations_to_run)

def state_amber_handler(iterations_to_run):
    print("AMBER -> ", end="")
    iterations_to_run = iterations_to_run - 1
    if iterations_to_run > 0:
        return ("RED", iterations_to_run)
    else:
        return ("END", iterations_to_run)


class simpleFSM:
    def __init__(self) -> None:
        self.handlers = {}
        
    def add_state(self, name, handler):
        self.handlers[name] = handler
        
    def run(self, startingState, iterations_to_run):
        handler = self.handlers[startingState]
        while True:
            (newState, iterations_to_run) = handler(iterations_to_run)
            if newState == "END":
                print("END")
                break
            handler = self.handlers[newState]


fsm = simpleFSM()
fsm.add_state("START", state_start_handler)

fsm.add_state("RED", state_red_handler)
fsm.add_state("RED_AMBER", state_red_amber_handler)
fsm.add_state("GREEN", state_green_handler)
fsm.add_state("AMBER", state_amber_handler)

iterations_to_run = 2

fsm.run("START", iterations_to_run)
