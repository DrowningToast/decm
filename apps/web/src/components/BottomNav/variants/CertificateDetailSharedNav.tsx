import { LockKeyhole, LockKeyholeOpen, Eye, EyeOff, Copy } from "lucide-react";
import { useState, useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { Switch } from "@/components/ui/switch";
import { CertificateDetailNav } from "./CertificateDetailNav";
import { useCertificateDetailsSharedNavStore } from "../stores/certificates";
import { toast } from "sonner";
import type { ClassValue } from "clsx";

interface CertificateDetailSharedNavProps {
    className?: ClassValue;
}

export const CertificateDetailSharedNav = (props: CertificateDetailSharedNavProps) => {
    const { t } = useTranslation();
    const { shareableHandle, isPublished, onChangePublish, isPasswordProtected } =
        useCertificateDetailsSharedNavStore();
    const [localPublish, setLocalPublish] = useState(isPublished);
    const toastIdRef = useRef<string | number | null>(null);
    const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    useEffect(() => {
        setLocalPublish(isPublished);
    }, [isPublished]);

    const handlePublishChange = (value: boolean) => {
        setLocalPublish(value);

        if (debounceTimerRef.current) clearTimeout(debounceTimerRef.current);
        if (toastIdRef.current) toast.dismiss(toastIdRef.current);

        toastIdRef.current = toast.loading(
            value
                ? t("participant.certificates.detail.makingVisible", "Making it visible...")
                : t("participant.certificates.detail.makingHidden", "Making it hidden..."),
        );

        debounceTimerRef.current = setTimeout(() => {
            onChangePublish(value);
            if (toastIdRef.current) {
                toast.dismiss(toastIdRef.current);
                toastIdRef.current = null;
            }
        }, 3000);
    };

    const handleCopyCode = async () => {
        if (!shareableHandle) return;
        await navigator.clipboard.writeText(shareableHandle);
        toast.success(t("participant.certificates.detail.codeCopied", "Code copied"));
    };

    const handleCopyUrl = async () => {
        if (!shareableHandle) return;
        const url = `${window.location.origin}/verify?handle=${shareableHandle}`;
        await navigator.clipboard.writeText(url);
        toast.success(t("participant.certificates.detail.urlCopied", "Link copied"));
    };

    return (
        <CertificateDetailNav overrideShowCreateShareButton={false} {...props}>
            {/* Copy code button */}
            {shareableHandle && (
                <button
                    onClick={handleCopyCode}
                    className="cursor-pointer flex items-center justify-center gap-2 px-3 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                    aria-label={t("participant.certificates.detail.copyCode", "Copy code")}
                >
                    <Copy className="w-4 h-4 text-background-alt flex-shrink-0" />
                    <span className="text-xs text-background-alt whitespace-nowrap">
                        {t("participant.certificates.detail.shareCode", "Share code")}
                    </span>
                </button>
            )}

            {/* Copy shareable URL */}
            {shareableHandle && (
                <button
                    onClick={handleCopyUrl}
                    className="cursor-pointer flex items-center justify-center gap-2 px-3 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                    aria-label={t(
                        "participant.certificates.detail.copyShareableUrl",
                        "Copy shareable URL",
                    )}
                >
                    <Copy className="w-4 h-4 text-background-alt flex-shrink-0" />
                    <span className="text-xs text-background-alt whitespace-nowrap">
                        {t("participant.certificates.detail.copyShareableUrl", "Copy URL")}
                    </span>
                </button>
            )}

            {/* Publish toggle */}
            <div className="flex items-center gap-2 px-3 h-10 bg-white rounded-[10px] flex-shrink-0">
                {localPublish ? (
                    <Eye className="w-5 h-5 text-background-alt flex-shrink-0" />
                ) : (
                    <EyeOff className="w-5 h-5 text-background-alt flex-shrink-0" />
                )}
                <Switch
                    checked={localPublish}
                    onCheckedChange={handlePublishChange}
                    aria-label={t("participant.certificates.detail.publish", "Publish")}
                />
            </div>

            {/* Password — lock icon reflects protection state */}
            <button
                onClick={() => {}}
                className="cursor-pointer flex items-center justify-center w-10 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                aria-label={t("participant.certificates.detail.editPassword", "Edit password")}
            >
                {isPasswordProtected ? (
                    <LockKeyhole className="w-5 h-5 text-background-alt" />
                ) : (
                    <LockKeyholeOpen className="w-5 h-5 text-background-alt" />
                )}
            </button>
        </CertificateDetailNav>
    );
};
