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
        grantHostRole(msg.sender);
        
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

    function grantHostRole(address host) public onlyHostOrAdmin {
        if (host == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _grantRole(Constants.HOST_ROLE, host);
    }

    function checkIsHost(address addr) public view returns (bool) {
        return hasRole(Constants.HOST_ROLE, addr);
    }

    function checkIsIssuer(address addr) public view returns (bool) {
        return hasRole(Constants.ISSUER_ROLE, addr);
    }
    
    function checkIsParticipant(address addr) public view returns (bool) {
        return hasRole(Constants.PARTICIPANT_ROLE, addr);
    }


    function checkIsHostOrAdmin(address addr) public view returns (bool) {
        return
            checkIsHost(addr) || DECM_ACCESS_MANAGER.checkIsAdmin(addr);
    }

    // Modifier
    modifier onlyHostOrAdmin() {
        bool hasHostRole = checkIsHost(msg.sender);
        bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(msg.sender);
        if (!hasHostRole && !hasAdminRole) {
            revert EventAccessManager__NotHostOrAdmin();
        }
        _;
    }

    modifier onlyParticipant() {
        if (!checkIsParticipant(msg.sender)) {
            revert EventAccessManager__NotParticipant();
        }
        _;
    }

    modifier onlyHostOrAdminOrParticipant() {
        bool hasHostRole = checkIsHost(msg.sender);
        bool hasParticipantRole = checkIsParticipant(msg.sender);
         bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(msg.sender);
        if (!hasHostRole && !hasAdminRole && !hasParticipantRole) {
            revert EventAccessManager__NotHostOrAdminOrParticipant();
        }
        _;
    }
}
