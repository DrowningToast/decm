// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library CertificateVCStructs {
    struct Header {
        string context;
        string credentialType;
        string id;
        string issuer;
        string issuanceDate;
    }

    struct Data {
        string eventName;
        string eventDescription;
        string certificateTokenId;
        string certificateId;
        string userId;
        string issuerId;
        string issuedAt;
        address[] issuerAddresses;
        address receiverAddress;
        string encryptedUserData;
        string backendEncryptedUserData;
        string certificateTitle;
        string certificateSubtitle;
    }

    struct HostProof {
        string signature;
        string publicKey;
    }

    struct IssuerProof {
        string issuerSignature;
        string issuerPublicKey;
    }

    struct Proof {
        string encryptedByUserRawData;
        string encryptedByBackendRawData;
        string hash;
        string signMessage;
        HostProof host;
        IssuerProof[] issuers;
    }

    struct CertificateVcData {
        Header header;
        Data data;
        Proof proof;
    }
}
