import { Helmet } from 'react-helmet-async';

interface FaviconHelmetProps {
    title?: string;
    description?: string;
}

export const FaviconHelmet = ({
    title = 'DECM - Decentralized Event Management',
    description = 'Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification'
}: FaviconHelmetProps) => {
    return (
        <Helmet>
            {/* Basic meta tags */}
            <title>{title}</title>
            <meta name="description" content={description} />

            {/* Standard favicon */}
            <link rel="icon" href="/favicon.ico" />

            {/* PNG favicons with sizes */}
            <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
            <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
            <link rel="icon" type="image/png" sizes="192x192" href="/android-chrome-192x192.png" />
            <link rel="icon" type="image/png" sizes="512x512" href="/android-chrome-512x512.png" />

            {/* Apple Touch Icon */}
            <link rel="apple-touch-icon" href="/apple-touch-icon.png" />

            {/* Web App Manifest */}
            <link rel="manifest" href="/site.webmanifest" />

            {/* Theme color for mobile browsers */}
            <meta name="theme-color" content="#0F1012" />

            {/* Additional meta tags for better SEO */}
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <meta httpEquiv="X-UA-Compatible" content="ie=edge" />
        </Helmet>
    );
};
