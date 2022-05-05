#! /bin/bash
geth/geth --networkid 9984 --datadir "ethereum/node3/data" --bootnodes enode://91ae5082e37df8407d5ed51bdb1bfdd084e07cdb2d9edeb0fac8572abb0bd7c971f1328f5bb573ae5258f62e5b410cf26aacc3d71c9016fd0a22dc9f0e815d63@127.0.0.1:0?discport=30301 --port 30304 --ipcdisable --syncmode full --allow-insecure-unlock --unlock 0xa5d0d24c0efbb343a8bde0b6051892c1cb404bf9 --password "ethereum/node3/password.txt" --mine console;
