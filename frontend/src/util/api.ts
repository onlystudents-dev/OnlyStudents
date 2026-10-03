export function postJSON(url: string, body?: unknown) {
    return fetch(url, {
        method: "POST",
        headers: body === undefined ? undefined : { "Content-Type": "application/json" },
        body: body === undefined ? undefined : JSON.stringify(body),
    })
}

export function getRetryAfter(response: Response): number {
    const raw = response.headers.get("Retry-After")
    if (!raw) return -1
    const seconds = Number(raw)
    return Number.isNaN(seconds) ? -1 : seconds
}
