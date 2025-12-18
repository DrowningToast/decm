import { EntitySystemStatus, type EntitySystemStatusSchedule } from "@decm/api";
import type { SystemStatus, SystemStatusSchedule } from "./SystemStatusService";

export const mapEntitySystemStatusToSystemStatus = (
    entityStatus: EntitySystemStatus,
): SystemStatus => {
    switch (entityStatus) {
        case EntitySystemStatus.SystemStatusMaintenance:
            return "maintenance";
        case EntitySystemStatus.SystemStatusOperating:
            return "operating";
        default:
            console.error(`Invalid entity system status: ${entityStatus}`);
            throw new Error(`Invalid entity system status: ${entityStatus}`);
    }
};

export const mapEntitySystemStatusScheduleToSystemStatusSchedule = (
    entity: EntitySystemStatusSchedule,
): SystemStatusSchedule => {
    return {
        id: entity.id,
        orderId: entity.order_id,
        startTime: new Date(entity.start_time),
        plannedEndTime: entity.planned_end_time ? new Date(entity.planned_end_time) : undefined,
        status: mapEntitySystemStatusToSystemStatus(entity.status),
        isPlanned: entity.is_planned,
        createdAt: new Date(entity.created_at),
        updatedAt: new Date(entity.updated_at),
        deletedAt: entity.deleted_at ? new Date(entity.deleted_at) : undefined,
    };
};
