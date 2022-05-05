const ethers = require('ethers');
const crypto = require('crypto');
const fs = require("fs");
const Web3 = require('web3');
const keythereum = require("keythereum");
const web3 = new Web3('http://localhost:8547');

const accountMap = new Map();
const balanceMap = new Map();

const PASSWORD = "helloworld";
const KEYSTORE = `${__dirname}/node1_keystore.json`;
const keyObject = JSON.parse(fs.readFileSync(KEYSTORE, {encoding: "utf8"}));
const privateKey = keythereum.recover(PASSWORD, keyObject).toString("hex");
const MAIN_ADDR = '0x'+keyObject.address;

function afterGenerate() {
    fs.readFile(`${__dirname}/accounts.json`, (err, data) => {
        if (err) return err;
        let val = data.toString();
        val = JSON.parse(val);

        for (let [key, value] of Object.entries(val)) {
            accountMap.set(key, value);
        }
        // gets balance for all N addresses generated, (will always be 0 if new address)
        async function getBalance() {
            for (const pair of accountMap) {
                const [addr, privateKey] = pair;
                const raw_bal = await web3.eth.getBalance(addr);
                const bal = web3.utils.fromWei(raw_bal,'ether');
                balanceMap.set(addr, bal);
                console.log(addr, "balance is:", bal, "ETH");
            }
            console.log("Total addresses:", accountMap.size);
        }

        /**
         * Once addresses are initially generated, they are filled with 1,000,000 ETH
         * for testing
         *
         * @returns {Promise<void>}
         */
        async function firstSendEth() {
            let txNonce = 0;
            for (const pair of accountMap) {
                const [addr, p] = pair;
                const createTransaction = await web3.eth.accounts.signTransaction(
                    {
                        from: MAIN_ADDR,
                        to: addr,
                        value: web3.utils.toWei('1000000', 'ether'),
                        gas: "21000",
                        nonce: txNonce
                    },
                    privateKey
                );
                txNonce++;
                const receipt = await web3.eth.sendSignedTransaction(createTransaction.rawTransaction);
                console.log(txNonce, "address:", addr, "been credited, tx_hash:", receipt.transactionHash);
            }
        }
        //firstSendEth();
        getBalance();
    });
}
afterGenerate();

/**
 * generateAddresses will generate n random public private key pairs
 */
function generateAddresses(n) {
    for (let i = 0; i < n; i++) {
        const id = crypto.randomBytes(32).toString('hex');
        const privateKey = '0x' + id;
        const wallet = new ethers.Wallet(privateKey);
        accountMap.set(wallet.address, privateKey);
    }

    const jsonMap = Object.fromEntries(accountMap);

    fs.writeFile(`${__dirname}/accounts.json`, JSON.stringify(jsonMap), err => {
        if (err) {
            console.error(err)
            return err;
        }
        console.log("accounts saved");
        //file written successfully
    });
}
generateAddresses(50);

