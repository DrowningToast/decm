// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

library StringUtils {
    using Strings for string;

    struct SignMessageStruct {
       address callerAddress;
       address contractAddress;
       uint256 deadlineBlock;
    }

    function toSignMessageStruct(string memory input) internal pure returns (SignMessageStruct memory) {
        (string memory part1, string memory part2, string memory part3) = splitSignMessage(input);
        return SignMessageStruct({
            callerAddress: toAddress(part1),
            contractAddress: toAddress(part2),
            deadlineBlock: toUint256(part3)
        });
    }

    function toAddress(string memory s) public pure returns (address) {
        return s.parseAddress();
    }

    function toUint256(string memory s) internal pure returns (uint256) {
        return s.parseUint(); 
    }

    function splitSignMessage(string memory input) internal pure returns (
        string memory part1,
        string memory part2,
        string memory part3
    ) {
        bytes memory inputBytes = bytes(input);
        uint256 firstComma = 0;
        uint256 secondComma = 0;
        bool foundFirst = false;
        
        for (uint256 i = 0; i < inputBytes.length; i++) {
            if (inputBytes[i] == ",") {
                if (!foundFirst) {
                    firstComma = i;
                    foundFirst = true;
                } else {
                    secondComma = i;
                    break;
                }
            }
        }

        require(foundFirst && secondComma > 0, "Invalid Sign Message");

        part1 = substring(input, 0, firstComma);
        part2 = substring(input, firstComma + 1, secondComma);
        part3 = substring(input, secondComma + 1, inputBytes.length);
    }

    function substring(string memory str, uint256 startIndex, uint256 endIndex) 
        internal pure 
        returns (string memory) 
    {
        bytes memory strBytes = bytes(str);
        bytes memory result = new bytes(endIndex - startIndex);
        
        for (uint256 i = startIndex; i < endIndex; i++) {
            result[i - startIndex] = strBytes[i];
        }
        
        return string(result);
    }
}