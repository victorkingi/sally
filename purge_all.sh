#! /bin/bash

[[ -d ethereum ]] && rm -r -f ethereum;
[[ -f accounts.txt ]] && [[ -f enode.txt ]] && rm accounts.txt enode.txt;

for i in 1 2 3 4 5
do
    [[ -f "start_node_$i.sh" ]] && rm "start_node_$i.sh";

done
