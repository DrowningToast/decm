// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IDecmAccessManager {
    function checkIsAdmin(address addr) external view returns (bool);

    function grantAdminRole(address admin) external;

    function revokeAdminRole(address admin) external;
}
