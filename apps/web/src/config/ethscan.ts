import { polygon, sepolia } from "@reown/appkit/networks";

interface Value {
    explorerUrl: `https://${string}.etherscan.io`;
}

export const EthscanConfig = {
    [polygon.id]: {
        explorerUrl: "https://polygon.etherscan.io",
    },
    [sepolia.id]: {
        explorerUrl: "https://sepolia.etherscan.io",
    },
} satisfies Record<number, Value>;
