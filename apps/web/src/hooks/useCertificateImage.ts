import { useState, useEffect } from "react";
import { certificateService } from "@/services/services";
import type { CertificateImage } from "@/services/CertificateService/mapper";

interface UseCertificateImageOptions {
    certificateId?: string;
    enabled?: boolean;
}

interface UseCertificateImageReturn {
    imageUrl: string | null;
    imageData: CertificateImage | null;
    isLoading: boolean;
    error: Error | null;
}

/**
 * Hook to fetch certificate image with authentication
 * Uses the CertificateService layer to fetch images with proper type mapping
 */
export const useCertificateImage = ({
    certificateId,
    enabled = true,
}: UseCertificateImageOptions): UseCertificateImageReturn => {
    const [imageData, setImageData] = useState<CertificateImage | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        if (!certificateId || !enabled) {
            setImageData(null);
            return;
        }

        let isMounted = true;
        let currentImageUrl: string | null = null;
        const controller = new AbortController();

        const fetchImage = async () => {
            setIsLoading(true);
            setError(null);

            try {
                const certificateImage =
                    await certificateService.getCertificateImage(certificateId);
                currentImageUrl = certificateImage.url;

                if (isMounted) {
                    setImageData(certificateImage);
                }
            } catch (err) {
                if (isMounted && err instanceof Error && err.name !== "AbortError") {
                    setError(err);
                    console.error("Failed to fetch certificate image:", err);
                }
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        };

        fetchImage();

        // Cleanup: revoke object URL when component unmounts or dependencies change
        return () => {
            isMounted = false;
            controller.abort();
            if (currentImageUrl) {
                URL.revokeObjectURL(currentImageUrl);
            }
        };
    }, [certificateId, enabled]);

    return {
        imageUrl: imageData?.url ?? null,
        imageData,
        isLoading,
        error,
    };
};
