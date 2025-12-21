import { coreApiClient } from "@/lib/api/api";
import { mapSystemStatusViewModel } from "./mapper";
import type { SystemStatusSchedule } from "../SystemStatusService/SystemStatusService";

export interface GasPriceInfo {
    currentGasPriceGwei: number;
    baseFeeGwei: number;
    priorityFeeGwei: number;
    softCapPriceGwei: number;
    hardCapPriceGwei: number;
    softSafetyMargin: number;
    hardSafetyMargin: number;
}

export interface SystemStatusViewModel {
    gasPrice: GasPriceInfo;
    latestSchedule?: SystemStatusSchedule | null;
}

export { type SystemStatusSchedule };

class BlockchainService {
    async getSystemStatusViewModel(): Promise<SystemStatusViewModel> {
        const response = await coreApiClient.v1.getSystemStatusViewmodel();
        return mapSystemStatusViewModel(response);
    }
}

export const blockchainService = new BlockchainService();
