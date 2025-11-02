// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

contract TestUtils is Test {
    using Strings for uint256;

    // Mock addresses for testing
    address public constant ADMIN = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;
    address public constant HOST = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address public constant CALLER = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
    address public constant PARTICIPANT = 0x90F79bf6EB2c4f870365E785982E1f101E93b906;
    address public constant ISSUER = 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65;

    // Private keys for the mock addresses (for signing)
    uint256 public constant ADMIN_PRIVATE_KEY = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
    uint256 public constant HOST_PRIVATE_KEY = 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d;
    uint256 public constant CALLER_PRIVATE_KEY = 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a;
    uint256 public constant PARTICIPANT_PRIVATE_KEY = 0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6;
    uint256 public constant ISSUER_PRIVATE_KEY = 0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a;

    // Mock sign messages
    string public constant MOCK_ADMIN_MESSAGE = "MOCK_ADMIN_MESSAGE";
    string public constant MOCK_HOST_MESSAGE = "MOCK_HOST_MESSAGE";
    string public constant MOCK_CALLER_MESSAGE = "MOCK_CALLER_MESSAGE";
    string public constant MOCK_PARTICIPANT_MESSAGE = "MOCK_PARTICIPANT_MESSAGE";
    string public constant MOCK_ISSUER_MESSAGE = "MOCK_ISSUER_MESSAGE";

    /**
     * @dev Creates a sign message with the format: ADDRESS,CONTRACT_ADDRESS,DEADLINE_BLOCK
     * @param signer The address of the signer
     * @return The formatted sign message string
     */
    function getMockSignMessage(
        address signer
    ) internal pure returns (string memory) {
        if (signer == ADMIN) {
            return MOCK_ADMIN_MESSAGE;
        } else if (signer == HOST) {
            return MOCK_HOST_MESSAGE;
        } else if (signer == CALLER) {
            return MOCK_CALLER_MESSAGE;
        } else if (signer == PARTICIPANT) {
            return MOCK_PARTICIPANT_MESSAGE;
        } else if (signer == ISSUER) {
            return MOCK_ISSUER_MESSAGE;
        }
        return "";
    }

    /**
     * @dev Signs a message with the provided private key
     * @param message The message to sign
     * @param privateKey The private key to sign with
     * @return The signature
     */
    function signMessage(string memory message, uint256 privateKey) internal pure returns (bytes memory) {
        // Hash the message first
        bytes32 messageHash = keccak256(abi.encodePacked(message));
        
        // Sign the hash
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKey, messageHash);
        return abi.encodePacked(r, s, v);
    } 

    /**
     * @dev Creates a signed message for a specific role
     * @param role The role address (HOST, ADMIN, etc.)
     * @param contractAddress The contract address
     * @param privateKey The private key of the role
     * @param deadlineBlocks The number of blocks until deadline (default: 12 blocks ≈ 1 minute)
     * @return The sign message and signature
     */
    function createSignedMessageForRole(
        address role,
        address contractAddress,
        uint256 privateKey,
        uint256 deadlineBlocks
    ) internal returns (string memory, bytes memory) {
        uint256 deadlineBlock = block.number + deadlineBlocks;
        string memory message = getMockSignMessage(role);
        bytes memory signature = signMessage(message, privateKey);
        return (message, signature);
    }

    /**
     * @dev Creates a signed message for a specific role with default deadline (12 blocks)
     * @param role The role address (HOST, ADMIN, etc.)
     * @param contractAddress The contract address
     * @param privateKey The private key of the role
     * @return The sign message and signature
     */
    function createSignedMessageForRole(
        address role,
        address contractAddress,
        uint256 privateKey
    ) internal returns (string memory, bytes memory) {
        return createSignedMessageForRole(role, contractAddress, privateKey, 12);
    }

    /**
     * @dev Creates an expired signed message for testing
     * @param role The role address (HOST, ADMIN, etc.)
     * @param contractAddress The contract address
     * @param privateKey The private key of the role
     * @return The sign message and signature
     */
    function createExpiredSignedMessage(
        address role,
        address contractAddress,
        uint256 privateKey
    ) internal returns (string memory, bytes memory) {
        uint256 pastBlock = block.number - 1;
        string memory message = getMockSignMessage(role);
        bytes memory signature = signMessage(message, privateKey);
        return (message, signature);
    }

    /**
     * @dev Creates a signed message with invalid contract address for testing
     * @param role The role address (HOST, ADMIN, etc.)
     * @param privateKey The private key of the role
     * @return The sign message and signature
     */
    function createInvalidContractSignedMessage(
        address role,
        uint256 privateKey
    ) internal returns (string memory, bytes memory) {
        address invalidContract = address(0x999);
        uint256 deadlineBlock = block.number + 12;
        string memory message = getMockSignMessage(role);
        bytes memory signature = signMessage(message, privateKey);
        return (message, signature);
    }

    /**
     * @dev Advances the blockchain by a specific number of blocks
     * @param numBlocks The number of blocks to advance
     */
    function advanceBlocks(uint256 numBlocks) internal {
        vm.roll(block.number + numBlocks);
    }

    /**
     * @dev Sets up the test environment with mock addresses
     */
    function setupMockAddresses() internal {
        vm.label(ADMIN, "Admin");
        vm.label(HOST, "Host");
        vm.label(CALLER, "Caller");
        vm.label(PARTICIPANT, "Participant");
        vm.label(ISSUER, "Issuer");
    }
}
