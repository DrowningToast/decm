import { useState, useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { AxiosError } from "axios";
import { ToastFromAxiosError } from "@/common/Err";
import { certificateService } from "@/services/services";
import { useCertificateShareImage } from "@/hooks/useCertificateShareImage";
import { PublicNavbar } from "@/components/layouts/navigations/PublicNavbar";
import { ChevronLeft } from "lucide-react";
import { Typography } from "@/components/typography/typography";
import type {
    GetCertificateShareDataResult,
    CertificateShareViewStatus,
} from "@/services/CertificateService/mapper";
import { CertificateDataTable } from "@/components/pages/Verify/CertificateDataTable";
import { HowItWorks } from "@/components/pages/Verify/HowItWorks";
import { CertificateArea } from "@/components/pages/Verify/CertificateArea";
import { PasswordUnlockArea } from "@/components/pages/Verify/PasswordUnlockArea";
import { VerifySearchForm } from "@/components/pages/Verify/VerifySearchForm";

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

    const { isLoading, isFetching, isError, error, refetch } = useQuery({
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
        gcTime: 0,
        staleTime: 30_000,
    });

    useEffect(() => {
        if (isError && error instanceof AxiosError) {
            ToastFromAxiosError(t, error);
        }
    }, [isError, error, t]);

    const errorStatus = error instanceof AxiosError ? error.response?.status : undefined;

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const trimmed = inputCode.trim();
        if (!trimmed) return;
        if (trimmed === handle) {
            // Same handle — reset state and force a refetch
            setShareStatus(null);
            setShareData(null);
            void refetch();
        } else {
            setSearchParams({ handle: trimmed });
        }
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
                    <CertificateArea
                        handle={handle}
                        shareStatus={shareStatus}
                        shareData={shareData}
                        isLoading={isLoading}
                        isError={isError}
                        errorStatus={errorStatus}
                        imageUrl={imageUrl ?? undefined}
                        isImageLoading={isImageLoading}
                    />

                    {/* On-chain data table */}
                    {shareStatus === "READY" && shareData && (
                        <CertificateDataTable data={shareData} />
                    )}

                    {/* Input */}
                    {handle && shareStatus === "PASSWORD_LOCKED" && !shareData ? (
                        <PasswordUnlockArea
                            password={password}
                            onPasswordChange={setPassword}
                            passwordError={passwordError}
                            onUnlock={() => void handleUnlock()}
                            isUnlocking={isUnlocking}
                        />
                    ) : (
                        <VerifySearchForm
                            inputCode={inputCode}
                            onInputCodeChange={setInputCode}
                            isFetching={isFetching}
                            onSubmit={handleSubmit}
                        />
                    )}

                    {/* How it works */}
                    <HowItWorks />
                </div>
            </main>
        </div>
    );
}
