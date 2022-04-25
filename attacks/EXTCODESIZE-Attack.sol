// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;


contract EXTAttack {

    uint256 number;

    // provide a smart contract address and how many times to run the loop
    // I found uint times = 34499; as the maximum solidity will let me use before throwing an out of gas error
    function doAttack(address _code, uint times) public {
        uint size;
        for (uint i = 0; i < times; ++i) {
            assembly {
                size := extcodesize(_code)
            }
        }
        number = size;
    }

    /**
     * @dev Return value
     * @return value of 'number'
     */
    function retrieve() public view returns (uint256){
        return number;
    }
}
