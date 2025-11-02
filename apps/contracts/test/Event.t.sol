// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Event} from "../src/contracts/event/Event.sol";
import {TestEventAccessManager} from "./utils/TestEventAccessManager.sol";
import {MockDecmAccessManager} from "./utils/MockDecmAccessManager.sol";
import {TestUtils} from "./utils/TestUtils.sol";

contract EventTest is TestUtils {
    Event public eventContract;
    TestEventAccessManager public eventAccessManager;
    MockDecmAccessManager public mockDecmAccessManager;

    string constant EVENT_NAME = "Test Event";
    string constant EVENT_DESCRIPTION = "This is a test event";
    uint256 constant SEATS_COUNT = 100;
    string constant NEW_EVENT_NAME = "Updated Event";
    string constant NEW_EVENT_DESCRIPTION = "This is an updated test event";
    uint256 constant NEW_SEATS_COUNT = 200;

    event RemovedParticipant(address indexed participant);
    event AddedParticipant(address indexed participant);
    event ParticipantSigned(address indexed participant, bytes32 signature);
    event EventConfirmed();
    event EventUpdated(
        string eventName,
        string eventDescription,
        uint256 seatsCount,
        Event.EventStatus eventStatus
    );

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
                           CONSTRUCTOR TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_DecmAccessManagerAddressIsZero() public {
        vm.expectRevert();
        new Event(
            address(0),
            EVENT_NAME,
            EVENT_DESCRIPTION,
            SEATS_COUNT,
            HOST
        );
    }

    function test_RevertWhen_EventNameIsEmpty() public {
        vm.expectRevert(Event.Event__InvalidEventName.selector);
        new Event(
            address(mockDecmAccessManager),
            "",
            EVENT_DESCRIPTION,
            SEATS_COUNT,
            HOST
        );
    }

    function test_WhenConstructorParametersAreValid() public {
        assertEq(eventContract.eventName(), EVENT_NAME, "Event name should be set correctly");
        assertEq(eventContract.eventDescription(), EVENT_DESCRIPTION, "Event description should be set correctly");
        assertEq(eventContract.seatsCount(), SEATS_COUNT, "Seats count should be set correctly");
        assertEq(eventContract.currentSeatsCount(), 0, "Current seats count should be 0");
        assertEq(uint256(eventContract.eventStatus()), uint256(Event.EventStatus.ACTIVE), "Event status should be ACTIVE");
        assertTrue(eventContract.checkIsHost(HOST), "Host should have host role");
    }

    /*//////////////////////////////////////////////////////////////
                           UPDATE EVENT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_UpdateEventWithInvalidSignature() public {
        (string memory signMessage, bytes memory signature) = createInvalidContractSignedMessage(HOST, HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.updateEvent(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_UpdateEventWithExpiredSignature() public {
        (string memory signMessage, bytes memory signature) = createExpiredSignedMessage(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.updateEvent(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_UpdateEventWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.updateEvent(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_UpdateEventWithEmptyName() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__InvalidEventName.selector);
        eventContract.updateEvent(
            "",
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.ACTIVE,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_UpdateEventWithReducedSeatsCount() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__CannotReduceSeatsCount.selector);
        eventContract.updateEvent(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            SEATS_COUNT - 1,
            Event.EventStatus.ACTIVE,
            signMessage,
            signature
        );
    }

    function test_WhenUpdateEventWithHostPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit EventUpdated(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.INACTIVE
        );
        
        vm.prank(CALLER);
        eventContract.updateEvent(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.INACTIVE,
            signMessage,
            signature
        );
        
        assertEq(eventContract.eventName(), NEW_EVENT_NAME, "Event name should be updated");
        assertEq(eventContract.eventDescription(), NEW_EVENT_DESCRIPTION, "Event description should be updated");
        assertEq(eventContract.seatsCount(), NEW_SEATS_COUNT, "Seats count should be updated");
        assertEq(uint256(eventContract.eventStatus()), uint256(Event.EventStatus.INACTIVE), "Event status should be updated");
    }

    function test_WhenUpdateEventWithAdminPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventContract), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.updateEvent(
            NEW_EVENT_NAME,
            NEW_EVENT_DESCRIPTION,
            NEW_SEATS_COUNT,
            Event.EventStatus.INACTIVE,
            signMessage,
            signature
        );
        
        assertEq(eventContract.eventName(), NEW_EVENT_NAME, "Event name should be updated");
        assertEq(eventContract.eventDescription(), NEW_EVENT_DESCRIPTION, "Event description should be updated");
        assertEq(eventContract.seatsCount(), NEW_SEATS_COUNT, "Seats count should be updated");
        assertEq(uint256(eventContract.eventStatus()), uint256(Event.EventStatus.INACTIVE), "Event status should be updated");
    }

    /*//////////////////////////////////////////////////////////////
                        ADD PARTICIPANT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_AddParticipantWithInvalidSignature() public {
        (string memory signMessage, bytes memory signature) = createInvalidContractSignedMessage(HOST, HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.addParticipant(
            PARTICIPANT,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_AddParticipantWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.addParticipant(
            PARTICIPANT,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_AddParticipantWithZeroAddress() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__AddressCannotBeZero.selector);
        eventContract.addParticipant(
            address(0),
            signMessage,
            signature
        );
    }

    function test_RevertWhen_AddParticipantWhenSeatsCountReached() public {
        // Create an event with only 1 seat
        vm.prank(ADMIN);
        Event smallEvent = new Event(
            address(mockDecmAccessManager),
            "Small Event",
            "Small event description",
            1,
            HOST
        );
        
        // Add the first participant
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(smallEvent), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        smallEvent.addParticipant(PARTICIPANT, signMessage1, signature1);
        
        // Try to add a second participant
        address secondParticipant = makeAddr("secondParticipant");
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(smallEvent), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__SeatsCountReached.selector);
        smallEvent.addParticipant(secondParticipant, signMessage2, signature2);
    }

    function test_RevertWhen_AddParticipantWhoIsAlreadyJoined() public {
        // First add participant
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, signMessage1, signature1);
        
        // Try to add the same participant again
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__ParticipantIsAlreadyJoined.selector);
        eventContract.addParticipant(PARTICIPANT, signMessage2, signature2);
    }

    function test_WhenAddParticipantWithHostPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit AddedParticipant(PARTICIPANT);
        
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, signMessage, signature);
        
        assertEq(eventContract.currentSeatsCount(), 1, "Current seats count should be 1");
        assertTrue(eventContract.checkIsParticipant(PARTICIPANT), "Participant should have participant role");
        
        address[] memory participants = eventContract.getParticipants();
        assertEq(participants.length, 1, "Participants array should have 1 element");
        assertEq(participants[0], PARTICIPANT, "First participant should be the added participant");
    }

    /*//////////////////////////////////////////////////////////////
                        REMOVE PARTICIPANT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_RemoveParticipantWithInvalidSignature() public {
        (string memory signMessage, bytes memory signature) = createInvalidContractSignedMessage(HOST, HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.removeParticipant(
            PARTICIPANT,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_RemoveParticipantWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.removeParticipant(
            PARTICIPANT,
            signMessage,
            signature
        );
    }

    function test_RevertWhen_RemoveParticipantWithZeroAddress() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__AddressCannotBeZero.selector);
        eventContract.removeParticipant(
            address(0),
            signMessage,
            signature
        );
    }

    function test_RevertWhen_RemoveParticipantWhoIsNotJoined() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__ParticipantIsNotJoined.selector);
        eventContract.removeParticipant(
            PARTICIPANT,
            signMessage,
            signature
        );
    }

    function test_WhenRemoveParticipantWithHostPermission() public {
        // First add participant
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, signMessage1, signature1);
        
        // Then remove participant
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit RemovedParticipant(PARTICIPANT);
        
        vm.prank(CALLER);
        eventContract.removeParticipant(PARTICIPANT, signMessage2, signature2);
        
        assertEq(eventContract.currentSeatsCount(), 0, "Current seats count should be 0");
        assertFalse(eventContract.checkIsParticipant(PARTICIPANT), "Participant should no longer have participant role");
        
        address[] memory participants = eventContract.getParticipants();
        assertEq(participants.length, 0, "Participants array should be empty");
    }

    /*//////////////////////////////////////////////////////////////
                           LEAVE EVENT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_LeaveEventWithInvalidSignature() public {
        (string memory signMessage, bytes memory signature) = createInvalidContractSignedMessage(PARTICIPANT, PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.leaveEvent(signMessage, signature);
    }

    function test_RevertWhen_LeaveEventWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.leaveEvent(signMessage, signature);
    }

    function test_RevertWhen_LeaveEventWhenNotJoined() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__ParticipantIsNotJoined.selector);
        eventContract.leaveEvent(signMessage, signature);
    }

    function test_WhenLeaveEventWithParticipantPermission() public {
        // First add participant
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, signMessage1, signature1);
        
        // Then participant leaves
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit RemovedParticipant(PARTICIPANT);
        
        vm.prank(CALLER);
        eventContract.leaveEvent(signMessage2, signature2);
        
        assertEq(eventContract.currentSeatsCount(), 0, "Current seats count should be 0");
        assertFalse(eventContract.checkIsParticipant(PARTICIPANT), "Participant should no longer have participant role");
        
        address[] memory participants = eventContract.getParticipants();
        assertEq(participants.length, 0, "Participants array should be empty");
    }

    /*//////////////////////////////////////////////////////////////
                          CONFIRM EVENT TESTS
    //////////////////////////////////////////////////////////////*/

    function test_RevertWhen_ConfirmEventWithInvalidSignature() public {
        (string memory signMessage, bytes memory signature) = createInvalidContractSignedMessage(HOST, HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.confirmEvent(signMessage, signature);
    }

    function test_RevertWhen_ConfirmEventWithoutPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventContract), PARTICIPANT_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert();
        eventContract.confirmEvent(signMessage, signature);
    }

    function test_RevertWhen_ConfirmEventWhenClosed() public {
        // First close the event
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.confirmEvent(signMessage1, signature1);
        
        // Try to confirm again
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__CantConfirmEvent.selector);
        eventContract.confirmEvent(signMessage2, signature2);
    }

    function test_RevertWhen_ConfirmEventWhenInactive() public {
        // First set event to inactive
        (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.updateEvent(
            EVENT_NAME,
            EVENT_DESCRIPTION,
            SEATS_COUNT,
            Event.EventStatus.INACTIVE,
            signMessage1,
            signature1
        );
        
        // Try to confirm
        (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.prank(CALLER);
        vm.expectRevert(Event.Event__CantConfirmEvent.selector);
        eventContract.confirmEvent(signMessage2, signature2);
    }

    function test_WhenConfirmEventWithHostPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        
        vm.expectEmit(true, true, true, true);
        emit EventConfirmed();
        
        vm.prank(CALLER);
        eventContract.confirmEvent(signMessage, signature);
        
        assertEq(uint256(eventContract.eventStatus()), uint256(Event.EventStatus.CLOSED), "Event status should be CLOSED");
    }

    function test_WhenConfirmEventWithAdminPermission() public {
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventContract), ADMIN_PRIVATE_KEY);
        
        vm.prank(CALLER);
        eventContract.confirmEvent(signMessage, signature);
        
        assertEq(uint256(eventContract.eventStatus()), uint256(Event.EventStatus.CLOSED), "Event status should be CLOSED");
    }

    /*//////////////////////////////////////////////////////////////
                           GETTER FUNCTION TESTS
    //////////////////////////////////////////////////////////////*/

    function test_GetEventName() public {
        assertEq(eventContract.getEventName(), EVENT_NAME, "Should return correct event name");
    }

    function test_GetEventDescription() public {
        assertEq(eventContract.getEventDescription(), EVENT_DESCRIPTION, "Should return correct event description");
    }

    function test_GetParticipants() public {
        // Initially empty
        address[] memory participants = eventContract.getParticipants();
        assertEq(participants.length, 0, "Participants array should be empty initially");
        
        // Add a participant
        (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventContract), HOST_PRIVATE_KEY);
        vm.prank(CALLER);
        eventContract.addParticipant(PARTICIPANT, signMessage, signature);
        
        // Check again
        participants = eventContract.getParticipants();
        assertEq(participants.length, 1, "Participants array should have 1 element");
        assertEq(participants[0], PARTICIPANT, "First participant should be the added participant");
    }
}
