# A state machine that checks if a floating point number is valid
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
        print(self.state, "-> ", end="")
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
                        if len(cargo) > 0:
                            print(self.state, "-> ", end="")
                        else:
                            print(self.state)
                        break
                
                if not found:
                    print(self.state)
                    break
                
            else:
                break
                



fsm = complexFSM()
fsm.add_state("START", "0123456789", "SECOND_DIGIT_ONWARDS*")
fsm.add_state("START", "-", "AFTER_MINUS")
fsm.add_state("AFTER_MINUS", "0123456789",  "SECOND_DIGIT_ONWARDS*")
fsm.add_state("SECOND_DIGIT_ONWARDS*", "0123456789", "SECOND_DIGIT_ONWARDS*")
fsm.add_state("SECOND_DIGIT_ONWARDS*", ".", "AFTER_DOT")
fsm.add_state("AFTER_DOT", "0123456789", "MANTISSA*")
fsm.add_state("MANTISSA*", "0123456789", "MANTISSA*")

fsm.run("-22.a0")
