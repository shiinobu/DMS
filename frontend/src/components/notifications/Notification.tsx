"use client";

import { useEffect } from "react";
import { NotificationItem } from "@/types/notification";

interface NotificationProps {
    notifications: NotificationItem[];
    onClose: (id: string) => void;
}

interface NotificationCardProps {
    item: NotificationItem;
    onClose: (id: string) => void;
}

function NotificationCard({
    item,
    onClose,
}: NotificationCardProps) {
    const { id, event } = item;

    useEffect(() => {
        const timer = setTimeout(() => {
            onClose(id);
        }, 20000);

        return () => {
            clearTimeout(timer);
        };
    }, [id, onClose]);

    const isOnline = event.status === "ONLINE";

    return (
        <div className="notification-slide-in">
            <div className="flex w-[320px] items-center gap-3 rounded-lg border border-gray-200 bg-white px-4 py-3 shadow-lg">
                <span
                    className={`h-3 w-3 shrink-0 rounded-full ${isOnline
                            ? "bg-green-500"
                            : "bg-red-500"
                        }`}
                />
                <div className="flex-1 text-sm text-gray-700">
                    Device{" "}
                    <span className="font-semibold text-gray-900">
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
                <button
                    type="button"
                    onClick={() => onClose(id)}
                    className="shrink-0 text-gray-400 transition hover:text-gray-700 cursor-pointer"
                    aria-label="Close notification"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        className="h-5 w-5"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M6 6l12 12"
                        />

                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M18 6L6 18"
                        />
                    </svg>
                </button>
            </div>
        </div>
    );
}

export default function Notification({
    notifications,
    onClose,
}: NotificationProps) {
    if (notifications.length === 0) {
        return null;
    }

    return (
        <div className="fixed right-6 top-6 z-50 flex flex-col gap-3">
            {notifications.map((item) => (
                <NotificationCard
                    key={item.id}
                    item={item}
                    onClose={onClose}
                />
            ))}
        </div>
    );
}