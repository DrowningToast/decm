// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library TicketVCStructs {
    struct TicketVcData {
        string eventName;
        string eventDescription;
        string ticketId; // References to offchain.
        string userId; // References to offchain.
        string issuerId; // References to offchain.
        address issuerAddress;
        address receiverAddress;
        uint256 issuedAt;
    }
}
