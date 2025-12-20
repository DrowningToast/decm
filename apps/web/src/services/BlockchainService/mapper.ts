import type { SystemStatusSystemStatusViewModel } from "@decm/api";
import type { SystemStatusViewModel, GasPriceInfo } from "./BlockchainService";
import { mapEntitySystemStatusScheduleToSystemStatusSchedule } from "../SystemStatusService/mapper";

export const mapGasPriceResponse = (
    data: SystemStatusSystemStatusViewModel["gas_price"],
): GasPriceInfo => {
    return {
        currentGasPriceGwei: data.current_gas_price_gwei,
        baseFeeGwei: data.base_fee_gwei,
        priorityFeeGwei: data.priority_fee_gwei,
        softCapPriceGwei: data.soft_cap_price_gwei,
        hardCapPriceGwei: data.hard_cap_price_gwei,
        softSafetyMargin: data.soft_safety_margin,
        hardSafetyMargin: data.hard_safety_margin,
    };
};

export const mapSystemStatusViewModel = (
    response: SystemStatusSystemStatusViewModel,
): SystemStatusViewModel => {
    return {
        gasPrice: mapGasPriceResponse(response.gas_price),
        latestSchedule: response.latest_schedule
            ? mapEntitySystemStatusScheduleToSystemStatusSchedule(response.latest_schedule)
            : null,
    };
};
