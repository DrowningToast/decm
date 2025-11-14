import { Link } from "@/router";
import { Typography } from "@/components/typography/typography";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export const AccountAlreadyExistsPage = () => {
    const { t } = useTranslation();

    return (
        <div className="min-h-screen bg-background flex items-center justify-center px-6">
            <div className="max-w-md w-full text-center space-y-8">
                {/* Icon */}
                <div className="flex justify-center">
                    <div className="w-24 h-24 md:w-32 md:h-32 rounded-full bg-primary/10 flex items-center justify-center">
                        <Typography variant="header" tag="div" className="text-5xl md:text-6xl">
                            👋
                        </Typography>
                    </div>
                </div>

                {/* Content */}
                <div className="space-y-4">
                    <Typography
                        variant="header"
                        tag="h1"
                        className="text-3xl md:text-4xl font-bold text-foreground"
                    >
                        {t("accountExists.title")}
                    </Typography>

                    <Typography variant="text" tag="p" className="text-lg text-muted-foreground">
                        {t("accountExists.description")}
                    </Typography>
                </div>

                {/* CTA */}
                <div className="pt-4">
                    <Link to="/signin" className="block">
                        <Button size="lg" variant="primary" className="w-full">
                            <Typography
                                variant="text"
                                tag="span"
                                className="font-medium text-secondary"
                            >
                                {t("accountExists.signInButton")}
                            </Typography>
                        </Button>
                    </Link>
                </div>
            </div>
        </div>
    );
};
