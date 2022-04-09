#! /bin/bash

zipFile=$(ls geth);
tar -xzvf  "geth/${zipFile}" --directory geth --strip-components=1;
mkdir "ethereum";

> accounts.txt;

for i in 1 2 3 4 5
do
    mkdir "ethereum/node$i"; touch "ethereum/node$i/password.txt";
    echo "helloworld" >> "ethereum/node$i/password.txt";
    geth/geth --datadir "ethereum/node$i/data" account new;
    OUTPUT=$(ls "ethereum/node$i/data/keystore");
    cat "ethereum/node$i/data/keystore/${OUTPUT}" | jq '. | .address' | tr -d '"' >> "accounts.txt";

done
