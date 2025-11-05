// // SPDX-License-Identifier: MIT
// pragma solidity ^0.8.20;

// import {EventAccessManager} from "../src/contracts/event/EventAccessManager.sol";
// import {TestEventAccessManager} from "./utils/TestEventAccessManager.sol";
// import {MockDecmAccessManager} from "./utils/MockDecmAccessManager.sol";
// import {TestUtils} from "./utils/TestUtils.sol";

// contract EventAccessManagerTest is TestUtils {
//     TestEventAccessManager public eventAccessManager;
//     MockDecmAccessManager public mockDecmAccessManager;

//     event IssuerRoleGranted(address indexed issuer, address indexed granter);
//     event IssuerRoleRevoked(address indexed issuer, address indexed revoker);
//     event ParticipantRoleGranted(address indexed participant, address indexed granter);
//     event ParticipantRoleRevoked(address indexed participant, address indexed revoker);
//     event HostRoleGranted(address indexed host, address indexed granter);
//     event MsgSenderAllowed(address indexed sender, address indexed granter);
//     event MsgSenderDisallowed(address indexed sender, address indexed revoker);

//     function setUp() public {
//         setupMockAddresses();
        
//         // Deploy MockDecmAccessManager with admin
//         address[] memory admins = new address[](1);
//         admins[0] = ADMIN;
//         vm.prank(ADMIN);
//         mockDecmAccessManager = new MockDecmAccessManager(admins);
        
//         // Deploy EventAccessManager with host
//         vm.prank(ADMIN);
//         eventAccessManager = new TestEventAccessManager(address(mockDecmAccessManager), HOST);
//     }

//     /*//////////////////////////////////////////////////////////////
//                            CONSTRUCTOR TESTS
//     //////////////////////////////////////////////////////////////*/

//     function test_RevertWhen_DecmAccessManagerAddressIsZero() public {
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccessManagerCannotBeZeroAddress.selector);
//         new EventAccessManager(address(0), HOST);
//     }

//     function test_RevertWhen_HostAddressIsZero() public {
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccountCannotBeZeroAddress.selector);
//         new EventAccessManager(address(mockDecmAccessManager), address(0));
//     }

//     function test_WhenConstructorParametersAreValid() public view {
//         assertTrue(eventAccessManager.checkIsHost(HOST), "Host should have host role");
//         assertTrue(eventAccessManager.allowedMsgSenders(ADMIN), "Deployer should be allowed msg sender");
//     }

//     /*//////////////////////////////////////////////////////////////
//                            ROLE GRANTING TESTS
//     //////////////////////////////////////////////////////////////*/

//     function test_RevertWhen_GrantIssuerRoleWithZeroAddress() public {
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccountCannotBeZeroAddress.selector);
//         eventAccessManager.grantIssuerRole(address(0), signMessage, signature);
//     }

//     function test_RevertWhen_GrantIssuerRoleWithoutPermission() public {
//         address newIssuer = makeAddr("newIssuer");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(PARTICIPANT, address(eventAccessManager), PARTICIPANT_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert(EventAccessManager.EventAccessManager__NotHostOrAdmin.selector);
//         eventAccessManager.grantIssuerRole(newIssuer, signMessage, signature);
//     }

//     function test_WhenGrantIssuerRoleWithHostPermission() public {
//         address newIssuer = makeAddr("newIssuer");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit IssuerRoleGranted(newIssuer, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantIssuerRole(newIssuer, signMessage, signature);
        
//         assertTrue(eventAccessManager.checkIsIssuer(newIssuer), "New issuer should have issuer role");
//     }

//     function test_WhenGrantIssuerRoleWithAdminPermission() public {
//         address newIssuer = makeAddr("newIssuer");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventAccessManager), ADMIN_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit IssuerRoleGranted(newIssuer, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantIssuerRole(newIssuer, signMessage, signature);
        
//         assertTrue(eventAccessManager.checkIsIssuer(newIssuer), "New issuer should have issuer role");
//     }

//     function test_RevertWhen_GrantParticipantRoleWithZeroAddress() public {
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccountCannotBeZeroAddress.selector);
//         eventAccessManager.grantParticipantRole(address(0), signMessage, signature);
//     }

//     function test_WhenGrantParticipantRoleWithHostPermission() public {
//         address newParticipant = makeAddr("newParticipant");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit ParticipantRoleGranted(newParticipant, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantParticipantRole(newParticipant, signMessage, signature);
        
