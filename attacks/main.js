const ethers = require("ethers");
const url = "ws://127.0.0.1:8546";

const init = function () {
    const customWsProvider = new ethers.providers.WebSocketProvider(url);

    customWsProvider.on("pending", (tx) => {
        customWsProvider.getTransaction(tx).then(function (transaction) {
            console.log(transaction);
            if (transaction.gasPrice._isBigNumber) {
                console.log("GAS PRICE:", transaction.gasPrice.toBigInt().toString());
            } else {
                console.log("NOT BIG GAS PRICE:", transaction.gasPrice._hex);
            }
            if (transaction.value._isBigNumber) {
                console.log("VALUE:", ethers.BigNumber.from(transaction.value._hex).toBigInt());
            } else {
                console.log("NOT BIG VALUE:", transaction.value);
            }
        });
    });

    customWsProvider._websocket.on("error", async () => {
        console.log(`Unable to connect to  retrying in 3s...`);
        setTimeout(init, 3000);
    });
    customWsProvider._websocket.on("close", async (code) => {
        console.log(
            `Connection lost with code ${code}! Attempting reconnect in 3s...`
        );
        customWsProvider._websocket.terminate();
        setTimeout(init, 3000);
    });
};

init();
