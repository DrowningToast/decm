// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {EventAccessManager} from "./EventAccessManager.sol";
import {Event} from "./Event.sol";
import {CertificateVCStructs} from "../../libraries/CertificateVCStructs.sol";
import {Constants} from "../constants/Constants.s.sol";

contract EventCertificate is ERC721 {
    // Contracts
    EventAccessManager public immutable EVENT_ACCESS_MANAGER;
    Event public immutable EVENT;

    // Enums
    enum CertificateStatus {
        VALID,
        REVOKED
    } 

    enum CertificateType {
        DEFAULT
    }

    // State Variables
    uint256 private tokenCounter;

    // Errors
    error EventCertificate__NotHostOrAdmin();
    error EventCertificate__TokenIdOutOfBounds();
    error EventCertificate__NotHost();

    // Mappings
    mapping(uint256 => CertificateVCStructs.CertificateVcData) private tokenIdToUri;
    mapping(uint256 => CertificateStatus) private tokenIdToStatus;

    constructor(
        address eventAccessManagerAddr,
        address eventAddr
    ) ERC721(Constants.EVENT_CERTIFICATE_NAME, Constants.EVENT_CERTIFICATE_SYMBOL) {
        EVENT_ACCESS_MANAGER = EventAccessManager(eventAccessManagerAddr);
        EVENT = Event(eventAddr);
        tokenCounter = 0;
    }

    function mintNft(
        address to,
        string memory userId,
        string memory certificateId,
        string memory issuerId
    ) public {
        bool isAllow = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(msg.sender);
        if (!isAllow) revert EventCertificate__NotHostOrAdmin();

        _safeMint(to, tokenCounter);

        tokenIdToStatus[tokenCounter] = CertificateStatus.VALID;
        tokenIdToUri[tokenCounter] = CertificateVCStructs.CertificateVcData({
            eventName: EVENT.getEventName(),
            eventDescription: EVENT.getEventDescription(),
            certificateId: certificateId,
            userId: userId,
            issuerId: issuerId,
            issuedAt: block.timestamp,
            issuerAddress: msg.sender,
            receiverAddress: to
        });

        tokenCounter++;
    }

    function tokenUri(uint256 tokenId) public view returns (string memory) {
        if (tokenId >= tokenCounter) revert EventCertificate__TokenIdOutOfBounds();

        CertificateVCStructs.CertificateVcData memory vc = tokenIdToUri[tokenId];
        CertificateStatus status = tokenIdToStatus[tokenId];
        CertificateType certificateType = CertificateType.DEFAULT;
        CertificateVCStructs.CertificateVcData memory certificateVcData = CertificateVCStructs
            .CertificateVcData({
                eventName: vc.eventName,
                eventDescription: vc.eventDescription,
                certificateId: vc.certificateId,
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
                    '"id": "', certificateVcData.certificateId, '",',
                    '"type": ["VerifiableCredential","EventCertificate"],',
                    '"issuer": "', certificateVcData.issuerId, '",',
                    '"issuanceDate": "', certificateVcData.issuedAt, '",',
                    '"credentialSubject": {',
                        '"eventName": "', certificateVcData.eventName, '",',
                        '"eventDescription": "', certificateVcData.eventDescription, '",',
                        '"certificateId": "', certificateVcData.certificateId, '",',
                        '"userId": "', certificateVcData.userId, '",',
                        '"issuerId": "', certificateVcData.issuerId, '",',
                        '"issuedAt": "', certificateVcData.issuedAt, '",',
                        '"issuerAddress": "', certificateVcData.issuerAddress, '",',
                        '"receiverAddress": "', certificateVcData.receiverAddress, '"',
                        '"status": "', status, '"',
                        '"type": "', certificateType, '"',
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

    function revokeCertificate(uint256 tokenId) public {
        bool isAllow = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(msg.sender);
        if (!isAllow) revert EventCertificate__NotHostOrAdmin();
        tokenIdToStatus[tokenId] = CertificateStatus.REVOKED;
    }
}