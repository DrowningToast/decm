import { useState, useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { certificateService } from "@/services/services";
import { useCertificateShareImage } from "@/hooks/useCertificateShareImage";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
    Award,
    ChevronLeft,
    ChevronRight,
    FileCheck2,
    Loader2,
    Search,
    Share2,
    ShieldCheck,
    CheckCircle2,
    XCircle,
    Copy,
    Check,
} from "lucide-react";
import { Typography } from "@/components/typography/typography";
import type {
    GetCertificateShareDataResult,
    CertificateShareViewStatus,
} from "@/services/CertificateService/mapper";

function truncate(value: string, maxLen = 20): string {
    if (value.length <= maxLen) return value;
    return `${value.slice(0, 10)}...${value.slice(-8)}`;
}

function CopyButton({ value }: { value: string }) {
    const [copied, setCopied] = useState(false);

    const handleCopy = () => {
        void navigator.clipboard.writeText(value).then(() => {
            setCopied(true);
            setTimeout(() => setCopied(false), 1500);
        });
    };

    return (
        <button
            onClick={handleCopy}
            className="shrink-0 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            aria-label="Copy to clipboard"
        >
            {copied ? <Check className="w-3 h-3 text-green-500" /> : <Copy className="w-3 h-3" />}
        </button>
    );
}

function Row({
    label,
    value,
    mono = false,
    copyValue,
}: {
    label: string;
    value: React.ReactNode;
    mono?: boolean;
    copyValue?: string;
}) {
    return (
        <tr className="border-b border-muted/10 last:border-0">
            <td className="py-2 pr-4 align-top w-1/3">
                <Typography variant="text" tag="span" color="muted" className="text-xs">
                    {label}
                </Typography>
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
                    {copyValue && <CopyButton value={copyValue} />}
                </div>
            </td>
        </tr>
    );
}

