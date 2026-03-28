import { Info, AlertTriangle } from "lucide-react";
import { Typography } from "@/components/typography/typography";
import { Tooltip, TooltipTrigger, TooltipContent } from "@/components/ui/tooltip";
import { VerifiedBadge } from "./VerifiedBadge";
import { CopyButton } from "./CopyButton";

export const SIGNATURE_TOOLTIP = `Verification uses Ethereum's personal_sign scheme.

1. The sign message (contract address + certificate hash) is hashed with the Ethereum prefix "\x19Ethereum Signed Message:\n".
2. The signer's address is mathematically recovered from the signature using elliptic curve recovery (ecrecover).
3. If the recovered address matches the claimed public key, the signature is marked Verified.

This proves the certificate was authorised by the claimed key without any network call.`;

export function truncate(value: string, maxLen = 20): string {
    if (value.length <= maxLen) return value;
    return `${value.slice(0, 10)}...${value.slice(-8)}`;
}

export function Row({
    label,
    value,
    mono = false,
    copyValue,
    verified,
    labelTooltip,
    warningTooltip,
}: {
    label: string;
    value: React.ReactNode;
    mono?: boolean;
    copyValue?: string;
    verified?: boolean | null;
    labelTooltip?: string;
    warningTooltip?: string;
}) {
    return (
        <tr className="border-b border-muted/10 last:border-0">
            <td className="py-2 pr-4 align-top w-1/3">
                <div className="flex items-center gap-1">
                    <Typography variant="text" tag="span" color="muted" className="text-xs">
                        {label}
                    </Typography>
                    {labelTooltip && (
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Info className="w-3 h-3 text-muted-foreground/50 hover:text-muted-foreground cursor-default shrink-0" />
                            </TooltipTrigger>
                            <TooltipContent
                                side="right"
                                className="max-w-72 whitespace-pre-line text-xs leading-relaxed"
                            >
                                {labelTooltip}
                            </TooltipContent>
                        </Tooltip>
                    )}
                </div>
            </td>
            <td className="py-2 align-top">
                <div className="flex items-start gap-1.5">
                    {typeof value === "string" ? (
                        <Typography
                            variant="text"
                            tag="span"
                            color="foreground"
                            className={`text-xs break-all ${mono ? "font-mono" : ""}`}
                        >
                            {value}
                        </Typography>
                    ) : (
                        value
                    )}
                    {verified !== undefined &&
                        (verified === false && warningTooltip ? (
                            <Tooltip>
                                <TooltipTrigger className="inline-flex cursor-default">
                                    <AlertTriangle className="w-3.5 h-3.5 text-yellow-500 shrink-0" />
                                </TooltipTrigger>
                                <TooltipContent
                                    side="right"
                                    className="max-w-64 text-xs leading-relaxed"
                                >
                                    {warningTooltip}
                                </TooltipContent>
                            </Tooltip>
                        ) : (
                            <VerifiedBadge verified={verified ?? null} />
                        ))}
                    {copyValue && <CopyButton value={copyValue} />}
                </div>
            </td>
        </tr>
    );
}
