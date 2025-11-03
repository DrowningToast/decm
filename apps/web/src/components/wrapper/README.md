# Wallet Connection Wrapper

This package contains wallet connection monitoring utilities for the DECM application.

## Components

### WalletConnectionWrapper

A simple wrapper component that monitors wallet connection status and logs to the console.

**Usage:**

```tsx
import { WalletConnectionWrapper } from "@/components/wrapper/WalletConnectionWrapper";

export const App = () => {
    return (
        <WalletConnectionWrapper>
            <YourAppContent />
        </WalletConnectionWrapper>
    );
};
```

## Context & Hooks

### WalletProvider + useWallet Hook

A more feature-rich solution that provides wallet state to the entire component tree via context.

**Setup (in \_app.tsx):**

```tsx
import { WalletProvider } from "@/context/WalletContext";

export const Layout = () => {
    return (
        <AppKitProvider>
            <WalletProvider>
                <YourApp />
            </WalletProvider>
        </AppKitProvider>
    );
};
```

**Usage in any component:**

```tsx
import { useWallet } from "@/hooks/useWallet";

export const MyComponent = () => {
    const { address, isConnected, publicClient, chainId } = useWallet();

    if (!isConnected) {
        return <div>Please connect your wallet</div>;
    }

    return (
        <div>
            <p>Connected: {address}</p>
            <p>Chain: {chainId}</p>
            <p>Public Client: {publicClient?.key}</p>
        </div>
    );
};
```

## Logged Information

When a wallet connects, the following information is logged to the console:

```
=== WALLET CONNECTED ===
📍 Connected Address: 0x1234...5678
🔗 Chain ID: 1
⛓️ Public Client: PublicClient {...}
📊 Public Client Details: {
  mode: "publicClient",
  transport: { type: "http", url: "(HTTP)" },
  chain: { name: "Ethereum", id: 1 },
  key: "public"
}
========================
```

## Public Client

The `publicClient` from the wallet context is a Viem `PublicClient` instance that provides:

- `mode`: Type of client ("publicClient")
- `transport`: Configuration of the RPC transport (HTTP, WebSocket, etc.)
- `chain`: Information about the connected blockchain
- `key`: Identifier for the client instance

You can use this public client for on-chain reads:

```tsx
const { publicClient, address } = useWallet();

// Example: Get account balance
const balance = await publicClient?.getBalance({ address });

// Example: Get block number
const blockNumber = await publicClient?.getBlockNumber();
```

## Integration Points

The WalletProvider is integrated at the top level of the application tree:

```
_app.tsx (Layout)
  ├── QueryClientProvider
  │   └── ErrorBoundary
  │       ├── Toaster
  │       └── AppKitProvider
  │           └── WalletProvider ← You are here
  │               └── main
  │                   ├── HelmetProvider
  │                   └── AuthProvider
  │                       └── Routes (Outlet)
```

This ensures that wallet connection is available throughout the entire application.
