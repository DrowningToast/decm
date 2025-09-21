// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {EventAccessManager} from "./EventAccessManager.sol";
import {Event} from "./Event.sol";
import {Constants} from "../constants/Constants.s.sol";
import {TicketVCStructs} from "../../libraries/TicketVCStructs.sol";

contract EventTicket is ERC721 {
    // Contracts
    EventAccessManager public immutable EVENT_ACCESS_MANAGER;
    Event public immutable EVENT;

    // Enums
    enum TicketStatus {
        ACTIVE,
        INACTIVE
    }

    // State Variables
    uint256 private tokenCounter;

    // Errors
    error EventTicket__NotHostOrAdmin();
    error EventTicket__TokenIdOutOfBounds();
    error EventTicket__NotHost();

    // Mappings
    mapping(uint256 => TicketVCStructs.TicketVcData) private tokenIdToUri;
    mapping(uint256 => TicketStatus) private tokenIdToStatus;

    constructor(
        address eventAccessManagerAddr,
        address eventAddr
    ) ERC721(Constants.EVENT_TICKET_NAME, Constants.EVENT_TICKET_SYMBOL) {
        EVENT_ACCESS_MANAGER = EventAccessManager(eventAccessManagerAddr);
        EVENT = Event(eventAddr);
        tokenCounter = 0;
    }

    function mintNft(
        address to,
        string memory userId, // Participant's userId (Offchain ID)
        string memory ticketId, // Ticket's id (Offchain ID)
        string memory issuerId // Issuer's id (Offchain ID)
    ) public {
        bool isAllow = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(msg.sender);
        if (!isAllow) revert EventTicket__NotHostOrAdmin();

        _safeMint(to, tokenCounter);

        tokenIdToStatus[tokenCounter] = TicketStatus.ACTIVE;
        tokenIdToUri[tokenCounter] = TicketVCStructs.TicketVcData({
            eventName: EVENT.getEventName(),
            eventDescription: EVENT.getEventDescription(),
            ticketId: ticketId,
            userId: userId,
            issuerId: issuerId,
            issuedAt: block.timestamp,
            issuerAddress: msg.sender,
            receiverAddress: to
        });

        tokenCounter++;
    }

    function tokenUri(
        uint256 tokenId
    ) public view returns (string memory) {
        if (tokenId >= tokenCounter) revert EventTicket__TokenIdOutOfBounds();

        TicketVCStructs.TicketVcData memory vc = tokenIdToUri[tokenId];
        TicketStatus status = tokenIdToStatus[tokenId];

        TicketVCStructs.TicketVcData memory ticketVcData = TicketVCStructs
            .TicketVcData({
                eventName: vc.eventName,
                eventDescription: vc.eventDescription,
                ticketId: vc.ticketId,
                userId: vc.userId,
                issuerId: vc.issuerId,
                issuedAt: vc.issuedAt,
                issuerAddress: vc.issuerAddress,
                receiverAddress: vc.receiverAddress
            });

            string memory json = string(
            abi.encodePacked(
                "{",
                    '"@context": ["https://www.w3.org/2018/credentials/v1"],',
                    '"id": "', ticketVcData.ticketId, '",',
                    '"type": ["VerifiableCredential","EventTicket"],',
                    '"issuer": "', ticketVcData.issuerId, '",',
                    '"issuanceDate": "', ticketVcData.issuedAt, '",',
                    '"credentialSubject": {',
                        '"eventName": "', ticketVcData.eventName, '",',
                        '"eventDescription": "', ticketVcData.eventDescription, '",',
                        '"ticketId": "', ticketVcData.ticketId, '",',
                        '"userId": "', ticketVcData.userId, '",',
                        '"issuerId": "', ticketVcData.issuerId, '",',
                        '"issuedAt": "', ticketVcData.issuedAt, '",',
                        '"issuerAddress": "', ticketVcData.issuerAddress, '",',
                        '"receiverAddress": "', ticketVcData.receiverAddress, '"',
                        '"status": "', status, '"',
                    "}",
                "}"
            )
        );

        return string(
            abi.encodePacked(
                "data:application/json;utf8,",
                json
            )
        );
    }
}
