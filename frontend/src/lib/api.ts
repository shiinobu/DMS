const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export interface ApiResponse<T> {
    success: boolean;
    message?: string;
    data: T;
}

export interface LoginResponse {
    token: string;
    user: {
        id: number;
        username: string;
        role: string;
    };
}

export async function login(
    username: string,
    password: string
): Promise<LoginResponse> {
    const response = await fetch(
        `${API_URL}/auth/login`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username,
                password,
            }),
        }
    );

    const result = await response.json();

    if (!response.ok) {
        throw new Error(
            result.message || "Login failed"
        );
    }

    return result.data;
}

function getAccessToken(): string | null {
    if (typeof window === "undefined") {
        return null;
    }
    return localStorage.getItem(
        "access_token"
    );
}

function handleUnauthorized() {
    if (typeof window !== "undefined") {
        localStorage.removeItem(
            "access_token"
        );
        localStorage.removeItem("user");
        window.location.href = "/login";
    }
}

export async function apiFetch<T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<T> {
    const token = getAccessToken();
    const headers = new Headers(options.headers);

    headers.set("Content-Type", "application/json");

    if (token) {
        headers.set(
            "Authorization",
            `Bearer ${token}`
        );
    }

    const response = await fetch(
        `${API_URL}${endpoint}`,
        {
            ...options,
            headers,
        }
    );

    if (response.status === 401) {
        handleUnauthorized();
        throw new Error("Unauthorized");
    }

    const result = await response.json();

    if (!response.ok) {
        throw new Error(
            result.message || "Request failed"
        );
    }

    return result;
}

export async function apiDownload(
    endpoint: string
): Promise<Blob> {
    const token = getAccessToken();
    const headers = new Headers();

    if (token) {
        headers.set(
            "Authorization",
            `Bearer ${token}`
        );
    }

    const response = await fetch(
        `${API_URL}${endpoint}`,
        {
            method: "GET",
            headers,
        }
    );

    if (response.status === 401) {
        handleUnauthorized();
        throw new Error("Unauthorized");
    }

    if (!response.ok) {
        throw new Error(
            "Failed to download file"
        );
    }

    return response.blob();
}