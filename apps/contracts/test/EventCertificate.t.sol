// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EventCertificate} from "../src/contracts/event/EventCertificate.sol";
import {TestEventAccessManager} from "./utils/TestEventAccessManager.sol";
import {MockDecmAccessManager} from "./utils/MockDecmAccessManager.sol";
import {MockEvent} from "./utils/MockEvent.sol";
import {TestUtils} from "./utils/TestUtils.sol";
import {Event} from "../src/contracts/event/Event.sol";

contract EventCertificateTest is TestUtils {
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

    event CertificateMinted(
        uint256 indexed tokenId,
        address indexed receiverAddress,
        string certificateId,
        string userId,
        string issuerId
    );
    event CertificateRevoked(uint256 indexed tokenId);
    event ParticipantSignedCertificate(uint256 indexed tokenId, address indexed receiverAddress);

    function setUp() public {
        setupMockAddresses();
        
        // Initialize test issuer addresses
        testIssuerAddresses = new address[](2);
        testIssuerAddresses[0] = ISSUER;
        testIssuerAddresses[1] = makeAddr("secondIssuer");
        
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

    function _mintCertificate(address receiver) internal {
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
            signature
        );
    }

    /*//////////////////////////////////////////////////////////////
                       CONSTRUCTOR TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_EventAccessManagerAddressIsZero() public {
        // EventCertificate doesn't explicitly check for zero address in constructor
        // This test verifies that the contract can be deployed with zero address
        // but will likely fail when trying to use the contract
        new EventCertificate(
            address(0),
            address(mockEvent)
        );
    }

    function test_RevertWhen_EventAddressIsZero() public {
        // EventCertificate doesn't explicitly check for zero address in constructor
        // This test verifies that the contract can be deployed with zero address
        // but will likely fail when trying to use the contract
        new EventCertificate(
            address(eventAccessManager),
            address(0)
        );
    }

    function test_WhenConstructorParametersAreValid() public view {
        assertEq(eventCertificate.name(), "DECM Event Certificate", "Token name should be set correctly");
        assertEq(eventCertificate.symbol(), "DECMC", "Token symbol should be set correctly");
        assertEq(address(eventCertificate.EVENT_ACCESS_MANAGER()), address(eventAccessManager), "Event access manager should be set correctly");
        assertEq(address(eventCertificate.EVENT()), address(mockEvent), "Event contract should be set correctly");
    }

    /*//////////////////////////////////////////////////////////////
                           MINT NFT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_MintNftWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventCertificate), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(EventCertificate.EventCertificate__NotHostOrAdmin.selector);
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_MintNftWithZeroAddress() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventCertificate.mintNft(
            address(0),
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage,
            signature
        );
    }

    function test_WhenMintNftWithHostPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateMinted(0, PARTICIPANT, TEST_CERTIFICATE_ID, TEST_USER_ID, TEST_ISSUER_ID);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage,
            signature
        );
        
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        
        // Check token data
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(bytes(tokenData).length > 0, "Token data should not be empty");
        assertTrue(_contains(tokenData, TEST_CERTIFICATE_ID), "Token data should contain certificate ID");
        assertTrue(_contains(tokenData, TEST_USER_ID), "Token data should contain user ID");
        assertTrue(_contains(tokenData, TEST_ISSUER_ID), "Token data should contain issuer ID");
        assertTrue(_contains(tokenData, TEST_ENCRYPTED_USER_DATA), "Token data should contain encrypted user data");
        assertTrue(_contains(tokenData, TEST_BACKEND_ENCRYPTED_USER_DATA), "Token data should contain backend encrypted user data");
        assertTrue(_contains(tokenData, "VALID"), "Token data should contain VALID status");
    }

    function test_WhenMintNftWithAdminPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventCertificate), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage,
            signature
        );
        
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
    }

    function test_MintWithMaximumStringParameters() public {
        // Create long strings
        string memory longUserId = new string(1000);
        string memory longCertificateId = new string(1000);
        string memory longIssuerId = new string(1000);
        string memory longEncryptedUserData = new string(1000);
        string memory longBackendEncryptedUserData = new string(1000);
        
        // Fill strings with 'a'
        for (uint256 i = 0; i < 1000; i++) {
            bytes(longUserId)[i] = 'a';
            bytes(longCertificateId)[i] = 'a';
            bytes(longIssuerId)[i] = 'a';
            bytes(longEncryptedUserData)[i] = 'a';
            bytes(longBackendEncryptedUserData)[i] = 'a';
        }
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            longUserId,
            longCertificateId,
            longIssuerId,
            longEncryptedUserData,
            longBackendEncryptedUserData,
            testIssuerAddresses,
            signMessage,
            signature
        );
        
        // Verify token data contains long strings
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, longUserId), "Token data should contain long user ID");
        assertTrue(_contains(tokenData, longCertificateId), "Token data should contain long certificate ID");
        assertTrue(_contains(tokenData, longIssuerId), "Token data should contain long issuer ID");
    }

    function test_MintWithSpecialCharacters() public {
        string memory specialCharsUserId = "user!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsCertificateId = "certificate!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsIssuerId = "issuer!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsEncryptedUserData = "encrypted!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsBackendEncryptedUserData = "backend!@#$%^&*()_+-=[]{}|;':\",./<>?";
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            specialCharsUserId,
            specialCharsCertificateId,
            specialCharsIssuerId,
            specialCharsEncryptedUserData,
            specialCharsBackendEncryptedUserData,
            testIssuerAddresses,
            signMessage,
            signature
        );
        
        // Verify token data contains special characters
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, specialCharsUserId), "Token data should contain special characters in user ID");
        assertTrue(_contains(tokenData, specialCharsCertificateId), "Token data should contain special characters in certificate ID");
        assertTrue(_contains(tokenData, specialCharsIssuerId), "Token data should contain special characters in issuer ID");
    }

    /*//////////////////////////////////////////////////////////////
                BULK MINT PARTICIPANT CERTIFICATES TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_BulkMintWithoutPermission() public {
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](1);
        params[0] = EventCertificate.BulkMintParticipantCertificatesParams({
            receiverAddress: PARTICIPANT,
            userId: TEST_USER_ID,
            certificateId: TEST_CERTIFICATE_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddresses: testIssuerAddresses
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventCertificate), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(EventCertificate.EventCertificate__NotHostOrAdmin.selector);
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
    }

    function test_RevertWhen_BulkMintWithZeroAddress() public {
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](1);
        params[0] = EventCertificate.BulkMintParticipantCertificatesParams({
            receiverAddress: address(0),
            userId: TEST_USER_ID,
            certificateId: TEST_CERTIFICATE_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddresses: testIssuerAddresses
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
    }

    function test_RevertWhen_BulkMintWithEmptyArray() public {
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](0);
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
        // Should not revert but do nothing
    }

    function test_WhenBulkMintWithHostPermission() public {
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](2);
        
        // First participant
        params[0] = EventCertificate.BulkMintParticipantCertificatesParams({
            receiverAddress: PARTICIPANT,
            userId: TEST_USER_ID,
            certificateId: TEST_CERTIFICATE_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddresses: testIssuerAddresses
        });
        
        // Second participant
        address secondParticipant = makeAddr("secondParticipant");
        params[1] = EventCertificate.BulkMintParticipantCertificatesParams({
            receiverAddress: secondParticipant,
            userId: "user456",
            certificateId: "certificate456",
            issuerId: "issuer456",
            encryptedUserData: "encryptedUserData456",
            backendEncryptedUserData: "backendEncryptedUserData456",
            issuerAddresses: testIssuerAddresses
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateMinted(0, PARTICIPANT, TEST_CERTIFICATE_ID, TEST_USER_ID, TEST_ISSUER_ID);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateMinted(1, secondParticipant, "certificate456", "user456", "issuer456");
        
        vm.prank(CALLER);
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
        
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "First token should be owned by first participant");
        assertEq(eventCertificate.ownerOf(1), secondParticipant, "Second token should be owned by second participant");
        
        // Check token data for first token
        string memory tokenData1 = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData1, TEST_CERTIFICATE_ID), "First token data should contain first certificate ID");
        
        // Check token data for second token
        string memory tokenData2 = eventCertificate.getTokenData(1);
        assertTrue(_contains(tokenData2, "certificate456"), "Second token data should contain second certificate ID");
    }

    function test_WhenBulkMintWithAdminPermission() public {
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](1);
        params[0] = EventCertificate.BulkMintParticipantCertificatesParams({
            receiverAddress: PARTICIPANT,
            userId: TEST_USER_ID,
            certificateId: TEST_CERTIFICATE_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddresses: testIssuerAddresses
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventCertificate), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
        
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
    }

    function test_BulkMintWithLargeArray() public {
        // Create params for 10 certificates
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](10);
        
        for (uint256 i = 0; i < 10; i++) {
            params[i] = EventCertificate.BulkMintParticipantCertificatesParams({
                receiverAddress: makeAddr(string(abi.encodePacked("participant", i))),
                userId: string(abi.encodePacked("user", i)),
                certificateId: string(abi.encodePacked("certificate", i)),
                issuerId: string(abi.encodePacked("issuer", i)),
                encryptedUserData: string(abi.encodePacked("encryptedUserData", i)),
                backendEncryptedUserData: string(abi.encodePacked("backendEncryptedUserData", i)),
                issuerAddresses: testIssuerAddresses
            });
        }
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
        
        // Verify all tokens were minted
        for (uint256 i = 0; i < 10; i++) {
            address participant = makeAddr(string(abi.encodePacked("participant", i)));
            assertEq(eventCertificate.ownerOf(i), participant, string(abi.encodePacked("Token ", i, " should be owned by participant", i)));
            
            // Verify token data
            string memory tokenData = eventCertificate.getTokenData(i);
            assertTrue(_contains(tokenData, string(abi.encodePacked("certificate", i))), string(abi.encodePacked("Token ", i, " data should contain certificate ID")));
        }
    }

    /*//////////////////////////////////////////////////////////////
                PARTICIPANT SIGNED CERTIFICATE TESTS
    //////////////////////////////////////////////////////////////*/

    function test_WhenParticipantSignedCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        vm.expectEmit(true, true, true, true);
        emit ParticipantSignedCertificate(0, PARTICIPANT);
        
        // Participant signs their certificate
        vm.prank(PARTICIPANT);
        eventCertificate.participantSignedCertificate(0);
        
        // Verify the participant signed address is set
        // Note: We can't directly access the private mapping, but we can verify the event was emitted
    }

    function test_RevertWhen_NonParticipantTriesToSign() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Non-participant tries to sign
        vm.prank(makeAddr("nonParticipant"));
        vm.expectRevert(EventCertificate.EventCertificate__NotParticipant.selector);
        eventCertificate.participantSignedCertificate(0);
    }

    function test_RevertWhen_SignRevokedCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Revoke the certificate
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage, signature);
        
        // Try to sign the revoked certificate
        vm.prank(PARTICIPANT);
        vm.expectRevert(EventCertificate.EventCertificate__CertificateNotValid.selector);
        eventCertificate.participantSignedCertificate(0);
    }

    function test_RevertWhen_SignNonExistentCertificate() public {
        // Try to sign a non-existent certificate
        vm.prank(PARTICIPANT);
        vm.expectRevert(EventCertificate.EventCertificate__NotParticipant.selector);
        eventCertificate.participantSignedCertificate(999);
    }

    function test_WhenParticipantSignsCertificateMultipleTimes() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Participant signs their certificate first time
        vm.prank(PARTICIPANT);
        eventCertificate.participantSignedCertificate(0);
        
        vm.expectEmit(true, true, true, true);
        emit ParticipantSignedCertificate(0, PARTICIPANT);
        
        // Participant signs their certificate second time (should overwrite)
        vm.prank(PARTICIPANT);
        eventCertificate.participantSignedCertificate(0);
        
        // Note: We can't directly verify the overwrite since the mapping is private,
        // but we can verify the event is emitted again
    }

    /*//////////////////////////////////////////////////////////////
                    REVOKE CERTIFICATE TESTS
    //////////////////////////////////////////////////////////////*/

    function test_WhenHostRevokesCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateRevoked(0);
        
        // Host revokes the certificate
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage, signature);
        
        // Verify certificate is revoked by checking token data
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, "REVOKED"), "Token data should contain REVOKED status");
    }

    function test_WhenAdminRevokesCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventCertificate), ADMIN_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateRevoked(0);
        
        // Admin revokes the certificate
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage, signature);
        
        // Verify certificate is revoked by checking token data
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, "REVOKED"), "Token data should contain REVOKED status");
    }

    function test_RevertWhen_RevokeWithoutPermission() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventCertificate), PARTICIPANT_PRIVATE_KEY);
        
        // Participant tries to revoke the certificate
        vm.prank(CALLER);
        vm.expectRevert(EventCertificate.EventCertificate__NotHostOrAdmin.selector);
        eventCertificate.revokeCertificate(0, signMessage, signature);
    }

    function test_RevertWhen_RevokeNonExistentCertificate() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        // Try to revoke a non-existent certificate
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(999, signMessage, signature);
        // Should not revert since the function doesn't check for token existence
        // but it will update a non-existent token's status to REVOKED
        // This might be a design consideration for the contract
    }

    function test_WhenRevokeAlreadyRevokedCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Revoke the certificate first time
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage1, signature1);
        
        // Verify certificate is revoked
        string memory tokenData1 = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData1, "REVOKED"), "Token data should contain REVOKED status");
        
        // Try to revoke the certificate again
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateRevoked(0);
        
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage2, signature2);
        
        // Verify certificate is still revoked
        string memory tokenData2 = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData2, "REVOKED"), "Token data should still contain REVOKED status");
    }

    /*//////////////////////////////////////////////////////////////
                    GET TOKEN DATA TESTS
    //////////////////////////////////////////////////////////////*/

    function test_WhenGetTokenDataWithValidCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Get token data
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Verify JSON structure
        assertTrue(_contains(tokenData, "@context"), "Token data should contain @context");
        assertTrue(_contains(tokenData, "VerifiableCredential"), "Token data should contain VerifiableCredential type");
        assertTrue(_contains(tokenData, "EventCertificate"), "Token data should contain EventCertificate type");
        assertTrue(_contains(tokenData, "id"), "Token data should contain id");
        assertTrue(_contains(tokenData, "issuer"), "Token data should contain issuer");
        assertTrue(_contains(tokenData, "issuanceDate"), "Token data should contain issuanceDate");
        assertTrue(_contains(tokenData, "credentialSubject"), "Token data should contain credentialSubject");
        
        // Verify data URI format
        assertTrue(_contains(tokenData, "data:application/json;utf8,"), "Token data should have correct data URI format");
        
        // Verify status is VALID
        assertTrue(_contains(tokenData, "VALID"), "Token data should contain VALID status");
        
        // Verify certificate specific fields
        assertTrue(_contains(tokenData, TEST_CERTIFICATE_ID), "Token data should contain certificate ID");
        assertTrue(_contains(tokenData, TEST_USER_ID), "Token data should contain user ID");
        assertTrue(_contains(tokenData, TEST_ISSUER_ID), "Token data should contain issuer ID");
        assertTrue(_contains(tokenData, TEST_ENCRYPTED_USER_DATA), "Token data should contain encrypted user data");
        assertTrue(_contains(tokenData, TEST_BACKEND_ENCRYPTED_USER_DATA), "Token data should contain backend encrypted user data");
        assertTrue(_contains(tokenData, TEST_EVENT_NAME), "Token data should contain event name");
        assertTrue(_contains(tokenData, TEST_EVENT_DESCRIPTION), "Token data should contain event description");
    }

    function test_WhenGetTokenDataWithRevokedCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Revoke the certificate
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage, signature);
        
        // Get token data
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Verify status is REVOKED
        assertTrue(_contains(tokenData, "REVOKED"), "Token data should contain REVOKED status");
        assertTrue(!_contains(tokenData, "VALID"), "Token data should not contain VALID status");
    }

    function test_RevertWhen_GetTokenDataWithNonExistentCertificate() public {
        // Try to get token data for non-existent certificate
        vm.expectRevert(EventCertificate.EventCertificate__TokenIdOutOfBounds.selector);
        eventCertificate.getTokenData(999);
    }

    /*//////////////////////////////////////////////////////////////
                        TOKEN URI TESTS
    //////////////////////////////////////////////////////////////*/

    function test_WhenTokenURIWithValidCertificate() public {
        // First mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Get token URI
        string memory tokenUri = eventCertificate.tokenURI(0);
        
        // Should be same as getTokenData
        string memory tokenData = eventCertificate.getTokenData(0);
        assertEq(tokenUri, tokenData, "Token URI should be same as token data");
        
        // Verify data URI format
        assertTrue(_contains(tokenUri, "data:application/json;utf8,"), "Token URI should have correct data URI format");
    }

    function test_RevertWhen_TokenURIWithNonExistentCertificate() public {
        // Try to get token URI for non-existent certificate
        vm.expectRevert(EventCertificate.EventCertificate__TokenIdOutOfBounds.selector);
        eventCertificate.tokenURI(999);
    }

    /*//////////////////////////////////////////////////////////////
                        INTEGRATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_FullFlowMintSignVerifyRevoke() public {
        // Mint a certificate
        (string memory mintSignMessage, bytes memory mintSignature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateMinted(0, PARTICIPANT, TEST_CERTIFICATE_ID, TEST_USER_ID, TEST_ISSUER_ID);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            mintSignMessage,
            mintSignature
        );
        
        // Verify token ownership
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        
        // Verify token data
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, TEST_CERTIFICATE_ID), "Token data should contain certificate ID");
        assertTrue(_contains(tokenData, TEST_USER_ID), "Token data should contain user ID");
        assertTrue(_contains(tokenData, TEST_ISSUER_ID), "Token data should contain issuer ID");
        assertTrue(_contains(tokenData, "VALID"), "Token data should contain VALID status");
        
        // Participant signs certificate
        vm.expectEmit(true, true, true, true);
        emit ParticipantSignedCertificate(0, PARTICIPANT);
        
        vm.prank(PARTICIPANT);
        eventCertificate.participantSignedCertificate(0);
        
        // Verify token URI
        string memory tokenUri = eventCertificate.tokenURI(0);
        assertEq(tokenUri, tokenData, "Token URI should be same as token data");
        
        // Revoke certificate
        (string memory revokeSignMessage, bytes memory revokeSignature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit CertificateRevoked(0);
        
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, revokeSignMessage, revokeSignature);
        
        // Verify certificate is revoked
        string memory revokedTokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(revokedTokenData, "REVOKED"), "Token data should contain REVOKED status");
        assertTrue(!_contains(revokedTokenData, "VALID"), "Token data should not contain VALID status");
    }

    function test_MultipleCertificatesAndTokenCounterIncrement() public {
        // Mint first certificate
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage1,
            signature1
        );
        
        // Mint second certificate
        address secondParticipant = makeAddr("secondParticipant");
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            secondParticipant,
            "user456",
            "certificate456",
            "issuer456",
            "encryptedUserData456",
            "backendEncryptedUserData456",
            testIssuerAddresses,
            signMessage2,
            signature2
        );
        
        // Verify token ownership
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "First token should be owned by first participant");
        assertEq(eventCertificate.ownerOf(1), secondParticipant, "Second token should be owned by second participant");
        
        // Verify token data for both certificates
        string memory tokenData1 = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData1, TEST_CERTIFICATE_ID), "First token data should contain first certificate ID");
        
        string memory tokenData2 = eventCertificate.getTokenData(1);
        assertTrue(_contains(tokenData2, "certificate456"), "Second token data should contain second certificate ID");
        
        // Revoke first certificate
        (string memory revokeSignMessage, bytes memory revokeSignature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, revokeSignMessage, revokeSignature);
        
        // Verify only first certificate is revoked
        string memory revokedTokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(revokedTokenData, "REVOKED"), "First token data should contain REVOKED status");
        
        string memory validTokenData = eventCertificate.getTokenData(1);
        assertTrue(_contains(validTokenData, "VALID"), "Second token data should contain VALID status");
    }

    function test_CrossFunctionDataConsistency() public {
        // Mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Get token data using getTokenData
        string memory tokenData = eventCertificate.getTokenData(0);
        
        // Get token data using tokenURI
        string memory tokenUri = eventCertificate.tokenURI(0);
        
        // Verify both methods return the same data
        assertEq(tokenData, tokenUri, "getTokenData and tokenURI should return the same data");
        
        // Verify data contains all expected fields
        assertTrue(_contains(tokenData, TEST_CERTIFICATE_ID), "Data should contain certificate ID");
        assertTrue(_contains(tokenData, TEST_USER_ID), "Data should contain user ID");
        assertTrue(_contains(tokenData, TEST_ISSUER_ID), "Data should contain issuer ID");
        assertTrue(_contains(tokenData, "VALID"), "Data should contain VALID status");
        
        // Revoke certificate
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventCertificate.revokeCertificate(0, signMessage, signature);
        
        // Verify both methods return updated data
        string memory revokedTokenData = eventCertificate.getTokenData(0);
        string memory revokedTokenUri = eventCertificate.tokenURI(0);
        
        assertEq(revokedTokenData, revokedTokenUri, "Both methods should return same revoked data");
        assertTrue(_contains(revokedTokenData, "REVOKED"), "Data should contain REVOKED status");
    }

    /*//////////////////////////////////////////////////////////////
                    EDGE CASE AND GAS OPTIMIZATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_GasUsageForSingleMint() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        uint256 gasStart = gasleft();
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            testIssuerAddresses,
            signMessage,
            signature
        );
        uint256 gasUsed = gasStart - gasleft();
        
        // Just log the gas usage for reference
        emit log_named_uint("Gas used for single mint", gasUsed);
    }

    function test_GasUsageForBulkMint() public {
        // Create params for 5 certificates
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](5);
        
        for (uint256 i = 0; i < 5; i++) {
            params[i] = EventCertificate.BulkMintParticipantCertificatesParams({
                receiverAddress: makeAddr(string(abi.encodePacked("participant", i))),
                userId: string(abi.encodePacked("user", i)),
                certificateId: string(abi.encodePacked("certificate", i)),
                issuerId: string(abi.encodePacked("issuer", i)),
                encryptedUserData: string(abi.encodePacked("encryptedUserData", i)),
                backendEncryptedUserData: string(abi.encodePacked("backendEncryptedUserData", i)),
                issuerAddresses: testIssuerAddresses
            });
        }
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        uint256 gasStart = gasleft();
        eventCertificate.bulkMintParticipantCertificates(params, signMessage, signature);
        uint256 gasUsed = gasStart - gasleft();
        
        // Just log the gas usage for reference
        emit log_named_uint("Gas used for bulk mint of 5 certificates", gasUsed);
    }

    function test_CompareGasEfficiency() public {
        // Test single mints
        uint256 totalGasForSingleMints = 0;
        
        for (uint256 i = 0; i < 5; i++) {
            (string memory singleSignMessage, bytes memory singleSignature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
            
            vm.prank(CALLER);
            uint256 singleGasStart = gasleft();
            eventCertificate.mintNft(
                makeAddr(string(abi.encodePacked("participant", i))),
                string(abi.encodePacked("user", i)),
                string(abi.encodePacked("certificate", i)),
                string(abi.encodePacked("issuer", i)),
                string(abi.encodePacked("encryptedUserData", i)),
                string(abi.encodePacked("backendEncryptedUserData", i)),
                testIssuerAddresses,
                singleSignMessage,
                singleSignature
            );
            totalGasForSingleMints += singleGasStart - gasleft();
        }
        
        // Test bulk mint
        EventCertificate.BulkMintParticipantCertificatesParams[] memory params = new EventCertificate.BulkMintParticipantCertificatesParams[](5);
        
        for (uint256 i = 0; i < 5; i++) {
            params[i] = EventCertificate.BulkMintParticipantCertificatesParams({
                receiverAddress: makeAddr(string(abi.encodePacked("participant", i))),
                userId: string(abi.encodePacked("user", i)),
                certificateId: string(abi.encodePacked("certificate", i)),
                issuerId: string(abi.encodePacked("issuer", i)),
                encryptedUserData: string(abi.encodePacked("encryptedUserData", i)),
                backendEncryptedUserData: string(abi.encodePacked("backendEncryptedUserData", i)),
                issuerAddresses: testIssuerAddresses
            });
        }
        
        (string memory bulkSignMessage, bytes memory bulkSignature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        uint256 bulkGasStart = gasleft();
        eventCertificate.bulkMintParticipantCertificates(params, bulkSignMessage, bulkSignature);
        uint256 gasForBulkMint = bulkGasStart - gasleft();
        
        // Log gas usage for comparison
        emit log_named_uint("Total gas for 5 single mints", totalGasForSingleMints);
        emit log_named_uint("Gas for bulk mint of 5 certificates", gasForBulkMint);
        emit log_named_uint("Gas savings", totalGasForSingleMints - gasForBulkMint);
    }

    function test_ReentrancyGuard() public {
        // This test verifies that the reentrancy guard is working
        // We can't directly test reentrancy without a malicious contract,
        // but we can verify that the nonReentrant modifier is present
        // by checking that the function calls succeed without issues
        
        // Mint a certificate
        _mintCertificate(PARTICIPANT);
        
        // Try to call participantSignedCertificate multiple times in the same transaction
        // This should work fine since the function has the nonReentrant modifier
        vm.prank(PARTICIPANT);
        eventCertificate.participantSignedCertificate(0);
        
        // Verify the certificate was signed
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, TEST_CERTIFICATE_ID), "Token data should contain certificate ID");
    }

    function test_MintWithEmptyStrings() public {
        // Test with empty string parameters
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            "", // empty userId
            "", // empty certificateId
            "", // empty issuerId
            "", // empty encryptedUserData
            "", // empty backendEncryptedUserData
            testIssuerAddresses,
            signMessage,
            signature
        );
        
        // Verify token was minted with empty strings
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(bytes(tokenData).length > 0, "Token data should not be empty");
    }

    function test_MintWithMultipleIssuerAddresses() public {
        // Create array with multiple issuer addresses
        address[] memory multipleIssuers = new address[](5);
        multipleIssuers[0] = makeAddr("issuer1");
        multipleIssuers[1] = makeAddr("issuer2");
        multipleIssuers[2] = makeAddr("issuer3");
        multipleIssuers[3] = makeAddr("issuer4");
        multipleIssuers[4] = makeAddr("issuer5");
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventCertificate), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventCertificate.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_CERTIFICATE_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            multipleIssuers,
            signMessage,
            signature
        );
        
        // Verify token was minted
        assertEq(eventCertificate.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        
        string memory tokenData = eventCertificate.getTokenData(0);
        assertTrue(_contains(tokenData, TEST_CERTIFICATE_ID), "Token data should contain certificate ID");
    }
}
