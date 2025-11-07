// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";  
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {EventAccessManager} from "./EventAccessManager.sol";
import {Event} from "./Event.sol";
import {Constants} from "../constants/Constants.s.sol";
import {TicketVCStructs} from "../../libraries/TicketVCStructs.sol";
import {ThemisUtils} from "../../utils/ThemisUtils.sol";


contract EventTicket is ERC721, ThemisUtils, ReentrancyGuard {
    // Contracts
    EventAccessManager public immutable EVENT_ACCESS_MANAGER;
    Event public immutable EVENT;

    // Events
     event SignatureUsed(
        address indexed transactor,
        address indexed signer,
        address contractAddress,
        string functionName,
        string signedMessageDigest,
        bytes signature,
        uint256 timestamp
    );

    // Enums
    enum TicketStatus {
        ACTIVE,
        INACTIVE
    }

    // State Variables
    uint256 private tokenCounter;

    // Getter for tokenCounter (for testing)
    function getTokenCounter() external view returns (uint256) {
        return tokenCounter;
    }

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
    mapping(uint256 => TicketVCStructs.TicketVcData) private tokenIdToVcData;
    mapping(uint256 => TicketStatus) private tokenIdToStatus;

    function requireHostOrAdmin(address signer) private view {
        bool isAllowedMsgSender = EVENT_ACCESS_MANAGER.checkIsAllowedMsgSender();
        bool isHostOrAdmin = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(signer);
        if (!isHostOrAdmin && !isAllowedMsgSender) {
            revert EventTicket__NotHostOrAdmin();
        }
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
        address receiverAddress,
        string memory userId,
        string memory ticketId,
        string memory issuerId,
        string memory encryptedUserData,
        string memory backendEncryptedUserData,
        address issuerAddress,
        string memory signedMessageDigest,
        bytes memory signature
    ) external nonReentrant {
        address signer = recoverSigner(signedMessageDigest, signature);
        requireHostOrAdmin(signer);

        if (receiverAddress == address(0)) {
            revert EventTicket__InvalidReceiver();
        }

        uint256 tokenId = tokenCounter;
        tokenIdToStatus[tokenId] = TicketStatus.ACTIVE;

        TicketVCStructs.TicketVcData memory newTicketVcData = _buildTicketVcData(
            tokenId,
            receiverAddress,
            userId,
            ticketId,
            issuerId,
            encryptedUserData,
            backendEncryptedUserData,
            issuerAddress,
            block.timestamp
        );

        tokenIdToVcData[tokenId] = newTicketVcData;

        _safeMint(receiverAddress, tokenId);
        tokenCounter++;

        emit TicketMinted(tokenId, signer, receiverAddress, ticketId);

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "mintNft",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }

    struct BulkMintParticipantTicketsParams {
        address receiverAddress;
        string userId;
        string ticketId;
        string issuerId;
        string encryptedUserData;
        string backendEncryptedUserData;
        address issuerAddress; // Host Address
    }

    function bulkMintParticipantTickets(
        BulkMintParticipantTicketsParams[] memory params,
        string memory signedMessageDigest,
        bytes memory signature
    ) external nonReentrant {
        address signer = recoverSigner(signedMessageDigest, signature);
        requireHostOrAdmin(signer);

        for (uint256 i = 0; i < params.length; i++) {
            if (params[i].receiverAddress == address(0)) {
                revert EventTicket__InvalidReceiver();
            }
            
            uint256 tokenId = tokenCounter;

            TicketVCStructs.TicketVcData memory newTicketVcData = _buildTicketVcData(
                tokenId,
                params[i].receiverAddress,
                params[i].userId,
                params[i].ticketId,
                params[i].issuerId,
                params[i].encryptedUserData,
                params[i].backendEncryptedUserData,
                params[i].issuerAddress,
                block.timestamp
            );

            _safeMint(params[i].receiverAddress, tokenId);

            tokenIdToVcData[tokenId] = newTicketVcData;
            tokenIdToStatus[tokenId] = TicketStatus.ACTIVE;

            tokenCounter++;

            emit TicketMinted(tokenId, signer, params[i].receiverAddress, params[i].ticketId);
        }

        emit SignatureUsed( 
            msg.sender,
            signer,
            address(this),
            "bulkMintParticipantTickets",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }

    function getTokenData(
        uint256 tokenId
    ) public view returns (string memory) {
        if (tokenId >= tokenCounter) {
             revert EventTicket__TokenIdOutOfBounds();
        }

        TicketVCStructs.TicketVcData memory vc = tokenIdToVcData[tokenId];
        TicketStatus status = tokenIdToStatus[tokenId];

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
                        '"ticketTokenId": "', vc.ticketTokenId, '",',
                        '"ticketId": "', vc.ticketId, '",',
                        '"userId": "', vc.userId, '",',
                        '"issuerId": "', vc.issuerId, '",',
                        '"issuedAt": "', vc.issuedAt, '",',
                        '"issuerAddress": "', vc.issuerAddress, '",',
                        '"receiverAddress": "', vc.receiverAddress, '",',
                        '"encryptedUserData": "', vc.encryptedUserData, '",',
                        '"backendEncryptedUserData": "', vc.backendEncryptedUserData, '",',
                        '"status": "', statusString, '"',
                    "}",
                "}"
            )
        );

        return string(abi.encodePacked("data:application/json;utf8,", json));
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        if (tokenId >= tokenCounter) {
            revert EventTicket__TokenIdOutOfBounds();
        }

        return getTokenData(tokenId);
    }

    function _buildTicketVcData(
        uint256 tokenId,
        address receiverAddress,
        string memory userId,
        string memory ticketId,
        string memory issuerId,
        string memory encryptedUserData,
        string memory backendEncryptedUserData,
        address issuerAddress,
        uint256 issuedAt
    ) private view returns (TicketVCStructs.TicketVcData memory) {
        return TicketVCStructs.TicketVcData({
            eventName: EVENT.getEventName(),
            eventDescription: EVENT.getEventDescription(),
            ticketTokenId: tokenId,
            ticketId: ticketId,
            userId: userId,
            issuerId: issuerId,
            issuedAt: issuedAt,
            issuerAddress: issuerAddress,
            receiverAddress: receiverAddress,
            encryptedUserData: encryptedUserData,
            backendEncryptedUserData: backendEncryptedUserData
        });
    }
}
