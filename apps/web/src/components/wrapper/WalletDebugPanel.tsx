import { useWallet } from "@/hooks/useWallet";
import { Typography } from "@/components/typography/typography";
import { Copy, Check } from "lucide-react";
import { useState } from "react";

/**
 * WalletDebugPanel Component
 *
 * A debug panel that displays the currently connected wallet information
 * and public client details. Useful for development and debugging.
 *
 * Usage:
 * <WalletDebugPanel />
 */
export const WalletDebugPanel = () => {
    const {
        address,
        isConnected,
        isConnecting,
        isDisconnected,
        chainId,
        publicClient,
        isLoading,
        isReady,
    } = useWallet();
    const [copiedField, setCopiedField] = useState<string | null>(null);

    const copyToClipboard = (text: string, field: string) => {
        navigator.clipboard.writeText(text)
            .then(() => {
                setCopiedField(field);
                setTimeout(() => setCopiedField(null), 2000);
            })
            .catch((err) => {
                console.error('Failed to copy to clipboard:', err);
            });
    };

    if (isDisconnected) {
        return (
            <div className="p-4 bg-background border border-border rounded-lg">
                <Typography variant="text" tag="p" color="muted">
                    No wallet connected
                </Typography>
            </div>
        );
    }

    if (isLoading && !isReady) {
        return (
            <div className="p-4 bg-background border border-border rounded-lg">
                <Typography variant="text" tag="p" color="muted">
                    ⏸️ Initializing wallet...
                </Typography>
                {address && (
                    <Typography variant="text" tag="p" color="muted" className="text-xs mt-2">
                        Address detected: {address.slice(0, 6)}...{address.slice(-4)}
                    </Typography>
                )}
            </div>
        );
    }

    return (
        <div className="p-4 bg-background border border-border rounded-lg space-y-3">
            <Typography variant="text" tag="h3" color="foreground" className="font-semibold">
                💼 Wallet Debug Info
            </Typography>

            {/* Connection Status */}
            <div className="flex items-center gap-2">
                <Typography variant="text" tag="span" color="muted" className="text-sm">
                    Status:
                </Typography>
                {isLoading && !isReady && (
                    <Typography variant="text" tag="span" color="muted" className="text-sm">
                        ⏸️ Initializing...
                    </Typography>
                )}
                {isConnecting && (
                    <Typography variant="text" tag="span" color="muted" className="text-sm">
                        ⏳ Connecting...
                    </Typography>
                )}
                {isReady && isConnected && (
                    <Typography
                        variant="text"
                        tag="span"
                        color="primary"
                        className="text-sm font-medium"
                    >
                        ✅ Ready
                    </Typography>
                )}
                {isConnected && !isReady && (
                    <Typography variant="text" tag="span" color="muted" className="text-sm">
                        🟡 Connected (waiting for public client)
                    </Typography>
                )}
            </div>

            {/* Address */}
            {address && (
                <div className="space-y-1">
                    <Typography variant="text" tag="p" color="muted" className="text-xs">
                        Address:
                    </Typography>
                    <div className="flex items-center gap-2 p-2 bg-foreground/5 rounded border border-border">
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-xs flex-1 font-mono break-all"
                        >
                            {address}
                        </Typography>
                        <button
                            onClick={() => copyToClipboard(address, "address")}
                            className="p-1.5 hover:bg-foreground/10 rounded transition-colors flex-shrink-0"
                            title="Copy address"
                        >
                            {copiedField === "address" ? (
                                <Check className="w-4 h-4 text-green-500" />
                            ) : (
                                <Copy className="w-4 h-4 text-muted-foreground" />
                            )}
                        </button>
                    </div>
                </div>
            )}

            {/* Chain ID */}
            {chainId && (
                <div>
                    <Typography variant="text" tag="p" color="muted" className="text-xs">
                        Chain ID: <span className="text-foreground">{chainId}</span>
                    </Typography>
                </div>
            )}

            {/* Public Client Details */}
            {publicClient && (
                <div className="space-y-2 pt-2 border-t border-border">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs font-semibold"
                    >
                        Public Client Details:
                    </Typography>

                    <div className="space-y-1 text-xs">
                        {/* Removed mode property - not available in viem >=1.2.0 (see https://github.com/wevm/viem/releases/tag/v1.2.0). Update this if/when mode returns. */}

                        <div className="flex justify-between items-center">
                            <Typography variant="text" tag="span" color="muted">
                                Key:
                            </Typography>
                            <Typography variant="text" tag="span" className="font-mono">
                                {publicClient.key}
                            </Typography>
                        </div>

                        {publicClient.chain && (
                            <>
                                <div className="flex justify-between items-center">
                                    <Typography variant="text" tag="span" color="muted">
                                        Chain:
                                    </Typography>
                                    <Typography variant="text" tag="span">
                                        {publicClient.chain.name} (ID: {publicClient.chain.id})
                                    </Typography>
                                </div>
                            </>
                        )}

                        {publicClient.transport && (
                            <div className="flex justify-between items-center">
                                <Typography variant="text" tag="span" color="muted">
                                    Transport:
                                </Typography>
                                <Typography variant="text" tag="span" className="font-mono">
                                    {publicClient.transport.type || "unknown"}
                                </Typography>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* Raw Public Client */}
            <details className="pt-2 border-t border-border">
                <summary className="cursor-pointer text-xs text-muted-foreground hover:text-muted">
                    View Raw Public Client
                </summary>
                <pre className="mt-2 p-2 bg-foreground/5 rounded border border-border overflow-auto max-h-40 text-xs">
                    {JSON.stringify(
                        publicClient,
                        (_key, value) => {
                            // Filter out functions and circular references
                            if (typeof value === "function") return "[Function]";
                            return value;
                        },
                        2,
                    )}
                </pre>
            </details>
        </div>
    );
};
