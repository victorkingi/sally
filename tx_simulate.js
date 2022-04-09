// generate 100 random addresses, for each, first send 10000 eth to all from main address
// after send 0.5 eth to random addresses from the 100
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:8548"));

/*
const accountMap = new Map();
const balanceMap = new Map();
const keythereum = require("keythereum");
const fs = require("fs");

fs.readFile(`${__dirname}/accounts.json`, (err, data) => {
    if (err) return err;
    const val = data.toString();
    const jsObj = JSON.parse(val);
    for (const x of Object.entries(jsObj)) {
        accountMap.set(x[0], x[1]);
    }

    // gets balance for all N addresses generated, (will always be 0 if new address)
    async function getBalance() {
        for (const pair of accountMap) {
            const [addr, privateKey] = pair;
            const raw_bal = await web3.eth.getBalance(addr);
            const bal = web3.utils.fromWei(raw_bal,'ether');
            balanceMap.set(addr, bal);
        }
    }
    //getBalance().then(() => console.log(balanceMap));
});
*/

//send a tx from an address
async function sendTx(from, to, prKey, nonce) {
    const createTransaction = await web3.eth.accounts.signTransaction(
        {
            from,
            to,
            value: web3.utils.toWei('0.5', 'ether'),
            gas: "21000",
            nonce
        },
        prKey
    );
    nonce++;
    const receipt = await web3.eth.sendSignedTransaction(createTransaction.rawTransaction);
    console.log("nonce:", nonce, "address:", from, "tx_hash:", receipt.transactionHash);
}

//continously send a tx given an address
async function sim(from, prKey, to) {
    let nonce = 0;
    while (true) {
        await sendTx(from, to, prKey, nonce);
    }
}
const from = process.argv[2];
const prKey = process.argv[3];
const to = process.argv[4];
sim(from, prKey, to);
