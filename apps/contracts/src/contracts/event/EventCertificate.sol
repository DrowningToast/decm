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
    enum CertificateStatus {
        VALID,
        REVOKED
    }

    // State Variables
    uint256 private tokenCounter;

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

    function requireHostOrAdmin(address signer, address msgSender) private view {
        bool isAllowedMsgSender = EVENT_ACCESS_MANAGER.checkIsAllowedMsgSender(msgSender);
        bool isHostOrAdmin = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(signer);
        if (!isHostOrAdmin && !isAllowedMsgSender) {
            require(false, "Not host or admin or allowed msg sender");
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
        address[] memory issuerAddresses,
        string memory signedMessageDigest,
        bytes memory signature,
        string memory hostSignature,
        string memory hostPublicKey,
        string memory signMessage,
        string memory userEncryptedProof,
        string memory backendEncryptedProof,
        string memory certificateTitle,
        string memory certificateSubtitle,
        CertificateVCStructs.IssuerProof[] memory issuerProofs
    ) external nonReentrant {
        address signer = recoverSigner(signedMessageDigest, signature);
        requireHostOrAdmin(signer, msg.sender);

        uint256 tokenId = tokenCounter;
        
        CertificateVCStructs.CertificateVcData
            memory newTokenData = _buildCertificateVcDataWithProof(
                tokenId,
                receiverAddress,
                userId,
                certificateId,
                issuerId,
                encryptedUserData,
                backendEncryptedUserData,
                issuerAddresses,
                block.timestamp,
                hostSignature,
                hostPublicKey,
                signMessage,
                userEncryptedProof,
                backendEncryptedProof,
                issuerProofs,
                certificateTitle,
                certificateSubtitle
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

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "mintNft",
            signedMessageDigest,
            signature,
            block.timestamp
        );

        tokenCounter++;
    }

    function participantSignedCertificate(
        uint256 tokenId
    ) external nonReentrant {
        if (tokenIdToStatus[tokenId] != CertificateStatus.VALID) {
            require(false, "Certificate not valid");
        }

        CertificateVCStructs.CertificateVcData memory vc = tokenIdToData[tokenId];

        bool isAllowToSign = vc.data.receiverAddress == msg.sender;
        if (!isAllowToSign) {
            require(false, "Not participant");
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
        address[] issuerAddresses;
        CertificateVCStructs.IssuerProof[] issuerProofs;
        string signMessage;
        string userEncryptedProof;
        string backendEncryptedProof;
        string certificateTitle;
        string certificateSubtitle;
    }

    function getTokenData(
        uint256 tokenId
    ) public view returns (string memory) {
        if (tokenId >= tokenCounter) {
            require(false, "Token id out of bounds");
        }
            
        CertificateVCStructs.CertificateVcData memory vc = tokenIdToData[
            tokenId
        ];
        CertificateStatus status = tokenIdToStatus[tokenId];
        string memory statusString = status == CertificateStatus.VALID
            ? "VALID"
            : "REVOKED";

        // Build the VC in the new format
        string memory json = string(
            abi.encodePacked(
                "{",
                    '"header": {',
                        '"@context": ["https://www.w3.org/2018/credentials/v1"],',
                        '"type": ["VerifiableCredential","EventCertificate"],',
                        '"id": "', vc.header.id, '",',
                        '"issuer": "', vc.header.issuer, '",',
                        '"issuanceDate": "', vc.header.issuanceDate, '"',
                    "},",
                    '"data": {',
                        '"eventName": "', vc.data.eventName, '",',
                        '"eventDescription": "', vc.data.eventDescription, '",',
                        '"certificateTokenId": "', vc.data.certificateTokenId, '",',
                        '"certificateId": "', vc.data.certificateId, '",',
                        '"userId": "', vc.data.userId, '",',
                        '"issuerId": "', vc.data.issuerId, '",',
                        '"issuedAt": "', vc.data.issuedAt, '",',
                        '"issuerAddresses": "', _addressArrayToString(vc.data.issuerAddresses), '",',
                        '"receiverAddress": "', _addressToString(vc.data.receiverAddress), '",',
                        '"encryptedUserData": "', vc.data.encryptedUserData, '",',
                        '"backendEncryptedUserData": "', vc.data.backendEncryptedUserData, '",',
                        '"certificateTitle": "', vc.data.certificateTitle, '",',
                        '"certificateSubtitle": "', vc.data.certificateSubtitle, '",',
                        '"status": "', statusString, '"',
                    "},",
                    '"proof": {',
                        '"encryptedByUserRawData": "', vc.proof.encryptedByUserRawData, '",',
                        '"encryptedByBackendRawData": "', vc.proof.encryptedByBackendRawData, '",',
                        '"hash": "', vc.proof.hash, '",',
                        '"signMessage": "', vc.proof.signMessage, '",',
                        '"host": {',
                            '"signature": "', vc.proof.host.signature, '",',
                            '"publicKey": "', vc.proof.host.publicKey, '"',
                        "},",
                        '"issuers": [',
                            _issuerProofsArrayToString(vc.proof.issuers),
                        "]",
                    "}",
                "}"
            )
        );

        return string(abi.encodePacked("data:application/json;utf8,", json));
    }

    function revokeCertificate(uint256 tokenId, string memory signedMessageDigest, bytes memory signature) external nonReentrant {
        address signer = recoverSigner(signedMessageDigest, signature);
        requireHostOrAdmin(signer, msg.sender);

        tokenIdToStatus[tokenId] = CertificateStatus.REVOKED;
        emit CertificateRevoked(tokenId);

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "revokeCertificate",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }


    function tokenURI(uint256 tokenId) public
        view
        override
        returns (string memory) {
        if (tokenId >= tokenCounter) {
            require(false, "Token id out of bounds");
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
        string memory backendEncryptedUserData,
        address[] memory issuerAddresses,
        uint256 issuedAt,
        string memory certificateTitle,
        string memory certificateSubtitle
    ) private view returns (CertificateVCStructs.CertificateVcData memory) {
        // Create empty issuer proofs array for now
        CertificateVCStructs.IssuerProof[] memory emptyIssuerProofs = new CertificateVCStructs.IssuerProof[](0);
        
        return
            CertificateVCStructs.CertificateVcData({
                header: CertificateVCStructs.Header({
                    context: "https://www.w3.org/2018/credentials/v1",
                    credentialType: "VerifiableCredential,EventCertificate",
                    id: certificateId,
                    issuer: issuerId,
                    issuanceDate: _uint256ToString(issuedAt)
                }),
                data: CertificateVCStructs.Data({
                    eventName: EVENT.getEventName(),
                    eventDescription: EVENT.getEventDescription(),
                    certificateTokenId: _uint256ToString(tokenId),
                    certificateId: certificateId,
                    userId: userId,
                    issuerId: issuerId,
                    issuedAt: _uint256ToString(issuedAt),
                    issuerAddresses: issuerAddresses,
                    receiverAddress: receiverAddress,
                    encryptedUserData: encryptedUserData,
                    backendEncryptedUserData: backendEncryptedUserData,
                    certificateTitle: certificateTitle,
                    certificateSubtitle: certificateSubtitle
                }),
                proof: CertificateVCStructs.Proof({
                    encryptedByUserRawData: "",
                    encryptedByBackendRawData: "",
                    hash: "",
                    signMessage: "",
                    host: CertificateVCStructs.HostProof({
                        signature: "",
                        publicKey: ""
                    }),
                    issuers: emptyIssuerProofs
                })
            });
    }

    function _buildCertificateVcDataWithProof(
        uint256 tokenId,
        address receiverAddress,
        string memory userId,
        string memory certificateId,
        string memory issuerId,
        string memory encryptedUserData,
        string memory backendEncryptedUserData,
        address[] memory issuerAddresses,
        uint256 issuedAt,
        string memory hostSignature,
        string memory hostPublicKey,
        string memory signMessage,
        string memory userEncryptedProof,
        string memory backendEncryptedProof,
        CertificateVCStructs.IssuerProof[] memory issuerProofs,
        string memory certificateTitle,
        string memory certificateSubtitle
    ) private view returns (CertificateVCStructs.CertificateVcData memory) {
        return
            CertificateVCStructs.CertificateVcData({
                header: CertificateVCStructs.Header({
                    context: "https://www.w3.org/2018/credentials/v1",
                    credentialType: "VerifiableCredential,EventCertificate",
                    id: certificateId,
                    issuer: issuerId,
                    issuanceDate: _uint256ToString(issuedAt)
                }),
                data: CertificateVCStructs.Data({
                    eventName: EVENT.getEventName(),
                    eventDescription: EVENT.getEventDescription(),
                    certificateTokenId: _uint256ToString(tokenId),
                    certificateId: certificateId,
                    userId: userId,
                    issuerId: issuerId,
                    issuedAt: _uint256ToString(issuedAt),
                    issuerAddresses: issuerAddresses,
                    receiverAddress: receiverAddress,
                    encryptedUserData: encryptedUserData,
                    backendEncryptedUserData: backendEncryptedUserData,
                    certificateTitle: certificateTitle,
                    certificateSubtitle: certificateSubtitle
                }),
                proof: CertificateVCStructs.Proof({
                    encryptedByUserRawData: userEncryptedProof,
                    encryptedByBackendRawData: backendEncryptedProof,
                    hash: "",
                    signMessage: signMessage,
                    host: CertificateVCStructs.HostProof({
                        signature: hostSignature,
                        publicKey: hostPublicKey
                    }),
                    issuers: issuerProofs
                })
            });
    }

    function _addressToString(address addr) private pure returns (string memory) {
        bytes memory data = abi.encodePacked(addr);
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < data.length; i++) {
            str[2 + i * 2] = alphabet[uint256(uint8(data[i] >> 4))];
            str[3 + i * 2] = alphabet[uint256(uint8(data[i] & 0x0f))];
        }
        return string(str);
    }

    function _addressArrayToString(address[] memory addresses) private pure returns (string memory) {
        if (addresses.length == 0) {
            return "[]";
        }

        string memory result = "[";
        for (uint256 i = 0; i < addresses.length; i++) {
            result = string(abi.encodePacked(result, "\"", _addressToString(addresses[i]), "\""));
            if (i < addresses.length - 1) {
                result = string(abi.encodePacked(result, ","));
            }
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    function _issuerProofsArrayToString(CertificateVCStructs.IssuerProof[] memory issuers) private pure returns (string memory) {
        if (issuers.length == 0) {
            return "[]";
        }

        string memory result = "";
        for (uint256 i = 0; i < issuers.length; i++) {
            result = string(abi.encodePacked(
                result,
                "{",
                    '"issuerSignature": "', issuers[i].issuerSignature, '",',
                    '"issuerPublicKey": "', issuers[i].issuerPublicKey, '"',
                "}"
            ));
            if (i < issuers.length - 1) {
                result = string(abi.encodePacked(result, ","));
            }
        }
        return result;
    }

    function _uint256ToString(uint256 value) private pure returns (string memory) {
        if (value == 0) {
            return "0";
        }
        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }
        return string(buffer);
    }
}
