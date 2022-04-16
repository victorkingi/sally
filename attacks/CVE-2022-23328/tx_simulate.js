const Web3 = require('web3');
const web3 = new Web3('http://localhost:8547');

//send a tx from an address
async function sendTx(from, to, prKey) {
    let nonce = await web3.eth.getTransactionCount(from);
    const createTransaction = await web3.eth.accounts.signTransaction(
        {
            from,
            to,
            value: web3.utils.toWei('0.05', 'ether'),
            gas: "21000",
            nonce
        },
        prKey
    );
    const receipt = await web3.eth.sendSignedTransaction(createTransaction.rawTransaction);
    console.log("nonce:", nonce, "address:", from, "tx_hash:", receipt.transactionHash);
    return receipt.transactionHash;
}

//continously send a tx given an address
async function sim(from, prKey, to) {
    while (true) {
        const hash = await sendTx(from, to, prKey);
        if (!hash) throw new Error("Function errored: "+hash);
    }
}
const from = process.argv[2];
const prKey = process.argv[3];
const to = process.argv[4];
sim(from, prKey, to);
