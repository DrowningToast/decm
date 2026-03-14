import { KeyRound, Globe, Copy } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import { Typography } from "@/components/typography/typography";
import { Switch } from "@/components/ui/switch";
import { CertificateDetailNav } from "./CertificateDetailNav";
import { useCertificateDetailsSharedNavStore } from "../stores/certificates";
import type { ClassValue } from "clsx";

interface CertificateDetailSharedNavProps {
    className?: ClassValue;
}

export const CertificateDetailSharedNav = (props: CertificateDetailSharedNavProps) => {
    const { t } = useTranslation();
    const [isPublish, setIsPublish] = useState(false);
    const { shareableHandle } = useCertificateDetailsSharedNavStore();

    return (
        <CertificateDetailNav overrideShowCreateShareButton={false} {...props}>
            {/* Shareable handle display */}
            {shareableHandle && (
                <div className="flex items-center gap-2 px-3 h-10 bg-white rounded-[10px] flex-shrink-0 min-w-0">
                    <Typography
                        variant="text"
                        tag="span"
                        color="background-alt"
                        className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap truncate"
                    >
                        {shareableHandle}
                    </Typography>
                    <button
                        onClick={() => {}}
                        className="cursor-pointer flex-shrink-0"
                        aria-label={t("participant.certificates.detail.copyCode", "Copy code")}
                    >
                        <Copy className="w-4 h-4 text-background-alt" />
                    </button>
                </div>
            )}

            {/* Publish toggle */}
            <div className="flex items-center gap-2 px-3 h-10 bg-white rounded-[10px] flex-shrink-0">
                <Globe className="w-4 h-4 text-background-alt flex-shrink-0" />
                <Typography
                    variant="text"
                    tag="span"
                    color="background-alt"
                    className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap"
                >
                    {t("participant.certificates.detail.publish", "Publish")}
                </Typography>
                <Switch
                    checked={isPublish}
                    onCheckedChange={setIsPublish}
                    aria-label={t("participant.certificates.detail.publish", "Publish")}
                />
            </div>

            {/* Copy shareable URL button */}
            <button
                onClick={() => {}}
                className="cursor-pointer flex items-center justify-center gap-2 px-4 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0"
                aria-label={t(
                    "participant.certificates.detail.copyShareableUrl",
                    "Copy shareable URL",
                )}
            >
                <Copy className="w-5 h-5 text-background-alt" />
                <Typography
                    variant="text"
                    tag="span"
                    color="background-alt"
                    className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap"
                >
                    {t("participant.certificates.detail.copyShareableUrl", "Copy shareable URL")}
                </Typography>
            </button>

            {/* Edit password button */}
            <button
                onClick={() => {}}
                className={cn(
                    "cursor-pointer flex items-center justify-center gap-2 px-4 h-10 bg-white rounded-[10px] hover:bg-white/90 transition-colors flex-shrink-0",
                )}
                aria-label={t("participant.certificates.detail.editPassword", "Edit password")}
            >
                <KeyRound className="w-5 h-5 text-background-alt" />
                <Typography
                    variant="text"
                    tag="span"
                    color="background-alt"
                    className="text-xs font-normal leading-normal tracking-[0.06px] whitespace-nowrap"
                >
                    {t("participant.certificates.detail.editPassword", "Edit password")}
                </Typography>
            </button>
        </CertificateDetailNav>
    );
};
