import { toast } from "react-toastify"
import { fromResponse } from "./language"

export function postJSON(url: string, body?: unknown) {
    return fetch(url, {
        method: "POST",
        ...(body !== undefined && {
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
        }),
    })
}

export function fetchWithHeaders(url: string) {
    const school = localStorage.getItem("school_id");
    const child = localStorage.getItem("child_id");
    const headers: Record<string, string> = {};
    if (school) headers["X-School"] = school;
    if (child) headers["X-Child"] = child;
    return fetch(url, { headers });
}

export function getRetryAfter(response: Response): number {
    const raw = response.headers.get("Retry-After")
    if (!raw) return -1
    const seconds = Number(raw)
    return Number.isNaN(seconds) ? -1 : seconds
}

export async function readJSON<T>(response: Response): Promise<T> {
    return await response.json() as T
}

export async function fetchInto<T>(
    api: string,
    method?: React.Dispatch<React.SetStateAction<T>>
) {
    const response = await fetchWithHeaders(api)

    if (!response.ok) {
        toast.error(await fromResponse(response))
        return
    }

    const json = await readJSON<T>(response)

    if (method) method(json)
    return json
}

export type ApiError = { error?: string }
