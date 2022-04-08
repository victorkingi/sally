#! /bin/bash

mkdir "ethereum";
cd  "ethereum";

for i in 1 2 3 4 5
do
    mkdir "node$i"; cd "node$i"; touch password.txt;
    echo "helloworld" >> password.txt;
    cd ..; cd ..;
    geth/geth --datadir "ethereum/node$i/data" account new;
    cd "ethereum";

done
