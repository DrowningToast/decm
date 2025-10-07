// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library TicketVCStructs {
    struct TicketVcData {
        string eventName;
        string eventDescription;
        uint256 ticketTokenId; // Onchain NFT Token ID.
        string ticketId; // References to offchain.
        string userId; // References to offchain.
        string issuerId; // References to offchain.
        address issuerAddress;
        address receiverAddress; // Participant's address.
        uint256 issuedAt;
        string encryptedUserData; // Decrypt using User PK.
        string backendEncryptedUserData; // Decrypt using Backend PK.
    }
}
