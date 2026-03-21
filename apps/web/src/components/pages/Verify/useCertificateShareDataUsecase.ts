import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { AxiosError } from "axios";
import { toast } from "sonner";
import { ToastFromAxiosError } from "@/common/Err";
import { certificateService } from "@/services/services";
import type {
    GetCertificateShareDataResult,
    CertificateShareViewStatus,
} from "@/services/CertificateService/mapper";

export const useCertificateShareDataUsecase = (handle: string | null) => {
    const { t } = useTranslation();

    const [shareStatus, setShareStatus] = useState<CertificateShareViewStatus | null>(null);
    const [shareData, setShareData] = useState<GetCertificateShareDataResult | null>(null);
    const [isUnlocking, setIsUnlocking] = useState(false);
    const [passwordError, setPasswordError] = useState("");

    useEffect(() => {
        setShareStatus(null);
        setShareData(null);
        setPasswordError("");
    }, [handle]);

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
                    toast.info(t("certificateVerify.passwordProtected"));
                    return null;
                }
                if (err instanceof AxiosError) {
                    if (err.response?.status === 404) {
                        toast.error(t("certificateVerify.notFound"));
                    } else {
                        ToastFromAxiosError(t, err);
                    }
                }
                throw err;
            }
        },
        enabled: !!handle,
        retry: false,
        gcTime: 0,
        staleTime: 30_000,
    });

    const handleUnlock = async (password: string) => {
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
            toast.success(t("certificateVerify.unlockSuccess"));
        } catch {
            setPasswordError(t("certificateVerify.incorrectPassword"));
            toast.error(t("certificateVerify.incorrectPassword"));
        } finally {
            setIsUnlocking(false);
        }
    };

    return {
        shareStatus,
        shareData,
        isLoading,
        isFetching,
        isError,
        error,
        refetch,
        isUnlocking,
        passwordError,
        handleUnlock,
    };
};
