import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import { languages, type Language } from "@/lib/i18n";
import { Check } from "lucide-react";
import { cn } from "@/lib/utils";

export const PreferencePage = () => {
    const { t, i18n } = useTranslation();
    const currentLanguage = (i18n.language as Language) || "en";

    const handleLanguageChange = (lng: Language) => {
        i18n.changeLanguage(lng);
    };

    return (
        <section className="relative z-10 w-full">
            <div className="relative w-full overflow-hidden">
                {/* Background image */}
                <div className="absolute bottom-0 right-0 w-[424px] h-[424px] md:w-[500px] md:h-[500px] opacity-20 pointer-events-none">
                    <img
                        src="/assets/scale.webp"
                        alt=""
                        className="w-full h-full object-cover object-center"
                    />
                </div>

                {/* Main content */}
                <div className="relative z-10 w-full max-w-[1384px] mx-auto px-4 md:px-16 py-4 md:pt-16 md:pb-24 flex flex-col gap-y-4 md:gap-y-6">
                    {/* Header section */}
                    <div className="flex flex-col gap-1.5">
                        <Typography
                            variant="header"
                            tag="h1"
                            color="primary"
                            className="text-[28px]/[34px] [text-shadow:rgba(255,255,255,0.2)_0px_0px_4px] font-header"
                        >
                            {t("preference.title")}
                        </Typography>
                        <Typography
                            variant="text"
                            tag="p"
                            color="muted"
                            className="text-base/[22px] [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                        >
                            {t("preference.description")}
                        </Typography>
                    </div>

                    {/* Preferences Section */}
                    <div className="flex flex-col gap-y-6">
                        {/* Language Preference */}
                        <div className="flex flex-col gap-y-4">
                            <div className="flex flex-col gap-y-1">
                                <Typography
                                    variant="text"
                                    tag="h2"
                                    color="foreground"
                                    className="text-lg md:text-xl font-semibold [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                                >
                                    {t("preference.language.title")}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    color="muted"
                                    className="text-sm md:text-base [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                                >
                                    {t("preference.language.description")}
                                </Typography>
                            </div>

                            {/* Language Options */}
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                {Object.entries(languages).map(([code, { label, flag }]) => (
                                    <button
                                        key={code}
                                        onClick={() => handleLanguageChange(code as Language)}
                                        className={cn(
                                            "relative flex items-center gap-3 p-4 rounded-lg border transition-all duration-200 text-left",
                                            "hover:border-primary hover:bg-primary/5",
                                            code === currentLanguage
                                                ? "border-primary bg-primary/10"
                                                : "border-border bg-background/50",
                                        )}
                                        aria-label={`${t("preference.language.select")} ${label}`}
                                    >
                                        {/* Flag */}
                                        <span className="text-3xl flex-shrink-0">{flag}</span>

                                        {/* Language Info */}
                                        <div className="flex-1 min-w-0">
                                            <Typography
                                                variant="text"
                                                tag="p"
                                                color={
                                                    code === currentLanguage
                                                        ? "primary"
                                                        : "foreground"
                                                }
                                                className="font-medium"
                                            >
                                                {label}
                                            </Typography>
                                            <Typography
                                                variant="text"
                                                tag="p"
                                                color="muted"
                                                className="text-xs md:text-sm"
                                            >
                                                {code.toUpperCase()}
                                            </Typography>
                                        </div>

                                        {/* Check Icon */}
                                        {code === currentLanguage && (
                                            <div className="absolute right-4 flex-shrink-0">
                                                <Check className="w-5 h-5 text-primary" />
                                            </div>
                                        )}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {/* Additional Preferences Section (Future) */}
                        <div className="flex flex-col gap-y-4 pt-6 border-t border-border">
                            <Typography
                                variant="text"
                                tag="h2"
                                color="muted"
                                className="text-sm [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px]"
                            >
                                {t("preference.moreComingSoon")}
                            </Typography>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};
