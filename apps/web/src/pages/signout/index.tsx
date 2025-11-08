import { useSignout } from "@/components/useSignout";
import { useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Typography } from "@/components/typography/typography";
import { delay } from "@/lib/utils";

const SignoutPage = () => {
    const { signout, isPending, error } = useSignout();
    const { t } = useTranslation();
    const hasSignedOut = useRef(false);

    useEffect(() => {
        const init = async () => {
            if (hasSignedOut.current) return;
            hasSignedOut.current = true;

            try {
                await delay(2000);
                await signout();
                window.location.href = "/";
            } catch (error) {
                if (error instanceof Error) {
                    toast.error(error.message);
                }
            }
        };
        init();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    if (isPending) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <Typography variant="text" tag="p">
                    {t("common.loading")}
                </Typography>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <Typography variant="text" tag="p" color="muted">
                    {t("common.error")}: {error.message}
                </Typography>
            </div>
        );
    }

    return (
        <div className="flex items-center justify-center min-h-screen">
            <Typography variant="header" tag="h1">
                {t("common.signingOut")}
            </Typography>
        </div>
    );
};

export default SignoutPage;
