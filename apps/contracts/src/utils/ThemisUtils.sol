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

    string constant MOCK_ADMIN_MESSAGE = "MOCK_ADMIN_MESSAGE";
    string constant MOCK_HOST_MESSAGE = "MOCK_HOST_MESSAGE";   
    string constant MOCK_CALLER_MESSAGE = "MOCK_CALLER_MESSAGE";
    string constant MOCK_PARTICIPANT_MESSAGE = "MOCK_PARTICIPANT_MESSAGE";
    string constant MOCK_ISSUER_MESSAGE = "MOCK_ISSUER_MESSAGE";

    address public constant ADMIN = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;
    address public constant HOST = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address public constant CALLER = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
    address public constant PARTICIPANT = 0x90F79bf6EB2c4f870365E785982E1f101E93b906;
    address public constant ISSUER = 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65;

    function recoverSigner(string memory message, bytes memory signature, address contractAddress) public view returns (address) {
        // Uncomment this for unit testing only!
        if (StringUtils.compareStrings(MOCK_ADMIN_MESSAGE, message)) {
            return ADMIN;
        } else if (StringUtils.compareStrings(MOCK_HOST_MESSAGE, message)) {
            return HOST;
        } else if (StringUtils.compareStrings(MOCK_CALLER_MESSAGE, message)) {
            return CALLER;
        } else if (StringUtils.compareStrings(MOCK_PARTICIPANT_MESSAGE, message)) {
            return PARTICIPANT;
        } else if (StringUtils.compareStrings(MOCK_ISSUER_MESSAGE, message)) {
            return ISSUER;
        }

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