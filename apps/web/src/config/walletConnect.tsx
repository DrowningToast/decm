import { createAppKit } from "@reown/appkit/react";

import { WagmiProvider } from "wagmi";
import { polygon, sepolia, type AppKitNetwork } from "@reown/appkit/networks";
import { WagmiAdapter } from "@reown/appkit-adapter-wagmi";
import { env } from "./env";

// 0. Setup

const projectId = env.VITE_WALLETCONNECT_PROJECT_ID;

const metadata = {
    name: "Themis",
    description: "Themis",
    url: import.meta.env.DEV ? "http://localhost:3000" : "https://themis.it.kmitl.ac.th.", // origin must match your domain & /logo.svg']
    icons: ["/logo.svg"],
};

const isDev = import.meta.env.DEV;

export const currentNetwork = (isDev ? sepolia : polygon) satisfies AppKitNetwork;
// 3. Set the networks
const networks: [AppKitNetwork, ...AppKitNetwork[]] = [currentNetwork];

// 4. Create Wagmi Adapter
const wagmiAdapter = new WagmiAdapter({
    networks,
    projectId,
});

// 5. Create modal
createAppKit({
    adapters: [wagmiAdapter],
    networks,
    projectId,
    metadata,
    features: {
        analytics: true, // Optional - defaults to your Cloud configuration
        socials: [],
        email: false,
        swaps: false,
        allWallets: true,
    },
    themeMode: "dark",
    themeVariables: {
        "--w3m-color-mix": "var(--primary)",
    },
    debug: import.meta.env.DEV,
    enableWalletGuide: true,
});

export const wagmiConfig = wagmiAdapter.wagmiConfig;

export function AppKitProvider({ children }: { children: React.ReactNode }) {
    return <WagmiProvider config={wagmiAdapter.wagmiConfig}>{children}</WagmiProvider>;
}
