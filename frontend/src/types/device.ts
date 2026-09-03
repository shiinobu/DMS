export type DeviceStatus = "ONLINE" | "OFFLINE";

export interface Device {
    id: number;
    device_id: string;
    device_name: string;
    serial_number: string;
    os_version?: string | null;
    ip_address?: string | null;
    location?: string | null;
    status: DeviceStatus;
    last_seen?: string | null;
    last_online_at?: string | null;
    last_offline_at?: string | null;
    created_at: string;
    updated_at: string;
}

export interface DeviceStatusChangedEvent {
    type: "DEVICE_STATUS_CHANGED";
    device_id: string;
    status: DeviceStatus;
    device: Device;
}