# A finite state machine that checks if a floating point number is valid
# states having * at the end denote valid end states i.e. MANTISSA*

from collections import namedtuple

tok_next_t = namedtuple('tok_next_t', ['tokens', 'next_state'])


class complexFSM:
    def __init__(self) -> None:
        self.tok_next = {}
        self.state = "START"
        
    def add_state(self, name, tokens, next_state):
        # tok_next is a list of tok_next_t namedtuples
        # i.e. a list of the valid tokens and transition states
        
        if name not in self.tok_next:
            self.tok_next[name] = []
        self.tok_next[name].append(tok_next_t(tokens, next_state))
    
        
    def run(self, cargo):
        print("->", self.state, "-> ", end="")
        while True:
            if len(cargo) > 0:
                # Get next char & remove it from string
                token = cargo[0]
                cargo = cargo[1:]
                found = False
                for tn in self.tok_next[self.state]:
                    if token in tn.tokens:
                        # Token is valid, jump to next state
                        self.state = tn.next_state
                        found = True
                        print(self.state, "-> ", end="")
                        break
                
                if not found:
                    print("INVALID INPUT")
                    break
                
            else:
                if self.state.endswith('*'):
                    print()
                    break
                else:
                    print("INVALID END STATE")
                    break
                



fsm = complexFSM()
fsm.add_state("START", "0123456789", "WHOLE_NUMBER*")
fsm.add_state("START", "-", "SIGN")
fsm.add_state("SIGN", "0123456789",  "WHOLE_NUMBER*")
fsm.add_state("WHOLE_NUMBER*", "0123456789", "WHOLE_NUMBER*")
fsm.add_state("WHOLE_NUMBER*", ".", "POINT")
fsm.add_state("POINT", "0123456789", "MANTISSA*")
fsm.add_state("MANTISSA*", "0123456789", "MANTISSA*")

fsm.run("-7.0")
