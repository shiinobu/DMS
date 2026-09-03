export interface LastOnline {
    device_id: string;
    device_name: string;
    last_online_at: string;
}

export interface LastOffline {
    device_id: string;
    device_name: string;
    last_offline_at: string;
}

export interface ReportSummary {
    total_devices: number;
    online_devices: number;
    offline_devices: number;
    last_online?: LastOnline | null;
    last_offline?: LastOffline | null;
}