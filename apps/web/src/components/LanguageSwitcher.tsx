import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { languages, type Language } from "@/lib/i18n";
import { Typography } from "@/components/typography/typography";

export function LanguageSwitcher() {
    const { i18n } = useTranslation();

    const changeLanguage = (lng: Language) => {
        i18n.changeLanguage(lng);
    };

    const currentLanguage = (i18n.language as Language) || "en";
    const currentLangData = languages[currentLanguage] || languages.en;

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button size="sm" className="gap-2">
                    <Typography variant="text" tag="span">
                        {currentLangData.flag}
                    </Typography>
                    <Typography variant="text" tag="span" className="hidden sm:inline-block">
                        {currentLangData.label}
                    </Typography>
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                {Object.entries(languages).map(([code, { label, flag }]) => (
                    <DropdownMenuItem
                        key={code}
                        onClick={() => changeLanguage(code as Language)}
                        className="gap-2 cursor-pointer"
                    >
                        <Typography variant="text" tag="span">
                            {flag}
                        </Typography>
                        <Typography variant="text" tag="span">
                            {label}
                        </Typography>
                        {code === currentLanguage && (
                            <Typography variant="text" tag="span" className="ml-auto text-xs">
                                ✓
                            </Typography>
                        )}
                    </DropdownMenuItem>
                ))}
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
