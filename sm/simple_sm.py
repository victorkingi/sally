def state_start_handler(cargo):
    print("-> START -> ", end="")
    return ("RED", cargo)

def state_red_handler(cargo):
    print("RED -> ", end="")
    return ("RED_AMBER", cargo)

def state_red_amber_handler(cargo):
    print("RED & AMBER -> ", end="")
    return ("GREEN", cargo)

def state_green_handler(cargo):
    print("GREEN -> ", end="")
    return ("AMBER", cargo)

def state_amber_handler(cargo):
    print("AMBER -> ", end="")
    cargo = cargo - 1
    if cargo > 0:
        return ("RED", cargo)
    else:
        return ("END", cargo)


class simpleFSM:
    def __init__(self) -> None:
        self.handlers = {}
        
    def add_state(self, name, handler):
        self.handlers[name] = handler
        
    def run(self, startingState, cargo):
        handler = self.handlers[startingState]
        while True:
            (newState, cargo) = handler(cargo)
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
