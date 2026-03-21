import { FileCheck2, Share2, ShieldCheck, ChevronRight } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";

export function HowItWorks() {
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
