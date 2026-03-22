import { useState, useEffect, useMemo, Fragment } from "react";
import { ChevronRight, CheckCircle2, XCircle, Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";
import { verifyProof, type ProofVerificationResult } from "@/lib/certificate/verifySignature";
import { verifyCertificateHash } from "@/lib/certificate/verifyHash";
import { useOnChainCertificate } from "@/hooks/useOnChainCertificate";
import { Typography } from "@/components/typography/typography";
import type { GetCertificateShareDataResult } from "@/services/CertificateService/mapper";
import { Row, truncate, SIGNATURE_TOOLTIP } from "./Row";
import { CopyButton } from "./CopyButton";
import { EthExplorerLink, EthNftExplorerLink } from "@/components/common/EthscanLink";

export function CertificateDataTable({ data }: { data: GetCertificateShareDataResult }) {
    const { t } = useTranslation();
    const { payload, contract, decryptedUserData, decryptedCertificateData } = data;
    const { header, data: cert, proof } = payload;
    const [techOpen, setTechOpen] = useState(false);
    const [proofOpen, setProofOpen] = useState(false);
    const [chainOpen, setChainOpen] = useState(false);
    const [attendeeOpen, setAttendeeOpen] = useState(false);
    const [verification, setVerification] = useState<ProofVerificationResult | null>(null);

    const {
        data: onChain,
        isLoading: isChainLoading,
        isError: isChainError,
    } = useOnChainCertificate({
        contractAddress: contract.eventCertificateContractAddress,
        tokenId: contract.certificateTokenId,
        enabled: chainOpen,
    });

    useEffect(() => {
        let cancelled = false;
        verifyProof(proof)
            .then((result) => {
                if (!cancelled) {
                    if (!result.host) {
                        console.error("[CertificateDataTable] Host signature verification failed", {
                            proof,
                        });
                    }
                    result.issuers.forEach((ok, i) => {
                        if (!ok) {
                            console.error(
                                `[CertificateDataTable] Issuer ${i + 1} signature verification failed`,
                                { issuer: proof.issuers[i] },
                            );
                        }
                    });
                    setVerification(result);
                }
            })
            .catch((err) => {
                if (!cancelled) {
                    console.error("[CertificateDataTable] verifyProof threw unexpectedly", err);
                    setVerification(null);
                }
            });
        return () => {
            cancelled = true;
        };
    }, [proof]);

    const hashVerified = useMemo(() => {
        if (!decryptedCertificateData) return undefined;
        const name = decryptedCertificateData.name ?? "";
        // name is stored as "{firstName} {lastName}" — split to match computeCertificateHash input
        const spaceIdx = name.indexOf(" ");
        const firstName = spaceIdx >= 0 ? name.slice(0, spaceIdx) : name;
        const lastName = spaceIdx >= 0 ? name.slice(spaceIdx + 1) : "";
        const result = verifyCertificateHash(
            {
                firstName,
                lastName,
                academicInstitution: decryptedCertificateData.academic_institution,
                certificateTitle: decryptedCertificateData.certificate_title,
                certificateSubtitle: decryptedCertificateData.certificate_subtitle,
            },
            proof.hash,
        );
        if (!result) {
            const csv = [
                name,
                decryptedCertificateData.academic_institution ?? "",
                decryptedCertificateData.certificate_title ?? "",
                decryptedCertificateData.certificate_subtitle ?? "",
            ].join(",");
            console.error(
                "[CertificateDataTable] Hash verification failed — recomputed hash does not match proof.hash",
                {
                    proof_hash: proof.hash,
                    recomputed_csv: csv,
                    decryptedCertificateData,
                },
            );
        }
        return result;
    }, [decryptedCertificateData, proof.hash]);

    const issuedAtDate = cert.issuedAt
        ? new Date(Number(cert.issuedAt) * 1000).toLocaleString()
        : "-";

    let parsedSignMessage: { eventContractAddress?: string; receivers?: string[] } = {};
    try {
        parsedSignMessage = JSON.parse(proof.signMessage) as typeof parsedSignMessage;
    } catch {
        /* ignore */
    }

    const isValid = cert.status === "VALID";

    // Cross-verify attendee name and academic_institution against decryptedCertificateData
    const attendeeFullName = decryptedUserData
        ? `${decryptedUserData.first_name ?? ""} ${decryptedUserData.last_name ?? ""}`.trim()
        : null;

    // Cross-verify attendee name and academic_institution against decryptedCertificateData.
    // proof.hash covers certificate record fields (imported by the host), not the attendee profile,
    // so we can only cross-check these two decrypted sources against each other — not directly against the hash.
    const nameVerified: boolean | undefined =
        attendeeFullName !== null && decryptedCertificateData?.name !== undefined
            ? attendeeFullName.toLowerCase() === (decryptedCertificateData.name ?? "").toLowerCase()
            : undefined;

    const academicInstitutionVerified: boolean | undefined =
        decryptedUserData && decryptedCertificateData?.academic_institution !== undefined
            ? (decryptedUserData.academic_institution ?? "").toLowerCase() ===
              (decryptedCertificateData.academic_institution ?? "").toLowerCase()
            : undefined;

    return (
        <div className="flex flex-col gap-4">
            {/* Certificate Info */}
            <div className="rounded-xl border border-muted/20 bg-muted/5 overflow-hidden">
                <div className="px-4 py-3 border-b border-muted/15">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-xs uppercase tracking-widest"
                    >
                        {t("certificateVerify.table.certDataHeading")}
                    </Typography>
                </div>
                <div className="px-4 py-2">
                    <table className="w-full">
                        <tbody>
                            <Row
                                label={t("certificateVerify.table.event")}
                                value={cert.eventName}
                            />
                            {cert.eventDescription && (
                                <Row
                                    label={t("certificateVerify.table.eventDescription")}
                                    value={cert.eventDescription}
                                />
                            )}
                            {decryptedCertificateData?.name && (
                                <Row
                                    label={t("certificateVerify.table.receiverName")}
                                    value={decryptedCertificateData.name}
                                />
                            )}
                            <Row
                                label={t("certificateVerify.table.title")}
                                value={cert.certificateTitle}
                            />
                            {cert.certificateSubtitle && (
                                <Row
                                    label={t("certificateVerify.table.subtitle")}
                                    value={cert.certificateSubtitle}
                                />
                            )}
                            <Row
                                label={t("certificateVerify.table.issuedAt")}
                                value={issuedAtDate}
                            />
                            <Row
                                label={t("certificateVerify.table.status")}
                                value={
                                    <div className="flex items-center gap-1.5">
                                        {isValid ? (
                                            <CheckCircle2 className="w-3.5 h-3.5 text-green-500 shrink-0" />
                                        ) : (
                                            <XCircle className="w-3.5 h-3.5 text-destructive shrink-0" />
                                        )}
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            color={isValid ? "foreground" : "destructive"}
                                            className="text-xs font-medium"
                                        >
                                            {cert.status}
                                        </Typography>
                                    </div>
                                }
                            />
                        </tbody>
                    </table>
                </div>

                {/* Technical details — collapsible */}
                <button
                    onClick={() => setTechOpen((o) => !o)}
                    className="w-full px-4 py-2.5 flex items-center justify-between cursor-pointer hover:bg-muted/10 transition-colors border-t border-muted/15"
                >
                    <Typography
                        variant="text"
                        tag="span"
                        color="muted"
                        className="text-xs uppercase tracking-widest"
                    >
                        {t("certificateVerify.table.technicalDetails")}
                    </Typography>
                    <ChevronRight
                        className={`w-3.5 h-3.5 text-muted-foreground transition-transform duration-200 ${techOpen ? "rotate-90" : ""}`}
                    />
                </button>
                {techOpen && (
                    <div className="border-t border-muted/15 px-4 py-2">
                        <table className="w-full">
                            <tbody>
                                <Row
                                    label={t("certificateVerify.table.tokenId")}
                                    value={`#${cert.certificateTokenId}`}
                                    mono
                                />
                                <Row
                                    label={t("certificateVerify.table.certificateId")}
                                    value={truncate(cert.certificateId, 26)}
                                    mono
                                    copyValue={cert.certificateId}
                                />
                                <Row
                                    label={t("certificateVerify.table.issuerAddress")}
                                    value={
                                        <EthExplorerLink
                                            address={cert.issuerAddresses}
                                            className="text-xs font-mono"
                                        >
                                            {truncate(cert.issuerAddresses, 26)}
                                        </EthExplorerLink>
                                    }
                                    mono
                                    copyValue={cert.issuerAddresses}
                                />
                                <Row
                                    label={t("certificateVerify.table.issuerId")}
                                    value={truncate(cert.issuerId, 26)}
                                    mono
                                    copyValue={cert.issuerId}
                                />
                                <Row
                                    label={t("certificateVerify.table.receiverAddress")}
                                    value={
                                        <EthExplorerLink
                                            address={cert.receiverAddress}
                                            className="text-xs font-mono"
                                        >
                                            {truncate(cert.receiverAddress, 26)}
                                        </EthExplorerLink>
                                    }
                                    mono
                                    copyValue={cert.receiverAddress}
                                />
                                <Row
                                    label={t("certificateVerify.table.userId")}
                                    value={truncate(cert.userId, 26)}
                                    mono
                                    copyValue={cert.userId}
                                />
                                <Row
                                    label={t("certificateVerify.table.encryptedUserData")}
                                    value={truncate(cert.encryptedUserData, 30)}
                                    mono
                                    copyValue={cert.encryptedUserData}
                                />
                                <Row
                                    label={t("certificateVerify.table.backendEncryptedUserData")}
                                    value={truncate(cert.backendEncryptedUserData, 30)}
                                    mono
                                    copyValue={cert.backendEncryptedUserData}
                                />
                            </tbody>
                        </table>
                    </div>
                )}

                {/* Direct On-Chain Query — collapsible, lazy */}
                <button
                    onClick={() => setChainOpen((o) => !o)}
                    className="w-full px-4 py-2.5 flex items-center justify-between cursor-pointer hover:bg-muted/10 transition-colors border-t border-muted/15"
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
                        className={`w-3.5 h-3.5 text-muted-foreground transition-transform duration-200 ${chainOpen ? "rotate-90" : ""}`}
                    />
                </button>
                {chainOpen &&
                    (() => {
                        let parsedTokenData: { hash?: string } | null = null;
                        if (onChain?.tokenData) {
                            try {
                                const json = onChain.tokenData.replace(/^data:[^,]+,/, "");
                                parsedTokenData = JSON.parse(json) as { hash?: string };
                            } catch {
                                /* ignore */
                            }
                        }
                        const hashMatch =
                            parsedTokenData?.hash !== undefined
                                ? parsedTokenData.hash.toLowerCase() === proof.hash.toLowerCase()
                                : null;
                        const ownerMatch =
                            onChain?.owner !== undefined
                                ? onChain.owner.toLowerCase() === cert.receiverAddress.toLowerCase()
                                : null;
                        return (
                            <div className="border-t border-muted/15 px-4 py-2">
                                {isChainLoading && (
                                    <div className="flex items-center gap-2 py-3">
                                        <Loader2 className="w-4 h-4 animate-spin text-muted-foreground" />
                                        <Typography
                                            variant="text"
                                            tag="span"
                                            color="muted"
                                            className="text-xs"
                                        >
                                            {t("certificateVerify.onChain.loading")}
                                        </Typography>
                                    </div>
                                )}
                                {isChainError && (
                                    <Typography
                                        variant="text"
                                        tag="p"
                                        color="destructive"
                                        className="text-xs py-3"
                                    >
                                        {t("certificateVerify.onChain.error")}
                                    </Typography>
                                )}
                                {onChain && (
                                    <table className="w-full">
                                        <tbody>
                                            <Row
                                                label={t("certificateVerify.onChain.contract")}
                                                value={
                                                    <EthExplorerLink
                                                        address={
                                                            contract.eventCertificateContractAddress
                                                        }
                                                        className="text-xs font-mono"
                                                    >
                                                        {truncate(
                                                            contract.eventCertificateContractAddress,
                                                            26,
                                                        )}
                                                    </EthExplorerLink>
                                                }
                                                mono
                                                copyValue={contract.eventCertificateContractAddress}
                                            />
                                            <Row
                                                label={t("certificateVerify.table.tokenId")}
                                                value={
                                                    <EthNftExplorerLink
                                                        contractAddress={
                                                            contract.eventCertificateContractAddress
                                                        }
                                                        tokenId={contract.certificateTokenId}
                                                        className="text-xs font-mono"
                                                    >
                                                        {`#${contract.certificateTokenId}`}
                                                    </EthNftExplorerLink>
                                                }
                                                mono
                                            />
                                            <Row
                                                label={t("certificateVerify.onChain.owner")}
                                                value={truncate(onChain.owner, 26)}
                                                mono
                                                copyValue={onChain.owner}
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
                        );
                    })()}
            </div>

            {/* Attendee Data — shown only when decryptedUserData is present */}
            {decryptedUserData && (
                <div className="rounded-xl border border-muted/20 bg-muted/5 overflow-hidden">
                    <button
                        onClick={() => setAttendeeOpen((o) => !o)}
                        className="w-full px-4 py-3 flex items-center justify-between cursor-pointer hover:bg-muted/10 transition-colors"
                    >
                        <Typography
                            variant="text"
                            tag="span"
                            color="muted"
                            className="text-xs uppercase tracking-widest"
                        >
                            {t("certificateVerify.table.attendeeDataHeading")}
                        </Typography>
                        <ChevronRight
                            className={`w-3.5 h-3.5 text-muted-foreground transition-transform duration-200 ${attendeeOpen ? "rotate-90" : ""}`}
                        />
                    </button>
                    {attendeeOpen && (
                        <div className="border-t border-muted/15 px-4 py-2">
                            <p className="text-xs text-muted-foreground py-2 leading-relaxed">
                                {t("certificateVerify.table.attendeeDataNote")}
                            </p>
                            <table className="w-full">
                                <tbody>
                                    {attendeeFullName && (
                                        <Row
                                            label={t("certificateVerify.table.fullName")}
                                            value={attendeeFullName}
                                        />
                                    )}
                                    {decryptedUserData.academic_institution && (
                                        <Row
                                            label={t("certificateVerify.table.academicInstitution")}
                                            value={decryptedUserData.academic_institution}
                                        />
                                    )}
                                    {decryptedUserData.email && (
                                        <Row
                                            label={t("certificateVerify.table.email")}
                                            value={decryptedUserData.email}
                                        />
                                    )}
                                    {decryptedUserData.academic_email && (
                                        <Row
                                            label={t("certificateVerify.table.academicEmail")}
                                            value={decryptedUserData.academic_email}
                                        />
                                    )}
                                    {decryptedUserData.phone_number && (
                                        <Row
                                            label={t("certificateVerify.table.phoneNumber")}
                                            value={decryptedUserData.phone_number}
                                        />
                                    )}
                                    {decryptedUserData.address && (
                                        <Row
                                            label={t("certificateVerify.table.address")}
                                            value={decryptedUserData.address}
                                        />
                                    )}
                                    {decryptedUserData.bio && (
                                        <Row
                                            label={t("certificateVerify.table.bio")}
                                            value={decryptedUserData.bio}
                                        />
                                    )}
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            )}

            {/* On-chain Proof — collapsible */}
            <div className="rounded-xl border border-muted/20 bg-muted/5 overflow-hidden">
                <button
                    onClick={() => setProofOpen((o) => !o)}
                    className="w-full px-4 py-3 flex items-center justify-between cursor-pointer hover:bg-muted/10 transition-colors"
                >
                    <Typography
                        variant="text"
                        tag="span"
                        color="muted"
                        className="text-xs uppercase tracking-widest"
                    >
                        {t("certificateVerify.table.onChainProofHeading")}
                    </Typography>
                    <ChevronRight
                        className={`w-3.5 h-3.5 text-muted-foreground transition-transform duration-200 ${proofOpen ? "rotate-90" : ""}`}
                    />
                </button>
                {proofOpen && (
                    <div className="border-t border-muted/15 px-4 py-2">
                        <table className="w-full">
                            <tbody>
                                <Row
                                    label={t("certificateVerify.table.credentialType")}
                                    value={header.type.join(", ")}
                                />
                                <Row
                                    label={t("certificateVerify.table.issuanceDate")}
                                    value={new Date(
                                        Number(header.issuanceDate) * 1000,
                                    ).toLocaleString()}
                                />
                                <Row
                                    label={t("certificateVerify.table.hash")}
                                    value={truncate(proof.hash, 30)}
                                    mono
                                    copyValue={proof.hash}
                                    verified={hashVerified}
                                    labelTooltip={
                                        decryptedCertificateData
                                            ? `Hash is recomputed client-side as:\nkeccak256("{firstName} {lastName},{academicInstitution},{certificateTitle},{certificateSubtitle}")\n\nIf it matches the on-chain proof hash, the certificate data has not been tampered with.`
                                            : undefined
                                    }
                                />
                                {parsedSignMessage.eventContractAddress && (
                                    <Row
                                        label={t("certificateVerify.table.contractAddress")}
                                        value={truncate(parsedSignMessage.eventContractAddress, 26)}
                                        mono
                                        copyValue={parsedSignMessage.eventContractAddress}
                                    />
                                )}
                                <Row
                                    label={t("certificateVerify.table.signMessage")}
                                    value={
                                        <div className="flex items-start gap-1.5 w-full">
                                            <pre className="text-xs font-mono text-foreground bg-muted/10 rounded-md px-2 py-1.5 break-all whitespace-pre-wrap flex-1 max-h-28 overflow-y-auto">
                                                {JSON.stringify(parsedSignMessage, null, 2)}
                                            </pre>
                                            <CopyButton value={proof.signMessage} />
                                        </div>
                                    }
                                />
                                <Row
                                    label={t("certificateVerify.table.hostPublicKey")}
                                    value={truncate(proof.host.publicKey, 26)}
                                    mono
                                    copyValue={proof.host.publicKey}
                                />
                                <Row
                                    label={t("certificateVerify.table.hostSignature")}
                                    value={truncate(proof.host.signature, 30)}
                                    mono
                                    copyValue={proof.host.signature}
                                    verified={verification ? verification.host : null}
                                    labelTooltip={SIGNATURE_TOOLTIP}
                                />
                                {proof.issuers.map((issuer, i) => (
                                    <Fragment key={i}>
                                        <Row
                                            label={t("certificateVerify.table.issuerKey", {
                                                n: proof.issuers.length > 1 ? ` ${i + 1}` : "",
                                            })}
                                            value={truncate(issuer.issuerPublicKey, 26)}
                                            mono
                                            copyValue={issuer.issuerPublicKey}
                                        />
                                        <Row
                                            label={t("certificateVerify.table.issuerSignature", {
                                                n: proof.issuers.length > 1 ? ` ${i + 1}` : "",
                                            })}
                                            value={truncate(issuer.issuerSignature, 30)}
                                            mono
                                            copyValue={issuer.issuerSignature}
                                            verified={
                                                verification
                                                    ? (verification.issuers[i] ?? null)
                                                    : null
                                            }
                                            labelTooltip={SIGNATURE_TOOLTIP}
                                        />
                                    </Fragment>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
