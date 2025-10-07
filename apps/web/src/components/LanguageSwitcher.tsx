import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { languages, type Language } from '@/lib/i18n';

export function LanguageSwitcher() {
    const { i18n } = useTranslation();

    const changeLanguage = (lng: Language) => {
        i18n.changeLanguage(lng);
    };

    const currentLanguage = (i18n.language as Language) || 'en';
    const currentLangData = languages[currentLanguage] || languages.en;

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button size="sm" className="gap-2">
                    <span>{currentLangData.flag}</span>
                    <span className="hidden sm:inline-block">{currentLangData.label}</span>
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                {Object.entries(languages).map(([code, { label, flag }]) => (
                    <DropdownMenuItem
                        key={code}
                        onClick={() => changeLanguage(code as Language)}
                        className="gap-2 cursor-pointer"
                    >
                        <span>{flag}</span>
                        <span>{label}</span>
                        {code === currentLanguage && (
                            <span className="ml-auto text-xs">✓</span>
                        )}
                    </DropdownMenuItem>
                ))}
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