//         assertTrue(eventAccessManager.checkIsParticipant(newParticipant), "New participant should have participant role");
//     }

//     function test_RevertWhen_GrantHostRoleWithZeroAddress() public {
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccountCannotBeZeroAddress.selector);
//         eventAccessManager.grantHostRole(address(0), signMessage, signature);
//     }

//     function test_WhenGrantHostRoleWithHostPermission() public {
//         address newHost = makeAddr("newHost");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit HostRoleGranted(newHost, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantHostRole(newHost, signMessage, signature);
        
//         assertTrue(eventAccessManager.checkIsHost(newHost), "New host should have host role");
//     }

//     /*//////////////////////////////////////////////////////////////
//                            ROLE REVOKING TESTS
//     //////////////////////////////////////////////////////////////*/

//     function test_RevertWhen_RevokeIssuerRoleWithZeroAddress() public {
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccountCannotBeZeroAddress.selector);
//         eventAccessManager.revokeIssuerRole(address(0), signMessage, signature);
//     }

//     function test_WhenRevokeIssuerRoleWithHostPermission() public {
//         // First grant the role
//         address issuer = makeAddr("issuer");
//         (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantIssuerRole(issuer, signMessage1, signature1);
        
//         // Then revoke the role
//         (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit IssuerRoleRevoked(issuer, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.revokeIssuerRole(issuer, signMessage2, signature2);
        
//         assertFalse(eventAccessManager.checkIsIssuer(issuer), "Issuer should no longer have issuer role");
//     }

//     function test_WhenRevokeParticipantRoleWithHostPermission() public {
//         // First grant the role
//         address participant = makeAddr("participant");
//         (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantParticipantRole(participant, signMessage1, signature1);
        
//         // Then revoke the role
//         (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit ParticipantRoleRevoked(participant, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.revokeParticipantRole(participant, signMessage2, signature2);
        
//         assertFalse(eventAccessManager.checkIsParticipant(participant), "Participant should no longer have participant role");
//     }

//     /*//////////////////////////////////////////////////////////////
//                     MESSAGE SENDER ALLOWLIST TESTS
//     //////////////////////////////////////////////////////////////*/

//     function test_RevertWhen_AddAllowedMsgSenderWithZeroAddress() public {
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventAccessManager), ADMIN_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert(EventAccessManager.EventAccessManager__AccountCannotBeZeroAddress.selector);
//         eventAccessManager.addAllowedMsgSender(address(0), signMessage, signature);
//     }

//     function test_RevertWhen_AddAllowedMsgSenderWithoutAdminRole() public {
//         address newSender = makeAddr("newSender");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         vm.expectRevert();
//         eventAccessManager.addAllowedMsgSender(newSender, signMessage, signature);
//     }

//     function test_WhenAddAllowedMsgSenderWithAdminPermission() public {
//         address newSender = makeAddr("newSender");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(ADMIN, address(eventAccessManager), ADMIN_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit MsgSenderAllowed(newSender, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.addAllowedMsgSender(newSender, signMessage, signature);
        
//         assertTrue(eventAccessManager.allowedMsgSenders(newSender), "New sender should be allowed");
//     }

//     function test_WhenRemoveAllowedMsgSenderWithAdminPermission() public {
//         // First add a sender
//         address sender = makeAddr("sender");
//         (string memory signMessage1, bytes memory signature1) = createSignedMessageForRole(ADMIN, address(eventAccessManager), ADMIN_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.addAllowedMsgSender(sender, signMessage1, signature1);
        
//         // Then remove the sender
//         (string memory signMessage2, bytes memory signature2) = createSignedMessageForRole(ADMIN, address(eventAccessManager), ADMIN_PRIVATE_KEY);
        
//         vm.expectEmit(true, true, true, true);
//         emit MsgSenderDisallowed(sender, CALLER);
        
//         vm.prank(CALLER);
//         eventAccessManager.removeAllowedMsgSender(sender, signMessage2, signature2);
        
//         assertFalse(eventAccessManager.allowedMsgSenders(sender), "Sender should no longer be allowed");
//     }

//     /*//////////////////////////////////////////////////////////////
//                            VIEW FUNCTION TESTS
//     //////////////////////////////////////////////////////////////*/

//     function test_CheckIsHost() public view {
//         assertTrue(eventAccessManager.checkIsHost(HOST), "Host should be identified as host");
//         assertFalse(eventAccessManager.checkIsHost(PARTICIPANT), "Participant should not be identified as host");
//     }

