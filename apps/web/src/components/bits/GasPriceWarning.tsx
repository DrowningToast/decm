import { Typography } from "@/components/typography/typography";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useSystemStatus } from "@/hooks/useSystemStatus";
import { AlertTriangle, AlertCircle } from "lucide-react";
import { useTranslation } from "react-i18next";

type WarningLevel = "none" | "soft" | "hard";

const getWarningLevel = (currentPrice: number, softCap: number, hardCap: number): WarningLevel => {
    if (hardCap > 0 && currentPrice >= hardCap) {
        return "hard";
    }
    if (softCap > 0 && currentPrice >= softCap) {
        return "soft";
    }
    return "none";
};

export const GasPriceWarning = () => {
    const { t } = useTranslation();
    const { data: systemStatus } = useSystemStatus();

    if (!systemStatus) {
        return null;
    }

    const { gasPrice } = systemStatus;
    const warningLevel = getWarningLevel(
        gasPrice.currentGasPriceGwei,
        gasPrice.softCapPriceGwei,
        gasPrice.hardCapPriceGwei,
    );

    if (warningLevel === "none") {
        return null;
    }

    const isHardWarning = warningLevel === "hard";
    const warningText = isHardWarning
        ? t("blockchain.gasPrice.veryHigh", "Very High Gas Price")
        : t("blockchain.gasPrice.high", "High Gas Price");

    const tooltipContent = isHardWarning
        ? t(
              "blockchain.gasPrice.hardCapTooltip",
              "Network gas prices have exceeded the maximum threshold. Transactions may experience significant delays or fail entirely. We recommend waiting for gas prices to decrease before performing blockchain operations. Please be patient.",
          )
        : t(
              "blockchain.gasPrice.softCapTooltip",
              "Network gas prices are elevated. Transactions may take longer than usual to complete. Consider waiting for gas prices to decrease if your transaction is not time-sensitive.",
          );

    return (
        <Tooltip>
            <TooltipTrigger asChild>
                <div className="flex items-center gap-1.5 cursor-help">
                    {isHardWarning ? (
                        <AlertCircle className="size-3 md:size-4 text-red-500 flex-shrink-0 animate-pulse" />
                    ) : (
                        <AlertTriangle className="size-3 md:size-4 text-yellow-500 flex-shrink-0" />
                    )}
                    <Typography
                        variant="text"
                        tag="span"
                        className={`text-[9px] md:text-xs font-medium ${
                            isHardWarning ? "text-red-500" : "text-yellow-500"
                        }`}
                    >
                        {warningText}
                    </Typography>
                </div>
            </TooltipTrigger>
            <TooltipContent side="bottom" className="max-w-xs">
                <Typography variant="text" tag="p" className="text-xs text-balance">
                    {tooltipContent}
                </Typography>
            </TooltipContent>
        </Tooltip>
    );
};
