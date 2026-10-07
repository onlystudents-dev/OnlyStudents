import {getLanguage, type Language} from "./language.ts";
import type {Me} from "../types/api.ts";

export function getTimeFormat(language?: Language): boolean {
    if (!language) return false

    const hourCycle = new Intl.DateTimeFormat(language.key, {
        hour: "numeric",
    }).resolvedOptions().hourCycle

    return hourCycle === "h23" || hourCycle === "h24"
}

export function formatSecondsToHourAndMinute(seconds: number, me: Me | null = null) {
    const language = getLanguage(me)?.key || "en-US"

    return new Intl.DateTimeFormat(language, {
        hour: "numeric",
        minute: "2-digit",
        timeZone: "UTC",
        hourCycle: me?.preferences.time_format || undefined
    }).format(new Date(seconds * 1000))
}

export function formatUnixDate(unix: number, locale: string): string {
    return new Intl.DateTimeFormat(locale, {
        month: "short",
        day: "numeric",
    }).format(new Date(unix * 1000))
}

export function getWeekRange(date: Date) {
    const day = date.getDay()
    const diffToMonday = day === 0 ? -6 : 1 - day

    const monday = new Date(date)
    monday.setDate(date.getDate() + diffToMonday)
    monday.setHours(0, 0, 0, 0)

    const sunday = new Date(monday)
    sunday.setDate(monday.getDate() + 6)

    const toISODate = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`

    return { start: toISODate(monday), end: toISODate(sunday), monday, mondayTime: Math.floor(monday.getTime() / 1000) }
}
