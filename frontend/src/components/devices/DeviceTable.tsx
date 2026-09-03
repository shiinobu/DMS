"use client";

import { Device } from "@/types/device";

interface DeviceTableProps {
    devices: Device[];
    onEdit: (device: Device) => void;
    onDelete: (device: Device) => void;
}

function StatusBadge({
    status,
}: {
    status: Device["status"];
}) {
    const isOnline = status === "ONLINE";

    return (
        <span
            className={`inline-flex rounded-full px-2.5 py-1 text-xs font-medium gap-2 items-center ${isOnline
                    ? "bg-green-100 text-green-700"
                    : "bg-red-100 text-red-700"
                }`}
        >
            <span
                className={`h-2 w-2 rounded-full ${isOnline
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

    return new Date(lastSeen).toLocaleString();
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
                                            onClick={() => onEdit(device)}
                                            className="rounded-lg border px-3 py-1.5 text-xs font-medium hover:bg-gray-50"
                                        >
                                            Edit
                                        </button>

                                        <button
                                            onClick={() => onDelete(device)}
                                            className="rounded-lg bg-red-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-red-700"
                                        >
                                            Delete
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