import {match} from "@formatjs/intl-localematcher";
import type {Me} from "../types/api.ts";
import {type ApiError, getRetryAfter, readJSON} from "./api.ts";

export type Language = {
    key: string,
    name: string,
    emoji: string,
}

export const languages: Language[] = [
    {"key": "en-US", "name": "English (US)", "emoji": "\uD83C\uDDFA\uD83C\uDDF8"},
    {"key": "hu-HU", "name": "Magyar", "emoji": "\uD83C\uDDED\uD83C\uDDFA"}
]

type LanguageFile = {
    messages: Record<string, string>,
    errors: Record<string, string>,
}

let language: LanguageFile = { messages: {}, errors: {} }

export async function fetchLanguage(me: Me | null) {
    const response = await fetch(`/assets/languages/${getLanguage(me)?.key ?? "en-US"}.json`);
    language = await readJSON<LanguageFile>(response);
}

export function getLanguage(me: Me | null, auto?: boolean) {
    const stored = !auto ? me?.preferences.lang : null;

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

    const json = await readJSON<ApiError | null>(response).catch(() => null)

    const val = json?.error ? language.errors[json.error] : undefined
    return val ? fill(val, args) : JSON.stringify(json)
}
