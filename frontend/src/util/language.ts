import i18n from "i18next";
import {initReactI18next} from "react-i18next";
import {match} from "@formatjs/intl-localematcher";
import type {Me} from "../types/api.ts";
import {type ApiError, getRetryAfter} from "./api.tsx";

declare global {
    interface Window { __INITIAL_LANGUAGE__?: LanguageFile, __INITIAL_LANG__?: string }
}

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

type Vars = Record<string, string | number>

function positional(args: string[]): Vars {
    return Object.fromEntries(args.map((arg, i) => [String(i + 1), arg]))
}

export async function fetchLanguage(me: Me | null) {
    const lng = getLanguage(me)?.key ?? "en-US"
    const {messages, errors} = await getLanguageFile(lng)

    if (!i18n.isInitialized) {
        await i18n.use(initReactI18next).init({
            lng,
            fallbackLng: "en-US",
            ns: ["translation", "errors"],
            defaultNS: "translation",
            resources: {[lng]: {translation: messages, errors}},
            keySeparator: false,
            nsSeparator: false,
            interpolation: {escapeValue: false},
            parseMissingKeyHandler: () => "",
            returnEmptyString: true,
        })
        return
    }

    i18n.addResourceBundle(lng, "translation", messages, true, true)
    i18n.addResourceBundle(lng, "errors", errors, true, true)
    await i18n.changeLanguage(lng)
}

async function getLanguageFile(key: string): Promise<LanguageFile> {
    if (window.__INITIAL_LANG__ === key && window.__INITIAL_LANGUAGE__) {
        return window.__INITIAL_LANGUAGE__
    }

    const response = await fetch(`/assets/languages/${key}.json`)
    return await response.json() as LanguageFile
}

export function getLanguage(me: Me | null, auto?: boolean) {
    const stored = !auto ? me?.preferences.lang : null

    return languages.find(o => o.key === match(
        stored ? [stored] : navigator.languages,
        languages.map((o) => o.key),
        "en-US"
    ))
}

export function t(key: string, vars?: Vars) {
    return i18n.t(key, vars)
}

export function getKey(key: string, ...args: string[]) {
    const vars = positional(args)
    if (!i18n.exists(key)) {
        return i18n.t(key, {ns: "errors", ...vars})
    }
    return t(key, vars)
}


export async function fromResponse(response: Response, ...args: string[]) {
    if (response.status === 429) {
        args.shift()
        args[0] = String(getRetryAfter(response))
    }

    const json = await response.json().catch(() => null) as ApiError | null

    if (json?.error && i18n.exists(json.error, {ns: "errors"})) {
        return i18n.t(json.error, {ns: "errors", ...positional(args)})
    }
    return JSON.stringify(json)
}
