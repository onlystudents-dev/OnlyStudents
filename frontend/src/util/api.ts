import { toast } from "react-toastify"
import { fromResponse } from "./language"
import React from "preact/compat";

export function postJSON(url: string, body?: unknown) {
    return fetch(url, {
        method: "POST",
        ...(body !== undefined && {
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
        }),
    })
}

export function fetchWithSchool(url: string, init?: RequestInit) {
    const school = localStorage.getItem("school_id");
    const child = localStorage.getItem("child_id");
    const headers = new Headers(init?.headers || {});
    if (school) headers.set("X-School", school);
    if (child) headers.set("X-Child", child);
    return fetch(url, { ...init, headers });
}

export function getRetryAfter(response: Response): number {
    const raw = response.headers.get("Retry-After")
    if (!raw) return -1
    const seconds = Number(raw)
    return Number.isNaN(seconds) ? -1 : seconds
}

export async function fetchInto<T>(api: string, method?: React.Dispatch<React.SetStateAction<T>>) {
    const response = await fetchWithSchool(api)

    if (!response.ok) {
        toast.error(await fromResponse(response))
        return
    }

    const json = await response.json() as T

    if (method) method(json)
    return json
}

export type ApiError = { error?: string }
