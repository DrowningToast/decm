import { useEffect, useRef, useState } from "react";
import { certificateService } from "@/services/services";
import type { CertificateImage } from "@/services/CertificateService/mapper";

interface UseCertificateShareImageOptions {
    handle?: string;
    password?: string;
    enabled?: boolean;
}

interface UseCertificateShareImageReturn {
    imageUrl: string | null;
    imageData: CertificateImage | null;
    isLoading: boolean;
    error: Error | null;
}

/**
 * Hook to fetch a certificate image via a public share handle.
 * Mirrors the useCertificateImage pattern but calls the unauthenticated
 * certificate-shares image endpoint, forwarding an optional password for
 * password-protected shares.
 */
export const useCertificateShareImage = ({
    handle,
    password,
    enabled = true,
}: UseCertificateShareImageOptions): UseCertificateShareImageReturn => {
    const [imageData, setImageData] = useState<CertificateImage | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    const imageUrlRef = useRef<string | null>(null);

    useEffect(() => {
        if (!handle || !enabled) {
            if (imageUrlRef.current) {
                URL.revokeObjectURL(imageUrlRef.current);
                imageUrlRef.current = null;
            }
            setImageData(null);
            return;
        }

        let isMounted = true;
        const controller = new AbortController();
        let currentImageUrl: string | null = null;

        const fetchImage = async () => {
            setIsLoading(true);
            setError(null);

            try {
                const certificateImage = await certificateService.getCertificateShareImage(
                    handle,
                    password,
                );

                if (isMounted) {
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
                }
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        };

        fetchImage();

        return () => {
            isMounted = false;
            controller.abort();
            const urlToRevoke = currentImageUrl || imageUrlRef.current;
            if (urlToRevoke) {
                URL.revokeObjectURL(urlToRevoke);
                imageUrlRef.current = null;
            }
        };
    }, [handle, password, enabled]);

    return {
        imageUrl: imageData?.url ?? null,
        imageData,
        isLoading,
        error,
    };
};
