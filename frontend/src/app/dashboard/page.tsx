"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Device, DeviceStatusChangedEvent } from "@/types/device";
import { NotificationItem } from "@/types/notification";
import { apiFetch, apiDownload, ApiResponse } from "@/lib/api";
import { connectWebSocket } from "@/lib/websocket";
import { ReportSummary } from "@/types/report";

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
    const [report, setReport] = useState<ReportSummary | null>(null);
    const [notifications, setNotifications] = useState<NotificationItem[]>([]);

    const total = devices.length;
    const online = devices.filter((device) => device.status === "ONLINE").length;
    const offline = devices.filter((device) => device.status === "OFFLINE").length;

    const loadDevices = useCallback(async () => {
        try {
            setLoading(true);
            setError("");
            const result = await apiFetch<ApiResponse<Device[]>>("/devices");
            setDevices(result.data);
        } catch (error) {
            setError(error instanceof Error ? error.message : "Failed to load devices");
        } finally {
            setLoading(false);
        }
    }, []);
    const loadReport = useCallback(async () => {
        try {
            const result = await apiFetch<ApiResponse<ReportSummary>>("/reports/summary");
            setReport(result.data);
        } catch (error) {
            console.error(
                "Failed to load report:",
                error
            );
        }
    }, []);
    const handleDeviceStatusChanged = useCallback(
        (event: DeviceStatusChangedEvent) => {
            setDevices((current) =>
                current.map((device) =>
                    device.device_id === event.device_id ? {
                        ...device,
                        status: event.status,
                    } : device
                )
            );
            setNotifications((current) => [
                ...current,
                {
                    id: crypto.randomUUID(),
                    event,
                },
            ]);
            loadReport();
        },
        [loadReport]
    );
    const removeNotification = useCallback((id: string) => {
        setNotifications((current) =>
            current.filter(
                (notification) =>
                    notification.id !== id
            )
        );
    }, []);

    useEffect(() => {
        loadDevices();
        loadReport();
        const disconnect = connectWebSocket(
            handleDeviceStatusChanged
        );
        return () => {
            disconnect();
        };
    }, [loadDevices, loadReport, handleDeviceStatusChanged]);

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

    function handleLogout() {
        localStorage.removeItem("access_token");
        localStorage.removeItem("user");
        router.replace("/login");
    }

    async function handleSuccess() {
        closeModal();
        await Promise.all([
            loadDevices(), 
            loadReport()
        ]);
    }
    async function handleExportCSV() {
        try {
            const blob = await apiDownload(
                "/reports/devices/export"
            );
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement("a");

            link.href = url;
            link.download = `dms-devices-${new Date().toISOString().slice(0, 10)}.csv`;

            document.body.appendChild(link);

            link.click();
            link.remove();

            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error(
                "Export failed:",
                error
            );
        }
    }

    return (
        <main className="min-h-screen bg-gray-400">
            <Notification notifications={notifications} onClose={removeNotification} />
            <header className="border-b bg-white">
                <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
                    <div>
                        <h1 className="text-xl font-bold">
                            Device Management System
                        </h1>
                        <p className="text-md text-gray-700">
                            Dashboard
                        </p>
                    </div>
                    <div className="flex items-center gap-3">
                        <button
                            type="button" onClick={handleExportCSV}
                            className="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-semibold text-gray-700 shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-green-200 hover:bg-green-50 hover:text-green-600 hover:shadow-md active:translate-y-0">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth="2"
                                stroke="currentColor"
                                className="h-4 w-4"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12M12 16.5V3"
                                />
                            </svg>
                            EXPORT
                        </button>
                        <button
                            onClick={handleLogout}
                            className="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-semibold text-gray-700 shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-red-200 hover:bg-red-50 hover:text-red-600 hover:shadow-md active:translate-y-0">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth="2"
                                stroke="currentColor"
                                className="h-4 w-4"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6A2.25 2.25 0 0 0 5.25 5.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3-6 3 3m0 0-3 3m3-3H9"
                                />
                            </svg>
                            LOGOUT
                        </button>
                    </div>
                </div>
            </header>
            <div className="mx-auto max-w-7xl space-y-6 p-6">
                <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                    <div className="flex items-center justify-between rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-2">
                            <span className="h-4 w-4 rounded-full bg-gray-400" />
                            <p className="text-lg font-bold">
                                DEVICE
                            </p>
                        </div>
                        <p className="text-2xl font-bold text-gray-900">
                            {total}
                        </p>
                    </div>
                    <div className="flex items-center justify-between rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-2">
                            <span className="h-4 w-4 rounded-full bg-green-500" />
                            <p className="text-lg font-bold">
                                ONLINE
                            </p>
                        </div>
                        <p className="text-2xl font-bold text-green-600">
                            {online}
                        </p>
                    </div>
                    <div className="flex items-center justify-between rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-2">
                            <span className="h-4 w-4 rounded-full bg-red-500" />
                            <p className="text-lg font-bold">
                                OFFLINE
                            </p>
                        </div>
                        <p className="text-2xl font-bold text-red-600">
                            {offline}
                        </p>
                    </div>
                </div>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                    <div className="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-5">
                            <div className="flex shrink-0 items-center gap-2">
                                <p className="text-xl font-bold">
                                    LAST ONLINE
                                </p>
                            </div>
                            {report?.last_online ? (
                                <div className="ml-auto min-w-0 text-right">
                                    <p className="text-md font-semibold">
                                        {report.last_online.device_id}
                                    </p>
                                    <p className="text-xs text-gray-500">
                                        {report.last_online.device_name}
                                    </p>
                                    <p className="mt-1 text-[11px] text-gray-400">
                                        {new Date(
                                            report.last_online.last_online_at
                                        ).toLocaleDateString("id-ID", {
                                            day: "numeric",
                                            month: "short",
                                            year: "numeric",
                                            timeZone: "Asia/Jakarta",
                                        })}{" • "}
                                        {new Date(
                                            report.last_online.last_online_at
                                        ).toLocaleTimeString("en-US", {
                                            hour: "numeric",
                                            minute: "2-digit",
                                            second: "2-digit",
                                            hour12: true,
                                            timeZone: "Asia/Jakarta",
                                        })}
                                    </p>
                                </div>
                            ) : (
                                <p className="text-xs text-gray-400">
                                    No data
                                </p>
                            )}
                        </div>
                    </div>
                    <div className="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
                        <div className="flex items-center gap-5">
                            <div className="flex shrink-0 items-center gap-2">
                                <p className="text-xl font-bold">
                                    LAST OFFLINE
                                </p>
                            </div>
                            {report?.last_offline ? (
                                <div className="ml-auto min-w-0 text-right">
                                    <p className="text-md font-semibold">
                                        {report.last_offline.device_id}
                                    </p>
                                    <p className="text-xs text-gray-500">
                                        {report.last_offline.device_name}
                                    </p>
                                    <p className="mt-1 text-[11px] text-gray-400">
                                        {new Date(
                                            report.last_offline.last_offline_at
                                        ).toLocaleDateString("id-ID", {
                                            day: "numeric",
                                            month: "short",
                                            year: "numeric",
                                            timeZone: "Asia/Jakarta",
                                        })}{" • "}
                                        {new Date(
                                            report.last_offline.last_offline_at
                                        ).toLocaleTimeString("en-US", {
                                            hour: "numeric",
                                            minute: "2-digit",
                                            second: "2-digit",
                                            hour12: true,
                                            timeZone: "Asia/Jakarta",
                                        })}
                                    </p>
                                </div>
                            ) : (
                                <p className="text-xs text-gray-400">
                                    No data
                                </p>
                            )}
                        </div>
                    </div>
                </div>
                <div className="flex items-center justify-between">
                    <div>
                        <h2 className="text-xl font-semibold">
                            Devices
                        </h2>
                        <p className="text-sm">
                            Manage registered devices
                        </p>
                    </div>
                    <button
                        onClick={openCreateModal}
                        className="inline-flex cursor-pointer items-center gap-2 rounded-xl bg-gray-900 px-8 py-2.5 text-sm font-semibold text-white shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:bg-gray-800 hover:shadow-md active:translate-y-0"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth="2"
                            stroke="currentColor"
                            className="h-4 w-4"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M12 4.5v15m7.5-7.5h-15"
                            />
                        </svg>
                        ADD
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