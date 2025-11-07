// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {Event} from "../src/contracts/event/Event.sol";
import {EventAccessManager} from "../src/contracts/event/EventAccessManager.sol";

contract DeployEvent is Script {
    // Errors
    error DeployEvent__DecmAccessManagerCannotBeZeroAddress();

    function run() public returns (Event) {
        address decmAccessManagerAddress = vm.envAddress("DECM_ACCESS_MANAGER_ADDRESS");
        address hostAddress = vm.envAddress("HOST_ADDRESS");

        string memory eventName = "TEST_EVENT";
        string memory eventDescription = "TEST_EVENT_DESCRIPTION";
        uint256 seatsCount = 100;

        if (decmAccessManagerAddress == address(0)) {
            revert DeployEvent__DecmAccessManagerCannotBeZeroAddress();
        }
 
        vm.startBroadcast();
        EventAccessManager eventAccessManager = new EventAccessManager(decmAccessManagerAddress, hostAddress);
        address eventAccessManagerAddress = address(eventAccessManager);
        Event eventContract = new Event(
            eventAccessManagerAddress,
            eventName,
            eventDescription,
            seatsCount
        );
        vm.stopBroadcast();

        return eventContract;
    }
}