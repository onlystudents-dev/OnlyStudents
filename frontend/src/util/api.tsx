import { toast } from "react-toastify"
import { fromResponse } from "./language"
import React from "preact/compat";
import {createRoot} from "preact/compat/client";
import type {ComponentChild} from "preact";
import Login from "../auth/login/login.tsx";
import {defaultChild, defaultSchool} from "../navbar/navbar.tsx";
import type {FetchMe} from "../types/props.ts";
import RateLimit from "./ratelimit.tsx";

export function postJSON(url: string, body?: unknown) {
    return fetch(url, {
        method: "POST",
        ...(body !== undefined && {
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
        }),
    })
}

let shown: boolean = false

export async function fetchWithSchool(fetchMe: FetchMe, url: string, init?: RequestInit) {
    const school = localStorage.getItem("school_id")
    const child = localStorage.getItem("child_id")
    const headers = new Headers(init?.headers || {})
    if (school) headers.set("X-School", school)
    if (child) headers.set("X-Child", child)
    let response = await fetch(url, { ...init, headers })
    switch (response.status) {
        case 401: {
            const id = localStorage.getItem("user")
            showElement()
            while (shown) {
                await new Promise(resolve => setTimeout(resolve, 1000))
            }
            const me = await fetchMe()
            if (me) {
                if (`${me.account_id}${me.role}` !== id) {
                    location.reload()
                    return response
                }

                defaultSchool(me, defaultChild(me))
            }
            response = await fetch(url, { ...init, headers })
            break
        }
        case 429: {
            showElement(Number(response.headers.get("Retry-After")))
            while (shown) {
                await new Promise(resolve => setTimeout(resolve, 1000))
            }
            response = await fetch(url, { ...init, headers })
            break
        }
    }
    return response
}

let root: {render(children: ComponentChild): void, unmount(): void} | null = null

export function showElement(ratelimit?: number) {
    shown = true
    if (!root) {
        const container = document.createElement('div')
        document.body.appendChild(container)
        root = createRoot(container)
    }

    if (ratelimit !== undefined) {
        root.render(
            <RateLimit retry={ratelimit} expire={hideElement} z={200} />,
        )
    } else {
        root.render(
            <Login finished={hideElement} z={200}/>
        )
    }
}

export function hideElement() {
    root?.unmount()
    root = null
    shown = false
}

export function getRetryAfter(response: Response): number {
    const raw = response.headers.get("Retry-After")
    if (!raw) return -1
    const seconds = Number(raw)
    return Number.isNaN(seconds) ? -1 : seconds
}

export async function fetchInto<T>(fetchMe: FetchMe, api: string, method?: React.Dispatch<React.SetStateAction<T>>) {
    const response = await fetchWithSchool(fetchMe, api)

    if (!response.ok) {
        toast.error(await fromResponse(response))
        return
    }

    const json = await response.json() as T

    if (method) method(json)
    return json
}

export type ApiError = { error?: string }
