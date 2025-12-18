import { coreApiClient, type CoreApiType } from "@/lib/api/api";
import { mapEntitySystemStatusScheduleToSystemStatusSchedule } from "./mapper";

export type SystemStatus = "maintenance" | "operating";

export interface SystemStatusSchedule {
    id: number;
    orderId: number;
    startTime: Date;
    plannedEndTime?: Date;
    status: SystemStatus;
    isPlanned: boolean;
    createdAt: Date;
    updatedAt: Date;
    deletedAt?: Date;
}

interface GetSchedulesBetweenParams {
    startTime: number;
    endTime: number;
}

export class SystemStatusService {
    private _coreApi: CoreApiType;

    constructor(coreApi: CoreApiType) {
        this._coreApi = coreApi;
    }

    public async getLatestSchedules(pageSize: number): Promise<SystemStatusSchedule[]> {
        const response = await this._coreApi.v1.getLatestSchedules({ page_size: pageSize });
        return response.schedules?.map(mapEntitySystemStatusScheduleToSystemStatusSchedule) ?? [];
    }

    public async getSchedulesBetween(
        params: GetSchedulesBetweenParams,
    ): Promise<SystemStatusSchedule[]> {
        const response = await this._coreApi.v1.getSchedulesBetween({
            start_time: params.startTime,
            end_time: params.endTime,
        });
        return response.schedules?.map(mapEntitySystemStatusScheduleToSystemStatusSchedule) ?? [];
    }

    public async getClosestIncomingSchedule(): Promise<SystemStatusSchedule | undefined> {
        const response = await this._coreApi.v1.getClosestIncomingSchedule();
        return response.schedule
            ? mapEntitySystemStatusScheduleToSystemStatusSchedule(response.schedule)
            : undefined;
    }
}

export const defaultSystemStatusService = new SystemStatusService(coreApiClient);
