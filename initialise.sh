#!/bin/bash

for i in 1 2 3 4 5
do
    geth/geth --datadir "ethereum/node$i/data" init "ethereum/year3project.json";

done

