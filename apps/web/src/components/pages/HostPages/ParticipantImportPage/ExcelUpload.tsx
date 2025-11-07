import { useRef } from "react";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";

interface ExcelUploadProps {
    onFileSelect: (file: File) => void;
    selectedFile: File | null;
    disabled?: boolean;
}

export const ExcelUpload = ({ onFileSelect, selectedFile, disabled = false }: ExcelUploadProps) => {
    const { t } = useTranslation();
    const fileInputRef = useRef<HTMLInputElement | null>(null);

    const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            // Validate file type
            const validTypes = [
                "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", // .xlsx
                "application/vnd.ms-excel", // .xls
            ];

            if (!validTypes.includes(file.type)) {
                // TODO: Show error message
                return;
            }

            // Validate file size (max 10MB)
            const maxSize = 10 * 1024 * 1024; //10MB in bytes
            if (file.size > maxSize) {
                // TODO: Show error message
                return;
            }

            onFileSelect(file);
        }
    };

    const handleClick = () => {
        fileInputRef.current?.click();
    };

    return (
        <div className="border-2 border-dashed border-back rounded-lg p-8 text-center">
            <input
                ref={fileInputRef}
                type="file"
                id="excel-upload"
                accept=".xlsx,.xls"
                onChange={handleFileSelect}
                disabled={disabled}
                className="hidden"
            />
            <div
                onClick={handleClick}
                className={`cursor-pointer ${disabled ? "opacity-50 pointer-events-none" : ""}`}
            >
                <div className="flex flex-col items-center">
                    <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="h-6 w-6 text-primary"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                            />
                        </svg>
                    </div>
                    <Typography
                        variant="text"
                        tag="p"
                        color="background-alt"
                        className="font-medium mb-2"
                    >
                        {selectedFile ? selectedFile.name : t("participantImport.selectFile")}
                    </Typography>
                    <Typography variant="text" tag="p" color="background" className="text-sm">
                        {t("participantImport.fileFormats")}
                    </Typography>
                    {selectedFile && (
                        <Typography
                            variant="text"
                            tag="p"
                            color="background"
                            className="text-xs mt-2"
                        >
                            {t("participantImport.fileSize")}:{" "}
                            {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                        </Typography>
                    )}
                </div>
            </div>
        </div>
    );
};
