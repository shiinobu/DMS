import { DeviceStatusChangedEvent } from "@/types/device";

export interface NotificationItem {
    id: string;
    event: DeviceStatusChangedEvent;
}