import React from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "@/components/typography/typography";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import type { CertificateFontFamily } from "@/services/CertificateService/CertificateService";

interface FontSettingRowProps {
    label: string;
    familyId: string;
    familyValue: number | undefined;
    weightValue: number | undefined;
    defaultFontFamilyId: number;
    fontFamilies: CertificateFontFamily[];
    availableWeights: number[];
    onFamilyChange: (value: string) => void;
    onWeightChange: (value: string) => void;
}

export const FontSettingRow: React.FC<FontSettingRowProps> = ({
    label,
    familyId,
    familyValue,
    weightValue,
    defaultFontFamilyId,
    fontFamilies,
    availableWeights,
    onFamilyChange,
    onWeightChange,
}) => {
    const { t } = useTranslation();

    return (
        <div className="grid grid-cols-[1fr_1.5fr_0.75fr] items-center gap-2 rounded-md px-1 py-1 hover:bg-muted/30">
            <Typography variant="text" tag="span" className="truncate text-sm font-medium">
                {label}
            </Typography>
            <Select
                value={familyValue?.toString() || defaultFontFamilyId.toString()}
                onValueChange={onFamilyChange}
            >
                <SelectTrigger id={familyId} className="h-8 text-xs">
                    <SelectValue
                        placeholder={t(
                            "certificateSettings.fontSettings.selectFont",
                            "Select font",
                        )}
                    />
                </SelectTrigger>
                <SelectContent>
                    {fontFamilies
                        .filter((font) => font?.id != null)
                        .map((font) => (
                            <SelectItem key={font.id} value={font.id.toString()}>
                                {font.font_family_name}
                            </SelectItem>
                        ))}
                </SelectContent>
            </Select>
            <Select value={weightValue?.toString() || "700"} onValueChange={onWeightChange}>
                <SelectTrigger id={`${familyId}-weight`} className="h-8 text-xs">
                    <SelectValue
                        placeholder={t(
                            "certificateSettings.fontSettings.selectWeight",
                            "Select weight",
                        )}
                    />
                </SelectTrigger>
                <SelectContent>
                    {availableWeights
                        .filter((weight) => weight != null)
                        .map((weight) => (
                            <SelectItem key={weight} value={weight.toString()}>
                                {weight}
                            </SelectItem>
                        ))}
                </SelectContent>
            </Select>
        </div>
    );
};
