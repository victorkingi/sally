#! /bin/bash

rm -r -f ethereum accounts.txt enode.txt;

for i in 1 2 3 4 5
do
    rm "start_node_$i.sh";

done
