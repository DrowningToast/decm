// // SPDX-License-Identifier: MIT
// pragma solidity ^0.8.20;

// import { AccessControl } from "@openzeppelin/contracts/access/AccessControl.sol"; 
// import { DecmAccessManager } from "../../src/contracts/decm/DecmAccessManager.sol";
// import { Constants } from "../../src/contracts/constants/Constants.s.sol";
// import { ThemisUtils } from "../../src/utils/ThemisUtils.sol";
// import { console } from "forge-std/console.sol";

// contract TestEventAccessManager is AccessControl, ThemisUtils {
//     using Constants for *;

//     // Contracts
//     DecmAccessManager public immutable DECM_ACCESS_MANAGER;

//     // Events
//     event IssuerRoleGranted(address indexed issuer, address indexed granter);
//     event IssuerRoleRevoked(address indexed issuer, address indexed revoker);
//     event ParticipantRoleGranted(address indexed participant, address indexed granter);
//     event ParticipantRoleRevoked(address indexed participant, address indexed revoker);
//     event HostRoleGranted(address indexed host, address indexed granter);
//     event MsgSenderAllowed(address indexed sender, address indexed granter);
//     event MsgSenderDisallowed(address indexed sender, address indexed revoker);

//     // Errors
//     error EventAccessManager__AccessManagerCannotBeZeroAddress();
//     error EventAccessManager__AccountCannotBeZeroAddress();
//     error EventAccessManager__NotHostOrAdmin();
//     error EventAccessManager__NotParticipant();
//     error EventAccessManager__NotHostOrAdminOrParticipant();
//     error EventAccessManager__NotAllowedMsgSender();
//     error EventAccessManager__NotAdmin();
    
//     // States
//     mapping(address => bool) public allowedMsgSenders;

//     constructor(address decmAccessManagerAddr, address hostAddress) {
//         if (decmAccessManagerAddr == address(0)) {
//             revert EventAccessManager__AccessManagerCannotBeZeroAddress();
//         }

//         DECM_ACCESS_MANAGER = DecmAccessManager(decmAccessManagerAddr);

//         allowedMsgSenders[msg.sender] = true;
//         _grantRole(Constants.HOST_ROLE, hostAddress);

//         emit MsgSenderAllowed(msg.sender, msg.sender);
//         emit HostRoleGranted(hostAddress, msg.sender);
//     }

//     // External wrapper functions for testing
//     function grantIssuerRole(address issuer, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _grantIssuerRole(issuer, signer);
//     }

//     function revokeIssuerRole(address issuer, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _revokeIssuerRole(issuer, signer);
//     }

//     function grantParticipantRole(address participant, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _grantParticipantRole(participant, signer);
//     }

//     function revokeParticipantRole(address participant, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _revokeParticipantRole(participant, signer);
//     }

//     function grantHostRole(address host, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _grantHostRole(host, signer);
//     }

//     function addAllowedMsgSender(address sender, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _addAllowedMsgSender(sender, signer);
//     }

//     function removeAllowedMsgSender(address sender, string memory signMessage, bytes memory signature) external {
//         address signer = recoverSigner(signMessage, signature, address(this));
//         _removeAllowedMsgSender(sender, signer);
//     }

//     // Internal functions (copied from original)
//     function _grantIssuerRole(address issuer, address signer) internal {
//         requireHostOrAdmin(signer);

//         if (issuer == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }
//         _grantRole(Constants.ISSUER_ROLE, issuer);
//         emit IssuerRoleGranted(issuer, msg.sender);
//     }

//     function _revokeIssuerRole(address issuer, address signer) internal {
//         requireHostOrAdmin(signer);

//         if (issuer == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }
//         _revokeRole(Constants.ISSUER_ROLE, issuer);
//         emit IssuerRoleRevoked(issuer, msg.sender);
//     }

//     function _grantParticipantRole(address participant, address signer) internal {
//         requireHostOrAdmin(signer);

//         if (participant == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }
//         _grantRole(Constants.PARTICIPANT_ROLE, participant);
//         emit ParticipantRoleGranted(participant, msg.sender);
//     }

//     function _revokeParticipantRole(address participant, address signer) internal {
//         requireHostOrAdminOrParticipant(signer);

//         if (participant == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }
//         _revokeRole(Constants.PARTICIPANT_ROLE, participant);
//         emit ParticipantRoleRevoked(participant, msg.sender);
//     }

//     function _grantHostRole(address host, address signer) internal {
//         requireHostOrAdmin(signer);

//         if (host == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }
//         _grantRole(Constants.HOST_ROLE, host);
//         emit HostRoleGranted(host, msg.sender);
//     }

//     function _addAllowedMsgSender(address sender, address signer) internal {
//         requireAdmin(signer);

//         if (sender == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }

//         allowedMsgSenders[sender] = true;
//         emit MsgSenderAllowed(sender, msg.sender);
//     }

//     function _removeAllowedMsgSender(address sender, address signer) internal {
//         requireAdmin(signer);

//         if (sender == address(0)) {
//             revert EventAccessManager__AccountCannotBeZeroAddress();
//         }

//         allowedMsgSenders[sender] = false;
//         emit MsgSenderDisallowed(sender, msg.sender);
//     }

//     // View functions
//     function checkIsHost(address addr) public view returns (bool) {
//         return hasRole(Constants.HOST_ROLE, addr);
//     }

//     function checkIsIssuer(address addr) public view returns (bool) {
//         return hasRole(Constants.ISSUER_ROLE, addr);
//     }

//     function checkIsParticipant(address addr) public view returns (bool) {
//         return hasRole(Constants.PARTICIPANT_ROLE, addr);
//     }

//     function checkIsHostOrAdmin(address addr) public view returns (bool) {
//         return checkIsHost(addr) || DECM_ACCESS_MANAGER.checkIsAdmin(addr);
//     }

//     function checkIsAllowedMsgSender() public view returns (bool) {
//         return allowedMsgSenders[msg.sender];
//     }

//     function requireAllowedMsgSender() public view {
//         if (!checkIsAllowedMsgSender()) {
//             revert EventAccessManager__NotAllowedMsgSender();
//         }
//     }

//     function requireAdmin(address addr) public view {
//         if (!DECM_ACCESS_MANAGER.checkIsAdmin(addr)) {
//             revert EventAccessManager__NotAdmin();
//         }
//     }

//     function requireHostOrAdmin(address addr) public view {
//         console.log("addr", addr);

//         if (!checkIsHostOrAdmin(addr)) {
//             revert EventAccessManager__NotHostOrAdmin();
//         }
//     }

//     function requireHostOrAdminOrParticipant(address addr) public view {
//         bool hasHostRole = checkIsHost(addr);
//         bool hasParticipantRole = checkIsParticipant(addr);
//         bool hasAdminRole = DECM_ACCESS_MANAGER.checkIsAdmin(addr);
//         bool isAllowed = checkIsAllowedMsgSender();
//         if (!hasHostRole && !hasAdminRole && !hasParticipantRole && !isAllowed) {
//             revert EventAccessManager__NotHostOrAdminOrParticipant();
//         }
//     }

//     function requireParticipant(address addr) public view {
//         bool hasParticipantRole = checkIsParticipant(addr);
//         bool isAllowed = checkIsAllowedMsgSender();
//         if (!hasParticipantRole && !isAllowed) {
//             revert EventAccessManager__NotParticipant();
//         }
//     }
// }
