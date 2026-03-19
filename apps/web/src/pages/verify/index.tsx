import { useState, useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { certificateService } from "@/services/services";
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
} from "lucide-react";
import { Typography } from "@/components/typography/typography";
import type { GetCertificateShareDataResult } from "@/services/CertificateService/mapper";

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
    const [unlockedData, setUnlockedData] = useState<GetCertificateShareDataResult | null>(null);
    const [isUnlocking, setIsUnlocking] = useState(false);

    // Sync input and reset unlock state whenever the handle in the URL changes
    useEffect(() => {
        setInputCode(handle ?? "");
        setPassword("");
        setPasswordError("");
        setUnlockedData(null);
    }, [handle]);

    const {
        data: statusData,
        isLoading,
        isError,
    } = useQuery({
        queryKey: ["certificateShareStatus", handle],
        queryFn: () => certificateService.getCertificateShareStatus(handle!),
        enabled: !!handle,
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
            const result = await certificateService.getCertificateShareDataWithPassword(
                handle!,
                password,
            );
            setUnlockedData(result);
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

        if (statusData?.status === "PASSWORD_LOCKED" && !unlockedData) {
            return <EmptyCertificateFrame />;
        }

        if (statusData?.status === "VALID_BUT_PENDING") {
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

        if (statusData?.status === "READY" && statusData?.certificate) {
            return (
                <AspectFrame>
                    <div className="absolute inset-0 rounded-xl border border-primary/30 bg-primary/5 flex items-center justify-center">
                        <Typography
                            variant="header"
                            tag="h2"
                            color="foreground"
                            className="text-xl font-header"
                        >
                            {statusData.certificate.certificateTitle}
                        </Typography>
                    </div>
                </AspectFrame>
            );
        }

        return <EmptyCertificateFrame />;
    };

    const renderInputArea = () => {
        if (handle && statusData?.status === "PASSWORD_LOCKED" && !unlockedData) {
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

                    {/* Input */}
                    {renderInputArea()}

                    {/* How it works */}
                    <HowItWorks />
                </div>
            </main>
        </div>
    );
}
