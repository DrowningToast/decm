// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library Constants {
    bytes32 public constant HOST_ROLE = keccak256("HOST_ROLE");
    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant PARTICIPANT_ROLE = keccak256("PARTICIPANT_ROLE");

    string public constant EVENT_TICKET_NAME = "DECM Event Ticket";
    string public constant EVENT_TICKET_SYMBOL = "DECMT";

    string public constant EVENT_CERTIFICATE_NAME = "DECM Event Certificate";
    string public constant EVENT_CERTIFICATE_SYMBOL = "DECMC";
}
