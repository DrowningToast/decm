# Smart Contract Go Bindings

This directory contains Go bindings for the DECM platform smart contracts, generated using `abigen` from the Ethereum Go implementation.

## Contracts

The following smart contracts have Go bindings generated:

1. **Event** - Main event management contract
2. **EventAccessManager** - Access control for events
3. **EventTicket** - NFT ticket functionality
4. **EventCertificate** - Digital certificate functionality
5. **DecmAccessManager** - Platform-wide access management

## Directory Structure

```
contracts/
├── event/           # Event contract bindings
├── accessmanager/   # EventAccessManager bindings
├── ticket/          # EventTicket bindings
├── certificate/     # EventCertificate bindings
├── decm/            # DecmAccessManager bindings
└── contracts.go     # Unified import file
```

## Installation and Setup

### Prerequisites

1. Go 1.24.1 or higher
2. Ethereum client (Geth, Infura, or Alchemy)
3. Private key for contract deployment/interaction

### Importing in Your Go Code

```go
import (
    "apps/backend/contracts"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)
```

## Usage Examples

### 1. Connecting to Ethereum

```go
// Create an Ethereum client
client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")
if err != nil {
    log.Fatalf("Failed to connect to the Ethereum client: %v", err)
}
defer client.Close()

// Create a transactor
privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
if err != nil {
    log.Fatalf("Failed to parse private key: %v", err)
}

auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // 1 for Mainnet
if err != nil {
    log.Fatalf("Failed to create authorized transactor: %v", err)
}
```

### 2. Deploying a Contract

```go
// Deploy Event contract
address, tx, _, err := contracts.DeployEvent(auth, client, "Event Name", "Event Description", 100)
if err != nil {
    log.Fatalf("Failed to deploy Event contract: %v", err)
}

fmt.Printf("Contract deployed at: %s\n", address.Hex())
fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
```

### 3. Interacting with an Existing Contract

```go
// Connect to an existing contract
contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
eventContract, err := contracts.NewEvent(contractAddress, client)
if err != nil {
    log.Fatalf("Failed to instantiate Event contract: %v", err)
}

// Read data from the contract
eventName, err := eventContract.EventName(nil)
if err != nil {
    log.Fatalf("Failed to get event name: %v", err)
}

fmt.Printf("Event Name: %s\n", eventName)
```

### 4. Writing to a Contract

```go
// Update event details
auth.Value = big.NewInt(0) // Set to 0 for read-only operations
tx, err := eventContract.UpdateEvent(auth, "New Event Name", "New Description", 150)
if err != nil {
    log.Fatalf("Failed to update event: %v", err)
}

fmt.Printf("Transaction submitted: %s\n", tx.Hash().Hex())

// Wait for transaction confirmation
receipt, err := bind.WaitMined(context.Background(), client, tx)
if err != nil {
    log.Fatalf("Failed to wait for transaction: %v", err)
}

fmt.Printf("Transaction confirmed in block: %d\n", receipt.BlockNumber)
```

## Contract-Specific Examples

### Event Contract

```go
// Create a new event
address, tx, _, err := contracts.DeployEvent(
    auth,
    client,
    "Tech Conference 2024",
    "Annual technology conference",
    500, // seats count
)

// Add a participant
tx, err := eventContract.AddParticipant(auth, common.HexToAddress("0xParticipantAddress"))

// Confirm the event
tx, err := eventContract.ConfirmEvent(auth)
```

### EventTicket Contract

```go
// Deploy ticket contract
address, tx, _, err := contracts.DeployEventTicket(
    auth,
    client,
    common.HexToAddress("0xEventContractAddress"),
    "Conference Ticket",
    "TICKET",
    "https://metadata-uri.com",
)

// Mint a ticket to a participant
tx, err := ticketContract.MintTicket(
    auth,
    common.HexToAddress("0xParticipantAddress"),
    1, // token ID
    "ipfs://QmTokenMetadataHash"
)
```

### EventCertificate Contract

