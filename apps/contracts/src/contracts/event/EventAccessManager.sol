// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {DecmAccessManager} from "../decm/DecmAccessManager.sol";
import {Constants} from "../constants/Constants.s.sol";

contract EventAccessManager is AccessControl {
    using Constants for *;

    // Contracts
    DecmAccessManager public immutable DECM_ACCESS_MANAGER;

    // Errors
    error EventAccessManager__AccessManagerCannotBeZeroAddress();
    error EventAccessManager__AccountCannotBeZeroAddress();
    error EventAccessManager__NotHostOrAdmin();
    error EventAccessManager__NotParticipant();
    error EventAccessManager__NotHostOrAdminOrParticipant();

    constructor(address decmAccessManagerAddr) {
        if (decmAccessManagerAddr == address(0)) {
            revert EventAccessManager__AccessManagerCannotBeZeroAddress();
        }
        DECM_ACCESS_MANAGER = DecmAccessManager(decmAccessManagerAddr);
        _grantRole(Constants.HOST_ROLE, msg.sender);
    }

    function grantIssuerRole(address issuer) external onlyHostOrAdmin {
        if (issuer == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _grantRole(Constants.ISSUER_ROLE, issuer);
    }

    function revokeIssuerRole(address issuer) external onlyHostOrAdmin {
        if (issuer == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _revokeRole(Constants.ISSUER_ROLE, issuer);
    }

    function grantParticipantRole(address participant) public onlyHostOrAdmin {
        if (participant == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _grantRole(Constants.PARTICIPANT_ROLE, participant);
    }

    function revokeParticipantRole(address participant) public onlyHostOrAdminOrParticipant {
        if (participant == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _revokeRole(Constants.PARTICIPANT_ROLE, participant);
    }

    function checkIsHost(address addr) external view returns (bool) {
        return DECM_ACCESS_MANAGER.checkIsHost(addr);
    }

    function checkIsIssuer(address addr) external view returns (bool) {
        return DECM_ACCESS_MANAGER.checkIsIssuer(addr);
    }

    function checkIsParticipant(address addr) external view returns (bool) {
        return DECM_ACCESS_MANAGER.checkIsParticipant(addr);
    }

    function checkIsHostOrAdmin(address addr) external view returns (bool) {
        return
            DECM_ACCESS_MANAGER.checkIsHost(addr) || DECM_ACCESS_MANAGER.checkIsAdmin(addr);
    }

    // Modifier
    modifier onlyHostOrAdmin() {
        bool hasHostRole = DECM_ACCESS_MANAGER.checkIsHost(msg.sender);
        bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(msg.sender);
        if (!hasHostRole && !hasAdminRole) {
            revert EventAccessManager__NotHostOrAdmin();
        }
        _;
    }

    modifier onlyParticipant() {
        if (!DECM_ACCESS_MANAGER.checkIsParticipant(msg.sender)) {
            revert EventAccessManager__NotParticipant();
        }
        _;
    }

    modifier onlyHostOrAdminOrParticipant() {
        bool hasHostRole = DECM_ACCESS_MANAGER.checkIsHost(msg.sender);
        bool hasParticipantRole = DECM_ACCESS_MANAGER.checkIsParticipant(msg.sender);
         bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(msg.sender);
        if (!hasHostRole && !hasAdminRole && !hasParticipantRole) {
            revert EventAccessManager__NotHostOrAdminOrParticipant();
        }
        _;
    }
}
