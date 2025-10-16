import { Typography } from '@/components/typography/typography';
import { Button } from '@/components/ui/button';
import { Link } from '@/router';
import { useTranslation } from 'react-i18next';

interface ErrorProps {
    title?: string
    description?: string
}

export const ErrorPage: React.FC<ErrorProps> = ({ title, description }) => {
    const { t } = useTranslation();

    const handleRefresh = () => {
        window.location.reload();
    };

    return (
        <div className="min-h-screen bg-background flex items-center justify-center px-6">
            <div className="max-w-2xl w-full text-center space-y-8">
                {/* Error Icon/Symbol */}
                <div className="space-y-4">
                    <div className="flex justify-center">
                        <div className="w-32 h-32 md:w-40 md:h-40 rounded-full bg-destructive/10 flex items-center justify-center">
                            <Typography
                                variant="header"
                                tag="div"
                                className="text-6xl md:text-8xl text-destructive"
                            >
                                ⚠️
                            </Typography>
                        </div>
                    </div>

                    <Typography
                        variant="header"
                        tag="h1"
                        className="text-3xl md:text-5xl font-bold text-foreground"
                    >
                        {title || t('error.heading')}
                    </Typography>

                    <Typography
                        variant="text"
                        tag="p"
                        className="text-lg text-muted-foreground max-w-md mx-auto"
                    >
                        {description || t('error.description')}
                    </Typography>
                </div>

                {/* Primary CTAs */}
                <div className="pt-4 flex flex-col sm:flex-row gap-4 justify-center items-center">
                    <Link to="/">
                        <Button size="lg" variant="primary" className="min-w-[180px]">
                            <Typography variant="text" tag="span" className="font-medium text-secondary">
                                {t('error.backHome')}
                            </Typography>
                        </Button>
                    </Link>

                    <Button
                        size="lg"
                        variant="secondary-light"
                        className="min-w-[180px]"
                        onClick={handleRefresh}
                    >
                        <Typography variant="text" tag="span" className="font-medium text-secondary-foreground">
                            {t('error.tryAgain')}
                        </Typography>
                    </Button>
                </div>

                {/* Suggestions */}
                <div className="pt-8 space-y-4">
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-sm text-muted-foreground font-medium"
                    >
                        {t('error.suggestions.title')}
                    </Typography>

                    <div className="flex flex-wrap justify-center gap-3">
                        <Link to="/">
                            <Button variant="secondary-light" size="sm">
                                <Typography variant="text" tag="span" className="text-secondary-foreground">
                                    {t('error.suggestions.home')}
                                </Typography>
                            </Button>
                        </Link>

                        <Button
                            variant="secondary-light"
                            size="sm"
                            onClick={handleRefresh}
                        >
                            <Typography variant="text" tag="span" className="text-secondary-foreground">
                                {t('error.suggestions.refresh')}
                            </Typography>
                        </Button>

                        <Link to="/signup">
                            <Button variant="secondary-light" size="sm">
                                <Typography variant="text" tag="span" className="text-secondary-foreground">
                                    {t('error.suggestions.signup')}
                                </Typography>
                            </Button>
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
};

