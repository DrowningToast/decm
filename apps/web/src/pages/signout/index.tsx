import { useSignout } from "@/components/useSignout";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Typography } from "@/components/typography/typography";

const SignoutPage = () => {
    const { signout, isPending, error } = useSignout();
    const navigate = useNavigate();
    const { t } = useTranslation();

    useEffect(() => {
        const init = async () => {
            try {
                await signout();
                navigate("/");
            } catch (error) {
                if (error instanceof Error) {
                    toast.error(error.message);
                }
            }
        };
        init();
    }, [navigate, signout]);

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
                {t("common.loading")}
            </Typography>
        </div>
    );
};

export default SignoutPage;
