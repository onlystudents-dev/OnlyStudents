export function postJSON(url: string, body?: unknown) {
    return fetch(url, {
        method: "POST",
        ...(body !== undefined && {
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
        }),
    })
}

/** Seconds from a 429 response's Retry-After header, or -1 when missing/invalid. */
export function getRetryAfter(response: Response): number {
    const raw = response.headers.get("Retry-After")
    if (!raw) return -1
    const seconds = Number(raw)
    return Number.isNaN(seconds) ? -1 : seconds
}

/**
 * The one place where untyped `response.json()` becomes a caller-declared type.
 * The server is trusted to match `T`; validate with a guard instead if that's ever not good enough.
 */
export async function readJSON<T>(response: Response): Promise<T> {
    return await response.json() as T
}

export type ApiError = { error?: string }
