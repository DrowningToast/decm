// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {IDecmAccessManager} from "../../interfaces/IDecmAccessManager.sol";

contract DecmAccessManager is AccessControl, IDecmAccessManager {
    // Errors
    error DecmAccessManager__AdminCannotBeZeroAddress();

    // Roles
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
        _grantRole(ADMIN_ROLE, admin);
    }

    function revokeAdminRole(
        address admin
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        _revokeRole(ADMIN_ROLE, admin);
    }

    function checkIsAdmin(address addr) external view returns (bool) {
        return hasRole(ADMIN_ROLE, addr);
    }
}
