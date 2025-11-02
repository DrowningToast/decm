// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {Constants} from "../../src/contracts/constants/Constants.s.sol";

contract MockDecmAccessManager is AccessControl {
    using Constants for *;

    // Events
    event AdminGranted(address indexed admin, address indexed granter);
    event AdminRevoked(address indexed admin, address indexed revoker);

    // Errors
    error DecmAccessManager__AdminCannotBeZeroAddress();

    bytes32 public constant ADMIN_ROLE = DEFAULT_ADMIN_ROLE;

    constructor(address[] memory initialAdmins) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);

        for (uint256 i = 0; i < initialAdmins.length; i++) {
            if (initialAdmins[i] == address(0)) {
                revert DecmAccessManager__AdminCannotBeZeroAddress();
            }
            _grantRole(ADMIN_ROLE, initialAdmins[i]);
        }
    }

    function grantAdminRole(
        address admin
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (admin == address(0)) {
            revert DecmAccessManager__AdminCannotBeZeroAddress();
        }
        _grantRole(ADMIN_ROLE, admin);
        emit AdminGranted(admin, msg.sender);
    }

    function revokeAdminRole(
        address admin
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (admin == address(0)) {
            revert DecmAccessManager__AdminCannotBeZeroAddress();
        }
        _revokeRole(ADMIN_ROLE, admin);
        emit AdminRevoked(admin, msg.sender);
    }

    function checkIsAdmin(address addr) external view returns (bool) {
        return hasRole(ADMIN_ROLE, addr);
    }
}
