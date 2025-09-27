// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {EventAccessManager} from "./EventAccessManager.sol";

contract Event is EventAccessManager {
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
    error Event__EventIsInactive();
    error Event__EventIsClosed();
    error Event__SeatsCountReached();
    error Event__InvalidEventName();
    error Event__CannotReduceSeatsCount();
    error Event__ParticipantIsNotApprovedOrJoined();
    error Event__SeatCountInvariantViolated(uint256 currentSeatsCount, uint256 participantsLength);

    // Events
    event ParticipantLeftEvent(address indexed participant);
    event EventConfirmed();
    event EventStatusUpdated(EventStatus indexed eventStatus);
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

    // Mappings
    mapping(address => ParticipantStatus) participantToStatus;
    
    // Arrays for enumeration
    address[] private participants;
    mapping(address => bool) private isParticipant;
    mapping(address => uint256) private participantIndex;

    constructor(
        address decmAccessManagerAddr,
        address[] memory initialIssuers,
        string memory _eventName,
        string memory _eventDescription,
        uint256 _seatsCount
    ) EventAccessManager(decmAccessManagerAddr) {
        _validateEventName(_eventName);

        for (uint256 i = 0; i < initialIssuers.length; i++) {
            if (initialIssuers[i] == address(0)) {
                revert EventAccessManager__AccountCannotBeZeroAddress();
            }
            _grantRole(ISSUER_ROLE, initialIssuers[i]);
        }

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
        _validateEventName(_eventName);

        if (_seatsCount < seatsCount) {
            revert Event__CannotReduceSeatsCount();
        }

        eventName = _eventName;
        eventDescription = _eventDescription;
        seatsCount = _seatsCount;
        eventStatus = _eventStatus;
        
        emit EventUpdated(
            _eventName,
            _eventDescription,
            _seatsCount,
            eventStatus
        );
    }

    function leaveEvent() external onlyParticipant {
        _removeApprovedParticipant(msg.sender);
    }

    function removeParticipant(address participant) external onlyHostOrAdmin {
        _removeApprovedParticipant(participant);
    }

    function confirmEvent() public onlyHostOrAdmin {
        if (eventStatus == EventStatus.INACTIVE) {
            revert Event__EventIsInactive();
        }

        if (eventStatus == EventStatus.CLOSED) {
            revert Event__EventIsClosed();
        }

        eventStatus = EventStatus.CLOSED;
        emit EventConfirmed();
    }

    function setEventStatus(EventStatus _eventStatus) external onlyHostOrAdmin {
        eventStatus = _eventStatus;
        emit EventStatusUpdated(_eventStatus);
    }
    
    function getEventName() external view returns (string memory) {
        return eventName;
    }

    function getEventDescription() external view returns (string memory) {
        return eventDescription;
    }

    function getEventSeatsCount() external view returns (uint256) {
        return seatsCount;
    }

    function getEventStatus() external view returns (EventStatus) {
        return eventStatus;
    }

    function getParticipantStatus(
        address participant
    ) external view returns (ParticipantStatus) {
        return participantToStatus[participant];
    }

    function getParticipants() external view returns (address[] memory) {
        return participants;
    }
    
    function getParticipantsCount() external view returns (uint256) {
        return participants.length;
    }
    
    function getParticipantAtIndex(uint256 index) external view returns (address) {
        require(index < participants.length, "Index out of bounds");
        return participants[index];
    }
    
    function _validateEventName(string memory _eventName) private pure {
        if (bytes(_eventName).length == 0) {
            revert Event__InvalidEventName();
        }
    }

    function _addParticipant(address participant) private {
        if (!isParticipant[participant]) {
            participantIndex[participant] = participants.length;
            participants.push(participant);
            isParticipant[participant] = true;
        }
    }
    
    function _removeParticipant(address participant) private {
        if (isParticipant[participant]) {
            uint256 index = participantIndex[participant];
            uint256 lastIndex = participants.length - 1;

            if (index != lastIndex) {
                address lastParticipant = participants[lastIndex];
                participants[index] = lastParticipant;
                participantIndex[lastParticipant] = index;
            }

            participants.pop();
            delete isParticipant[participant];
            delete participantIndex[participant];
        }
    }

    function _removeApprovedParticipant(address participant) private {
        if (participantToStatus[participant] != ParticipantStatus.APPROVED) {
            revert Event__ParticipantIsNotApprovedOrJoined();
        }

        delete participantToStatus[participant];
        _removeParticipant(participant);
        revokeParticipantRole(participant);
        _decrementSeatCount();

        emit ParticipantLeftEvent(participant);
    }

    function _decrementSeatCount() private {
        if (currentSeatsCount == 0) {
            revert Event__SeatCountInvariantViolated(currentSeatsCount, participants.length);
        }

        unchecked {
            currentSeatsCount--;
        }

        if (currentSeatsCount != participants.length) {
            revert Event__SeatCountInvariantViolated(currentSeatsCount, participants.length);
        }
    }
}
