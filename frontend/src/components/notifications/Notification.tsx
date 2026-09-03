"use client";

import { useEffect } from "react";
import { DeviceStatusChangedEvent } from "@/types/device";

interface NotificationProps {
    event: DeviceStatusChangedEvent | null;
    onClose: () => void;
}

export default function Notification({
    event,
    onClose,
}: NotificationProps) {
    useEffect(() => {
        if (!event) {
            return;
        }

        const timer = setTimeout(() => {
            onClose();
        }, 10000 );

        return () => {
            clearTimeout(timer);
        };
    }, [event, onClose]);

    if (!event) {
        return null;
    }

    const isOnline = event.status === "ONLINE";

    return (
        <div className="fixed right-6 top-6 z-50">
            <div className="flex w-[320px] items-center gap-3 rounded-lg border border-gray-200 bg-white px-4 py-3 shadow-lg">
                {/* Status Dot */}
                <span
                    className={`h-3 w-3 shrink-0 rounded-full ${isOnline ? "bg-green-500" : "bg-red-500"
                        }`}
                />

                {/* Message */}
                <div className="flex-1 text-sm text-gray-700">
                    Device{" "}
                    <span className="font-semibold">
                        {event.device_id}
                    </span>{" "}
                    is{" "}
                    <span
                        className={`font-semibold ${isOnline
                                ? "text-green-600"
                                : "text-red-600"
                            }`}
                    >
                        {event.status}
                    </span>
                </div>

                {/* Close */}
                <button
                    type="button"
                    onClick={onClose}
                    className="text-lg leading-none text-gray-400 transition hover:text-gray-700"
                    aria-label="Close notification"
                >
                    ×
                </button>
            </div>
        </div>
    );
}