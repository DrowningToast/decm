// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

import {StringUtils} from "./StringUtils.sol";

contract ThemisUtils {
    error Themis__InvalidSignature();
    error Themis__SignatureExpired();
    error Themis__InvalidCaller();
    error Themis__InvalidContract();
    error Themis__InvalidSigner();

    function recoverSigner(string memory message, bytes memory signature, address contractAddress) public view returns (address) {
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        StringUtils.SignMessageStruct memory signMessage = StringUtils.toSignMessageStruct(message);

        // 1. Check deadline block
        if (block.number > signMessage.deadlineBlock) {
            revert Themis__SignatureExpired();
        }

        // 2. Check is valid contract
        if (signMessage.contractAddress != contractAddress) {
            revert Themis__InvalidContract();
        }

        address recoveredSigner = recoverSigner(messageHash, signature);

        // 3. Check is valid signer
        if (recoveredSigner != signMessage.signerAddress) {
            revert Themis__InvalidSigner();
        }

        return recoveredSigner;
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