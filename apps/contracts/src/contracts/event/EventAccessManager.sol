// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {DecmAccessManager} from "../decm/DecmAccessManager.sol";
import {Constants} from "../constants/Constants.s.sol";
import {ThemisUtils} from "../../utils/ThemisUtils.sol";

contract EventAccessManager is AccessControl, ThemisUtils {
    using Constants for *;

    // Contracts
    DecmAccessManager public immutable DECM_ACCESS_MANAGER;

    // Errors
    error EventAccessManager__AccessManagerCannotBeZeroAddress();
    error EventAccessManager__AccountCannotBeZeroAddress();
    error EventAccessManager__NotHostOrAdmin();
    error EventAccessManager__NotParticipant();
    error EventAccessManager__NotHostOrAdminOrParticipant();

    constructor(address decmAccessManagerAddr, address hostAddress) {
        if (decmAccessManagerAddr == address(0)) {
            revert EventAccessManager__AccessManagerCannotBeZeroAddress();
        }
        DECM_ACCESS_MANAGER = DecmAccessManager(decmAccessManagerAddr);

        _grantRole(Constants.HOST_ROLE, hostAddress);
    }

    function grantIssuerRole(address issuer, address signer) public {
        requireHostOrAdmin(signer);

        if (issuer == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _grantRole(Constants.ISSUER_ROLE, issuer);
    }

    function revokeIssuerRole(address issuer, address signer) public {
        requireHostOrAdmin(signer);

        if (issuer == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _revokeRole(Constants.ISSUER_ROLE, issuer);
    }

    function grantParticipantRole(address participant, address signer) public {
        requireHostOrAdmin(signer);

        if (participant == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _grantRole(Constants.PARTICIPANT_ROLE, participant);
    }

    function revokeParticipantRole(address participant, address signer) public {
        requireHostOrAdminOrParticipant(signer);

        if (participant == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _revokeRole(Constants.PARTICIPANT_ROLE, participant);
    }

    function grantHostRole(address host, address signer) public {
        requireHostOrAdmin(signer);

        if (host == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        _grantRole(Constants.HOST_ROLE, host);
    }

    function addAllowedMsgSender(address sender) public {
        requireAdmin(msg.sender);
        
        if (sender == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }

        allowedMsgSenders[sender] = true;
        emit MsgSenderAllowed(sender, msg.sender);
    }

    function removeAllowedMsgSender(address sender) public {
        requireAdmin(msg.sender);
        
        if (sender == address(0)) {
            revert EventAccessManager__AccountCannotBeZeroAddress();
        }
        
        allowedMsgSenders[sender] = false;
        emit MsgSenderDisallowed(sender, msg.sender);
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

    function requireHostOrAdmin(address addr) public view {
        if (!checkIsHostOrAdmin(addr)) {
            revert EventAccessManager__NotHostOrAdmin();
        }
    }

    function requireHostOrAdminOrParticipant(address addr) public view {
        bool hasHostRole = checkIsHost(addr);
        bool hasParticipantRole = checkIsParticipant(addr);
         bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(addr);
        if (!hasHostRole && !hasAdminRole && !hasParticipantRole) {
            revert EventAccessManager__NotHostOrAdminOrParticipant();
        }
    }

    function requireParticipant(address addr) public view {
        if (!checkIsParticipant(addr)) {
            revert EventAccessManager__NotParticipant();
        }
    }
}