```go
// Deploy certificate contract
address, tx, _, err := contracts.DeployEventCertificate(
    auth,
    client,
    common.HexToAddress("0xEventContractAddress"),
    "Completion Certificate",
    "CERT",
)

// Issue a certificate
tx, err := certificateContract.IssueCertificate(
    auth,
    common.HexToAddress("0xRecipientAddress"),
    1, // certificate ID
    "ipfs://QmCertificateMetadataHash"
)
```

## Error Handling

```go
// Handle common errors
if err != nil {
    if strings.Contains(err.Error(), "insufficient funds") {
        fmt.Println("Insufficient funds for gas")
    } else if strings.Contains(err.Error(), "gas required exceeds allowance") {
        fmt.Println("Gas limit too low")
    } else if strings.Contains(err.Error(), "nonce too low") {
        fmt.Println("Transaction nonce issue")
    } else {
        fmt.Printf("Transaction failed: %v\n", err)
    }
}
```

## Event Listening

```go
// Create a channel for events
eventCh := make(chan *contracts.EventEventUpdated)

// Subscribe to events
sub, err := eventContract.WatchEventUpdated(nil, eventCh)
if err != nil {
    log.Fatalf("Failed to subscribe to events: %v", err)
}

// Listen for events
for {
    select {
    case err := <-sub.Err():
        log.Fatalf("Subscription error: %v", err)
    case event := <-eventCh:
        fmt.Printf("Event updated: %s\n", event.EventName)
    }
}
```

## Gas Management

```go
// Estimate gas for a transaction
auth.GasLimit = uint64(300000) // Set gas limit
auth.GasPrice = big.NewInt(20000000000) // 20 Gwei

// Or use suggested gas price
gasPrice, err := client.SuggestGasPrice(context.Background())
if err != nil {
    log.Fatalf("Failed to get suggested gas price: %v", err)
}
auth.GasPrice = gasPrice
```

## Testing

```go
// For testing with a local Ganache network
client, err := ethclient.Dial("http://localhost:8545")
if err != nil {
    log.Fatalf("Failed to connect to Ganache: %v", err)
}

// Use a test account
privateKey, _ := crypto.HexToECDSA("YOUR_TEST_PRIVATE_KEY")
auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337)) // 1337 for Ganache
```

## Regenerating Bindings

If you update your smart contracts, regenerate the Go bindings:

```bash
# Build contracts and generate Go bindings
pnpm run contract:build-and-generate

# Or separately
pnpm run contract:build-only
pnpm run contract:generate-go
```

## Security Considerations

1. **Private Key Management**
    - Never hardcode private keys in production
    - Use environment variables or secure key management systems
    - Consider using hardware wallets for high-value operations

2. **Network Configuration**
    - Always verify you're connected to the correct network
    - Use appropriate chain IDs for each network
    - Test on testnets before mainnet deployment

3. **Transaction Security**
    - Set appropriate gas limits to avoid failed transactions
    - Use reasonable gas prices to ensure timely inclusion
    - Implement proper error handling and retry logic

## Troubleshooting

### Common Issues

1. **"no contract code at given address"**
    - Verify the contract address is correct
    - Check if you're connected to the right network

2. **"insufficient funds for gas"**
    - Ensure the account has enough ETH for gas
    - Check if gas price is too high

3. **"nonce too low" or "nonce too high"**
    - Transaction count is out of sync
    - Reset the account or wait for pending transactions

4. **"abi: cannot marshal"**
    - Check if the contract ABI matches the deployed contract
    - Regenerate bindings if contract was updated

### Debugging Tips

```go
// Enable detailed logging
log.SetFlags(log.LstdFlags | log.Lshortfile)

// Print transaction details before sending
fmt.Printf("Sending transaction to: %s\n", contractAddress.Hex())
fmt.Printf("Function: %s\n", "updateEvent")
fmt.Printf("Gas limit: %d\n", auth.GasLimit)
fmt.Printf("Gas price: %s\n", auth.GasPrice.String())
```

## Resources

- [Ethereum Go Documentation](https://geth.ethereum.org/docs/interface/ethereum)
- [Go Ethereum Contracts](https://geth.ethereum.org/docs/dapp/native-bindings)
- [Solidity Documentation](https://docs.soliditylang.org/)
- [DECM Smart Contracts](../contracts/src/contracts/)
