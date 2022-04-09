#! /bin/bash
geth/geth --networkid 9984 --datadir "ethereum/node4/data" --bootnodes enode://e53a1bd8134812cdbe8383bfcbca97cc99cb13fc99c8ba5bf45830081578021fd9e6947ee50e36391294bb0e6d4f74606db38a4786f8bf738f5d14ea17ed3e37@127.0.0.1:0?discport=30301 --port 30305 --ipcdisable --syncmode full --http --allow-insecure-unlock --http.corsdomain "*" --http.port 8550 --ws --ws.origins "*" --unlock 0x09af82e4b659d52d7fc5539e5246ca04f52f9cd4 --password "ethereum/node4/password.txt" --mine console
