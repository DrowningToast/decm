import { useTranslation } from "react-i18next";
import { type ParticipantSettingsData } from "@/lib/schemas/participantSettingsSchema";
import { Typography } from "@/components/typography/typography";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Lock, Globe, Mail } from "lucide-react";

interface RegistrationFormPreviewProps {
    /**
     * Participant settings configuration
     */
    settings: ParticipantSettingsData;
}

export const RegistrationFormPreview = ({ settings }: RegistrationFormPreviewProps) => {
    const { t } = useTranslation();

    const renderField = (
        fieldName: keyof Pick<
            ParticipantSettingsData,
            | "firstName"
            | "lastName"
            | "email"
            | "bio"
            | "phoneNumber"
            | "address"
            | "academicInstitution"
            | "academicEmail"
        >,
        label: string,
        type: "text" | "email" | "textarea" = "text",
    ) => {
        const requirement = settings[fieldName];

        if (requirement === "not_required") {
            return null;
        }

        const isRequired = requirement === "required";

        return (
            <div key={fieldName} className="space-y-2">
                <Label htmlFor={`preview-${fieldName}`}>
                    <Typography variant="text" tag="span" className="text-sm font-medium">
                        {label}
                        {isRequired && (
                            <Typography variant="text" tag="span" className="text-destructive ml-1">
                                *
                            </Typography>
                        )}
                    </Typography>
                </Label>
                {type === "textarea" ? (
                    <Textarea
                        id={`preview-${fieldName}`}
                        placeholder={`${t("participantSettings.preview.enter")} ${label.toLowerCase()}`}
                        disabled
                        className="opacity-60"
                    />
                ) : (
                    <Input
                        id={`preview-${fieldName}`}
                        type={type}
                        placeholder={`${t("participantSettings.preview.enter")} ${label.toLowerCase()}`}
                        disabled
                        className="opacity-60"
                    />
                )}
                {!isRequired && (
                    <Typography variant="text" tag="p" className="text-xs text-muted-foreground">
                        {t("participantSettings.preview.optional")}
                    </Typography>
                )}
            </div>
        );
    };

    return (
        <div className="space-y-6">
            {/* Preview Header */}
            <div className="space-y-2">
                <div className="flex items-center justify-between">
                    <Typography variant="header" tag="h2" className="text-xl font-bold">
                        {t("participantSettings.preview.title")}
                    </Typography>
                    <div className="flex items-center space-x-2">
                        {settings.eventType === "invite" ? (
                            <>
                                <Mail className="h-4 w-4 text-orange-500" />
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-sm text-orange-500 font-medium"
                                >
                                    {t("participantSettings.eventTypeInviteOnly")}
                                </Typography>
                            </>
                        ) : (
                            <>
                                <Lock className="h-4 w-4 text-muted-foreground" />
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-sm text-muted-foreground font-medium"
                                >
                                    {t("participantSettings.eventTypePrivate")}
                                </Typography>
                            </>
                        )}
                    </div>
                </div>
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t("participantSettings.preview.description")}
                </Typography>
            </div>

            {/* Preview Form */}
            {/* <div className="rounded-lg border bg-muted/20 p-6 space-y-6"> */}
            {/* Event Settings Info */}
            <div className="flex flex-wrap gap-4 pb-4 border-b">
                {settings.isBookingRequired && (
                    <div className="flex items-center space-x-2 px-3 py-1.5 rounded-md bg-background border">
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-xs font-medium text-primary"
                        >
                            {t("participantSettings.preview.bookingRequired")}
                        </Typography>
                    </div>
                )}
                {settings.finalCallRegistrationDate && (
                    <div className="flex items-center space-x-2 px-3 py-1.5 rounded-md bg-background border">
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-xs font-medium text-orange-500"
                        >
                            {t("participantSettings.finalCallRegistrationDate")}:{" "}
                            {settings.finalCallRegistrationDate.toLocaleDateString()}
                        </Typography>
                    </div>
                )}
                {settings.isTicketTransferable && (
                    <div className="flex items-center space-x-2 px-3 py-1.5 rounded-md bg-background border">
                        <Typography
                            variant="text"
                            tag="span"
                            className="text-xs font-medium text-primary"
                        >
                            {t("participantSettings.preview.transferable")}
                        </Typography>
                    </div>
                )}
            </div>

            {/* Basic Information Fields */}
            {(settings.firstName !== "not_required" ||
                settings.lastName !== "not_required" ||
                settings.email !== "not_required" ||
                settings.phoneNumber !== "not_required") && (
                <div className="space-y-4">
                    <Typography variant="header" tag="h3" className="text-sm font-semibold">
                        {t("participantSettings.basicInformation")}
                    </Typography>
                    <div className="space-y-4">
                        {renderField(
                            "firstName",
                            t("participantSettings.fields.firstName"),
                            "text",
                        )}
                        {renderField("lastName", t("participantSettings.fields.lastName"), "text")}
                        {renderField("email", t("participantSettings.fields.email"), "email")}
                        {renderField(
                            "phoneNumber",
                            t("participantSettings.fields.phoneNumber"),
                            "text",
                        )}
                    </div>
                </div>
            )}

            {/* Additional Information Fields */}
            {(settings.bio !== "not_required" || settings.address !== "not_required") && (
                <div className="space-y-4 pt-4 border-t">
                    <Typography variant="header" tag="h3" className="text-sm font-semibold">
                        {t("participantSettings.additionalInformation")}
                    </Typography>
                    <div className="space-y-4">
                        {renderField("bio", t("participantSettings.fields.bio"), "textarea")}
                        {renderField(
                            "address",
                            t("participantSettings.fields.address"),
                            "textarea",
                        )}
                    </div>
                </div>
            )}

            {/* Academic Information Fields */}
            {(settings.academicInstitution !== "not_required" ||
                settings.academicEmail !== "not_required") && (
                <div className="space-y-4 pt-4">
                    <Typography variant="header" tag="h3" className="text-sm font-semibold">
                        {t("participantSettings.academicInformation")}
                    </Typography>
                    <div className="space-y-4">
                        {renderField(
                            "academicInstitution",
                            t("participantSettings.fields.academicInstitution"),
                            "text",
                        )}
                        {renderField(
                            "academicEmail",
                            t("participantSettings.fields.academicEmail"),
                            "email",
                        )}
                    </div>
                </div>
            )}

            {/* Preview Submit Button */}
            <div className="pt-4">
                <Button disabled className="w-full opacity-60" variant="primary" size="lg">
                    <Typography variant="text" tag="span" className="font-medium">
                        {t("participantSettings.preview.register")}
                    </Typography>
                </Button>
            </div>

            {/* Preview Note */}
            <Typography
                variant="text"
                tag="p"
                className="text-xs text-center text-muted-foreground italic"
            >
                {t("participantSettings.preview.note")}
            </Typography>
            {/* </div> */}
        </div>
    );
};
