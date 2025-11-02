// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {DecmAccessManager} from "../src/contracts/decm/DecmAccessManager.sol";

contract DeployDecmAccessManager is Script {
    function run() public returns (DecmAccessManager) {
        address initialAdminAddress = vm.envAddress("INITIAL_ADMIN_ADDRESS");
        address[] memory initialAdmins = new address[](1);
        initialAdmins[0] = initialAdminAddress;

        vm.startBroadcast();
        DecmAccessManager decmAccessManager = new DecmAccessManager(initialAdmins);
        vm.stopBroadcast();

        return decmAccessManager;
    }
}