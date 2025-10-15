import { Typography } from "@/components/typography/typography"
import { useTranslation } from "react-i18next"
import { useNavigate } from "react-router-dom";

export const LogoutButton: React.FC = () => {
    const navigate = useNavigate();
    const { t } = useTranslation()

    const onLogout = () => {
        navigate("/");
    };

    return (
        <button
            type="button"
            onClick={onLogout}
            className="text-start h-[14.5px] inline-block"
        >
            <Typography
                variant="text"
                tag="span"
                color="background-alt"
                className="text-xs italic underline [text-shadow:rgba(255,255,255,0.3)_0px_0px_4px] hover:text-primary transition-colors"
            >
                {t("onboard.logout")}
            </Typography>
        </button>
    )
}