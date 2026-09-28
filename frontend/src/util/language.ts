import {match} from "@formatjs/intl-localematcher";
import type {Me} from "../types/api.ts";

export type Language = {
    key: string,
    name: string,
    emoji: string,
}

export const languages: Language[] = [
    {"key": "en-US", "name": "English (US)", "emoji": "\uD83C\uDDFA\uD83C\uDDF8"},
    {"key": "hu-HU", "name": "Magyar", "emoji": "\uD83C\uDDED\uD83C\uDDFA"}
]

let language: Record<string, Record<string, string>>

export async function fetchLanguage(me: Me | null) {
    const response = await fetch(`/assets/languages/${getLanguage(me)?.key}.json`);
    language = await response.json();
}

export function getLanguage(me: Me | null, auto?: boolean) {
    const stored = !auto ? me?.preferences?.lang : null;

    return languages.find(o => o.key === match(
        stored ? [stored] : navigator.languages,
        languages.map((o) => o.key),
        "en-US"
    ))
}

export function getKey(key: string, ...args: string[]) {
    let val = language.messages[key]
    if (!val) return ""
    let i = 0
    for (const arg of args) {
        i++
        val = val.replace(`%${i}`, arg)
    }
    return val
}

export async function fromResponse(response: Response, ...args: string[]) {
    if (response.status === 429) {
        let seconds = response.headers.get("Retry-After")
        if (!seconds) seconds = "-1"
        args.shift()
        args[0] = seconds
    }

    const json = await response.json()

    let val = language.errors[json?.error]
    if (!val) return JSON.stringify(json)
    let i = 0
    for (const arg of args) {
        i++
        val = val.replace(`%${i}`, arg)
    }

    return val
}
