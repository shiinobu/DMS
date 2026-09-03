"use client";

import { Device } from "@/types/device";

interface DeviceTableProps {
    devices: Device[];
    onEdit: (device: Device) => void;
    onDelete: (device: Device) => void;
}

function StatusBadge({
    status,
}:{
    status: Device["status"];
}) {
    const isOnline = status === "ONLINE";
    return (
        <span className={`inline-flex items-center gap-2 rounded-full px-2.5 py-1 text-xs font-medium ${isOnline
                ? "bg-green-100 text-green-700"
                : "bg-red-100 text-red-700"
            }`}
        >
            <span className={`h-2 w-2 rounded-full ${isOnline
                    ? "bg-green-500"
                    : "bg-red-500"
                }`}
            />
            {status}
        </span>
    );
}

function formatLastSeen(lastSeen?: string | null) {
    if (!lastSeen) {
        return "-";
    }
    return new Date(lastSeen).toLocaleString("id-ID", {
        timeZone: "Asia/Jakarta",
    });
}

export default function DeviceTable({
    devices,
    onEdit,
    onDelete,
}: DeviceTableProps) {
    if (devices.length === 0) {
        return (
            <div className="rounded-xl bg-white p-8 text-center text-gray-500 shadow">
                No devices found.
            </div>
        );
    }

    return (
        <div className="overflow-hidden rounded-xl bg-white shadow">
            <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                    <thead className="border-b bg-gray-50">
                        <tr>
                            <th className="px-4 py-3 font-semibold">
                                Device ID
                            </th>
                            <th className="px-4 py-3 font-semibold">
                                Name
                            </th>
                            <th className="px-4 py-3 font-semibold">
                                Serial Number
                            </th>
                            <th className="px-4 py-3 font-semibold">
                                IP Address
                            </th>
                            <th className="px-4 py-3 font-semibold">
                                Location
                            </th>
                            <th className="px-4 py-3 font-semibold">
                                Status
                            </th>
                            <th className="px-4 py-3 font-semibold">
                                Last Seen
                            </th>
                            <th className="px-4 py-3 text-center font-semibold">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="divide-y">
                        {devices.map((device) => (
                            <tr key={device.device_id}>
                                <td className="px-4 py-3 font-medium">
                                    {device.device_id}
                                </td>
                                <td className="px-4 py-3">
                                    {device.device_name}
                                </td>
                                <td className="px-4 py-3">
                                    {device.serial_number}
                                </td>
                                <td className="px-4 py-3">
                                    {device.ip_address || "-"}
                                </td>
                                <td className="px-4 py-3">
                                    {device.location || "-"}
                                </td>
                                <td className="px-4 py-3">
                                    <StatusBadge status={device.status} />
                                </td>
                                <td className="px-4 py-3">
                                    {formatLastSeen(device.last_seen)}
                                </td>
                                <td className="px-4 py-3">
                                    <div className="flex justify-end gap-2">
                                        <button
                                            type="button"
                                            onClick={() => onEdit(device)}
                                            aria-label={`Edit ${device.device_id}`}
                                            className="inline-flex cursor-pointer items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-black shadow-sm transition-all duration-200 hover:border-gray-500 hover:bg-gray-50 hover:text-gray-900 hover:shadow active:scale-95"
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                strokeWidth="1.8"
                                                stroke="currentColor"
                                                className="h-5 w-5"
                                            >
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Z"
                                                />
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    d="M19.5 7.5 16.5 4.5"
                                                />
                                            </svg>
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => onDelete(device)}
                                            aria-label={`Delete ${device.device_id}`}
                                            className="inline-flex cursor-pointer items-center gap-1.5 rounded-lg border border-red-200 bg-red-100 px-3 py-1.5 text-xs font-medium text-red-600 transition-all duration-200 hover:border-red-500 hover:bg-red-100 hover:text-red-800 hover:shadow-sm active:scale-95"
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                strokeWidth="1.8"
                                                stroke="currentColor"
                                                className="h-5 w-5"
                                            >
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673A2.25 2.25 0 0 1 15.916 21H8.084a2.25 2.25 0 0 1-2.244-1.327L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12.37.562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0C8.91 2.73 8 3.714 8 4.894v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
                                                />
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}