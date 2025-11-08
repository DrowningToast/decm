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
        bool isAllowedMsgSender = EVENT_ACCESS_MANAGER.checkIsAllowedMsgSender();
        bool isHostOrAdmin = EVENT_ACCESS_MANAGER.checkIsHostOrAdmin(signer);
        if (!isHostOrAdmin && !isAllowedMsgSender) {
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
        address[] memory issuerAddresses,
        string memory signedMessageDigest,
        bytes memory signature,
        string memory hostSignature,
        string memory hostPublicKey,
        CertificateVCStructs.IssuerProof[] memory issuerProofs
    ) external nonReentrant {
        address signer = recoverSigner(signedMessageDigest, signature);
        requireHostOrAdmin(signer);

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
                issuerProofs
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
            revert EventCertificate__CertificateNotValid();
        }

        CertificateVCStructs.CertificateVcData memory vc = tokenIdToData[tokenId];

        bool isAllowToSign = vc.data.receiverAddress == msg.sender;
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
        address[] issuerAddresses;
        CertificateVCStructs.IssuerProof[] issuerProofs;
    }

    function bulkMintParticipantCertificates(
        BulkMintParticipantCertificatesParams[] memory params,
        string memory signedMessageDigest,
        bytes memory signature,
        string memory hostSignature,
        string memory hostPublicKey
    ) external nonReentrant {
        address signer = recoverSigner(signedMessageDigest, signature);
        requireHostOrAdmin(signer);

        for (uint256 i = 0; i < params.length; i++) {
            uint256 tokenId = tokenCounter++;
            tokenIdToStatus[tokenId] = CertificateStatus.VALID;

            CertificateVCStructs.CertificateVcData
                memory newTokenData = _buildCertificateVcDataWithProof(
                    tokenId,
                    params[i].receiverAddress,
                    params[i].userId,
                    params[i].certificateId,
                    params[i].issuerId,
                    params[i].encryptedUserData,
                    params[i].backendEncryptedUserData,
                    params[i].issuerAddresses,
                    block.timestamp,
                    hostSignature,
                    hostPublicKey,
                    params[i].issuerProofs
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

         emit SignatureUsed(
                msg.sender,
                signer,
                address(this),
                "bulkMintParticipantCertificates",
                signedMessageDigest,
                signature,
                block.timestamp
            );
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
                        '"status": "', statusString, '"',
                    "},",
                    '"proof": {',
                        '"encryptedByUserRawData": "', vc.proof.encryptedByUserRawData, '",',
                        '"encryptedByBackendRawData": "', vc.proof.encryptedByBackendRawData, '",',
                        '"hash": "', vc.proof.hash, '",',
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
        requireHostOrAdmin(signer);

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
        string memory backendEncryptedUserData,
        address[] memory issuerAddresses,
        uint256 issuedAt
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
                    backendEncryptedUserData: backendEncryptedUserData
                }),
                proof: CertificateVCStructs.Proof({
                    encryptedByUserRawData: "",
                    encryptedByBackendRawData: "",
                    hash: "",
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
        CertificateVCStructs.IssuerProof[] memory issuerProofs
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
                    backendEncryptedUserData: backendEncryptedUserData
                }),
                proof: CertificateVCStructs.Proof({
                    encryptedByUserRawData: "",
                    encryptedByBackendRawData: "",
                    hash: "",
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
