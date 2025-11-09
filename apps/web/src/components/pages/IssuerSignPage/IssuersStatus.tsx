import { Typography } from "@/components/typography/typography";
import { IssuerStatusBadge } from "./IssuerStatusBadge";
import { isIssuerSigned } from "./issuerStateUtils";

interface IssuersStatusProps {
    issuers: Array<{
        id: string;
        issuer_credential_id: string;
        is_signed: number;
        // Add other issuer fields as needed
    }>;
    currentIssuerId?: string;
    className?: string;
}

export function IssuersStatus({ issuers, currentIssuerId, className = "" }: IssuersStatusProps) {
    return (
        <div className={`bg-[#1a1a1a] border border-[#333333] rounded-lg p-6 ${className}`}>
            <Typography variant="header" tag="h3" className="text-xl font-semibold text-white mb-4">
                Issuers Status
            </Typography>

            {issuers.length === 0 ? (
                <Typography
                    variant="text"
                    tag="p"
                    color="muted-foreground"
                    className="text-center py-4"
                >
                    No issuers found for this event.
                </Typography>
            ) : (
                <div className="space-y-3">
                    {issuers.map((issuer) => (
                        <div
                            key={issuer.id}
                            className={`flex items-center justify-between p-3 rounded-md border ${
                                issuer.issuer_credential_id === currentIssuerId
                                    ? "border-[#ff6a39] bg-[#ff6a3910]"
                                    : "border-[#333333] bg-[#2a2a2a]"
                            }`}
                        >
                            <div className="flex items-center space-x-3">
                                <div className="w-10 h-10 rounded-full bg-gray-600 flex items-center justify-center text-white font-semibold">
                                    {issuer.issuer_credential_id?.substring(0, 2).toUpperCase() ||
                                        "IS"}
                                </div>
                                <div>
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="font-medium text-white"
                                    >
                                        Issuer #{issuer.id?.substring(0, 8) || "Unknown"}
                                    </Typography>
                                    {issuer.issuer_credential_id === currentIssuerId && (
                                        <Typography
                                            variant="text"
                                            tag="p"
                                            className="text-sm text-[#ff6a39]"
                                        >
                                            You
                                        </Typography>
                                    )}
                                </div>
                            </div>

                            <div className="flex items-center space-x-2">
                                <IssuerStatusBadge isSigned={issuer.is_signed} />
                                {isIssuerSigned(issuer.is_signed) && (
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        className="text-xs text-green-500"
                                    >
                                        Completed
                                    </Typography>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
