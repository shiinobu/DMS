"use client";

import { use, useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Device, DeviceStatusChangedEvent } from "@/types/device";
import { apiFetch } from "@/lib/api";
import { connectWebSocket } from "@/lib/websocket";

import DeviceForm from "@/components/devices/DeviceForm";
import DeviceTable from "@/components/devices/DeviceTable";
import DeleteDeviceDialog from "@/components/devices/DeleteDeviceDialog";
import Notification from "@/components/notifications/Notification";

type ModalType = "create" | "edit" | "delete" | null;

export default function DashboardPage() {
    const router = useRouter();
    const [devices, setDevices] = useState<Device[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [modal, setModal] = useState<ModalType>(null);
    const [selectedDevice, setSelectedDevice] = useState<Device | null>(null);
    const [notification, setNotification] = useState<DeviceStatusChangedEvent | null>(null);

    const loadDevices = useCallback(async () => {
        try {
            setLoading(true);
            setError("");
            const result = await apiFetch<{
                success: boolean;
                data: Device[];
            }>("/devices");
            setDevices(result.data || []);
        } catch (error) {
            setError(error instanceof Error ? error.message : "Failed to load devices");
        } finally {
            setLoading(false);
        }
    },
        []
    );

    const handleDeviceStatusChanged = (
        event: DeviceStatusChangedEvent
    ) => {
        setDevices((current) =>
            current.map((device) => device.device_id === event.device_id ? event.device : device)
        );
        setNotification(event);
    };

    useEffect(() => {
        const accessToken = localStorage.getItem("access_token");
        if (!accessToken) {
            router.push("/login");
            return;
        }
        loadDevices();
    }, [router, loadDevices]);

    useEffect(() => {
        const disconnect = connectWebSocket(
            handleDeviceStatusChanged
        );

        return () => {
            disconnect();
        };
    }, []);

    function openCreateModal() {
        setSelectedDevice(null);
        setModal("create");
    }

    function openEditModal(device: Device) {
        setSelectedDevice(device);
        setModal("edit");
    }

    function openDeleteModal(device: Device) {
        setSelectedDevice(device);
        setModal("delete");
    }

    function closeModal() {
        setModal(null);
        setSelectedDevice(null);
    }

    async function handleSuccess() {
        closeModal();
        await loadDevices();
    }

    function handleLogout() {
        localStorage.removeItem("access_token");
        localStorage.removeItem("user");

        router.replace("/login");
    }

    const total = devices.length;

    const online = devices.filter(
        (device) => device.status === "ONLINE"
    ).length;

    const offline = devices.filter(
        (device) => device.status === "OFFLINE"
    ).length;

    return (
        <main className="min-h-screen bg-gray-400">
            <Notification event={notification} onClose={() => setNotification(null)}/>
            <header className="border-b bg-white">
                <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
                    <div>
                        <h1 className="text-xl font-bold">
                            Device Management System
                        </h1>
                        <p className="text-sm">
                            Dashboard
                        </p>
                    </div>
                    <button onClick={handleLogout} className="rounded-lg border px-4 py-2 text-sm hover:bg-gray-50 cursor-pointer">
                        Logout
                    </button>
                </div>
            </header>
            <div className="mx-auto max-w-7xl space-y-6 p-6">
                <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                    {/* Total Device */}
                    <div className="flex items-center justify-between rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-2">
                            <span className="h-4 w-4 rounded-full bg-gray-400" />
                            <p className="text-lg font-bold uppercase">
                                Device
                            </p>
                        </div>
                        <p className="text-2xl font-bold text-gray-900">
                            {total}
                        </p>
                    </div>
                    {/* Online */}
                    <div className="flex items-center justify-between rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-2">
                            <span className="h-4 w-4 rounded-full bg-green-500" />
                            <p className="text-lg font-bold uppercase">
                                Online
                            </p>
                        </div>
                        <p className="text-2xl font-bold text-green-600">
                            {online}
                        </p>
                    </div>
                    {/* Offline */}
                    <div className="flex items-center justify-between rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-2">
                            <span className="h-4 w-4 rounded-full bg-red-500" />
                            <p className="text-lg font-bold uppercase">
                                Offline
                            </p>
                        </div>
                        <p className="text-2xl font-bold text-red-600">
                            {offline}
                        </p>
                    </div>
                </div>
                <div className="flex items-center justify-between">
                    <div>
                        <h2 className="text-xl font-semibold">
                            Devices
                        </h2>
                        <p className="text-sm">
                            Manage registered devices.
                        </p>
                    </div>
                    <button
                        onClick={openCreateModal}
                        className="rounded-lg bg-black px-4 py-2 text-sm font-medium text-white cursor-pointer"
                    >
                        + Add Device
                    </button>
                </div>

                {error && (
                    <div className="rounded-lg bg-red-100 p-4 text-sm text-red-700">
                        {error}
                    </div>
                )}

                {loading ? (
                    <div className="rounded-xl bg-white p-8 text-center text-gray-500 shadow">
                        Loading devices...
                    </div>
                ) : (
                    <DeviceTable
                        devices={devices}
                        onEdit={openEditModal}
                        onDelete={openDeleteModal}
                    />
                )}
            </div>
            {modal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-6">
                    <div className="max-h-[90vh] w-full max-w-lg overflow-y-auto">
                        {modal === "create" && (
                            <DeviceForm
                                onSuccess={handleSuccess}
                                onCancel={closeModal}
                            />
                        )}

                        {modal === "edit" && selectedDevice && (
                            <DeviceForm
                                device={selectedDevice}
                                onSuccess={handleSuccess}
                                onCancel={closeModal}
                            />
                        )}

                        {modal === "delete" && selectedDevice && (
                            <DeleteDeviceDialog
                                device={selectedDevice}
                                onSuccess={handleSuccess}
                                onCancel={closeModal}
                            />
                        )}
                    </div>
                </div>
            )}
        </main>
    );
}