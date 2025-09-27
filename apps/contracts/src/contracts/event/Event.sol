// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EventAccessManager} from "./EventAccessManager.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";

contract Event is EventAccessManager {
    using ECDSA for bytes32;

    // Enums
    enum ParticipantStatus {
        PENDING,
        APPROVED,
        REJECTED,
        LEAVED
    }

    enum EventStatus {
        ACTIVE,
        INACTIVE,
        CLOSED
    }

    // Errors
    error Event__InvalidEventName();
    error Event__CannotReduceSeatsCount();
    error Event__SeatsCountReached();
    error Event__ParticipantIsNotJoined();
    error Event__ParticipantIsAlreadyJoined();
    error Event__AddressCannotBeZero();
    error Event__CantConfirmEvent(string message);

    event RemovedParticipant(address indexed participant);
    event AddedParticipant(address indexed participant);
    event ParticipantSigned(address indexed participant, bytes32 signature);
    event EventConfirmed();
    error Event__InvalidSignature();
    event EventUpdated(
        string eventName,
        string eventDescription,
        uint256 seatsCount,
        EventStatus eventStatus
    );

    // State Variables
    string public eventName;
    string public eventDescription;
    uint256 public seatsCount;
    uint256 public currentSeatsCount;
    EventStatus public eventStatus;

    // Participants Mappings
    mapping(address => ParticipantStatus) participantToStatus;
    mapping(address => bool) private isParticipant;
    mapping(address => uint256) private participantIndex;
    mapping(address => bytes32) private participantToSignature;
    address[] private participants;

    constructor(
        address decmAccessManagerAddr,
        string memory _eventName,
        string memory _eventDescription,
        uint256 _seatsCount
    ) EventAccessManager(decmAccessManagerAddr) {
        _validateEventName(_eventName);

        eventName = _eventName;
        eventDescription = _eventDescription;
        seatsCount = _seatsCount;
        currentSeatsCount = 0;
        eventStatus = EventStatus.ACTIVE;
    }

    function updateEvent(
        string memory _eventName,
        string memory _eventDescription,
        uint256 _seatsCount,
        EventStatus _eventStatus
    ) external onlyHostOrAdmin {
        // 1. Validate Event Name
        _validateEventName(_eventName);

        // 2. Validate Seats Count
        if (_seatsCount < seatsCount) {
            revert Event__CannotReduceSeatsCount();
        }

        // 3. Update Event
        eventName = _eventName;
        eventDescription = _eventDescription;
        seatsCount = _seatsCount;
        eventStatus = _eventStatus;

        // 4. Emit Event
        emit EventUpdated(
            _eventName,
            _eventDescription,
            _seatsCount,
            eventStatus
        );
    }

    function addParticipant(
        address participantAddress
    ) external onlyHostOrAdmin {
        // Pre Conditions
        if (participantAddress == address(0)) {
            revert Event__AddressCannotBeZero();
        }

        if (isParticipant[participantAddress]) {
            revert Event__ParticipantIsAlreadyJoined();
        }

        if (currentSeatsCount >= seatsCount) {
            revert Event__SeatsCountReached();
        }

        // 1. Validate Participant
        if (isParticipant[participantAddress]) {
            revert Event__ParticipantIsAlreadyJoined();
        }

        // 2. Add Participant
        _addParticipant(participantAddress);

        // 3. Emit Event
        emit AddedParticipant(participantAddress);
    }

    function leaveEvent(address participantAddress) external onlyParticipant {
        // 1. Validate Participant
        if (!isParticipant[participantAddress]) {
            revert Event__ParticipantIsNotJoined();
        }

        // 2. Remove Participant
        _removeParticipant(participantAddress);

        // 3. Emit Event
        emit RemovedParticipant(participantAddress);
    }

    function confirmEvent() external onlyHostOrAdmin {
        // Pre Conditions
        if (eventStatus == EventStatus.CLOSED) {
            revert Event__CantConfirmEvent("Event is closed");
        }

        if (eventStatus == EventStatus.INACTIVE) {
            revert Event__CantConfirmEvent("Event is inactive");
        }

        // 1. Update Event Status
        eventStatus = EventStatus.CLOSED;

        // 2. Emit Event
        emit EventConfirmed();
    }

    function setSigningMessage(
        address participant,
        bytes32 messageHash
    ) external {
        if (participant == address(0)) revert Event__AddressCannotBeZero();
        participantToSignature[participant] = messageHash;
        emit ParticipantSigned(participant, messageHash);
    }

    function getSigningMessage(
        address participant
    ) external view returns (bytes32) {
        return participantToSignature[participant];
    }

    function verifySignature(
        address participant,
        bytes calldata signature
    ) external view returns (address) {
        bytes32 msgHash = participantToSignature[participant];
        if (msgHash == bytes32(0)) revert Event__InvalidSignature();

        address recovered = ECDSA.recover(
            MessageHashUtils.toEthSignedMessageHash(msgHash),
            signature
        );

        if (recovered != participant) revert Event__InvalidSignature();

        return recovered;
    }

    function _addParticipant(address participantAddress) private {
        // 1. Check Participant Signature
        if (participantToSignature[participantAddress] == bytes32(0)) {
            revert Event__InvalidSignature();
        }

        // 2. Add Participant
        participantIndex[participantAddress] = participants.length;
        participants.push(participantAddress);
        isParticipant[participantAddress] = true;
        participantToStatus[participantAddress] = ParticipantStatus.APPROVED;

        // 2. Current SeatsCount Increment
        currentSeatsCount++;
    }

    function _removeParticipant(address participantAddress) private {
        // 1. Remove & Revoke Participant Role
        delete participantToStatus[participantAddress];
        _removeParticipantFromList(participantAddress);
        revokeParticipantRole(participantAddress);

        // 2. Current SeatsCount Decrement
        currentSeatsCount--;
    }

    function _removeParticipantFromList(address participantAddress) private {
        uint256 index = participantIndex[participantAddress];
        uint256 lastIndex = participants.length - 1;

        if (index != lastIndex) {
            // 1. Swap the participant with the last participant
            address lastParticipant = participants[lastIndex];
            participants[index] = lastParticipant;
            participantIndex[lastParticipant] = index;
        }

        // 2. Remove the participant from the list using Pop
        participants.pop();
        delete isParticipant[participantAddress];
        delete participantIndex[participantAddress];
    }

    function _validateEventName(string memory _eventName) private pure {
        if (bytes(_eventName).length == 0) {
            revert Event__InvalidEventName();
        }
    }

    // Getters
    function getEventName() external view returns (string memory) {
        return eventName;
    }

    function getEventDescription() external view returns (string memory) {
        return eventDescription;
    }
}
