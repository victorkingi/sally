#! /bin/bash

geth/bootnode -nodekey "ethereum/bnode/boot.key" -verbosity 7 -addr "127.0.0.1:30301" | head -n 1 >> "enode.txt"
