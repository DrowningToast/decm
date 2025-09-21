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
    error Event__EventIsNotActive();
    error Event__ParticipantIsAlreadyPending();
    error Event__ParticipantIsAlreadyApproved();
    error Event__ParticipantIsNotPending();
    error Event__ParticipantIsNotApproved();
    error Event__EventIsInactive();
    error Event__EventIsClosed();
    error Event__SeatsCountReached();
    error Event__InvalidEventName();
    error Event__CannotReduceSeatsCount();

    // Events
    event ParticipantRequestedJoinEvent(address indexed participant);
    event ParticipantApprovedJoinEvent(address indexed participant);
    event ParticipantRejectedJoinEvent(address indexed participant);
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

    constructor(
        address decmAccessManagerAddr,
        address[] memory initialIssuers,
        string memory _eventName,
        string memory _eventDescription,
        uint256 _seatsCount
    ) EventAccessManager(decmAccessManagerAddr) {
        if (bytes(_eventName).length == 0) {
            revert Event__InvalidEventName();
        }

        for (uint256 i = 0; i < initialIssuers.length; i++) {
            if (initialIssuers[i] == address(0)) {
                revert EventAccessManager__IssuerCannotBeZeroAddress();
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
        if (bytes(_eventName).length == 0) {
            revert Event__InvalidEventName();
        }

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

    function requestJoinEvent() external {
        if (eventStatus != EventStatus.ACTIVE) {
            revert Event__EventIsNotActive();
        }

        if (participantToStatus[msg.sender] == ParticipantStatus.PENDING) {
            revert Event__ParticipantIsAlreadyPending();
        }

        if (participantToStatus[msg.sender] == ParticipantStatus.APPROVED) {
            revert Event__ParticipantIsAlreadyApproved();
        }

        participantToStatus[msg.sender] = ParticipantStatus.PENDING;
        emit ParticipantRequestedJoinEvent(msg.sender);
    }

    function approveJoinRequest(address participant) external onlyHostOrAdmin {
        if (eventStatus != EventStatus.ACTIVE) {
            revert Event__EventIsNotActive();
        }

        if (participantToStatus[participant] != ParticipantStatus.PENDING) {
            revert Event__ParticipantIsNotPending();
        }

        if (currentSeatsCount >= seatsCount) {
            revert Event__SeatsCountReached();
        }

        participantToStatus[participant] = ParticipantStatus.APPROVED;
        currentSeatsCount++;
        grantParticipantRole(participant);

        emit ParticipantApprovedJoinEvent(participant);
    }

    function rejectJoinRequest(address participant) external onlyHostOrAdmin {
        if (eventStatus != EventStatus.ACTIVE) {
            revert Event__EventIsNotActive();
        }

        if (participantToStatus[participant] != ParticipantStatus.PENDING) {
            revert Event__ParticipantIsNotPending();
        }

        participantToStatus[participant] = ParticipantStatus.REJECTED;
        emit ParticipantRejectedJoinEvent(participant);
    }

    function leaveEvent() external onlyParticipant {
        if (participantToStatus[msg.sender] != ParticipantStatus.APPROVED) {
            revert Event__ParticipantIsNotApproved();
        }

        participantToStatus[msg.sender] = ParticipantStatus.LEAVED;
        revokeParticipantRole(msg.sender);

        if (currentSeatsCount > 0) {
            currentSeatsCount--;
        }

        emit ParticipantLeftEvent(msg.sender);
    }

    function removeParticipant(address participant) external onlyHostOrAdmin {
        if (participantToStatus[participant] != ParticipantStatus.APPROVED) {
            revert Event__ParticipantIsNotApproved();
        }

        participantToStatus[participant] = ParticipantStatus.LEAVED;
        revokeParticipantRole(participant);

        if (currentSeatsCount > 0) {
            currentSeatsCount--;
        }

        emit ParticipantLeftEvent(participant);
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
}
