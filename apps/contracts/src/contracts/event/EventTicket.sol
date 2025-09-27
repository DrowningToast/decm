// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {EventAccessManager} from "./EventAccessManager.sol";
import {Event} from "./Event.sol";
import {Constants} from "../constants/Constants.s.sol";
import {TicketVCStructs} from "../../libraries/TicketVCStructs.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

contract EventTicket is ERC721 {
    using Strings for uint256;
    using Strings for address;

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
    error EventTicket__AccessManagerCannotBeZeroAddress();
    error EventTicket__EventAddressCannotBeZeroAddress();
    error EventTicket__InvalidReceiver();

    // Events
    event TicketMinted(
        uint256 indexed tokenId,
        address indexed issuer,
        address indexed receiver,
        string ticketId
    );

    event TicketStatusUpdated(uint256 indexed tokenId, TicketStatus status);

    // Mappings
    mapping(uint256 => TicketVCStructs.TicketVcData) private ticketVcData;
    mapping(uint256 => TicketStatus) private ticketStatus;

    modifier onlyHostOrAdmin() {
        if (!EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(msg.sender)) {
            revert EventTicket__NotHostOrAdmin();
        }
        _;
    }

    constructor(
        address eventAccessManagerAddr,
        address eventAddr
    ) ERC721(Constants.EVENT_TICKET_NAME, Constants.EVENT_TICKET_SYMBOL) {
        if (eventAccessManagerAddr == address(0)) {
            revert EventTicket__AccessManagerCannotBeZeroAddress();
        }
        if (eventAddr == address(0)) {
            revert EventTicket__EventAddressCannotBeZeroAddress();
        }

        EVENT_ACCESS_MANAGER = EventAccessManager(eventAccessManagerAddr);
        EVENT = Event(eventAddr);
        tokenCounter = 0;
    }

    function mintNft(
        address to,
        string memory userId, // Participant's userId (Offchain ID)
        string memory ticketId, // Ticket's id (Offchain ID)
        string memory issuerId // Issuer's id (Offchain ID)
    ) public onlyHostOrAdmin {
        if (to == address(0)) {
            revert EventTicket__InvalidReceiver();
        }

        uint256 tokenId = tokenCounter;
        _safeMint(to, tokenId);

        ticketStatus[tokenId] = TicketStatus.ACTIVE;
        ticketVcData[tokenId] = TicketVCStructs.TicketVcData({
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

        emit TicketMinted(tokenId, msg.sender, to, ticketId);
    }

    function setTicketStatus(
        uint256 tokenId,
        TicketStatus newStatus
    ) external onlyHostOrAdmin {
        _requireExistingToken(tokenId);
        ticketStatus[tokenId] = newStatus;
        emit TicketStatusUpdated(tokenId, newStatus);
    }

    function tokenURI(uint256 tokenId)
        public
        view
        override
        returns (string memory)
    {
        TicketVCStructs.TicketVcData memory vc = _requireExistingToken(tokenId);

        TicketStatus status = ticketStatus[tokenId];
        string memory statusString = status == TicketStatus.ACTIVE
            ? "ACTIVE"
            : "INACTIVE";

            string memory json = string(
            abi.encodePacked(
                "{",
                    '"@context": ["https://www.w3.org/2018/credentials/v1"],',
                    '"id": "', vc.ticketId, '",',
                    '"type": ["VerifiableCredential","EventTicket"],',
                    '"issuer": "', vc.issuerId, '",',
                    '"issuanceDate": "', vc.issuedAt, '",',
                    '"credentialSubject": {',
                        '"eventName": "', vc.eventName, '",',
                        '"eventDescription": "',vc.eventDescription, '",',
                        '"ticketId": "', vc.ticketId, '",',
                        '"userId": "', vc.userId, '",',
                        '"issuerId": "', vc.issuerId, '",',
                        '"issuedAt": "', vc.issuedAt, '",',
                        '"issuerAddress": "', vc.issuerAddress, '",',
                        '"receiverAddress": "', vc.receiverAddress, '"',
                        '"status": "', statusString, '"',
                    "}",
                "}"
            )
        );

        return string(
            abi.encodePacked("data:application/json;utf8,", json)
        );
    }

    function getTicketStatus(uint256 tokenId)
        external
        view
        returns (TicketStatus)
    {
        _requireExistingToken(tokenId);
        return ticketStatus[tokenId];
    }

    function getTicketData(uint256 tokenId)
        external
        view
        returns (TicketVCStructs.TicketVcData memory)
    {
        return _requireExistingToken(tokenId);
    }

    function totalMinted() external view returns (uint256) {
        return tokenCounter;
    }

    function _requireExistingToken(uint256 tokenId)
        internal
        view
        returns (TicketVCStructs.TicketVcData memory)
    {
        if (tokenId >= tokenCounter) {
            revert EventTicket__TokenIdOutOfBounds();
        }
        return ticketVcData[tokenId];
    }
}
