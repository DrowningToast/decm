import { wagmiConfig } from "@/config/walletConnect";
import { getAccount, signMessage, type Config } from "@wagmi/core";

export class WalletNotConnectedError extends Error {
    constructor() {
        super("Wallet is not connected");
        this.name = "WalletNotConnectedError";
    }
}

export class WalletSigningRejectedError extends Error {
    constructor() {
        super("User rejected the signing request");
        this.name = "WalletSigningRejectedError";
    }
}

export interface IWalletSigningService {
    signMessage(message: string): Promise<string>;
    getConnectedAddress(): string | undefined;
    isConnected(): boolean;
}

export class WalletSigningService implements IWalletSigningService {
    private _wagmiConfig: Config;

    constructor(wagmiConfig: Config) {
        this._wagmiConfig = wagmiConfig;
    }

    public async signMessage(message: string): Promise<string> {
        const account = getAccount(this._wagmiConfig);
        if (!account.isConnected || !account.address) {
            throw new WalletNotConnectedError();
        }

        try {
            const signature = await signMessage(this._wagmiConfig, { message });
            return signature;
        } catch (error) {
            if (error instanceof Error && error.message.toLowerCase().includes("user rejected")) {
                throw new WalletSigningRejectedError();
            }
            throw error;
        }
    }

    public getConnectedAddress(): string | undefined {
        return getAccount(this._wagmiConfig).address;
    }

    public isConnected(): boolean {
        return getAccount(this._wagmiConfig).isConnected ?? false;
    }
}

export const defaultWalletSigningService = new WalletSigningService(wagmiConfig);
