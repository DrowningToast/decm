import { ChevronRight, Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { Row, truncate } from "./Row";

export function OnChainTable({
    contractAddress,
    tokenId,
    open,
    onToggle,
    onChainData,
    isLoading,
    isError,
    expectedHash,
    expectedOwner,
}: {
    contractAddress: string;
    tokenId: string;
    open: boolean;
    onToggle: () => void;
    onChainData: { owner: string; tokenData: string } | null;
    isLoading: boolean;
    isError: boolean;
    expectedHash: string;
    expectedOwner: string;
}) {
    const { t } = useTranslation();
    let parsedTokenData: { hash?: string } | null = null;
    if (onChainData?.tokenData) {
        try {
            const json = onChainData.tokenData.replace(/^data:[^,]+,/, "");
            parsedTokenData = JSON.parse(json) as { hash?: string };
        } catch {
            /* ignore */
        }
    }

    const hashMatch =
        parsedTokenData?.hash !== undefined
            ? parsedTokenData.hash.toLowerCase() === expectedHash.toLowerCase()
            : null;

    const ownerMatch =
        onChainData?.owner !== undefined
            ? onChainData.owner.toLowerCase() === expectedOwner.toLowerCase()
            : null;

    if (hashMatch === false) {
        console.error("[OnChainTable] On-chain hash does not match proof hash", {
            onChain: parsedTokenData?.hash,
            expected: expectedHash,
        });
    }
    if (ownerMatch === false) {
        console.error("[OnChainTable] On-chain owner does not match receiver address", {
            onChain: onChainData?.owner,
            expected: expectedOwner,
        });
    }

    return (
        <div className="rounded-xl border border-muted/20 bg-muted/5 overflow-hidden">
            <button
                onClick={onToggle}
                className="w-full px-4 py-3 flex items-center justify-between cursor-pointer hover:bg-muted/10 transition-colors"
            >
                <Typography
                    variant="text"
                    tag="span"
                    color="muted"
                    className="text-xs uppercase tracking-widest"
                >
                    {t("certificateVerify.onChain.heading")}
                </Typography>
                <ChevronRight
                    className={`w-3.5 h-3.5 text-muted-foreground transition-transform duration-200 ${open ? "rotate-90" : ""}`}
                />
            </button>
            {open && (
                <div className="border-t border-muted/15 px-4 py-2">
                    {isLoading && (
                        <div className="flex items-center gap-2 py-3">
                            <Loader2 className="w-4 h-4 animate-spin text-muted-foreground" />
                            <Typography variant="text" tag="span" color="muted" className="text-xs">
                                {t("certificateVerify.onChain.loading")}
                            </Typography>
                        </div>
                    )}
                    {isError && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="destructive"
                            className="text-xs py-3"
                        >
                            {t("certificateVerify.onChain.error")}
                        </Typography>
                    )}
                    {onChainData && (
                        <table className="w-full">
                            <tbody>
                                <Row
                                    label={t("certificateVerify.onChain.contract")}
                                    value={truncate(contractAddress, 26)}
                                    mono
                                    copyValue={contractAddress}
                                />
                                <Row
                                    label={t("certificateVerify.table.tokenId")}
                                    value={`#${tokenId}`}
                                    mono
                                />
                                <Row
                                    label={t("certificateVerify.onChain.owner")}
                                    value={truncate(onChainData.owner, 26)}
                                    mono
                                    copyValue={onChainData.owner}
                                    verified={ownerMatch}
                                />
                                {parsedTokenData?.hash && (
                                    <Row
                                        label={t("certificateVerify.table.hash")}
                                        value={truncate(parsedTokenData.hash, 30)}
                                        mono
                                        copyValue={parsedTokenData.hash}
                                        verified={hashMatch}
                                    />
                                )}
                            </tbody>
                        </table>
                    )}
                </div>
            )}
        </div>
    );
}
