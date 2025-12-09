// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract ThemisUtils {
    
    mapping(bytes => bool) public usedSignatures;

    /**
     * @dev Recovers the signer address from a message digest and signature
     * @param signedMessageDigest The message digest that was signed (as a string)
     * @param signature The signature bytes
     * @return The address of the signer
     */
    function recoverSigner(string memory signedMessageDigest, bytes memory signature) public returns (address) {
        // Hash the message with the Ethereum signed message prefix
        bytes32 ethSignedMessageHash = MessageHashUtils.toEthSignedMessageHash(bytes(signedMessageDigest));
        
        // Recover the signer address
        address signer = ECDSA.recover(ethSignedMessageHash, signature);
        
        if (signer == address(0)) {
            require(false, "Invalid signature");
        }
        
        if (usedSignatures[signature]) {
            require(false, "Signature already used");
        }

        usedSignatures[signature] = true;

        return signer;
    }
}