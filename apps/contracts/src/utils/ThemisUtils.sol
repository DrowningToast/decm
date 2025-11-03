// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract ThemisUtils {
    error Themis__InvalidSignature();

    function recoverSigner(string memory message, bytes memory signature) public pure returns (address) {
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        return recoverSigner(messageHash, signature);
    }

    function recoverSigner(bytes32 messageHash, bytes memory signature) public pure returns (address) {
        bytes32 ethSignedMessageHash = MessageHashUtils.toEthSignedMessageHash(messageHash);
        address signer = ECDSA.recover(ethSignedMessageHash, signature);
        
        if (signer == address(0)) {
            revert Themis__InvalidSignature();
        }
        
        return signer;
    }
}