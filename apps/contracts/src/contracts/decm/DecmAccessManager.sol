// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {Constants} from "../constants/Constants.s.sol";

contract DecmAccessManager is AccessControl {
    using Constants for *;

    // Events
    event AdminGranted(address indexed admin, address indexed granter);
    event AdminRevoked(address indexed admin, address indexed revoker);

    bytes32 public constant ADMIN_ROLE = DEFAULT_ADMIN_ROLE;

    // States
    mapping(address => bool) public allowedMsgSenders;

    constructor(address[] memory initialAdmins) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);

        for (uint256 i = 0; i < initialAdmins.length; i++) {
            if (initialAdmins[i] == address(0)) {
                require(false, "Admin cannot be zero address");
            }
            _grantRole(ADMIN_ROLE, initialAdmins[i]);
        }
    }

    function grantAdminRole(
        address admin
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (admin == address(0)) {
            require(false, "Admin cannot be zero address");
        }
        _grantRole(ADMIN_ROLE, admin);
        emit AdminGranted(admin, msg.sender);
    }

    function revokeAdminRole(
        address admin
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (admin == address(0)) {
            require(false, "Admin cannot be zero address");
        }
        _revokeRole(ADMIN_ROLE, admin);
        emit AdminRevoked(admin, msg.sender);
    }

    function checkIsAdmin(address addr) external view returns (bool) {
        return hasRole(ADMIN_ROLE, addr);
    }

    function addAllowedMsgSender(address msgSender) external onlyRole(DEFAULT_ADMIN_ROLE) {
        allowedMsgSenders[msgSender] = true;
    }

    function removeAllowedMsgSender(address msgSender) external onlyRole(DEFAULT_ADMIN_ROLE) {
        allowedMsgSenders[msgSender] = false;
    }

    function checkIsAllowedMsgSender(address addr) public view returns (bool) {
        return allowedMsgSenders[addr];
    } 
}
 