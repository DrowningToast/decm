// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EventAccessManager} from "./EventAccessManager.sol";
import {ThemisUtils} from "../../utils/ThemisUtils.sol";

contract Event is ThemisUtils {
    // Contracts
        EventAccessManager public immutable EVENT_ACCESS_MANAGER;

    // Enums
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
    error Event__InvalidSignature();
    error Event__AccessManagerCannotBeZeroAddress();

    event RemovedParticipant(address indexed participant);
    event AddedParticipant(address indexed participant);
    event ParticipantSigned(address indexed participant, bytes32 signature);
    event EventConfirmed();
    event EventUpdated(
        string eventName,
        string eventDescription,
        uint256 seatsCount,
        EventStatus eventStatus
    );
    event SignatureUsed(
        address transactor,
        address signer,
        address contractAddress,
        string functionName,
        string signedMessageDigest,
        bytes signature,
        uint256 timestamp
    );

    // State Variables
    string public eventName;
    string public eventDescription;
    uint256 public seatsCount;
    uint256 public currentSeatsCount;
    EventStatus public eventStatus;

    // Participants Mappings
    mapping(address => bool) private isParticipant;
    mapping(address => uint256) private participantIndex;
    address[] private participants;

    constructor(
        address eventAccessManagerAddr,
        string memory _eventName,
        string memory _eventDescription,
        uint256 _seatsCount
    ) {
        if (eventAccessManagerAddr == address(0)) {
            revert Event__AccessManagerCannotBeZeroAddress();
        }

        EVENT_ACCESS_MANAGER = EventAccessManager(eventAccessManagerAddr);

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
        EventStatus _eventStatus,
        string memory signedMessageDigest,
        bytes memory signature
    ) external {
        address signer = recoverSigner(signedMessageDigest, signature); 
        EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer); 

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
            _eventStatus
        );

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "updateEvent",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }

    function addParticipant(
        address participantAddress,
        string memory signedMessageDigest,
        bytes memory signature
    ) external {
        address signer = recoverSigner(signedMessageDigest, signature);
        EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer);

        // Pre Conditions
        if (participantAddress == address(0)) {
            revert Event__AddressCannotBeZero();
        }

        if (currentSeatsCount >= seatsCount) {
            revert Event__SeatsCountReached();
        }

        // 1. Validate Participant
        if (isParticipant[participantAddress]) {
            revert Event__ParticipantIsAlreadyJoined();
        }

        // 2. Add Participant
        _addParticipant(participantAddress, signer);

        // 3. Emit Event
        emit AddedParticipant(participantAddress);

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "addParticipant",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }

    function leaveEvent(
        string memory signedMessageDigest,
        bytes memory signature
    ) external {
        address signer = recoverSigner(signedMessageDigest, signature);
        EVENT_ACCESS_MANAGER.requireParticipant(signer);

        address participantAddress = signer;

        // 1. Validate Participant
        if (!isParticipant[participantAddress]) {
            revert Event__ParticipantIsNotJoined();
        }

        // 2. Remove Participant
        _removeParticipant(participantAddress, signer);

        // 3. Emit Event
        emit RemovedParticipant(participantAddress);

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "leaveEvent",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }

    function removeParticipant(address participantAddress, string memory signedMessageDigest, bytes memory signature) external {
        address signer = recoverSigner(signedMessageDigest, signature);
        EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer);

        // Pre Conditions
        if (participantAddress == address(0)) {
            revert Event__AddressCannotBeZero();
        }

        // 1. Validate Participant
        if (!isParticipant[participantAddress]) {
            revert Event__ParticipantIsNotJoined();
        }
        
        // 2. Remove Participant
        _removeParticipant(participantAddress, signer);

        // 3. Emit Event
        emit RemovedParticipant(participantAddress);
    }

    function confirmEvent(
        string memory signedMessageDigest,
        bytes memory signature
    ) external {
        address signer = recoverSigner(signedMessageDigest, signature);
        EVENT_ACCESS_MANAGER.requireHostOrAdmin(signer);

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

        emit SignatureUsed(
            msg.sender,
            signer,
            address(this),
            "confirmEvent",
            signedMessageDigest,
            signature,
            block.timestamp
        );
    }

    function _addParticipant(address participantAddress, address signer) private {
        // 1. Add Participant
        participantIndex[participantAddress] = participants.length;
        participants.push(participantAddress);
        isParticipant[participantAddress] = true;

        // 2. Grant Participant Role
        EVENT_ACCESS_MANAGER.grantParticipantRole(participantAddress, signer);

        // 3. Current SeatsCount Increment
        currentSeatsCount++;
    }

    function _removeParticipant(address participantAddress, address signer) private {
        // 1. Remove Participant 
        _removeParticipantFromList(participantAddress);

        // 2. Revoke Participant Role
        EVENT_ACCESS_MANAGER.revokeParticipantRole(participantAddress, signer);

        // 3. Current SeatsCount Decrement
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
            // revert Event__InvalidEventName();
            require(false, "Invalid event name");
        }
    }

    // Getters
    function getEventName() external view returns (string memory) {
        return eventName;
    }

    function getEventDescription() external view returns (string memory) {
        return eventDescription;
    }

    function getParticipants() external view returns (address[] memory) {
        return participants;
    }
}
