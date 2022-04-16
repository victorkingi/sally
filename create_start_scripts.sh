#! /bin/bash

enode=$(cat enode.txt)

# Node 5 will be the only node to run http and ws command option
addr_5=$(cat accounts.txt | sed 5!d accounts.txt)
touch "start_node_5.sh"; echo "#! /bin/bash" > "start_node_5.sh";
text="geth/geth --networkid 9984 --datadir \"ethereum/node5/data\" --bootnodes ${enode} --port 30306 --ipcdisable --syncmode full --http --allow-insecure-unlock --http.corsdomain \"*\" --http.port 8547 --http.api \"eth,net,web3,txpool\" --verbosity 4 --ws --ws.origins \"*\" --unlock 0x$addr_5 --password \"ethereum/node5/password.txt\" --mine console;"
echo ${text} >> "start_node_5.sh";
chmod ugo+x "start_node_5.sh";

for i in 1 2 3 4
do
    port=$(expr $i + 1);
    addr=$(cat accounts.txt | sed $i!d accounts.txt)
    touch "start_node_$i.sh"; echo "#! /bin/bash" > "start_node_$i.sh";

    # Node 2 will also be printing debug statements for testing purposes
    if [ $i -eq 2 ] 
    then 
        text="geth/geth --networkid 9984 --datadir \"ethereum/node$i/data\" --bootnodes ${enode} --port 3030$port --ipcdisable --verbosity 4 --syncmode full --allow-insecure-unlock --unlock 0x$addr --password \"ethereum/node$i/password.txt\" --mine console;"
        echo ${text} >> "start_node_$i.sh";
        chmod ugo+x "start_node_$i.sh"
    else
        text="geth/geth --networkid 9984 --datadir \"ethereum/node$i/data\" --bootnodes ${enode} --port 3030$port --ipcdisable --syncmode full --allow-insecure-unlock --unlock 0x$addr --password \"ethereum/node$i/password.txt\" --mine console;"
        echo ${text} >> "start_node_$i.sh";
        chmod ugo+x "start_node_$i.sh"
    fi

done