function CertificateDataTable({ data }: { data: GetCertificateShareDataResult }) {
    const { payload } = data;
    const { header, data: cert, proof } = payload;
    const [proofOpen, setProofOpen] = useState(false);

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
                        Certificate Data
                    </Typography>
                </div>
                <div className="px-4 py-2">
                    <table className="w-full">
                        <tbody>
                            <Row label="Event" value={cert.eventName} />
                            <Row label="Title" value={cert.certificateTitle} />
                            {cert.certificateSubtitle && (
                                <Row label="Subtitle" value={cert.certificateSubtitle} />
                            )}
                            <Row label="Token ID" value={`#${cert.certificateTokenId}`} mono />
                            <Row
                                label="Certificate ID"
                                value={truncate(cert.certificateId, 26)}
                                mono
                                copyValue={cert.certificateId}
                            />
                            <Row label="Issued At" value={issuedAtDate} />
                            <Row
                                label="Issuer Address"
                                value={truncate(cert.issuerAddresses, 26)}
                                mono
                                copyValue={cert.issuerAddresses}
                            />
                            <Row
                                label="Receiver Address"
                                value={truncate(cert.receiverAddress, 26)}
                                mono
                                copyValue={cert.receiverAddress}
                            />
                            <Row
                                label="Status"
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
            </div>

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
                        On-Chain Proof
                    </Typography>
                    <ChevronRight
                        className={`w-3.5 h-3.5 text-muted-foreground transition-transform duration-200 ${proofOpen ? "rotate-90" : ""}`}
                    />
                </button>
                {proofOpen && (
                    <div className="border-t border-muted/15 px-4 py-2">
                        <table className="w-full">
                            <tbody>
                                <Row label="Credential Type" value={header.type.join(", ")} />
                                <Row
                                    label="Issuance Date"
                                    value={new Date(
                                        Number(header.issuanceDate) * 1000,
                                    ).toLocaleString()}
                                />
                                <Row
                                    label="Hash"
                                    value={truncate(proof.hash, 30)}
                                    mono
                                    copyValue={proof.hash}
                                />
                                {parsedSignMessage.eventContractAddress && (
                                    <Row
                                        label="Contract Address"
                                        value={truncate(parsedSignMessage.eventContractAddress, 26)}
                                        mono
                                        copyValue={parsedSignMessage.eventContractAddress}
                                    />
                                )}
                                <Row
                                    label="Host Public Key"
                                    value={truncate(proof.host.publicKey, 26)}
                                    mono
                                    copyValue={proof.host.publicKey}
                                />
                                <Row
                                    label="Host Signature"
                                    value={truncate(proof.host.signature, 30)}
                                    mono
                                    copyValue={proof.host.signature}
                                />
                                {proof.issuers.map((issuer, i) => (
                                    <>
                                        <Row
                                            key={`key-${i}`}
                                            label={`Issuer ${proof.issuers.length > 1 ? i + 1 : ""} Key`.trim()}
                                            value={truncate(issuer.issuerPublicKey, 26)}
                                            mono
                                            copyValue={issuer.issuerPublicKey}
                                        />
                                        <Row
                                            key={`sig-${i}`}
                                            label={`Issuer ${proof.issuers.length > 1 ? i + 1 : ""} Signature`.trim()}
                                            value={truncate(issuer.issuerSignature, 30)}
                                            mono
                                            copyValue={issuer.issuerSignature}
                                        />
                                    </>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}

function HowItWorks() {
    const { t } = useTranslation();

    const steps = [
        {
            icon: <FileCheck2 className="w-6 h-6" />,
            title: t("certificateVerify.how.step1.title"),
            description: t("certificateVerify.how.step1.description"),
        },
        {
            icon: <Share2 className="w-6 h-6" />,
            title: t("certificateVerify.how.step2.title"),
            description: t("certificateVerify.how.step2.description"),
        },
        {
            icon: <ShieldCheck className="w-6 h-6" />,
            title: t("certificateVerify.how.step3.title"),
            description: t("certificateVerify.how.step3.description"),
        },
    ] as const;

    return (
        <div className="flex flex-col gap-4">
            <Typography
                variant="text"
                tag="p"
                color="muted"
                className="text-xs uppercase tracking-widest"
            >
                {t("certificateVerify.how.heading")}
            </Typography>

            {/* Mobile: vertical stack */}
            <div className="flex flex-col gap-3 md:hidden">
                {steps.map((step, i) => (
                    <div key={i} className="flex flex-col gap-3">
                        <div className="flex items-start gap-3 p-4 rounded-xl border border-muted/20 bg-muted/5">
                            <div className="shrink-0 w-10 h-10 rounded-lg bg-muted/15 flex items-center justify-center text-muted-foreground">
                                {step.icon}
                            </div>
                            <div className="flex flex-col gap-0.5">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="foreground"
                                    className="text-sm font-medium"
                                >
                                    {step.title}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-xs leading-relaxed"
                                >
                                    {step.description}
                                </Typography>
                            </div>
                        </div>
                        {i < steps.length - 1 && (
                            <div className="flex justify-center">
                                <div className="w-px h-4 bg-muted/30" />
                            </div>
                        )}
                    </div>
                ))}
            </div>

            {/* Desktop: horizontal flow */}
            <div className="hidden md:flex items-stretch gap-2">
                {steps.map((step, i) => (
                    <div key={i} className="flex items-stretch gap-2 flex-1">
                        <div className="flex-1 flex flex-col gap-3 p-5 rounded-xl border border-muted/20 bg-muted/5">
                            <div className="w-10 h-10 rounded-lg bg-muted/15 flex items-center justify-center text-muted-foreground shrink-0">
                                {step.icon}
                            </div>
                            <div className="flex flex-col gap-1">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="foreground"
                                    className="text-sm font-medium"
                                >
                                    {step.title}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-xs leading-relaxed"
                                >
                                    {step.description}
                                </Typography>
                            </div>
                        </div>
                        {i < steps.length - 1 && (
                            <div className="flex items-center self-center shrink-0">
                                <ChevronRight className="w-4 h-4 text-muted-foreground/40" />
                            </div>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
}

function AspectFrame({ children }: { children: React.ReactNode }) {
    return (
        <div className="w-full relative" style={{ paddingBottom: "75%" }}>
            <div className="absolute inset-0 rounded-xl overflow-hidden">{children}</div>
        </div>
    );
}

function EmptyCertificateFrame() {
    const { t } = useTranslation();
    return (
        <AspectFrame>
            <div className="absolute inset-0 rounded-xl border-2 border-dashed border-muted-foreground/20 bg-muted/10 flex flex-col items-center justify-center gap-3">
                <Award className="w-12 h-12 text-muted-foreground/30" />
                <Typography variant="text" tag="p" color="muted" className="text-sm">
                    {t("certificateVerify.emptyFrame")}
                </Typography>
            </div>
        </AspectFrame>
    );
}

export default function VerifyPage() {
    const { t } = useTranslation();
    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate();
    const handle = searchParams.get("handle");

    const [inputCode, setInputCode] = useState(handle ?? "");
    const [password, setPassword] = useState("");
    const [passwordError, setPasswordError] = useState("");
    const [shareStatus, setShareStatus] = useState<CertificateShareViewStatus | null>(null);
    const [shareData, setShareData] = useState<GetCertificateShareDataResult | null>(null);
    const [isUnlocking, setIsUnlocking] = useState(false);

    // Sync input and reset state whenever the handle in the URL changes
    useEffect(() => {
        setInputCode(handle ?? "");
        setPassword("");
        setPasswordError("");
        setShareStatus(null);
        setShareData(null);
    }, [handle]);

    const unlockedPassword = shareStatus === "READY" ? password || undefined : undefined;

    const { imageUrl, isLoading: isImageLoading } = useCertificateShareImage({
        handle: handle ?? undefined,
        password: unlockedPassword,
        enabled: shareStatus === "READY",
    });

    const { isLoading, isError } = useQuery({
        queryKey: ["certificateShareData", handle],
        queryFn: async () => {
            try {
                const result = await certificateService.getCertificateShareData(handle!);
                setShareData(result);
                setShareStatus("READY");
                return result;
            } catch (err: unknown) {
                const status = (err as { status?: number })?.status;
                if (status === 403) {
                    setShareStatus("PASSWORD_LOCKED");
                    return null;
                }
                throw err;
            }
        },
        enabled: !!handle,
        retry: false,
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const trimmed = inputCode.trim();
        if (!trimmed) return;
        setSearchParams({ handle: trimmed });
    };

    const handleUnlock = async () => {
        if (!password) {
            setPasswordError(t("certificateVerify.passwordRequired"));
            return;
        }
        setIsUnlocking(true);
        setPasswordError("");
        try {
            const result = await certificateService.getCertificateShareData(handle!, password);
            setShareData(result);
            setShareStatus("READY");
        } catch {
            setPasswordError(t("certificateVerify.incorrectPassword"));
        } finally {
            setIsUnlocking(false);
        }
    };

    const renderCertificateArea = () => {
        if (!handle) {
            return <EmptyCertificateFrame />;
        }

        if (isLoading) {
            return (
                <AspectFrame>
                    <div className="absolute inset-0 rounded-xl border border-muted/30 bg-muted/10 flex items-center justify-center">
                        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                    </div>
                </AspectFrame>
            );
        }

        if (isError) {
            return (
                <AspectFrame>
                    <div className="absolute inset-0 rounded-xl border-2 border-dashed border-destructive/20 bg-destructive/5 flex flex-col items-center justify-center gap-2">
                        <Typography variant="text" tag="p" color="destructive" className="text-sm">
                            {t("certificateVerify.notFound")}
                        </Typography>
                    </div>
                </AspectFrame>
            );
        }

        if (shareStatus === "PASSWORD_LOCKED" && !shareData) {
            return <EmptyCertificateFrame />;
        }

        if (shareStatus === "VALID_BUT_PENDING") {
            return (
                <AspectFrame>
                    <div className="absolute inset-0 rounded-xl border border-muted/30 bg-muted/10 flex flex-col items-center justify-center gap-2">
                        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                        <Typography variant="text" tag="p" color="muted" className="text-sm">
                            {t("certificateVerify.pending")}
                        </Typography>
                    </div>
                </AspectFrame>
            );
        }

        if (shareStatus === "READY" && shareData) {
            return (
                <AspectFrame>
                    {isImageLoading && (
                        <div className="absolute inset-0 rounded-xl border border-muted/30 bg-muted/10 flex items-center justify-center">
                            <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
                        </div>
                    )}
                    {imageUrl && (
                        <img
                            src={imageUrl}
                            alt={shareData.payload.data.certificateTitle ?? "Certificate"}
                            className="absolute inset-0 w-full h-full object-contain rounded-xl"
                        />
                    )}
                    {!isImageLoading && !imageUrl && (
                        <div className="absolute inset-0 rounded-xl border border-primary/30 bg-primary/5 flex items-center justify-center">
                            <Typography
                                variant="header"
                                tag="h2"
                                color="foreground"
                                className="text-xl font-header"
                            >
                                {shareData.payload.data.certificateTitle}
                            </Typography>
                        </div>
                    )}
                </AspectFrame>
            );
        }

        return <EmptyCertificateFrame />;
    };

    const renderInputArea = () => {
        if (handle && shareStatus === "PASSWORD_LOCKED" && !shareData) {
            return (
                <div className="flex flex-col gap-3">
                    <Typography
                        variant="text"
                        tag="p"
                        color="muted"
                        className="text-sm text-center"
                    >
                        {t("certificateVerify.passwordProtected")}
                    </Typography>
                    <div className="flex gap-2">
                        <Input
                            type="password"
                            placeholder={t("certificateVerify.passwordPlaceholder")}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            onKeyDown={(e) => e.key === "Enter" && void handleUnlock()}
                            aria-label={t("certificateVerify.passwordPlaceholder")}
                            className="flex-1"
                        />
                        <Button
                            variant="primary"
                            onClick={() => void handleUnlock()}
                            disabled={isUnlocking}
                            loading={isUnlocking}
                        >
                            {t("certificateVerify.unlockButton")}
                        </Button>
                    </div>
                    {passwordError && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="destructive"
                            className="text-xs"
                            role="alert"
                        >
                            {passwordError}
                        </Typography>
                    )}
                </div>
            );
        }

        return (
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
                <Typography variant="text" tag="p" color="muted" className="text-sm text-center">
                    {t("certificateVerify.enterCodeHint")}
                </Typography>
                <div className="flex items-stretch gap-2">
                    <Input
                        type="text"
                        placeholder={t("certificateVerify.codePlaceholder")}
                        value={inputCode}
                        onChange={(e) => setInputCode(e.target.value)}
                        aria-label={t("certificateVerify.codePlaceholder")}
                        className="flex-1 h-auto"
                    />
                    <Button type="submit" variant="primary" disabled={!inputCode.trim()}>
                        <Search className="w-4 h-4" />
                        {t("certificateVerify.verifyButton")}
                    </Button>
                </div>
            </form>
        );
    };

    return (
        <div className="min-h-screen flex flex-col">
            <PublicNavbar />
            <main className="flex-1 flex flex-col px-4 pt-20 pb-12">
                <div className="w-full max-w-2xl mx-auto flex flex-col gap-6">
                    {/* Back */}
                    <button
                        onClick={() => navigate(-1)}
                        className="cursor-pointer flex items-center gap-1 text-muted-foreground hover:text-foreground transition-colors w-fit"
                    >
                        <ChevronLeft className="w-4 h-4" />
                        <Typography variant="text" tag="span" color="muted" className="text-sm">
                            {t("certificateVerify.back")}
                        </Typography>
                    </button>

                    {/* Context */}
                    <div className="flex flex-col gap-1">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="foreground"
                            className="text-2xl font-header"
                        >
                            {t("certificateVerify.pageTitle")}
                        </Typography>
                        <Typography variant="text" tag="p" color="muted" className="text-sm">
                            {t("certificateVerify.pageDescription")}
                        </Typography>
                    </div>

                    {/* Certificate frame — 4:3 */}
                    {renderCertificateArea()}

                    {/* On-chain data table */}
                    {shareStatus === "READY" && shareData && (
                        <CertificateDataTable data={shareData} />
                    )}

                    {/* Input */}
                    {renderInputArea()}

                    {/* How it works */}
                    <HowItWorks />
                </div>
            </main>
        </div>
    );
}
