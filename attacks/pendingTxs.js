const Web3 = require('web3');
const web3 = new Web3('http://localhost:8547');
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
//getPool();


//web3.eth.getTransaction("0xa80890bfb40e3191445b4e8a00feefd2bb7e2fe1c9dc27a41536fa2b489f9f8e").then(console.log);
//web3.eth.getBlockTransactionCount('0x0fe36c2090a2913017aa8eecc65f899e751506a6b19756579388b56a7fac61d3').then(console.log);
