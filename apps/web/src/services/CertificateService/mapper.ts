export interface CertificateImage {
    url: string;
    blob: Blob;
    contentType: string;
}

/**
 * Converts a Blob to a CertificateImage object with object URL
 * @param blob - The image blob from the API
 * @returns CertificateImage with object URL
 */
export const mapBlobToCertificateImage = (blob: Blob): CertificateImage => {
    const url = URL.createObjectURL(blob);
    return {
        url,
        blob,
        contentType: blob.type || "image/png",
    };
};
