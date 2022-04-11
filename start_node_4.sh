#! /bin/bash
geth/geth --networkid 9984 --datadir "ethereum/node4/data" --bootnodes enode://ed03da9238538a424707e8962c756001e08065c897bc4ccc1ef4b93d49909b693063b584bd398d6f6e476b1e6e90b0f28369a349cd1f3c98e5ae092ad920e969@127.0.0.1:0?discport=30301 --port 30305 --ipcdisable --syncmode full --allow-insecure-unlock --unlock 0xd9d4baedfd98ddfda9bede8911467816f7103704 --password "ethereum/node4/password.txt" --mine console
