// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import { AccessControl } from "@openzeppelin/contracts/access/AccessControl.sol"; 
import { DecmAccessManager } from "../decm/DecmAccessManager.sol";
import { Constants } from "../constants/Constants.s.sol";
import { ThemisUtils } from "../../utils/ThemisUtils.sol";

contract EventAccessManager is AccessControl, ThemisUtils {
    using Constants for *;

    // Contracts
    DecmAccessManager public immutable DECM_ACCESS_MANAGER;

    // Events
    event IssuerRoleGranted(address indexed issuer, address indexed granter);
    event IssuerRoleRevoked(address indexed issuer, address indexed revoker);
    event ParticipantRoleGranted(address indexed participant, address indexed granter);
    event ParticipantRoleRevoked(address indexed participant, address indexed revoker);
    event HostRoleGranted(address indexed host, address indexed granter);
    event MsgSenderAllowed(address indexed sender, address indexed granter);
    event MsgSenderDisallowed(address indexed sender, address indexed revoker);

    constructor(address decmAccessManagerAddr, address hostAddress) {
        if (decmAccessManagerAddr == address(0)) {
            require(false, "Access manager cannot be zero address");
        }

        if (hostAddress == address(0)) {
            require(false, "Account cannot be zero");
        }

        DECM_ACCESS_MANAGER = DecmAccessManager(decmAccessManagerAddr);

        _grantRole(Constants.HOST_ROLE, hostAddress);

        emit MsgSenderAllowed(msg.sender, msg.sender);
        emit HostRoleGranted(hostAddress, msg.sender);
    }

    function grantIssuerRole(address issuer, address signer) public {
        requireHostOrAdmin(signer, msg.sender);

        if (issuer == address(0)) {
            require(false, "Account cannot be zero");
        }
        _grantRole(Constants.ISSUER_ROLE, issuer);
        emit IssuerRoleGranted(issuer, msg.sender);
    }

    function revokeIssuerRole(address issuer, address signer) public {
        requireHostOrAdmin(signer, msg.sender);

        if (issuer == address(0)) {
            require(false, "Account cannot be zero");
        }
        _revokeRole(Constants.ISSUER_ROLE, issuer);
        emit IssuerRoleRevoked(issuer, msg.sender);
    }

    function grantParticipantRole(address participant, address signer) public {
        requireHostOrAdmin(signer, msg.sender);

        if (participant == address(0)) {
            require(false, "Account cannot be zero");
        }

        _grantRole(Constants.PARTICIPANT_ROLE, participant);
        emit ParticipantRoleGranted(participant, msg.sender);
    }

    function revokeParticipantRole(address participant, address signer) public {
        requireHostOrAdminOrParticipant(signer, msg.sender);

        if (participant == address(0)) {
            require(false, "Account cannot be zero");
        }

        _revokeRole(Constants.PARTICIPANT_ROLE, participant);
        emit ParticipantRoleRevoked(participant, msg.sender);
    }

    function grantParticipantRoleUsingAllowedMsgSender(address participant, address msgSender) public {
        requireAllowedMsgSender(msgSender);

        if (participant == address(0)) {
            require(false, "Account cannot be zero");
        }

        _grantRole(Constants.PARTICIPANT_ROLE, participant);
        emit ParticipantRoleGranted(participant, msgSender);
    }

    function grantHostRole(address host, address signer) public {
        requireHostOrAdmin(signer, msg.sender);

        if (host == address(0)) {
            require(false, "Account cannot be zero");
        }

        _grantRole(Constants.HOST_ROLE, host);
        emit HostRoleGranted(host, msg.sender);
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
        return checkIsHost(addr) || DECM_ACCESS_MANAGER.checkIsAdmin(addr);
    }

    function checkIsAllowedMsgSender(address addr) public view returns (bool) {
        return DECM_ACCESS_MANAGER.checkIsAllowedMsgSender(addr);
    }

    function requireAllowedMsgSender(address addr) public view {
        if (!checkIsAllowedMsgSender(addr)) {
            require(false, "Not allowed msg sender");
        }
    }

    function requireHostOrAdmin(address addr, address msgSender) public view {
        bool isHostOrAdmin = checkIsHostOrAdmin(addr);
        bool isAllowedMsgSender = checkIsAllowedMsgSender(msgSender);
        if (!isHostOrAdmin && !isAllowedMsgSender) {
            require(false, "Not host or admin or allowed msg sender");
        }
    }

    function requireAdmin(address addr, address msgSender) public view {
        bool isAdmin = DECM_ACCESS_MANAGER.checkIsAdmin(addr);
        bool isAllowedMsgSender = checkIsAllowedMsgSender(msgSender);
        if (!isAdmin && !isAllowedMsgSender) {
            require(false, "Not admin or allowed msg sender");
        }
    }

    function requireHostOrAdminOrParticipant(address addr, address msgSender) public view {
        bool hasHostRole = checkIsHost(addr);
        bool hasParticipantRole = checkIsParticipant(addr);
        bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(addr);
        bool isAllowedMsgSender = checkIsAllowedMsgSender(msgSender);
        if (!hasHostRole && !hasAdminRole && !hasParticipantRole && !isAllowedMsgSender) {
            require(false, "Not host or admin or participant or allowed msg sender");
        }
    }

    function requireParticipant(address addr, address msgSender) public view {
        bool hasParticipantRole = checkIsParticipant(addr);
        bool isAllowedMsgSender = checkIsAllowedMsgSender(msgSender);
        if (!hasParticipantRole && !isAllowedMsgSender) {
            require(false, "Not participant or allowed msg sender");
        }
    }
}
