// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library Constants {
    enum Role {
        ADMIN,
        HOST,
        ISSUER,
        PARTICIPANT
    }

    bytes32 public constant ADMIN_ROLE =
        keccak256(abi.encodePacked(uint256(Role.ADMIN)));
    bytes32 public constant HOST_ROLE =
        keccak256(abi.encodePacked(uint256(Role.HOST)));
    bytes32 public constant ISSUER_ROLE =
        keccak256(abi.encodePacked(uint256(Role.ISSUER)));
    bytes32 public constant PARTICIPANT_ROLE =
        keccak256(abi.encodePacked(uint256(Role.PARTICIPANT)));

    string public constant EVENT_TICKET_NAME = "DECM Event Ticket";
    string public constant EVENT_TICKET_SYMBOL = "DECMT";

    string public constant EVENT_CERTIFICATE_NAME = "DECM Event Certificate";
    string public constant EVENT_CERTIFICATE_SYMBOL = "DECMC";
}
