// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {EventCertificate} from "../src/contracts/event/EventCertificate.sol";

contract DeployEventCertificate is Script {
    // Errors
    error DeployEventCertificate__DecmAccessManagerOrEventCannotBeZeroAddress();

    function run() public returns (EventCertificate) {
        address decmAccessManagerAddress = vm.envAddress("DECM_ACCESS_MANAGER_ADDRESS");
        address eventAddress = vm.envAddress("EVENT_ADDRESS");

        if (decmAccessManagerAddress == address(0) || eventAddress == address(0)) {
            revert DeployEventCertificate__DecmAccessManagerOrEventCannotBeZeroAddress();
        }

        vm.startBroadcast();
        EventCertificate eventCertificate = new EventCertificate(
            decmAccessManagerAddress,
            eventAddress
        );
        vm.stopBroadcast();

        return eventCertificate;
    }
}