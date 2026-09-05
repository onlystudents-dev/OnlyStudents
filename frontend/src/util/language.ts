let language: Record<string, Record<string, string>>
export async function fetchLanguage() {
    let response = await fetch(`/assets/languages/${localStorage.getItem("language") || navigator.language}.json`)
    if (response.status === 404) response = await fetch("/assets/languages/en-US.json")
    language = await response.json()
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