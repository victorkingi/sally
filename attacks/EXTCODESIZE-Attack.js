const Web3 = require('web3');
const web3 = new Web3('http://localhost:8547');
const { performance } = require('perf_hooks');

/**
 * This code will monitor block production rate while executing the EXTCODESIZE ATTACK
 *
 */

async function printBlockNumberWithTime() {
    let currentBlock = 0;
    let startTime = performance.now();
    while (true) {
        const blockNumber = await web3.eth.getBlockNumber();
        if ((blockNumber - currentBlock) === 1) {
            const endTime = performance.now();
            const diffSeconds = (endTime - startTime) / 1000;
            const prettyPrintTimeSeconds =  Math.floor(diffSeconds * 100) / 100;
            console.log(`Block time ${prettyPrintTimeSeconds} seconds`);
            startTime = performance.now();
            currentBlock = blockNumber;
        } else if ((blockNumber - currentBlock) > 1) {
            const endTime = performance.now();
            const diffSeconds = (endTime - startTime) / (1000 * (blockNumber - currentBlock));
            const prettyPrintTimeSeconds =  Math.floor(diffSeconds * 100) / 100;
            console.log(blockNumber - currentBlock, "blocks missed", `Block time ${prettyPrintTimeSeconds} seconds`);
            startTime = performance.now();
            currentBlock = blockNumber;
        }
    }
}
printBlockNumberWithTime();
