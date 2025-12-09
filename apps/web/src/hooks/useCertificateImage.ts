import { useState, useEffect, useRef } from "react";
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
    const imageUrlRef = useRef<string | null>(null);

    useEffect(() => {
        if (!certificateId || !enabled) {
            // Revoke previous URL if it exists
            if (imageUrlRef.current) {
                URL.revokeObjectURL(imageUrlRef.current);
                imageUrlRef.current = null;
            }
            setImageData(null);
            return;
        }

        let isMounted = true;
        const controller = new AbortController();
        // Store URL in closure so cleanup can access it
        let currentImageUrl: string | null = null;

        const fetchImage = async () => {
            setIsLoading(true);
            setError(null);

            try {
                const certificateImage =
                    await certificateService.getCertificateImage(certificateId);

                if (isMounted) {
                    // Revoke previous URL if it exists
                    if (imageUrlRef.current) {
                        URL.revokeObjectURL(imageUrlRef.current);
                    }
                    currentImageUrl = certificateImage.url;
                    imageUrlRef.current = currentImageUrl;
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
            // Use the closure variable first, then fall back to ref
            // This ensures we have access to the URL even if the ref is cleared
            const urlToRevoke = currentImageUrl || imageUrlRef.current;
            if (urlToRevoke) {
                URL.revokeObjectURL(urlToRevoke);
                imageUrlRef.current = null;
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
