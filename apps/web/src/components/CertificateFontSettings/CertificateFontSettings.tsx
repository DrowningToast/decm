import React from "react";
import { Typography } from "@/components/typography/typography";
import { Label } from "@/components/ui/label";
import { useCertificateFontSettings } from "./useCertificateFontSettings";
import { FontSettingRow } from "./FontSettingRow";
import type { CertificateFontSettingsProps } from "./types";

export const CertificateFontSettings: React.FC<CertificateFontSettingsProps> = (props) => {
    const {
        t,
        fontFamilies,
        isLoadingFonts,
        defaultFontFamilyId,
        visibleFields,
        hasAnyFields,
        getAvailableWeights,
    } = useCertificateFontSettings(props);

    if (isLoadingFonts) {
        return (
            <div className="space-y-4 rounded-lg border p-4">
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t("common.loading", "Loading...")}
                </Typography>
            </div>
        );
    }

    return (
        <div className="space-y-3 rounded-lg border p-3">
            <div className="space-y-0.5">
                <Typography variant="header" tag="h3" className="text-base font-semibold">
                    {t("certificateSettings.fontSettings.title", "Font Settings")}
                </Typography>
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t(
                        "certificateSettings.fontSettings.description",
                        "Configure font family and weight for each text field in the certificate template.",
                    )}
                </Typography>
            </div>

            {!hasAnyFields ? (
                <div className="rounded-md bg-muted/50 p-3">
                    <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                        {t(
                            "certificateSettings.fontSettings.noFieldsConfigured",
                            "No text fields have been positioned in the certificate template yet. Configure text positions first to manage font settings.",
                        )}
                    </Typography>
                </div>
            ) : (
                <div className="space-y-1">
                    {/* Column headers */}
                    <div className="grid grid-cols-[1fr_1.5fr_0.75fr] gap-2 px-1">
                        <Label className="text-xs text-muted-foreground">
                            {t("certificateSettings.fontSettings.field", "Field")}
                        </Label>
                        <Label className="text-xs text-muted-foreground">
                            {t("certificateSettings.fontSettings.fontFamily", "Font Family")}
                        </Label>
                        <Label className="text-xs text-muted-foreground">
                            {t("certificateSettings.fontSettings.fontWeight", "Weight")}
                        </Label>
                    </div>

                    {/* Font rows */}
                    {visibleFields.map((field) => (
                        <FontSettingRow
                            key={field.key}
                            label={field.label}
                            familyId={field.familyId}
                            familyValue={field.familyValue}
                            weightValue={field.weightValue}
                            defaultFontFamilyId={defaultFontFamilyId}
                            fontFamilies={fontFamilies}
                            availableWeights={getAvailableWeights(field.familyValue)}
                            onFamilyChange={field.onFamilyChange}
                            onWeightChange={field.onWeightChange}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};
