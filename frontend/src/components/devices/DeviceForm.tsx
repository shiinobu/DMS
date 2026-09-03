"use client";

import { FormEvent, useEffect, useState } from "react";
import { Device } from "@/types/device";
import { apiFetch } from "@/lib/api";

interface DeviceFormProps {
    device?: Device | null;
    onSuccess: () => void;
    onCancel: () => void;
}

interface DeviceFormData {
    device_id: string;
    device_name: string;
    serial_number: string;
    os_version: string;
    ip_address: string;
    location: string;
}

const emptyForm: DeviceFormData = {
    device_id: "",
    device_name: "",
    serial_number: "",
    os_version: "",
    ip_address: "",
    location: "",
};

export default function DeviceForm({
    device,
    onSuccess,
    onCancel,
}: DeviceFormProps) {
    const [form, setForm] = useState<DeviceFormData>(emptyForm);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const isEdit = Boolean(device);

    useEffect(() => {
        if (device) {
            setForm({
                device_id: device.device_id,
                device_name: device.device_name,
                serial_number: device.serial_number,
                os_version: device.os_version || "",
                ip_address: device.ip_address || "",
                location: device.location || "",
            });
        } else {
            setForm(emptyForm);
        }

        setError("");
    }, [device]);

    function handleChange(
        field: keyof DeviceFormData,
        value: string
    ) {
        setForm((current) => ({
            ...current,
            [field]: value,
        }));
    }

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>
    ) {
        event.preventDefault();

        setLoading(true);
        setError("");

        try {
            const payload = {
                device_id: form.device_id,
                device_name: form.device_name,
                serial_number: form.serial_number,
                os_version: form.os_version || null,
                ip_address: form.ip_address || null,
                location: form.location || null,
            };

            if (isEdit) {
                await apiFetch(
                    `/devices/${encodeURIComponent(form.device_id)}`,
                    {
                        method: "PUT",
                        body: JSON.stringify(payload),
                    }
                );
            } else {
                await apiFetch("/devices", {
                    method: "POST",
                    body: JSON.stringify(payload),
                });
            }

            onSuccess();
        } catch (error) {
            setError(
                error instanceof Error
                    ? error.message
                    : "Failed to save device"
            );
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="rounded-xl bg-white p-6 shadow">
            <div className="mb-6">
                <h2 className="text-xl font-semibold">
                    {isEdit ? "Edit Device" : "Add Device"}
                </h2>

                <p className="mt-1 text-sm text-gray-500">
                    {isEdit
                        ? "Update device information."
                        : "Register a new device."}
                </p>
            </div>

            {error && (
                <div className="mb-4 rounded-lg bg-red-100 p-3 text-sm text-red-700">
                    {error}
                </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label className="mb-1 block text-sm font-medium">
                        Device ID
                    </label>

                    <input
                        type="text"
                        value={form.device_id}
                        onChange={(event) =>
                            handleChange("device_id", event.target.value)
                        }
                        disabled={isEdit}
                        required
                        className="w-full rounded-lg border px-3 py-2 disabled:bg-gray-100"
                        placeholder="DMS-006"
                    />
                </div>

                <div>
                    <label className="mb-1 block text-sm font-medium">
                        Device Name
                    </label>

                    <input
                        type="text"
                        value={form.device_name}
                        onChange={(event) =>
                            handleChange("device_name", event.target.value)
                        }
                        required
                        className="w-full rounded-lg border px-3 py-2"
                        placeholder="Office PC 006"
                    />
                </div>

                <div>
                    <label className="mb-1 block text-sm font-medium">
                        Serial Number
                    </label>

                    <input
                        type="text"
                        value={form.serial_number}
                        onChange={(event) =>
                            handleChange("serial_number", event.target.value)
                        }
                        required
                        className="w-full rounded-lg border px-3 py-2"
                        placeholder="SN-DMS-006"
                    />
                </div>

                <div>
                    <label className="mb-1 block text-sm font-medium">
                        OS Version
                    </label>

                    <input
                        type="text"
                        value={form.os_version}
                        onChange={(event) =>
                            handleChange("os_version", event.target.value)
                        }
                        className="w-full rounded-lg border px-3 py-2"
                        placeholder="Windows 11 Pro"
                    />
                </div>

                <div>
                    <label className="mb-1 block text-sm font-medium">
                        IP Address
                    </label>

                    <input
                        type="text"
                        value={form.ip_address}
                        onChange={(event) =>
                            handleChange("ip_address", event.target.value)
                        }
                        className="w-full rounded-lg border px-3 py-2"
                        placeholder="192.168.1.106"
                    />
                </div>

                <div>
                    <label className="mb-1 block text-sm font-medium">
                        Location
                    </label>

                    <input
                        type="text"
                        value={form.location}
                        onChange={(event) =>
                            handleChange("location", event.target.value)
                        }
                        className="w-full rounded-lg border px-3 py-2"
                        placeholder="Office - Floor 2"
                    />
                </div>

                <div className="flex gap-3 pt-2">
                    <button
                        type="submit"
                        disabled={loading}
                        className="rounded-lg bg-black px-4 py-2 text-white disabled:opacity-50"
                    >
                        {loading
                            ? "Saving..."
                            : isEdit
                                ? "Update Device"
                                : "Create Device"}
                    </button>

                    <button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="rounded-lg border px-4 py-2"
                    >
                        Cancel
                    </button>
                </div>
            </form>
        </div>
    );
}