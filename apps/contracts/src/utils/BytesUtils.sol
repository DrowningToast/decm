// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;


library BytesUtils {
    function compareBytes(bytes memory a, bytes memory b) internal pure returns (bool) {
        return keccak256(a) == keccak256(b);
    }
}