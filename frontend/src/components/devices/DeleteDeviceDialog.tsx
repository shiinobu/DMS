"use client";

import { useState } from "react";
import { Device } from "@/types/device";
import { apiFetch } from "@/lib/api";

interface DeleteDeviceDialogProps {
    device: Device;
    onSuccess: () => void;
    onCancel: () => void;
}

export default function DeleteDeviceDialog({
    device,
    onSuccess,
    onCancel,
}: DeleteDeviceDialogProps) {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    async function handleDelete() {
        setLoading(true);
        setError("");

        try {
            await apiFetch(
                `/devices/${encodeURIComponent(device.device_id)}`,
                {
                    method: "DELETE",
                }
            );
            onSuccess();
        } catch (error) {
            setError(
                error instanceof Error
                    ? error.message
                    : "Failed to delete device"
            );
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
            <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
                <div className="flex items-start gap-4">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-red-50">
                        <svg
                            className="h-5 w-5 text-red-500"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M12 9v3.75m0 3.75h.008M10.29 3.86l-7.5 13A1.875 1.875 0 004.42 19.5h15.16a1.875 1.875 0 001.63-2.81l-7.5-13a1.875 1.875 0 00-3.42 0z"
                            />
                        </svg>
                    </div>
                    <div>
                        <h2 className="text-lg font-semibold text-gray-900">
                            Delete Device
                        </h2>
                        <p className="mt-1 text-sm text-gray-500">
                            Are you sure you want to delete this device?
                        </p>
                    </div>
                </div>
                <div className="mt-5 rounded-lg border border-gray-200 bg-gray-50 px-4 py-3">
                    <p className="text-sm font-semibold text-gray-900">
                        {device.device_id}
                    </p>
                    <p className="mt-0.5 text-sm text-gray-500">
                        {device.device_name}
                    </p>
                </div>
                <p className="mt-3 text-xs text-gray-400 text-right">
                    *This action cannot be undone!
                </p>
                <div className="mt-6 flex items-center justify-between">
                    <button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="rounded-lg border border-gray-300 px-5 py-2.5 text-sm font-semibold text-gray-700 transition hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-300 focus:ring-offset-2 cursor-pointer"
                    >
                        CANCEL
                    </button>
                    <button
                        type="button"
                        onClick={handleDelete}
                        disabled={loading}
                        className="rounded-lg bg-red-500 px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 cursor-pointer"
                    >
                        {loading ? "DELETING..." : "DELETE"}
                    </button>
                </div>
            </div>
        </div>
    );
}