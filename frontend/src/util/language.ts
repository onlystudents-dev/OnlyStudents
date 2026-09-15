import {match} from "@formatjs/intl-localematcher";

export type language = {
    key: string,
    name: string,
    emoji: string,
}

let languages: language[]
export async function fetchLanguages() {
    const response = await fetch("/assets/languages.json")
    languages = await response.json()
}

let language: Record<string, Record<string, string>>
export async function fetchLanguage() {
    let response = await fetch(`/assets/languages/${getLanguage()?.key}.json`)
    if (response.status === 404) response = await fetch("/assets/languages/en-US.json")
    language = await response.json()
}

export function getLanguage(auto?: boolean) {
    const stored = !auto ? localStorage.getItem("language") : null;

    return languages.find(o => o.key === match(
        stored ? [stored] : navigator.languages,
        languages.map((o) => o.key),
        "en-US"
    ))
}

export function getLanguages() {
    return languages
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