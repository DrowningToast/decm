/**
 * Certificate Image representation
 * Contains the blob URL and metadata for rendering certificate images
 */
export interface CertificateImage {
    /** Blob URL that can be used as src attribute */
    url: string;
    /** Certificate ID this image belongs to */
    certificateId: string;
    /** MIME type of the image */
    mimeType: string;
    /** Size of the image in bytes */
    size: number;
}

/**
 * Maps a Blob response to a CertificateImage object
 * Creates an object URL from the blob for rendering
 */
export const mapBlobToCertificateImage = (blob: Blob, certificateId: string): CertificateImage => {
    const url = URL.createObjectURL(blob);
    return {
        url,
        certificateId,
        mimeType: blob.type,
        size: blob.size,
    };
};
