#! /bin/bash

enode=$(cat enode.txt)

for i in 1 2 3 4 5
do
    port=$(expr $i + 1);
    http=$(expr 46 + $i);
    addr=$(cat accounts.txt | sed $i!d accounts.txt)
    touch "start_node_$i.sh"; echo "#! /bin/bash" > "start_node_$i.sh";
    text="geth/geth --networkid 9984 --datadir \"ethereum/node$i/data\" --bootnodes ${enode} --port 3030$port --ipcdisable --syncmode full --http --allow-insecure-unlock --http.corsdomain \"*\" --http.port 85$http --ws --ws.origins \"*\" --unlock 0x$addr --password \"ethereum/node$i/password.txt\" --mine console"
    echo ${text} >> "start_node_$i.sh";
    chmod ugo+x "start_node_$i.sh"

done
