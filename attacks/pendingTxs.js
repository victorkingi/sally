const Web3 = require('web3');
const web3 = new Web3('http://localhost:8548');
web3.eth.extend({   property: 'txpool',   methods: [{
        name: 'content',
        call: 'txpool_content'   },{
        name: 'inspect',
        call: 'txpool_inspect'   },{
        name: 'status',
        call: 'txpool_status'   }] });
let unconfirmedTxs = [];
let others = [];

async function getPool() {
    while (true) {
        unconfirmedTxs = [];
        others = [];
        const x = await web3.eth.txpool.content();
        const pending = x.pending;
        console.log(x);

        for (const tx of Object.entries(pending)) {
            for (const val of Object.entries(tx[1])) {
                if (val[1].to === "0xa7138fb2a194e312764fc1243f6b2eef4d87fc93") {
                    unconfirmedTxs.push(tx);
                    others.push(tx);
                } else {
                    others.push(tx);
                }
            }
        }
        console.log("TX pool size: Attack txs:", unconfirmedTxs.length, ", Total:", others.length);
    }
}
getPool();
