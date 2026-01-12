import { useState, useEffect } from "react";
import { env } from "@/config/env";

/**
 * Fetches the raw SVG content of an event's certificate template through the
 * backend proxy endpoint (`GET /api/v1/events/:eventId/config/certificate/template`).
 *
 * Using the backend as a proxy avoids the CORS issue that occurs when calling
 * the S3 presigned URL directly via fetch() — the browser blocks those cross-origin
 * requests unless the bucket has explicit CORS headers, whereas <img src> is exempt.
 *
 * Returns an empty string until the SVG is loaded or if loading fails, allowing
 * the caller to fall back to <img> rendering.
 */
export const useFetchedCertificateSvg = (
    eventId: string | undefined,
    hasNewFileUploaded: boolean,
): string => {
    const [fetchedSvg, setFetchedSvg] = useState<string>("");

    useEffect(() => {
        if (!eventId || hasNewFileUploaded) {
            setFetchedSvg("");
            return;
        }

        let cancelled = false;

        fetch(`${env.VITE_CORE_BACKEND_API}/api/v1/events/${eventId}/config/certificate/template`, {
            credentials: "include",
        })
            .then((res) => {
                if (!res.ok) throw new Error(`HTTP ${res.status}`);
                return res.text();
            })
            .then((text) => {
                if (cancelled) return;
                const trimmed = text.trimStart();
                if (trimmed.startsWith("<svg") || trimmed.startsWith("<?xml")) {
                    setFetchedSvg(text);
                }
            })
            .catch(() => {
                // Silently fall back to <img> rendering on failure
            });

        return () => {
            cancelled = true;
        };
    }, [eventId, hasNewFileUploaded]);

    return fetchedSvg;
};
