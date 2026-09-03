"use client";

import { DeviceStatusChangedEvent } from "@/types/device";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1/ws";
const RECONNECT_DELAY = 3000;

export function connectWebSocket(
    onMessage: (
        event: DeviceStatusChangedEvent
    ) => void,
    onConnectionChange?: (
        connected: boolean
    ) => void
) {
    let socket: WebSocket | null = null;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
    let stopped = false;

    function connect() {
        if (stopped) {
            return;
        }

        socket = new WebSocket(WS_URL);

        socket.onopen = () => {
            onConnectionChange?.(true);
        };

        socket.onmessage = (message) => {
            try {
                const event = JSON.parse(
                    message.data
                ) as DeviceStatusChangedEvent;
                if (
                    event.type ===
                    "DEVICE_STATUS_CHANGED"
                ) {
                    onMessage(event);
                }
            } catch (error) {
                console.error(
                    "Invalid WebSocket message:",
                    error
                );
            }
        };

        socket.onclose = () => {
            onConnectionChange?.(false);

            if (!stopped) {
                reconnectTimer =
                    setTimeout(
                        connect,
                        RECONNECT_DELAY
                    );
            }
        };

        socket.onerror = () => {};
    }

    connect();

    return () => {
        stopped = true;
        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null;
        }
        socket?.close();
        socket = null;
    };
}