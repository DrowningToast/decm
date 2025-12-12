import { Helmet } from "react-helmet-async";

interface CustomHelmetProps {
    title?: string;
    description?: string;
    faviconUrl?: string;
    appleTouchIconUrl?: string;
    androidIconUrls?: {
        icon192: string;
        icon512: string;
    };
    themeColor?: string;
    additionalMeta?: Array<{
        name?: string;
        property?: string;
        content: string;
    }>;
}

export const CustomHelmet = ({
    title = "DECM - Decentralized Event Management",
    description = "Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification",
    faviconUrl = "/favicon.ico",
    appleTouchIconUrl = "/apple-touch-icon.png",
    androidIconUrls = {
        icon192: "/android-chrome-192x192.png",
        icon512: "/android-chrome-512x512.png",
    },
    themeColor = "#ffffff",
    additionalMeta = [],
}: CustomHelmetProps) => {
    return (
        <Helmet>
            {/* Basic meta tags */}
            <title>{title}</title>
            <meta name="description" content={description} />

            {/* Custom favicon */}
            <link rel="icon" href={faviconUrl} />

            {/* PNG favicons with sizes */}
            <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
            <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
            <link rel="icon" type="image/png" sizes="192x192" href={androidIconUrls.icon192} />
            <link rel="icon" type="image/png" sizes="512x512" href={androidIconUrls.icon512} />

            {/* Apple Touch Icon */}
            <link rel="apple-touch-icon" href={appleTouchIconUrl} />

            {/* Web App Manifest */}
            <link rel="manifest" href="/site.webmanifest" />

            {/* Theme color */}
            <meta name="theme-color" content={themeColor} />

            {/* Additional meta tags */}
            {additionalMeta.map((meta, index) => (
                <meta
                    key={index}
                    {...(meta.name ? { name: meta.name } : {})}
                    {...(meta.property ? { property: meta.property } : {})}
                    content={meta.content}
                />
            ))}

            {/* Standard meta tags */}
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <meta httpEquiv="X-UA-Compatible" content="ie=edge" />
        </Helmet>
    );
};
