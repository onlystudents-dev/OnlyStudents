import {match} from "@formatjs/intl-localematcher";
import type {Me} from "../types/api.ts";
import {getRetryAfter} from "./api.ts";

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

function fill(template: string, args: string[]) {
    return args.reduce((val, arg, i) => val.replace(`%${i + 1}`, arg), template)
}

export function getKey(key: string, ...args: string[]) {
    const val = language.messages[key]
    return val ? fill(val, args) : ""
}

export async function fromResponse(response: Response, ...args: string[]) {
    if (response.status === 429) {
        args.shift()
        args[0] = String(getRetryAfter(response))
    }

    const json = await response.json()

    const val = language.errors[json?.error]
    return val ? fill(val, args) : JSON.stringify(json)
}
