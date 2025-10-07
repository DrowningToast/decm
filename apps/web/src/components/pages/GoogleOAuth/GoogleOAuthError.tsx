import { Typography } from '@/components/typography/typography';
import { useTranslation } from 'react-i18next';
import { Link } from '@/router';

interface GoogleOAuthErrorProps {
    errorType?: 'missingCode' | 'missingState' | 'invalidParams' | 'verificationFailed';
}

export const GoogleOAuthError = ({ errorType = 'invalidParams' }: GoogleOAuthErrorProps) => {
    const { t } = useTranslation();

    const getErrorMessage = () => {
        switch (errorType) {
            case 'missingCode':
                return t('oauth.google.error.missingCode');
            case 'missingState':
                return t('oauth.google.error.missingState');
            case 'verificationFailed':
                return t('oauth.google.error.verificationFailed');
            default:
                return t('oauth.google.error.invalidParams');
        }
    };

    const getDescription = () => {
        if (errorType === 'verificationFailed') {
            return t('oauth.google.error.verificationDescription');
        }
        return t('oauth.google.error.description');
    };

    const getReasons = () => {
        if (errorType === 'verificationFailed') {
            return [
                t('oauth.google.error.reasons.serverError'),
                t('oauth.google.error.reasons.networkError'),
                t('oauth.google.error.reasons.sessionExpired'),
            ];
        }
        return [
            t('oauth.google.error.reasons.expired'),
            t('oauth.google.error.reasons.cancelled'),
            t('oauth.google.error.reasons.invalid'),
        ];
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-background px-6">
            <div className="flex flex-col items-center gap-8 max-w-lg w-full">
                {/* Error Icon */}
                <div className="w-20 h-20 rounded-full bg-destructive/10 flex items-center justify-center">
                    <svg
                        className="w-10 h-10 text-destructive"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                        />
                    </svg>
                </div>

                {/* Error Title */}
                <div className="flex flex-col items-center gap-3 text-center">
                    <Typography
                        variant="header"
                        tag="h1"
                        className="text-2xl md:text-3xl font-semibold text-foreground"
                    >
                        {t('oauth.google.error.title')}
                    </Typography>

                    <Typography
                        variant="text"
                        tag="p"
                        className="text-base text-destructive font-medium"
                    >
                        {getErrorMessage()}
                    </Typography>
                </div>

                {/* Error Description */}
                <div className="flex flex-col gap-4 w-full bg-muted/50 rounded-lg p-6">
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-sm text-muted-foreground"
                    >
                        {getDescription()}
                    </Typography>

                    <ul className="list-disc list-inside space-y-2 ml-2">
                        {getReasons().map((reason, index) => (
                            <li key={index}>
                                <Typography
                                    variant="text"
                                    tag="span"
                                    className="text-sm text-muted-foreground"
                                >
                                    {reason}
                                </Typography>
                            </li>
                        ))}
                    </ul>
                </div>

                {/* Action Buttons */}
                <div className="flex flex-col gap-3 w-full">
                    <Link to="/signup" className="w-full">
                        <button className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg px-6 py-3 transition-colors">
                            <Typography variant="text" tag="span" className="font-medium">
                                {t('oauth.google.error.backToSignUp')}
                            </Typography>
                        </button>
                    </Link>

                    <Typography
                        variant="text"
                        tag="p"
                        className="text-xs text-center text-muted-foreground"
                    >
                        {t('oauth.google.error.tryAgain')}
                    </Typography>
                </div>
            </div>
        </div>
    );
};

