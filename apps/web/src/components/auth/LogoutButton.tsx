import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useDisconnect } from "wagmi";

interface LogoutButtonProps {
    type: "disconnect" | "signout";
}

export const LogoutButton: React.FC<LogoutButtonProps> = ({ type }) => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { disconnect } = useDisconnect();

    const handleClick = () => {
        if (type === "disconnect") {
            disconnect();
            navigate("/");
        } else {
            navigate("/signout");
        }
    };

    return (
        <Button
            type="button"
            onClick={handleClick}
            variant="secondary-light"
            size="xl"
            className="w-full"
        >
            {type === "disconnect" ? t("auth.disconnect") : t("auth.signOut")}
        </Button>
    );
};
