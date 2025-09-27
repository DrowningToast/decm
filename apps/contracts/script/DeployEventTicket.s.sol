// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {EventTicket} from "../src/contracts/event/EventTicket.sol";

contract DeployEventTicket is Script {
    // Errors
    error DeployEventTicket__DecmAccessManagerOrEventCannotBeZeroAddress();

    function run() public returns (EventTicket) {
        address decmAccessManagerAddress = vm.envAddress("DECM_ACCESS_MANAGER_ADDRESS");
        address eventAddress = vm.envAddress("EVENT_ADDRESS");

        if (decmAccessManagerAddress == address(0) || eventAddress == address(0)) {
            revert DeployEventTicket__DecmAccessManagerOrEventCannotBeZeroAddress();
        }

        vm.startBroadcast();
        EventTicket eventTicket = new EventTicket(
            decmAccessManagerAddress,
            eventAddress
        );
        vm.stopBroadcast();

        return eventTicket;
    }
}