import { useEffect, useState } from "react";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

interface GoogleMapsEmbedProps {
    query: string;
    className?: string;
    height?: string;
}

export const GoogleMapsEmbed = ({
    query,
    className = "",
    height = "400px",
}: GoogleMapsEmbedProps) => {
    const { t } = useTranslation();
    const [mapUrl, setMapUrl] = useState<string>("");
    const [hasApiKey, setHasApiKey] = useState<boolean>(true);
    const [hasError, setHasError] = useState<boolean>(false);

    useEffect(() => {
        if (query.trim()) {
            // Encode the query for URL
            const encodedQuery = encodeURIComponent(query);
            // Get API key from environment variable
            const apiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY;

            if (!apiKey) {
                console.error(
                    "Google Maps API key is missing. Please set VITE_GOOGLE_MAPS_API_KEY in your .env.local file",
                );
                setHasApiKey(false);
                setMapUrl("");
                return;
            }

            setHasApiKey(true);
            setHasError(false);
            // Generate Google Maps embed URL
            const url = `https://www.google.com/maps/embed/v1/place?key=${apiKey}&q=${encodedQuery}`;
            setMapUrl(url);
        } else {
            setMapUrl("");
            setHasError(false);
        }
    }, [query]);

    // Show placeholder when no query
    if (!query.trim()) {
        return (
            <div
                className={cn(
                    "flex items-center justify-center border border-[#D9D9D91A] rounded-lg bg-muted/10",
                    className,
                )}
                style={{ height }}
            >
                <Typography variant="text" tag="p" className="text-sm text-muted-foreground">
                    {t("events.form.googleMapPlaceholder")}
                </Typography>
            </div>
        );
    }

    // Show error when API key is missing
    if (!hasApiKey) {
        return (
            <div
                className={cn(
                    "flex flex-col items-center justify-center border border-destructive/50 rounded-lg bg-destructive/10 p-4",
                    className,
                )}
                style={{ height }}
            >
                <Typography
                    variant="text"
                    tag="p"
                    className="text-sm text-destructive font-medium mb-2"
                >
                    ⚠️ Google Maps API Key Missing
                </Typography>
                <Typography
                    variant="text"
                    tag="p"
                    className="text-xs text-muted-foreground text-center"
                >
                    Please set VITE_GOOGLE_MAPS_API_KEY in your .env.local file
                </Typography>
            </div>
        );
    }

    // Show error when map fails to load (403 error)
    if (hasError) {
        return (
            <div
                className={cn(
                    "flex flex-col items-center justify-center border border-destructive/50 rounded-lg bg-destructive/10 p-4",
                    className,
                )}
                style={{ height }}
            >
                <Typography
                    variant="text"
                    tag="p"
                    className="text-sm text-destructive font-medium mb-2"
                >
                    ⚠️ Map Load Error (403)
                </Typography>
                <Typography
                    variant="text"
                    tag="p"
                    className="text-xs text-muted-foreground text-center mb-2"
                >
                    The API key restrictions may be blocking this request.
                </Typography>
                <Typography
                    variant="text"
                    tag="p"
                    className="text-xs text-muted-foreground text-center"
                >
                    Please check API key restrictions in Google Cloud Console
                </Typography>
            </div>
        );
    }

    // Only render iframe if mapUrl is not empty
    if (!mapUrl) {
        return (
            <div
                className={cn(
                    "flex items-center justify-center border border-[#D9D9D91A] rounded-lg bg-muted/10",
                    className,
                )}
                style={{ height }}
            >
                <Typography
                    variant="text"
                    tag="p"
                    className="text-sm text-muted-foreground animate-pulse"
                >
                    Loading map...
                </Typography>
            </div>
        );
    }

    return (
        <div
            className={cn(
                "relative w-full overflow-hidden rounded-lg border border-[#D9D9D91A]",
                className,
            )}
            style={{ height }}
        >
            <iframe
                width="100%"
                height="100%"
                loading="lazy"
                allowFullScreen
                referrerPolicy="no-referrer-when-downgrade"
                src={mapUrl}
                title={`Google Maps - ${query}`}
                className="absolute inset-0 border-0"
                onError={() => setHasError(true)}
            />
        </div>
    );
};
