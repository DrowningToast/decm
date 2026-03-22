import { Search } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Typography } from "@/components/typography/typography";

export interface VerifySearchFormProps {
    inputCode: string;
    onInputCodeChange: (value: string) => void;
    isFetching: boolean;
    onSubmit: (e: React.FormEvent) => void;
}

export function VerifySearchForm({
    inputCode,
    onInputCodeChange,
    isFetching,
    onSubmit,
}: VerifySearchFormProps) {
    const { t } = useTranslation();

    return (
        <form onSubmit={onSubmit} className="flex flex-col gap-3">
            <Typography variant="text" tag="p" color="muted" className="text-sm text-center">
                {t("certificateVerify.enterCodeHint")}
            </Typography>
            <div className="flex items-stretch gap-2">
                <Input
                    type="text"
                    placeholder={t("certificateVerify.codePlaceholder")}
                    value={inputCode}
                    onChange={(e) => onInputCodeChange(e.target.value)}
                    aria-label={t("certificateVerify.codePlaceholder")}
                    className="flex-1 h-auto"
                />
                <Button
                    type="submit"
                    variant="primary"
                    disabled={!inputCode.trim() || isFetching}
                    loading={isFetching}
                >
                    {!isFetching && <Search className="w-4 h-4" />}
                    {t("certificateVerify.verifyButton")}
                </Button>
            </div>
        </form>
    );
}
