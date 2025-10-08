import { Typography } from '@/components/typography/typography';
import { useTranslation } from 'react-i18next';

export const GoogleOAuth = () => {
    const { t } = useTranslation();

    return (
        <div className="min-h-screen flex items-center justify-center bg-background">
            <div className="flex flex-col items-center gap-8 max-w-md px-6">
                {/* Loading Spinner */}
                <div className="relative w-16 h-16">
                    <div className="absolute top-0 left-0 w-full h-full border-4 border-primary/20 rounded-full" />
                    <div className="absolute top-0 left-0 w-full h-full border-4 border-primary border-t-transparent rounded-full animate-spin" />
                </div>

                {/* Title */}
                <div className="flex flex-col items-center gap-3 text-center">
                    <Typography
                        variant="header"
                        tag="h1"
                        className="text-2xl md:text-3xl font-semibold text-foreground"
                    >
                        {t('oauth.google.title')}
                    </Typography>

                    <Typography
                        variant="text"
                        tag="p"
                        className="text-base text-muted-foreground"
                    >
                        {t('oauth.google.subtitle')}
                    </Typography>
                </div>

                {/* Processing Text */}
                <div className="flex flex-col items-center gap-2">
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-sm text-muted-foreground animate-pulse"
                    >
                        {t('oauth.google.processing')}
                    </Typography>

                    <Typography
                        variant="text"
                        tag="p"
                        className="text-xs text-muted-foreground/60"
                    >
                        {t('oauth.google.redirect')}
                    </Typography>
                </div>
            </div>
        </div>
    );
};

