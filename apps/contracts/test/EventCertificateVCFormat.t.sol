// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EventCertificate} from "../src/contracts/event/EventCertificate.sol";
import {TestEventAccessManager} from "./utils/TestEventAccessManager.sol";
import {MockDecmAccessManager} from "./utils/MockDecmAccessManager.sol";
import {MockEvent} from "./utils/MockEvent.sol";
import {TestUtils} from "./utils/TestUtils.sol";
import {Event} from "../src/contracts/event/Event.sol";
import {CertificateVCStructs} from "../src/libraries/CertificateVCStructs.sol";

contract EventCertificateVCFormatTest is TestUtils {
    EventCertificate public eventCertificate;
    TestEventAccessManager public eventAccessManager;
    MockDecmAccessManager public mockDecmAccessManager;
    MockEvent public mockEvent;

    // Test constants
    string constant TEST_EVENT_NAME = "Test Event";
    string constant TEST_EVENT_DESCRIPTION = "This is a test event";
    uint256 constant TEST_SEATS_COUNT = 100;
    
    string constant TEST_USER_ID = "user123";
    string constant TEST_CERTIFICATE_ID = "certificate123";
    string constant TEST_ISSUER_ID = "issuer123";
    string constant TEST_ENCRYPTED_USER_DATA = "encryptedUserData123";
    string constant TEST_BACKEND_ENCRYPTED_USER_DATA = "backendEncryptedUserData123";
    
    // Test issuer addresses array
    address[] testIssuerAddresses;
    
    // Test proof data
    string constant TEST_HOST_SIGNATURE = "hostSignature123";
    string constant TEST_HOST_PUBLIC_KEY = "hostPublicKey123";
    string constant TEST_SIGN_MESSAGE = "signMessage123";
    string constant TEST_USER_ENCRYPTED_PROOF = "userEncryptedProof123";
    string constant TEST_BACKEND_ENCRYPTED_PROOF = "backendEncryptedProof123";
    
    // Test issuer proofs
    CertificateVCStructs.IssuerProof[] testIssuerProofs;

    function setUp() public {
        setupMockAddresses();
        
        // Initialize test issuer addresses
        testIssuerAddresses = new address[](2);
        testIssuerAddresses[0] = ISSUER;
        testIssuerAddresses[1] = makeAddr("secondIssuer");
        
        // Initialize test issuer proofs
        testIssuerProofs = new CertificateVCStructs.IssuerProof[](2);
        testIssuerProofs[0] = CertificateVCStructs.IssuerProof({
            issuerSignature: "issuerSignature1",
            issuerPublicKey: "issuerPublicKey1"
        });
        testIssuerProofs[1] = CertificateVCStructs.IssuerProof({
            issuerSignature: "issuerSignature2",
            issuerPublicKey: "issuerPublicKey2"
        });
        
        // Deploy MockDecmAccessManager with admin
        address[] memory admins = new address[](1);
        admins[0] = ADMIN;
        vm.prank(ADMIN);
        mockDecmAccessManager = new MockDecmAccessManager(admins);
        
        // Deploy TestEventAccessManager with host
        vm.prank(ADMIN);
        eventAccessManager = new TestEventAccessManager(address(mockDecmAccessManager), HOST);
        
        // Deploy MockEvent contract
        mockEvent = new MockEvent(
            TEST_EVENT_NAME,
            TEST_EVENT_DESCRIPTION,
            TEST_SEATS_COUNT,
            Event.EventStatus.ACTIVE
        );
        
        // Deploy EventCertificate contract
        vm.prank(ADMIN);
        eventCertificate = new EventCertificate(
            address(eventAccessManager),
            address(mockEvent)
        );
    }

    /*//////////////////////////////////////////////////////////////
                       HELPER FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    function _contains(string memory str, string memory substr) internal pure returns (bool) {
        bytes memory strBytes = bytes(str);
        bytes memory substrBytes = bytes(substr);
        
        if (substrBytes.length > strBytes.length) {
            return false;
        }
        
        for (uint256 i = 0; i <= strBytes.length - substrBytes.length; i++) {
            bool found = true;
            for (uint256 j = 0; j < substrBytes.length; j++) {
                if (strBytes[i + j] != substrBytes[j]) {
                    found = false;
                    break;
                }
            }
            if (found) {
                return true;
            }
        }
        
        return false;
    }

    function _mintCertificateWithProof(address receiver) internal {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            receiver,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage,
            signature,
            TEST_HOST_SIGNATURE,
            TEST_HOST_PUBLIC_KEY,
            TEST_SIGN_MESSAGE,
            TEST_USER_ENCRYPTED_PROOF,
            TEST_BACKEND_ENCRYPTED_PROOF,
            testIssuerProofs
        );
    }

    function _extractJsonSection(string memory json, string memory key) internal pure returns (string memory) {
        // Simple extraction of JSON section by key
        // This is a simplified approach for testing purposes
        bytes memory jsonBytes = bytes(json);
        bytes memory keyBytes = bytes(key);
        bytes memory result = new bytes(jsonBytes.length);
        
        uint256 resultIndex = 0;
        bool foundKey = false;
        uint256 i = 0;
        
        // Find the key in the JSON
        while (i < jsonBytes.length - keyBytes.length - 2) {
            bool isMatch = true;
            
            // Check for key with quotes
            if (jsonBytes[i] == '"' && jsonBytes[i + keyBytes.length + 1] == '"') {
                for (uint256 j = 0; j < keyBytes.length; j++) {
                    if (jsonBytes[i + 1 + j] != keyBytes[j]) {
                        isMatch = false;
                        break;
                    }
                }
                
                if (isMatch) {
                    foundKey = true;
                    i += keyBytes.length + 3; // Skip key and ":"
                    break;
                }
            }
            i++;
        }
        
        if (!foundKey) return "";
        
        // Extract the value (simplified - assumes no nested objects with the same key)
        uint256 braceCount = 0;
        bool inString = false;
        bool escapeNext = false;
        
        while (i < jsonBytes.length) {
            bytes1 char = jsonBytes[i];
            
            if (escapeNext) {
                escapeNext = false;
            } else if (char == '\\') {
                escapeNext = true;
            } else if (char == '"') {
                inString = !inString;
            } else if (!inString) {
                if (char == '{') {
                    braceCount++;
                } else if (char == '}') {
                    if (braceCount == 0) {
                        // End of object
                        result[resultIndex] = char;
                        resultIndex++;
                        break;
                    }
                    braceCount--;
                } else if (char == ',' && braceCount == 0) {
                    // End of value
                    break;
                }
            }
            
            result[resultIndex] = char;
            resultIndex++;
            i++;
        }
        
        // Trim the result
        bytes memory trimmedResult = new bytes(resultIndex);
        for (uint256 j = 0; j < resultIndex; j++) {
            trimmedResult[j] = result[j];
        }
        
        return string(trimmedResult);
    }

    /*//////////////////////////////////////////////////////////////
                    VC FORMAT VALIDATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_VCFormat_ContainsRequiredHeaderFields() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Remove data URI prefix
        string memory json = string(bytes(tokenData)[28:]);
        
        // Check header section exists
        assertTrue(_contains(json, "\"header\""), "VC should contain header section");
        
        // Check required header fields
        assertTrue(_contains(json, "\"@context\""), "VC should contain @context field");
        assertTrue(_contains(json, "\"https://www.w3.org/2018/credentials/v1\""), "VC should contain correct @context value");
        assertTrue(_contains(json, "\"type\""), "VC should contain type field");
        assertTrue(_contains(json, "\"VerifiableCredential\""), "VC should contain VerifiableCredential type");
        assertTrue(_contains(json, "\"EventCertificate\""), "VC should contain EventCertificate type");
        assertTrue(_contains(json, "\"id\""), "VC should contain id field");
        assertTrue(_contains(json, TEST_CERTIFICATE_ID), "VC should contain certificate ID");
        assertTrue(_contains(json, "\"issuer\""), "VC should contain issuer field");
        assertTrue(_contains(json, TEST_ISSUER_ID), "VC should contain issuer ID");
        assertTrue(_contains(json, "\"issuanceDate\""), "VC should contain issuanceDate field");
    }

    function test_VCFormat_ContainsRequiredDataFields() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Remove data URI prefix
        string memory json = string(bytes(tokenData)[28:]);
        
        // Check data section exists
        assertTrue(_contains(json, "\"data\""), "VC should contain data section");
        
        // Check required data fields
        assertTrue(_contains(json, "\"eventName\""), "VC should contain eventName field");
        assertTrue(_contains(json, TEST_EVENT_NAME), "VC should contain event name");
        assertTrue(_contains(json, "\"eventDescription\""), "VC should contain eventDescription field");
        assertTrue(_contains(json, TEST_EVENT_DESCRIPTION), "VC should contain event description");
        assertTrue(_contains(json, "\"certificateTokenId\""), "VC should contain certificateTokenId field");
        assertTrue(_contains(json, "\"certificateId\""), "VC should contain certificateId field");
        assertTrue(_contains(json, TEST_CERTIFICATE_ID), "VC should contain certificate ID");
        assertTrue(_contains(json, "\"userId\""), "VC should contain userId field");
        assertTrue(_contains(json, TEST_USER_ID), "VC should contain user ID");
        assertTrue(_contains(json, "\"issuerId\""), "VC should contain issuerId field");
        assertTrue(_contains(json, TEST_ISSUER_ID), "VC should contain issuer ID");
        assertTrue(_contains(json, "\"issuedAt\""), "VC should contain issuedAt field");
        assertTrue(_contains(json, "\"issuerAddresses\""), "VC should contain issuerAddresses field");
        assertTrue(_contains(json, "\"receiverAddress\""), "VC should contain receiverAddress field");
        assertTrue(_contains(json, "\"encryptedUserData\""), "VC should contain encryptedUserData field");
        assertTrue(_contains(json, TEST_ENCRYPTED_USER_DATA), "VC should contain encrypted user data");
        assertTrue(_contains(json, "\"backendEncryptedUserData\""), "VC should contain backendEncryptedUserData field");
        assertTrue(_contains(json, TEST_BACKEND_ENCRYPTED_USER_DATA), "VC should contain backend encrypted user data");
        assertTrue(_contains(json, "\"status\""), "VC should contain status field");
        assertTrue(_contains(json, "\"VALID\""), "VC should contain VALID status");
    }

    function test_VCFormat_ContainsRequiredProofFields() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Remove data URI prefix
        string memory json = string(bytes(tokenData)[28:]);
        
        // Check proof section exists
        assertTrue(_contains(json, "\"proof\""), "VC should contain proof section");
        
        // Check required proof fields
        assertTrue(_contains(json, "\"encryptedByUserRawData\""), "VC should contain encryptedByUserRawData field");
        assertTrue(_contains(json, TEST_USER_ENCRYPTED_PROOF), "VC should contain user encrypted proof");
        assertTrue(_contains(json, "\"encryptedByBackendRawData\""), "VC should contain encryptedByBackendRawData field");
        assertTrue(_contains(json, TEST_BACKEND_ENCRYPTED_PROOF), "VC should contain backend encrypted proof");
        assertTrue(_contains(json, "\"signMessage\""), "VC should contain signMessage field");
        assertTrue(_contains(json, TEST_SIGN_MESSAGE), "VC should contain sign message");
        assertTrue(_contains(json, "\"host\""), "VC should contain host section");
        assertTrue(_contains(json, "\"signature\""), "VC should contain host signature field");
        assertTrue(_contains(json, TEST_HOST_SIGNATURE), "VC should contain host signature");
        assertTrue(_contains(json, "\"publicKey\""), "VC should contain host public key field");
        assertTrue(_contains(json, TEST_HOST_PUBLIC_KEY), "VC should contain host public key");
        assertTrue(_contains(json, "\"issuers\""), "VC should contain issuers array");
        assertTrue(_contains(json, "\"issuerSignature\""), "VC should contain issuer signature field");
        assertTrue(_contains(json, "\"issuerPublicKey\""), "VC should contain issuer public key field");
    }

    function test_VCFormat_HasCorrectDataUriPrefix() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Check data URI prefix
        assertTrue(_contains(tokenData, "data:application/json;utf8,"), "VC should have correct data URI prefix");
    }

    function test_VCFormat_RevokedCertificateStatus() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        // Revoke the certificate
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage, signature);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Check status is REVOKED
        assertTrue(_contains(tokenData, "\"REVOKED\""), "Revoked VC should contain REVOKED status");
        assertTrue(!_contains(tokenData, "\"VALID\""), "Revoked VC should not contain VALID status");
    }

    function test_VCFormat_ValidCertificateStatus() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Check status is VALID
        assertTrue(_contains(tokenData, "\"VALID\""), "Valid VC should contain VALID status");
        assertTrue(!_contains(tokenData, "\"REVOKED\""), "Valid VC should not contain REVOKED status");
    }

    function test_VCFormat_JsonStructure() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Remove data URI prefix
        string memory json = string(bytes(tokenData)[28:]);
        
        // Check basic JSON structure
        assertTrue(_contains(json, "{"), "VC should start with opening brace");
        assertTrue(_contains(json, "}"), "VC should end with closing brace");
        
        // Check main sections
        assertTrue(_contains(json, "\"header\":"), "VC should have header section");
        assertTrue(_contains(json, "\"data\":"), "VC should have data section");
        assertTrue(_contains(json, "\"proof\":"), "VC should have proof section");
        
        // Check sections are properly separated
        uint256 headerCount = 0;
        uint256 dataCount = 0;
        uint256 proofCount = 0;
        
        bytes memory jsonBytes = bytes(json);
        for (uint256 i = 0; i < jsonBytes.length - 6; i++) {
            if (
                jsonBytes[i] == '"' &&
                jsonBytes[i + 1] == 'h' &&
                jsonBytes[i + 2] == 'e' &&
                jsonBytes[i + 3] == 'a' &&
                jsonBytes[i + 4] == 'd' &&
                jsonBytes[i + 5] == 'e' &&
                jsonBytes[i + 6] == 'r' &&
                jsonBytes[i + 7] == '"'
            ) {
                headerCount++;
            } else if (
                jsonBytes[i] == '"' &&
                jsonBytes[i + 1] == 'd' &&
                jsonBytes[i + 2] == 'a' &&
                jsonBytes[i + 3] == 't' &&
                jsonBytes[i + 4] == 'a' &&
                jsonBytes[i + 5] == '"'
            ) {
                dataCount++;
            } else if (
                jsonBytes[i] == '"' &&
                jsonBytes[i + 1] == 'p' &&
                jsonBytes[i + 2] == 'r' &&
                jsonBytes[i + 3] == 'o' &&
                jsonBytes[i + 4] == 'o' &&
                jsonBytes[i + 5] == 'f' &&
                jsonBytes[i + 6] == '"'
            ) {
                proofCount++;
            }
        }
        
        assertEq(headerCount, 1, "VC should have exactly one header section");
        assertEq(dataCount, 1, "VC should have exactly one data section");
        assertEq(proofCount, 1, "VC should have exactly one proof section");
    }

    function test_VCFormat_IssuerAddressesFormat() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Check issuer addresses are formatted as JSON array
        assertTrue(_contains(tokenData, "\"issuerAddresses\": \"["), "Issuer addresses should be formatted as array");
        assertTrue(_contains(tokenData, "]\""), "Issuer addresses array should be closed");
        
        // Check addresses are in quotes
        assertTrue(_contains(tokenData, string(abi.encodePacked("\"", _addressToString(testIssuerAddresses[0]), "\""))), "First issuer address should be in quotes");
        assertTrue(_contains(tokenData, string(abi.encodePacked("\"", _addressToString(testIssuerAddresses[1]), "\""))), "Second issuer address should be in quotes");
        
        // Check addresses are separated by comma
        assertTrue(_contains(tokenData, string(abi.encodePacked(_addressToString(testIssuerAddresses[0]), ","))), "Issuer addresses should be separated by comma");
    }

    function test_VCFormat_IssuerProofsFormat() public {
        _mintCertificateWithProof(PARTICIPANT);
        
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Check issuer proofs are formatted as JSON array
        assertTrue(_contains(tokenData, "\"issuers\": ["), "Issuer proofs should be formatted as array");
        assertTrue(_contains(tokenData, "]"), "Issuer proofs array should be closed");
        
        // Check each issuer proof has required fields
        assertTrue(_contains(tokenData, "\"issuerSignature\""), "Issuer proof should contain signature field");
        assertTrue(_contains(tokenData, "\"issuerPublicKey\""), "Issuer proof should contain public key field");
        
        // Check issuer proofs are separated by comma
        assertTrue(_contains(tokenData, "issuerSignature1"), "First issuer signature should be present");
        assertTrue(_contains(tokenData, "issuerSignature2"), "Second issuer signature should be present");
    }

    function _addressToString(address addr) internal pure returns (string memory) {
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
}
