import { useReadContracts } from "wagmi";
import { EventCertificateABI } from "@/contracts/EventCertificate";
import { wagmiConfig } from "@/config/walletConnect";

export interface OnChainCertificateData {
    owner: string;
    tokenData: string;
}

export function useOnChainCertificate({
    contractAddress,
    tokenId,
    enabled = true,
}: {
    contractAddress: string;
    tokenId: string;
    enabled?: boolean;
}) {
    const contract = {
        address: contractAddress as `0x${string}`,
        abi: EventCertificateABI,
    } as const;

    const tokenIdBigInt = BigInt(tokenId);

    const { data, isLoading, isError } = useReadContracts({
        config: wagmiConfig,
        contracts: [
            { ...contract, functionName: "ownerOf", args: [tokenIdBigInt] },
            { ...contract, functionName: "getTokenData", args: [tokenIdBigInt] },
        ],
        query: { enabled: enabled && !!contractAddress && !!tokenId },
    });

    if (!data) return { data: null, isLoading, isError };

    const [ownerResult, tokenDataResult] = data;

    if (ownerResult.status === "failure" || tokenDataResult.status === "failure") {
        return { data: null, isLoading, isError: true };
    }

    return {
        data: {
            owner: ownerResult.result as string,
            tokenData: tokenDataResult.result as string,
        } satisfies OnChainCertificateData,
        isLoading,
        isError,
    };
}
