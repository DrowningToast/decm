// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {IEventAccessManager} from "../../interfaces/IEventAccessManager.sol";
import {DecmAccessManager} from "../decm/DecmAccessManager.sol";
import {Constants} from "../constants/Constants.s.sol";

contract EventAccessManager is AccessControl, IEventAccessManager {
    using Constants for *;

    // Contracts
    DecmAccessManager public immutable DECM_ACCESS_MANAGER;

    // Errors
    error EventAccessManager__IssuerCannotBeZeroAddress();
    error EventAccessManager__NotHostOrAdmin();
    error EventAccessManager__NotParticipant();

    // Roles
    bytes32 public constant HOST_ROLE = Constants.HOST_ROLE;
    bytes32 public constant ISSUER_ROLE = Constants.ISSUER_ROLE;
    bytes32 public constant PARTICIPANT_ROLE = Constants.PARTICIPANT_ROLE;

    constructor(address decmAccessManagerAddr) {
        DECM_ACCESS_MANAGER = DecmAccessManager(decmAccessManagerAddr);
        _grantRole(HOST_ROLE, msg.sender);
    }

    function grantIssuerRole(address issuer) external onlyRole(HOST_ROLE) {
        _grantRole(ISSUER_ROLE, issuer);
    }

    function revokeIssuerRole(address issuer) external onlyRole(HOST_ROLE) {
        _revokeRole(ISSUER_ROLE, issuer);
    }

    function grantParticipantRole(
        address participant
    ) public onlyRole(HOST_ROLE) {
        _grantRole(PARTICIPANT_ROLE, participant);
    }

    function revokeParticipantRole(
        address participant
    ) public onlyRole(HOST_ROLE) {
        _revokeRole(PARTICIPANT_ROLE, participant);
    }

    function checkIsHost(address addr) external view returns (bool) {
        return hasRole(HOST_ROLE, addr);
    }

    function checkIsIssuer(address addr) external view returns (bool) {
        return hasRole(ISSUER_ROLE, addr);
    }

    function checkIsParticipant(address addr) external view returns (bool) {
        return hasRole(PARTICIPANT_ROLE, addr);
    }

    function checkIsHostOrAdmin(address addr) external view returns (bool) {
        return
            hasRole(HOST_ROLE, addr) || DECM_ACCESS_MANAGER.checkIsAdmin(addr);
    }

    // Modifier
    modifier onlyHostOrAdmin() {
        bool hasHostRole = hasRole(HOST_ROLE, msg.sender);
        bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(msg.sender);
        if (!hasHostRole && !hasAdminRole) {
            revert EventAccessManager__NotHostOrAdmin();
        }
        _;
    }

    modifier onlyParticipant() {
        if (!hasRole(PARTICIPANT_ROLE, msg.sender)) {
            revert EventAccessManager__NotParticipant();
        }
        _;
    }
}
