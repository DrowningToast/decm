import { CustomHelmet } from "@/components/providers/helmets/CustomHelmet";
import { useHelmet } from "@/hooks/useHelmet";
import { Typography } from "@/components/typography/typography";

export const LoginPage = () => {
    const helmetData = useHelmet({
        pageType: "login",
        title: "Sign In | DECM Platform",
        description: "Sign in to access your DECM account and manage your digital credentials",
    });

    return (
        <>
            <CustomHelmet
                title={helmetData.title}
                description={helmetData.description}
                themeColor={helmetData.themeColor}
                additionalMeta={[
                    { name: "robots", content: "noindex, nofollow" }, // Don't index login pages
                    { property: "og:type", content: "website" },
                    { property: "og:title", content: helmetData.title },
                    { property: "og:description", content: helmetData.description },
                ]}
            />
            <div>
                <Typography variant="header" tag="h1">
                    Login to DECM
                </Typography>
                <Typography variant="text" tag="p">
                    Sign in to your account
                </Typography>
                {/* Your login form would go here */}
            </div>
        </>
    );
};
