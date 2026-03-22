import { EthscanConfig } from "@/config/ethscan";
import { currentNetwork } from "@/config/walletConnect";
import { cn } from "@/lib/utils";
import type { ClassValue } from "clsx";
import { ExternalLink } from "lucide-react";

interface EthscanLinkProps extends React.PropsWithChildren {
    address: string;
    className?: ClassValue;
    ariaLabel?: string;
}

export const EthExplorerLink: React.FC<EthscanLinkProps> = ({
    address,
    className,
    children,
    ariaLabel,
}) => {
    const explorerUrl = EthscanConfig[currentNetwork.id].explorerUrl;

    return (
        <a
            href={`${explorerUrl}/address/${address}`}
            target="_blank"
            rel="noopener noreferrer"
            aria-label={ariaLabel}
        >
            <div className={cn("flex items-center gap-2", className)}>
                {children}
                <ExternalLink className="w-4 h-4 text-foreground-alt" />
            </div>
        </a>
    );
};

interface EthNftLinkProps extends React.PropsWithChildren {
    contractAddress: string;
    tokenId: string | number;
    className?: ClassValue;
    ariaLabel?: string;
}

export const EthNftExplorerLink: React.FC<EthNftLinkProps> = ({
    contractAddress,
    tokenId,
    className,
    children,
    ariaLabel,
}) => {
    const explorerUrl = EthscanConfig[currentNetwork.id].explorerUrl;

    return (
        <a
            href={`${explorerUrl}/nft/${contractAddress}/${tokenId}`}
            target="_blank"
            rel="noopener noreferrer"
            aria-label={ariaLabel}
        >
            <div className={cn("flex items-center gap-2", className)}>
                {children}
                <ExternalLink className="w-4 h-4 text-foreground-alt" />
            </div>
        </a>
    );
};
