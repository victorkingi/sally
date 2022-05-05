const fs = require('fs');
const Web3 = require('web3');
const web3 = new Web3('http://localhost:8547');


async function getBalance() {
    const raw_bal = await web3.eth.getBalance("0x2dd4aea78a11ab6efce6d7bdfdd5e2a82e9a09d9");
    const bal = web3.utils.fromWei(raw_bal,'ether');
    console.log("balance is:", bal, "ETH");
}
getBalance();