//     function test_CheckIsIssuer() public {
//         // First grant issuer role
//         address issuer = makeAddr("issuer");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantIssuerRole(issuer, signMessage, signature);
        
//         assertTrue(eventAccessManager.checkIsIssuer(issuer), "Issuer should be identified as issuer");
//         assertFalse(eventAccessManager.checkIsIssuer(PARTICIPANT), "Non-issuer should not be identified as issuer");
//     }

//     function test_CheckIsParticipant() public {
//         // First grant participant role
//         address participant = makeAddr("participant");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantParticipantRole(participant, signMessage, signature);
        
//         assertTrue(eventAccessManager.checkIsParticipant(participant), "Participant should be identified as participant");
//         assertFalse(eventAccessManager.checkIsParticipant(HOST), "Non-participant should not be identified as participant");
//     }

//     function test_CheckIsHostOrAdmin() public view {
//         assertTrue(eventAccessManager.checkIsHostOrAdmin(HOST), "Host should be identified as host or admin");
//         assertTrue(eventAccessManager.checkIsHostOrAdmin(ADMIN), "Admin should be identified as host or admin");
//         assertFalse(eventAccessManager.checkIsHostOrAdmin(PARTICIPANT), "Non-host/non-admin should not be identified as host or admin");
//     }

//     function test_CheckIsAllowedMsgSender() public {
//         vm.prank(ADMIN);
//         assertTrue(eventAccessManager.checkIsAllowedMsgSender(), "Admin should be allowed msg sender");
        
//         vm.prank(PARTICIPANT);
//         assertFalse(eventAccessManager.checkIsAllowedMsgSender(), "Non-allowed should not be allowed msg sender");
//     }

//     /*//////////////////////////////////////////////////////////////
//                            REQUIRE FUNCTION TESTS
//     //////////////////////////////////////////////////////////////*/

//     function test_RequireAllowedMsgSender() public {
//         vm.prank(ADMIN);
//         eventAccessManager.requireAllowedMsgSender(); // Should not revert
        
//         vm.prank(PARTICIPANT);
//         vm.expectRevert(EventAccessManager.EventAccessManager__NotAllowedMsgSender.selector);
//         eventAccessManager.requireAllowedMsgSender();
//     }

//     function test_RequireHostOrAdmin() public {
//         vm.prank(HOST);
//         eventAccessManager.requireHostOrAdmin(HOST); // Should not revert
        
//         vm.prank(ADMIN);
//         eventAccessManager.requireHostOrAdmin(ADMIN); // Should not revert
        
//         vm.prank(PARTICIPANT);
//         vm.expectRevert(EventAccessManager.EventAccessManager__NotHostOrAdmin.selector);
//         eventAccessManager.requireHostOrAdmin(PARTICIPANT);
//     }

//     function test_RequireHostOrAdminOrParticipant() public {
//         vm.prank(HOST);
//         eventAccessManager.requireHostOrAdminOrParticipant(HOST); // Should not revert
        
//         vm.prank(ADMIN);
//         eventAccessManager.requireHostOrAdminOrParticipant(ADMIN); // Should not revert
        
//         // First grant participant role
//         address participant = makeAddr("participant");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantParticipantRole(participant, signMessage, signature);
        
//         vm.prank(participant);
//         eventAccessManager.requireHostOrAdminOrParticipant(participant); // Should not revert
        
//         vm.prank(PARTICIPANT);
//         vm.expectRevert(EventAccessManager.EventAccessManager__NotHostOrAdminOrParticipant.selector);
//         eventAccessManager.requireHostOrAdminOrParticipant(PARTICIPANT);
//     }

//     function test_RequireParticipant() public {
//         // First grant participant role
//         address participant = makeAddr("participant");
//         (string memory signMessage, bytes memory signature) = createSignedMessageForRole(HOST, address(eventAccessManager), HOST_PRIVATE_KEY);
        
//         vm.prank(CALLER);
//         eventAccessManager.grantParticipantRole(participant, signMessage, signature);
        
//         vm.prank(participant);
//         eventAccessManager.requireParticipant(participant); // Should not revert
        
//         vm.prank(HOST);
//         vm.expectRevert(EventAccessManager.EventAccessManager__NotParticipant.selector);
//         eventAccessManager.requireParticipant(HOST);
//     }
// }
