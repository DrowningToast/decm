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
    error Themis__InvalidDigestHash();

    function hashEthereumMessage(string memory rawMessage) public pure returns (bytes32) {
        bytes memory message = abi.encodePacked(
            "\x19Ethereum Signed Message:\n",
            bytes(rawMessage).length,
            rawMessage
        );
        
        return keccak256(message);
    }
    

    function recoverSigner(bytes32 digestHash, string memory message, bytes memory signature, address contractAddress) public view returns (address) {
        bytes32 messageHash = hashEthereumMessage(message);
        StringUtils.SignMessageStruct memory signMessage = StringUtils.toSignMessageStruct(message);


        if (digestHash != messageHash) {
            revert Themis__InvalidDigestHash();
        }

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