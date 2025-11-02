// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {Event} from "../src/contracts/event/Event.sol";
import {TestEventAccessManager} from "./utils/TestEventAccessManager.sol";
import {MockDecmAccessManager} from "./utils/MockDecmAccessManager.sol";
import {TestUtils} from "./utils/TestUtils.sol";
import {Constants} from "../src/contracts/constants/Constants.s.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract SignatureVerificationTest is TestUtils {
    Event public eventContract;
    TestEventAccessManager public eventAccessManager;
    MockDecmAccessManager public mockDecmAccessManager;

    string constant EVENT_NAME = "Test Event";
    string constant EVENT_DESCRIPTION = "This is a test event";
    uint256 constant SEATS_COUNT = 100;

    function setUp() public {
        setupMockAddresses();
        
        // Deploy MockDecmAccessManager with admin
        address[] memory admins = new address[](1);
        admins[0] = ADMIN;
        vm.prank(ADMIN);
        mockDecmAccessManager = new MockDecmAccessManager(admins);
        
        // Deploy EventAccessManager with host
        vm.prank(ADMIN);
        eventAccessManager = new TestEventAccessManager(address(mockDecmAccessManager), HOST);
        
        // Deploy Event contract
        vm.prank(ADMIN);
        eventContract = new Event(
            address(mockDecmAccessManager),
            EVENT_NAME,
            EVENT_DESCRIPTION,
            SEATS_COUNT,
            HOST
        );
    }

    /*//////////////////////////////////////////////////////////////
                        SIGNATURE FORMAT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_SignMessageFormat() public {
        uint256 deadlineBlock = block.number + 12;
        string memory message = createSignMessage(HOST, address(eventContract), deadlineBlock);
        
        // Verify format: HOST_ADDRESS,CONTRACT_ADDRESS,DEADLINE_BLOCK
        string memory expected = string(
            abi.encodePacked(
                vm.toString(HOST),
                ",",
                vm.toString(address(eventContract)),
                ",",
                vm.toString(deadlineBlock)
            )
        );
        
        assertEq(message, expected, "Sign message format should be correct");
    }

    function test_RecoverSignerWithValidSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Recovered signer should match the original signer");
    }

    function test_RecoverSignerWithInvalidSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        // Corrupt the signature by changing one byte
        signature[0] = bytes1(uint8(signature[0]) ^ 0xff);
        
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, address(eventContract));
    }

    /*//////////////////////////////////////////////////////////////
                        DEADLINE BLOCK TESTS
    //////////////////////////////////////////////////////////////*/

    function test_ValidDeadlineBlock() public {
        // Create a signature with a future deadline
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY, 12);
        
        // Should not revert
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer with valid deadline");
    }

    function test_ExpiredDeadlineBlock() public {
        // Create a signature with a past deadline
        (string memory message, bytes memory signature) = createExpiredSignedMessage(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, address(eventContract));
    }

    function test_ExactlyCurrentBlock() public {
        // Create a signature with the current block as deadline
        string memory message = createSignMessage(HOST, address(eventContract), block.number);
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Should not revert since current block is still valid
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer with current block as deadline");
    }

    function test_OneMinuteDeadline() public {
        // Create a signature with a 1 minute deadline (approximately 12 blocks)
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY, 12);
        
        // Should not revert
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer with 1 minute deadline");
    }

    function test_AdvanceBlocksPastDeadline() public {
        // Create a signature with a 1 block deadline
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY, 1);
        
        // Should not revert initially
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer before deadline");
        
        // Advance blocks past the deadline
        advanceBlocks(2);
        
        // Should now revert
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, address(eventContract));
    }

    /*//////////////////////////////////////////////////////////////
                    CONTRACT ADDRESS VERIFICATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_ValidContractAddress() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        // Should not revert with correct contract address
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer with valid contract address");
    }

    function test_InvalidContractAddress() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        // Should revert with incorrect contract address
        address wrongContract = address(0x999);
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, wrongContract);
    }

    function test_SignatureWithDifferentContractAddress() public {
        // Create a signature for a different contract address
        address differentContract = address(0x123);
        string memory message = createSignMessage(HOST, differentContract, block.number + 12);
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Should revert when verifying with the actual contract address
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, address(eventContract));
    }

    /*//////////////////////////////////////////////////////////////
                        ROLE-BASED SIGNATURE TESTS
    //////////////////////////////////////////////////////////////*/

    function test_HostSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover host signer");
    }

    function test_AdminSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventContract), ADMIN_PRIVATE_KEY);
        
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, ADMIN, "Should recover admin signer");
    }

    function test_ParticipantSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, PARTICIPANT, "Should recover participant signer");
    }

    function test_CallerSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(CALLER, address(eventContract), CALLER_PRIVATE_KEY);
        
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, CALLER, "Should recover caller signer");
    }

    /*//////////////////////////////////////////////////////////////
                    SIGNATURE INTEGRATION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_UpdateEventWithHostSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.updateEvent(
            "Updated Event",
            "Updated Description",
            SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            message,
            signature
        );
        
        assertEq(eventContract.eventName(), "Updated Event", "Event should be updated with host signature");
    }

    function test_UpdateEventWithAdminSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventContract), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.updateEvent(
            "Updated Event",
            "Updated Description",
            SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            message,
            signature
        );
        
        assertEq(eventContract.eventName(), "Updated Event", "Event should be updated with admin signature");
    }

    function test_RevertWhen_UpdateEventWithParticipantSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.updateEvent(
            "Updated Event",
            "Updated Description",
            SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            message,
            signature
        );
    }

    function test_AddParticipantWithHostSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, message, signature);
        
        assertEq(eventContract.currentSeatsCount(), 1, "Participant should be added with host signature");
    }

    function test_LeaveEventWithParticipantSignature() public {
        // First add participant
        (string memory message1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, message1, signature1);
        
        // Then participant leaves
        (string memory message2, bytes memory signature2) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.leaveEvent(message2, signature2);
        
        assertEq(eventContract.currentSeatsCount(), 0, "Participant should be removed with participant signature");
    }

    function test_ConfirmEventWithHostSignature() public {
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.confirmEvent(message, signature);
        
        assertEq(uint256(eventContract.eventStatus()), uint256(Event.EventStatus.CLOSED), "Event should be confirmed with host signature");
    }

    /*//////////////////////////////////////////////////////////////
                    SIGNATURE EDGE CASE TESTS
    //////////////////////////////////////////////////////////////*/

    function test_SignatureWithZeroAddress() public {
        address zeroAddress = address(0);
        (string memory message, bytes memory signature) = createSignedMessageForRole(zeroAddress, address(eventContract), 0x1);
        
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, zeroAddress, "Should recover zero address signer");
    }

    function test_SignatureWithMaxDeadline() public {
        uint256 maxDeadline = type(uint256).max;
        string memory message = createSignMessage(HOST, address(eventContract), maxDeadline);
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Should not revert with max deadline
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer with max deadline");
    }

    function test_SignatureWithMinDeadline() public {
        uint256 minDeadline = 0;
        string memory message = createSignMessage(HOST, address(eventContract), minDeadline);
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Should revert with min deadline (past block)
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, address(eventContract));
    }

    function test_SignatureWithVeryShortDeadline() public {
        // Create a signature with a 1 block deadline
        (string memory message, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY, 1);
        
        // Should not revert initially
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        assertEq(recoveredSigner, HOST, "Should recover signer with very short deadline");
        
        // Advance 1 block
        advanceBlocks(1);
        
        // Should now revert
        vm.expectRevert();
        eventContract.recoverSigner(message, signature, address(eventContract));
    }

    function test_DebugSignatureVerification() public {
        // Create a simple message
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Sign it
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Try to recover
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        
        // Just assert without logging
        assertEq(recoveredSigner, HOST, "Recovered signer should match original signer");
    }

    function test_SimpleSignatureVerification() public {
        // Create a simple message
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Sign it
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Try to recover
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        
        assertEq(recoveredSigner, HOST, "Recovered signer should match original signer");
    }

    function test_DirectSignatureVerification() public {
        // Create a simple message
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Sign it
        bytes memory signature = signMessage(message, HOST_PRIVATE_KEY);
        
        // Try to recover
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        
        assertEq(recoveredSigner, HOST, "Recovered signer should match original signer");
    }

    function test_ManualSignatureVerification() public {
        // Create a message manually
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Sign it manually
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        bytes32 ethSignedMessageHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(HOST_PRIVATE_KEY, ethSignedMessageHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        
        // Try to recover
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        
        assertEq(recoveredSigner, HOST, "Recovered signer should match original signer");
    }

    function test_RawSignatureVerification() public {
        // Create a message manually
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Hash it directly
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        
        // Create Ethereum signed message hash
        bytes32 ethSignedMessageHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash));
        
        // Sign it
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(HOST_PRIVATE_KEY, ethSignedMessageHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        
        // Recover directly
        address recoveredSigner = ECDSA.recover(ethSignedMessageHash, signature);
        
        assertEq(recoveredSigner, HOST, "Recovered signer should match original signer");
    }

    function test_HashSignatureVerification() public {
        // Create a message manually
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Hash it manually
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        
        // Create Ethereum signed message hash
        bytes32 ethSignedMessageHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash));
        
        // Sign it
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(HOST_PRIVATE_KEY, ethSignedMessageHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        
        // Recover directly
        address recoveredSigner = ECDSA.recover(ethSignedMessageHash, signature);
        
        assertEq(recoveredSigner, HOST, "Recovered signer should match original signer");
    }

    function test_ContractSignatureVerification() public {
        // Create a message manually
        string memory message = "0x2345678901234567890123456789012345678901,0xf7c49BE3d09B504206a79Bf68AD8EB41f6dCD541,100";
        
        // Hash it directly
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        
        // Create Ethereum signed message hash
        bytes32 ethSignedMessageHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash));
        
        // Sign it
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(HOST_PRIVATE_KEY, ethSignedMessageHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        
        // Recover using contract method
        address recoveredSigner = eventContract.recoverSigner(message, signature, address(eventContract));
        
        // Recover using ECDSA directly
        address directRecoveredSigner = ECDSA.recover(ethSignedMessageHash, signature);
        
        console.log("Expected signer:", vm.toString(HOST));
        console.log("Contract recovered:", vm.toString(recoveredSigner));
        console.log("Direct recovered:", vm.toString(directRecoveredSigner));
        
        assertEq(recoveredSigner, directRecoveredSigner, "Contract recovery should match direct ECDSA recovery");
    }
}
