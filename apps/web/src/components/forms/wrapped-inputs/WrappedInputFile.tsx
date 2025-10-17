import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";
import { Upload, X, File as FileIcon } from "lucide-react";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

interface WrappedInputFileProps<T extends FieldValues> {
    name: Path<T>;
    control: Control<T>;
    label: string;
    accept?: string;
    required?: boolean;
    disabled?: boolean;
    className?: string;
    maxSize?: number; // in bytes
    previewClassName?: string;
}

// Separate component for file input with preview
const FileInputWithPreview = ({
    value,
    onChange,
    name,
    accept,
    disabled,
    className,
    maxSize,
    previewClassName,
    label,
    required,
    error,
}: {
    value: File | undefined;
    onChange: (file: File | undefined) => void;
    name: string;
    accept: string;
    disabled: boolean;
    className: string;
    maxSize: number;
    previewClassName: string;
    label: string;
    required: boolean;
    error?: { message?: string };
}) => {
    const { t } = useTranslation();
    const [preview, setPreview] = useState<string | null>(null);
    const [fileInfo, setFileInfo] = useState<{ name: string; size: string; type: string } | null>(
        null,
    );

    const formatFileSize = (bytes: number): string => {
        if (bytes === 0) return "0 Bytes";
        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i];
    };

    const isImageFile = (type: string): boolean => {
        return type.startsWith("image/");
    };

    // Update preview when file value changes
    useEffect(() => {
        if (value) {
            // Set file info
            setFileInfo({
                name: value.name,
                size: formatFileSize(value.size),
                type: value.type,
            });

            // Generate preview if image
            if (isImageFile(value.type)) {
                const reader = new FileReader();
                reader.onloadend = () => {
                    setPreview(reader.result as string);
                };
                reader.readAsDataURL(value);
            } else {
                setPreview(null);
            }
        } else {
            setPreview(null);
            setFileInfo(null);
        }
    }, [value]);

    const handleFileChange = (file: File | undefined) => {
        if (!file) {
            setPreview(null);
            setFileInfo(null);
            onChange(undefined);
            return;
        }

        // Set file info
        setFileInfo({
            name: file.name,
            size: formatFileSize(file.size),
            type: file.type,
        });

        // Generate preview if image
        if (isImageFile(file.type)) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        } else {
            setPreview(null);
        }

        onChange(file);
    };

    const handleRemoveFile = () => {
        setPreview(null);
        setFileInfo(null);
        onChange(undefined);
    };

    return (
        <div className="space-y-2">
            <Label htmlFor={name}>
                <Typography variant="text" tag="span" className="text-sm font-medium">
                    {label}
                </Typography>
                {required && <span className="text-destructive ml-1">*</span>}
            </Label>

            {!value ? (
                <div className="relative">
                    <input
                        type="file"
                        id={name}
                        accept={accept}
                        className="hidden"
                        onChange={(e) => handleFileChange(e.target.files?.[0])}
                        disabled={disabled}
                    />
                    <label
                        htmlFor={name}
                        className={cn(
                            "flex flex-col items-center justify-center w-full h-48 border-2 border-dashed rounded-lg cursor-pointer transition-colors",
                            "border-input bg-background hover:bg-accent/50",
                            disabled && "opacity-50 cursor-not-allowed",
                            className,
                        )}
                    >
                        <div className="flex flex-col items-center justify-center pt-5 pb-6">
                            <Upload className="w-10 h-10 mb-3 text-muted-foreground" />
                            <Typography
                                variant="text"
                                tag="p"
                                className="mb-2 text-sm text-muted-foreground"
                            >
                                <span className="font-semibold">{t("common.clickToUpload")}</span>
                            </Typography>
                            <Typography
                                variant="text"
                                tag="p"
                                className="text-xs text-muted-foreground"
                            >
                                {accept.includes("image")
                                    ? `PNG, JPG, or WebP (MAX. ${formatFileSize(maxSize)})`
                                    : `MAX. ${formatFileSize(maxSize)}`}
                            </Typography>
                        </div>
                    </label>
                </div>
            ) : (
                <div className="space-y-3">
                    {/* Image Preview or File Info */}
                    {preview ? (
                        <div
                            className={cn(
                                "relative w-full h-72 rounded-lg overflow-hidden border border-[#D9D9D91A]",
                                previewClassName,
                            )}
                        >
                            <img
                                src={preview}
                                alt={fileInfo?.name || "File preview"}
                                className="w-full h-full object-cover"
                            />
                        </div>
                    ) : fileInfo ? (
                        <div className="flex items-center gap-4 p-4 border border-[#D9D9D91A] rounded-lg bg-background">
                            <div className="flex-shrink-0">
                                <FileIcon className="w-10 h-10 text-muted-foreground" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-sm font-medium truncate"
                                >
                                    {fileInfo.name}
                                </Typography>
                                <Typography
                                    variant="text"
                                    tag="p"
                                    className="text-xs text-muted-foreground"
                                >
                                    {fileInfo.size}
                                </Typography>
                            </div>
                        </div>
                    ) : null}

                    {/* Action Buttons */}
                    <div className="flex gap-2">
                        <label htmlFor={`${name}-change`} className="flex-1">
                            <input
                                type="file"
                                id={`${name}-change`}
                                accept={accept}
                                className="hidden"
                                onChange={(e) => handleFileChange(e.target.files?.[0])}
                                disabled={disabled}
                            />
                            <Button
                                type="button"
                                variant="secondary-dark"
                                size="default"
                                disabled={disabled}
                                className="w-full cursor-pointer"
                                onClick={() => document.getElementById(`${name}-change`)?.click()}
                            >
                                <Upload className="h-4 w-4 mr-2" />
                                <Typography variant="text" tag="span" className="font-medium">
                                    {t("common.change")}
                                </Typography>
                            </Button>
                        </label>
                        <Button
                            type="button"
                            variant="secondary-dark"
                            size="default"
                            onClick={handleRemoveFile}
                            disabled={disabled}
                            className="flex-1"
                        >
                            <X className="h-4 w-4 mr-2" />
                            <Typography variant="text" tag="span" className="font-medium">
                                {t("common.remove")}
                            </Typography>
                        </Button>
                    </div>
                </div>
            )}

            {error && (
                <Typography
                    variant="text"
                    tag="p"
                    className="text-sm text-destructive"
                    role="alert"
                >
                    {t(error.message as string)}
                </Typography>
            )}
        </div>
    );
};

export const WrappedInputFile = <T extends FieldValues>({
    name,
    control,
    label,
    accept = "image/jpeg,image/jpg,image/png,image/webp",
    required = false,
    disabled = false,
    className = "",
    maxSize = 5 * 1024 * 1024, // 5MB default
    previewClassName = "",
}: WrappedInputFileProps<T>) => {
    return (
        <Controller
            name={name}
            control={control}
            render={({ field: { value, onChange }, fieldState: { error } }) => (
                <FileInputWithPreview
                    value={value}
                    onChange={onChange}
                    name={name}
                    accept={accept}
                    disabled={disabled}
                    className={className}
                    maxSize={maxSize}
                    previewClassName={previewClassName}
                    label={label}
                    required={required}
                    error={error}
                />
            )}
        />
    );
};
