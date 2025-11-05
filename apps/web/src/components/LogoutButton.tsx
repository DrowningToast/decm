import { Typography } from "@/components/typography/typography";
import { cn } from "@/lib/utils";
import { useNavigate } from "@/router";
import type { ClassValue } from "clsx";
import { useTranslation } from "react-i18next";

interface LogoutButtonProps {
    className?: ClassValue;
    // signout: logout from the app
    // disconnect: disconnect from the wallet + signout from the app
    // mostly change the label
    type: "signout" | "disconnect";
}

export const LogoutButton: React.FC<LogoutButtonProps> = ({ className, type }) => {
    const navigate = useNavigate();
    const { t } = useTranslation();

    const onLogout = async () => {
        navigate("/signout");
    };

    return (
        <button
            type="button"
            onClick={onLogout}
            className={cn("text-start h-[14.5px] inline-block", className)}
        >
            <Typography
                variant="text"
                tag="span"
                color="background-alt"
                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
            >
                {type === "signout" ? t("auth.signOut") : t("verify.disconnectLink")}
            </Typography>
        </button>
    );
};
