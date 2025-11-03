// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Event} from "../../src/contracts/event/Event.sol";

contract MockEvent {
    // State variables to simulate Event contract
    string public eventName;
    string public eventDescription;
    uint256 public seatsCount;
    uint256 public currentSeatsCount;
    Event.EventStatus public eventStatus;

    // Participants mappings
    mapping(address => bool) public isParticipant;
    mapping(address => uint256) public participantIndex;
    address[] public participants;

    constructor(
        string memory _eventName,
        string memory _eventDescription,
        uint256 _seatsCount,
        Event.EventStatus _eventStatus
    ) {
        eventName = _eventName;
        eventDescription = _eventDescription;
        seatsCount = _seatsCount;
        currentSeatsCount = 0;
        eventStatus = _eventStatus;
    }

    // View functions that EventTicket will call
    function getEventName() external view returns (string memory) {
        return eventName;
    }

    function getEventDescription() external view returns (string memory) {
        return eventDescription;
    }

    // Helper functions for testing
    function setEventName(string memory _eventName) external {
        eventName = _eventName;
    }

    function setEventDescription(string memory _eventDescription) external {
        eventDescription = _eventDescription;
    }

    function setEventStatus(Event.EventStatus _eventStatus) external {
        eventStatus = _eventStatus;
    }

    function addParticipant(address participant) external {
        if (!isParticipant[participant]) {
            participantIndex[participant] = participants.length;
            participants.push(participant);
            isParticipant[participant] = true;
            currentSeatsCount++;
        }
    }

    function removeParticipant(address participant) external {
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
            currentSeatsCount--;
        }
    }

    function getParticipants() external view returns (address[] memory) {
        return participants;
    }
}
