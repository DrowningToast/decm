// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IEventAccessManager {
    function checkIsHost(address addr) external view returns (bool);

    function checkIsIssuer(address addr) external view returns (bool);

    function checkIsParticipant(address addr) external view returns (bool);

    function grantIssuerRole(address issuer) external;

    function revokeIssuerRole(address issuer) external;

    function grantParticipantRole(address participant) external;

    function revokeParticipantRole(address participant) external;

    function checkIsHostOrAdmin(address addr) external view returns (bool);
}
