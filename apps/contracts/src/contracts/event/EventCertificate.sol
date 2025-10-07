// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {EventAccessManager} from "./EventAccessManager.sol";
import {Event} from "./Event.sol";
import {CertificateVCStructs} from "../../libraries/CertificateVCStructs.sol";
import {Constants} from "../constants/Constants.s.sol";
import {ThemisUtils} from "../../utils/ThemisUtils.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract EventCertificate is ERC721, ThemisUtils, ReentrancyGuard {
    // Contracts
    EventAccessManager public immutable EVENT_ACCESS_MANAGER;
    Event public immutable EVENT;

    // Enums
    enum CertificateStatus {
        VALID,
        REVOKED
    }

    // State Variables
    uint256 private tokenCounter;

    // Errors
    error EventCertificate__NotHostOrAdmin();
    error EventCertificate__TokenIdOutOfBounds();
    error EventCertificate__NotHost();
    error EventCertificate__CertificateNotValid();
    error EventCertificate__NotParticipant();

    // Events
    event CertificateRevoked(uint256 indexed tokenId);
    event ParticipantSignedCertificate(uint256 indexed tokenId, address indexed receiverAddress);
    event CertificateMinted(
        uint256 indexed tokenId,
        address indexed receiverAddress,
        string certificateId,
        string userId,
        string issuerId
    );

    // Mappings
    mapping(uint256 => CertificateVCStructs.CertificateVcData)
        private tokenIdToData;
    mapping(uint256 => CertificateStatus) private tokenIdToStatus;
    mapping(uint256 => address) private tokenIdToParticipantSignedAddress;

    function requireHostOrAdmin(address signer) private view {
        if (!EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(signer)) {
            revert EventCertificate__NotHostOrAdmin();
        }
    }
    
    constructor(
        address eventAccessManagerAddr,
        address eventAddr
    )
        ERC721(
            Constants.EVENT_CERTIFICATE_NAME,
            Constants.EVENT_CERTIFICATE_SYMBOL
        )
    {
        EVENT_ACCESS_MANAGER = EventAccessManager(eventAccessManagerAddr);
        EVENT = Event(eventAddr);
        tokenCounter = 0;
    }

    function mintNft(
        address receiverAddress,
        string memory userId,
        string memory certificateId,
        string memory issuerId,
        string memory encryptedUserData,
        string memory backendEncryptedUserData,
        string memory signMessage,
        bytes memory signature
    ) external nonReentrant {
        address signer = recoverSigner(signMessage, signature);
        requireHostOrAdmin(signer);

        uint256 tokenId = tokenCounter;
        CertificateVCStructs.CertificateVcData
            memory newTokenData = _buildCertificateVcData(
                tokenId,
                receiverAddress,
                userId,
                certificateId,
                issuerId,
                encryptedUserData,
                backendEncryptedUserData
            );

        _safeMint(receiverAddress, tokenId);
        tokenIdToData[tokenId] = newTokenData;
        tokenIdToStatus[tokenId] = CertificateStatus.VALID;

        emit CertificateMinted(
            tokenId,
            receiverAddress,
            certificateId,
            userId,
            issuerId
        );

        tokenCounter++;
    }

    function participantSignedCertificate(
        uint256 tokenId
    ) external nonReentrant {
        if (tokenIdToStatus[tokenId] != CertificateStatus.VALID) {
            revert EventCertificate__CertificateNotValid();
        }

        CertificateVCStructs.CertificateVcData memory vc = tokenIdToData[tokenId];

        bool isAllowToSign = vc.receiverAddress == msg.sender;
        if (!isAllowToSign) {
            revert EventCertificate__NotParticipant();
        }

        tokenIdToParticipantSignedAddress[tokenId] = msg.sender;

        emit ParticipantSignedCertificate(tokenId, msg.sender);
    }

    struct BulkMintParticipantCertificatesParams {
        address receiverAddress;
        string userId;
        string certificateId;
        string issuerId;
        string encryptedUserData;
        string backendEncryptedUserData;
    }

    function bulkMintParticipantCertificates(
        BulkMintParticipantCertificatesParams[] memory params,
        string memory signMessage,
        bytes memory signature
    ) external nonReentrant {
        address signer = recoverSigner(signMessage, signature);
        requireHostOrAdmin(signer);

        for (uint256 i = 0; i < params.length; i++) {
            uint256 tokenId = tokenCounter++;
            tokenIdToStatus[tokenId] = CertificateStatus.VALID;

            CertificateVCStructs.CertificateVcData
                memory newTokenData = _buildCertificateVcData(
                    tokenId,
                    params[i].receiverAddress,
                    params[i].userId,
                    params[i].certificateId,
                    params[i].issuerId,
                    params[i].encryptedUserData,
                    params[i].backendEncryptedUserData
                );

            tokenIdToData[tokenId] = newTokenData;

            emit CertificateMinted(
                tokenId,
                params[i].receiverAddress,
                params[i].certificateId,
                params[i].userId,
                params[i].issuerId
            );

            _safeMint(params[i].receiverAddress, tokenId);
        }
    }

    function getTokenData(
        uint256 tokenId
    ) public view returns (string memory) {
        if (tokenId >= tokenCounter) {
            revert EventCertificate__TokenIdOutOfBounds();
        }
            
        CertificateVCStructs.CertificateVcData memory vc = tokenIdToData[
            tokenId
        ];
        CertificateStatus status = tokenIdToStatus[tokenId];
        string memory statusString = status == CertificateStatus.VALID
            ? "VALID"
            : "REVOKED";

        string memory json = string(
            abi.encodePacked(
                "{",
                    '"@context": ["https://www.w3.org/2018/credentials/v1"],',
                    '"id": "', vc.certificateId, '",',
                    '"type": ["VerifiableCredential","EventCertificate"],',
                    '"issuer": "', vc.issuerId, '",',
                    '"issuanceDate": "', vc.issuedAt, '",',
                    '"credentialSubject": {',
                        '"eventName": "', vc.eventName, '",',
                        '"eventDescription": "', vc.eventDescription, '",',
                        '"certificateTokenId": "', vc.certificateTokenId, '",',
                        '"certificateId": "', vc.certificateId, '",',
                        '"userId": "', vc.userId, '",',
                        '"issuerId": "', vc.issuerId, '",',
                        '"issuedAt": "', vc.issuedAt, '",',
                        '"issuerAddress": "', vc.issuerAddress, '",',
                        '"receiverAddress": "', vc.receiverAddress, '",',
                        '"status": "', status, '",',
                        '"encryptedUserData": "', vc.encryptedUserData, '",',
                        '"backendEncryptedUserData": "', vc.backendEncryptedUserData, '",',
                        '"status": "', statusString, '"',
                    "}",
                "}"
            )
        );

        return string(abi.encodePacked("data:application/json;utf8,", json));
    }

    function revokeCertificate(uint256 tokenId, string memory signMessage, bytes memory signature) external nonReentrant {
        address signer = recoverSigner(signMessage, signature);
        requireHostOrAdmin(signer);

        tokenIdToStatus[tokenId] = CertificateStatus.REVOKED;
        emit CertificateRevoked(tokenId);
    }


    function tokenURI(uint256 tokenId) public
        view
        override
        returns (string memory) {
        if (tokenId >= tokenCounter) {
            revert EventCertificate__TokenIdOutOfBounds();
        }
        return getTokenData(tokenId);
    }

    function _buildCertificateVcData(
        uint256 tokenId,
        address receiverAddress,
        string memory userId,
        string memory certificateId,
        string memory issuerId,
        string memory encryptedUserData,
        string memory backendEncryptedUserData
    ) private view returns (CertificateVCStructs.CertificateVcData memory) {
        return
            CertificateVCStructs.CertificateVcData({
                eventName: EVENT.getEventName(),
                eventDescription: EVENT.getEventDescription(),
                certificateTokenId: tokenId,
                certificateId: certificateId,
                userId: userId,
                issuerId: issuerId,
                issuedAt: block.timestamp,
                issuerAddress: msg.sender,
                receiverAddress: receiverAddress,
                encryptedUserData: encryptedUserData,
                backendEncryptedUserData: backendEncryptedUserData
            });
    }

}
