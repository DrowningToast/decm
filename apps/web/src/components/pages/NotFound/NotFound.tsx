import { Typography } from '@/components/typography/typography';
import { Button } from '@/components/ui/button';
import { Link } from '@/router';
import { useTranslation } from 'react-i18next';

export const NotFound = () => {
    const { t } = useTranslation();

    return (
        <div className="min-h-screen bg-background flex items-center justify-center px-6">
            <div className="max-w-2xl w-full text-center space-y-8">
                {/* 404 Number */}
                <div className="space-y-4">
                    <Typography
                        variant="header"
                        tag="h1"
                        className="text-[120px] md:text-[180px] font-bold text-primary leading-none"
                    >
                        404
                    </Typography>

                    <Typography
                        variant="header"
                        tag="h2"
                        className="text-3xl md:text-5xl font-bold text-foreground"
                    >
                        {t('notFound.heading')}
                    </Typography>

                    <Typography
                        variant="text"
                        tag="p"
                        className="text-lg text-muted-foreground max-w-md mx-auto"
                    >
                        {t('notFound.description')}
                    </Typography>
                </div>

                {/* Primary CTA */}
                <div className="pt-4">
                    <Link to="/">
                        <Button size="lg" className="min-w-[200px]">
                            <Typography variant="text" tag="span" className="font-medium">
                                {t('notFound.backHome')}
                            </Typography>
                        </Button>
                    </Link>
                </div>

                {/* Suggestions */}
                <div className="pt-8 space-y-4">
                    <Typography
                        variant="text"
                        tag="p"
                        className="text-sm text-muted-foreground font-medium"
                    >
                        {t('notFound.suggestions.title')}
                    </Typography>

                    <div className="flex flex-wrap justify-center gap-3">
                        <Link to="/">
                            <Button variant="secondary-light" size="sm">
                                <Typography variant="text" tag="span">
                                    {t('notFound.suggestions.home')}
                                </Typography>
                            </Button>
                        </Link>

                        <Link to="/signup">
                            <Button variant="secondary-light" size="sm">
                                <Typography variant="text" tag="span">
                                    {t('notFound.suggestions.signup')}
                                </Typography>
                            </Button>
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
};

