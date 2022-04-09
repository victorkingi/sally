#!/bin/bash

> accounts.txt;
echo "------All nodes addresses-------------" >> "accounts.txt";

for i in 1 2 3 4 5
do
    geth/geth --datadir "ethereum/node$i/data" init "ethereum/year3project.json";
    OUTPUT=$(ls "ethereum/node$i/data/keystore");
    echo -n "node$i: 0x" >> "accounts.txt";
    cat "ethereum/node$i/data/keystore/${OUTPUT}" | jq '. | .address' | tr -d '"' >> "accounts.txt";

done

