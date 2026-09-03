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
        <div className="rounded-xl bg-white p-6 shadow">
            <h2 className="text-xl font-semibold">
                Delete Device
            </h2>

            <p className="mt-3 text-sm text-gray-600">
                Are you sure you want to delete:
            </p>

            <p className="mt-2 font-semibold">
                {device.device_id} - {device.device_name}
            </p>

            {error && (
                <div className="mt-4 rounded-lg bg-red-100 p-3 text-sm text-red-700">
                    {error}
                </div>
            )}

            <div className="mt-6 flex gap-3">
                <button
                    onClick={handleDelete}
                    disabled={loading}
                    className="rounded-lg bg-red-600 px-4 py-2 text-white disabled:opacity-50"
                >
                    {loading ? "Deleting..." : "Delete"}
                </button>

                <button
                    onClick={onCancel}
                    disabled={loading}
                    className="rounded-lg border px-4 py-2"
                >
                    Cancel
                </button>
            </div>
        </div>
    );
}