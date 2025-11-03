// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EventTicket} from "../src/contracts/event/EventTicket.sol";
import {TestEventAccessManager} from "./utils/TestEventAccessManager.sol";
import {MockDecmAccessManager} from "./utils/MockDecmAccessManager.sol";
import {MockEvent} from "./utils/MockEvent.sol";
import {TestUtils} from "./utils/TestUtils.sol";
import {Event} from "../src/contracts/event/Event.sol";

contract EventTicketTest is TestUtils {
    EventTicket public eventTicket;
    TestEventAccessManager public eventAccessManager;
    MockDecmAccessManager public mockDecmAccessManager;
    MockEvent public mockEvent;

    // Test constants
    string constant TEST_EVENT_NAME = "Test Event";
    string constant TEST_EVENT_DESCRIPTION = "This is a test event";
    uint256 constant TEST_SEATS_COUNT = 100;
    
    string constant TEST_USER_ID = "user123";
    string constant TEST_TICKET_ID = "ticket123";
    string constant TEST_ISSUER_ID = "issuer123";
    string constant TEST_ENCRYPTED_USER_DATA = "encryptedUserData123";
    string constant TEST_BACKEND_ENCRYPTED_USER_DATA = "backendEncryptedUserData123";

    event TicketMinted(
        uint256 indexed tokenId,
        address indexed issuer,
        address indexed receiver,
        string ticketId
    );
    event TicketStatusUpdated(uint256 indexed tokenId, EventTicket.TicketStatus status);

    function setUp() public {
        setupMockAddresses();
        
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
        
        // Deploy EventTicket contract
        vm.prank(ADMIN);
        eventTicket = new EventTicket(
            address(eventAccessManager),
            address(mockEvent)
        );
    }

    /*//////////////////////////////////////////////////////////////
                           CONSTRUCTOR TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_EventAccessManagerAddressIsZero() public {
        vm.expectRevert(EventTicket.EventTicket__AccessManagerCannotBeZeroAddress.selector);
        new EventTicket(
            address(0),
            address(mockEvent)
        );
    }

    function test_RevertWhen_EventAddressIsZero() public {
        vm.expectRevert(EventTicket.EventTicket__EventAddressCannotBeZeroAddress.selector);
        new EventTicket(
            address(eventAccessManager),
            address(0)
        );
    }

    function test_WhenConstructorParametersAreValid() public view {
        assertEq(eventTicket.name(), "DECM Event Ticket", "Token name should be set correctly");
        assertEq(eventTicket.symbol(), "DECMT", "Token symbol should be set correctly");
        assertEq(address(eventTicket.EVENT_ACCESS_MANAGER()), address(eventAccessManager), "Event access manager should be set correctly");
        assertEq(address(eventTicket.EVENT()), address(mockEvent), "Event contract should be set correctly");
    }

    /*//////////////////////////////////////////////////////////////
                           MINT NFT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_MintNftWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventTicket), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(EventTicket.EventTicket__NotHostOrAdmin.selector);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_MintNftWithZeroAddress() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(EventTicket.EventTicket__InvalidReceiver.selector);
        eventTicket.mintNft(
            address(0),
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
    }

    function test_WhenMintNftWithHostPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit TicketMinted(0, HOST, PARTICIPANT, TEST_TICKET_ID);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        assertEq(eventTicket.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        assertEq(eventTicket.getTokenCounter(), 1, "Token counter should be 1");
        
        // Check token data
        string memory tokenData = eventTicket.getTokenData(0);
        assertTrue(bytes(tokenData).length > 0, "Token data should not be empty");
        assertTrue(_contains(tokenData, TEST_TICKET_ID), "Token data should contain ticket ID");
        assertTrue(_contains(tokenData, TEST_USER_ID), "Token data should contain user ID");
        assertTrue(_contains(tokenData, TEST_ISSUER_ID), "Token data should contain issuer ID");
        assertTrue(_contains(tokenData, TEST_ENCRYPTED_USER_DATA), "Token data should contain encrypted user data");
        assertTrue(_contains(tokenData, TEST_BACKEND_ENCRYPTED_USER_DATA), "Token data should contain backend encrypted user data");
    }

    function test_WhenMintNftWithAdminPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventTicket), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        assertEq(eventTicket.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        assertEq(eventTicket.getTokenCounter(), 1, "Token counter should be 1");
    }

    /*//////////////////////////////////////////////////////////////
                    BULK MINT PARTICIPANT TICKETS TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_BulkMintWithoutPermission() public {
        EventTicket.BulkMintParticipantTicketsParams[] memory params = new EventTicket.BulkMintParticipantTicketsParams[](1);
        params[0] = EventTicket.BulkMintParticipantTicketsParams({
            receiverAddress: PARTICIPANT,
            userId: TEST_USER_ID,
            ticketId: TEST_TICKET_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddress: ISSUER
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventTicket), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(EventTicket.EventTicket__NotHostOrAdmin.selector);
        eventTicket.bulkMintParticipantTickets(params, signMessage, signature);
    }

    function test_RevertWhen_BulkMintWithZeroAddress() public {
        EventTicket.BulkMintParticipantTicketsParams[] memory params = new EventTicket.BulkMintParticipantTicketsParams[](1);
        params[0] = EventTicket.BulkMintParticipantTicketsParams({
            receiverAddress: address(0),
            userId: TEST_USER_ID,
            ticketId: TEST_TICKET_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddress: ISSUER
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(EventTicket.EventTicket__InvalidReceiver.selector);
        eventTicket.bulkMintParticipantTickets(params, signMessage, signature);
    }

    function test_WhenBulkMintWithHostPermission() public {
        EventTicket.BulkMintParticipantTicketsParams[] memory params = new EventTicket.BulkMintParticipantTicketsParams[](2);
        
        // First participant
        params[0] = EventTicket.BulkMintParticipantTicketsParams({
            receiverAddress: PARTICIPANT,
            userId: TEST_USER_ID,
            ticketId: TEST_TICKET_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddress: ISSUER
        });
        
        // Second participant
        address secondParticipant = makeAddr("secondParticipant");
        params[1] = EventTicket.BulkMintParticipantTicketsParams({
            receiverAddress: secondParticipant,
            userId: "user456",
            ticketId: "ticket456",
            issuerId: "issuer456",
            encryptedUserData: "encryptedUserData456",
            backendEncryptedUserData: "backendEncryptedUserData456",
            issuerAddress: ISSUER
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit TicketMinted(0, HOST, PARTICIPANT, TEST_TICKET_ID);
        
        vm.expectEmit(true, true, true, true);
        emit TicketMinted(1, HOST, secondParticipant, "ticket456");
        
        vm.prank(CALLER);
        eventTicket.bulkMintParticipantTickets(params, signMessage, signature);
        
        assertEq(eventTicket.ownerOf(0), PARTICIPANT, "First token should be owned by first participant");
        assertEq(eventTicket.ownerOf(1), secondParticipant, "Second token should be owned by second participant");
        assertEq(eventTicket.getTokenCounter(), 2, "Token counter should be 2");
        
        // Check token data for first token
        string memory tokenData1 = eventTicket.getTokenData(0);
        assertTrue(_contains(tokenData1, TEST_TICKET_ID), "First token data should contain first ticket ID");
        
        // Check token data for second token
        string memory tokenData2 = eventTicket.getTokenData(1);
        assertTrue(_contains(tokenData2, "ticket456"), "Second token data should contain second ticket ID");
    }

    function test_WhenBulkMintWithAdminPermission() public {
        EventTicket.BulkMintParticipantTicketsParams[] memory params = new EventTicket.BulkMintParticipantTicketsParams[](1);
        params[0] = EventTicket.BulkMintParticipantTicketsParams({
            receiverAddress: PARTICIPANT,
            userId: TEST_USER_ID,
            ticketId: TEST_TICKET_ID,
            issuerId: TEST_ISSUER_ID,
            encryptedUserData: TEST_ENCRYPTED_USER_DATA,
            backendEncryptedUserData: TEST_BACKEND_ENCRYPTED_USER_DATA,
            issuerAddress: ISSUER
        });
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventTicket), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.bulkMintParticipantTickets(params, signMessage, signature);
        
        assertEq(eventTicket.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        assertEq(eventTicket.getTokenCounter(), 1, "Token counter should be 1");
    }

    /*//////////////////////////////////////////////////////////////
                        GET TOKEN DATA TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_GetTokenDataWithOutOfBoundsTokenId() public {
        // First mint a token
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        // Try to get token data with out of bounds token ID
        vm.expectRevert(EventTicket.EventTicket__TokenIdOutOfBounds.selector);
        eventTicket.getTokenData(1);
    }

    function test_WhenGetTokenDataWithValidTokenId() public {
        // First mint a token
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        // Get token data
        string memory tokenData = eventTicket.getTokenData(0);
        
        // Verify JSON structure
        assertTrue(_contains(tokenData, "@context"), "Token data should contain @context");
        assertTrue(_contains(tokenData, "VerifiableCredential"), "Token data should contain VerifiableCredential type");
        assertTrue(_contains(tokenData, "EventTicket"), "Token data should contain EventTicket type");
        assertTrue(_contains(tokenData, "id"), "Token data should contain id");
        assertTrue(_contains(tokenData, "issuer"), "Token data should contain issuer");
        assertTrue(_contains(tokenData, "issuanceDate"), "Token data should contain issuanceDate");
        assertTrue(_contains(tokenData, "credentialSubject"), "Token data should contain credentialSubject");
        
        // Verify data URI format
        assertTrue(_contains(tokenData, "data:application/json;utf8,"), "Token data should have correct data URI format");
        
        // Verify status is ACTIVE
        assertTrue(_contains(tokenData, "ACTIVE"), "Token data should contain ACTIVE status");
    }

    /*//////////////////////////////////////////////////////////////
                        TOKEN URI TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_TokenURIWithOutOfBoundsTokenId() public {
        // First mint a token
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        // Try to get token URI with out of bounds token ID
        vm.expectRevert(EventTicket.EventTicket__TokenIdOutOfBounds.selector);
        eventTicket.tokenURI(1);
    }

    function test_WhenTokenURIWithValidTokenId() public {
        // First mint a token
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        // Get token URI
        string memory tokenUri = eventTicket.tokenURI(0);
        
        // Should be same as getTokenData
        string memory tokenData = eventTicket.getTokenData(0);
        assertEq(tokenUri, tokenData, "Token URI should be same as token data");
    }

    /*//////////////////////////////////////////////////////////////
                        INTEGRATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_FullFlowMintAndVerify() public {
        // Mint a ticket
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        
        // Verify token ownership
        assertEq(eventTicket.ownerOf(0), PARTICIPANT, "Token should be owned by participant");
        
        // Verify token data
        string memory tokenData = eventTicket.getTokenData(0);
        assertTrue(_contains(tokenData, TEST_TICKET_ID), "Token data should contain ticket ID");
        assertTrue(_contains(tokenData, TEST_USER_ID), "Token data should contain user ID");
        assertTrue(_contains(tokenData, TEST_ISSUER_ID), "Token data should contain issuer ID");
        
        // Verify token URI
        string memory tokenUri = eventTicket.tokenURI(0);
        assertEq(tokenUri, tokenData, "Token URI should be same as token data");
    }

    function test_MultipleMintsAndTokenCounterIncrement() public {
        // Mint first ticket
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage1,
            signature1
        );
        
        // Verify token counter is 1
        assertEq(eventTicket.getTokenCounter(), 1, "Token counter should be 1");
        
        // Mint second ticket
        address secondParticipant = makeAddr("secondParticipant");
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            secondParticipant,
            "user456",
            "ticket456",
            "issuer456",
            "encryptedUserData456",
            "backendEncryptedUserData456",
            ISSUER,
            signMessage2,
            signature2
        );
        
        // Verify token counter is 2
        assertEq(eventTicket.getTokenCounter(), 2, "Token counter should be 2");
        
        // Verify token ownership
        assertEq(eventTicket.ownerOf(0), PARTICIPANT, "First token should be owned by first participant");
        assertEq(eventTicket.ownerOf(1), secondParticipant, "Second token should be owned by second participant");
    }

    /*//////////////////////////////////////////////////////////////
                        EDGE CASE TESTS
    //////////////////////////////////////////////////////////////*/

    function test_MintWithMaximumStringParameters() public {
        // Create long strings
        string memory longUserId = new string(1000);
        string memory longTicketId = new string(1000);
        string memory longIssuerId = new string(1000);
        string memory longEncryptedUserData = new string(1000);
        string memory longBackendEncryptedUserData = new string(1000);
        
        // Fill strings with 'a'
        for (uint256 i = 0; i < 1000; i++) {
            bytes(longUserId)[i] = 'a';
            bytes(longTicketId)[i] = 'a';
            bytes(longIssuerId)[i] = 'a';
            bytes(longEncryptedUserData)[i] = 'a';
            bytes(longBackendEncryptedUserData)[i] = 'a';
        }
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            longUserId,
            longTicketId,
            longIssuerId,
            longEncryptedUserData,
            longBackendEncryptedUserData,
            ISSUER,
            signMessage,
            signature
        );
        
        // Verify token data contains long strings
        string memory tokenData = eventTicket.getTokenData(0);
        assertTrue(_contains(tokenData, longUserId), "Token data should contain long user ID");
        assertTrue(_contains(tokenData, longTicketId), "Token data should contain long ticket ID");
        assertTrue(_contains(tokenData, longIssuerId), "Token data should contain long issuer ID");
    }

    function test_MintWithSpecialCharacters() public {
        string memory specialCharsUserId = "user!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsTicketId = "ticket!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsIssuerId = "issuer!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsEncryptedUserData = "encrypted!@#$%^&*()_+-=[]{}|;':\",./<>?";
        string memory specialCharsBackendEncryptedUserData = "backend!@#$%^&*()_+-=[]{}|;':\",./<>?";
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventTicket.mintNft(
            PARTICIPANT,
            specialCharsUserId,
            specialCharsTicketId,
            specialCharsIssuerId,
            specialCharsEncryptedUserData,
            specialCharsBackendEncryptedUserData,
            ISSUER,
            signMessage,
            signature
        );
        
        // Verify token data contains special characters
        string memory tokenData = eventTicket.getTokenData(0);
        assertTrue(_contains(tokenData, specialCharsUserId), "Token data should contain special characters in user ID");
        assertTrue(_contains(tokenData, specialCharsTicketId), "Token data should contain special characters in ticket ID");
        assertTrue(_contains(tokenData, specialCharsIssuerId), "Token data should contain special characters in issuer ID");
    }

    /*//////////////////////////////////////////////////////////////
                        GAS OPTIMIZATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_GasUsageForSingleMint() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        uint256 gasStart = gasleft();
        eventTicket.mintNft(
            PARTICIPANT,
            TEST_USER_ID,
            TEST_TICKET_ID,
            TEST_ISSUER_ID,
            TEST_ENCRYPTED_USER_DATA,
            TEST_BACKEND_ENCRYPTED_USER_DATA,
            ISSUER,
            signMessage,
            signature
        );
        uint256 gasUsed = gasStart - gasleft();
        
        // Just log the gas usage for reference
        emit log_named_uint("Gas used for single mint", gasUsed);
    }

    function test_GasUsageForBulkMint() public {
        // Create params for 5 tickets
        EventTicket.BulkMintParticipantTicketsParams[] memory params = new EventTicket.BulkMintParticipantTicketsParams[](5);
        
        for (uint256 i = 0; i < 5; i++) {
            params[i] = EventTicket.BulkMintParticipantTicketsParams({
                receiverAddress: makeAddr(string(abi.encodePacked("participant", i))),
                userId: string(abi.encodePacked("user", i)),
                ticketId: string(abi.encodePacked("ticket", i)),
                issuerId: string(abi.encodePacked("issuer", i)),
                encryptedUserData: string(abi.encodePacked("encryptedUserData", i)),
                backendEncryptedUserData: string(abi.encodePacked("backendEncryptedUserData", i)),
                issuerAddress: ISSUER
            });
        }
        
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        uint256 gasStart = gasleft();
        eventTicket.bulkMintParticipantTickets(params, signMessage, signature);
        uint256 gasUsed = gasStart - gasleft();
        
        // Just log the gas usage for reference
        emit log_named_uint("Gas used for bulk mint of 5 tickets", gasUsed);
    }

    function test_CompareGasEfficiency() public {
        // Test single mints
        uint256 totalGasForSingleMints = 0;
        
        for (uint256 i = 0; i < 5; i++) {
            (string memory singleSignMessage, bytes memory singleSignature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
            
            vm.prank(CALLER);
            uint256 singleGasStart = gasleft();
            eventTicket.mintNft(
                makeAddr(string(abi.encodePacked("participant", i))),
                string(abi.encodePacked("user", i)),
                string(abi.encodePacked("ticket", i)),
                string(abi.encodePacked("issuer", i)),
                string(abi.encodePacked("encryptedUserData", i)),
                string(abi.encodePacked("backendEncryptedUserData", i)),
                ISSUER,
                singleSignMessage,
                singleSignature
            );
            totalGasForSingleMints += singleGasStart - gasleft();
        }
        
        // Test bulk mint
        EventTicket.BulkMintParticipantTicketsParams[] memory params = new EventTicket.BulkMintParticipantTicketsParams[](5);
        
        for (uint256 i = 0; i < 5; i++) {
            params[i] = EventTicket.BulkMintParticipantTicketsParams({
                receiverAddress: makeAddr(string(abi.encodePacked("participant", i))),
                userId: string(abi.encodePacked("user", i)),
                ticketId: string(abi.encodePacked("ticket", i)),
                issuerId: string(abi.encodePacked("issuer", i)),
                encryptedUserData: string(abi.encodePacked("encryptedUserData", i)),
                backendEncryptedUserData: string(abi.encodePacked("backendEncryptedUserData", i)),
                issuerAddress: ISSUER
            });
        }
        
        (string memory bulkSignMessage, bytes memory bulkSignature) = createSignedMessageForRole(HOST, address(eventTicket), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        uint256 bulkGasStart = gasleft();
        eventTicket.bulkMintParticipantTickets(params, bulkSignMessage, bulkSignature);
        uint256 gasForBulkMint = bulkGasStart - gasleft();
        
        // Log gas usage for comparison
        emit log_named_uint("Total gas for 5 single mints", totalGasForSingleMints);
        emit log_named_uint("Gas for bulk mint of 5 tickets", gasForBulkMint);
        emit log_named_uint("Gas savings", totalGasForSingleMints - gasForBulkMint);
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
}
