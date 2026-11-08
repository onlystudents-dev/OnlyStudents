import "./save.css";
import React from "preact/compat";
import Button from "../button/button.tsx";
import {getKey} from "../language.ts";

function isObject(v: unknown): v is Record<string, unknown> {
    return v !== null && typeof v === "object" && !Array.isArray(v)
}

function isEqual(a: unknown, b: unknown): boolean {
    if (Object.is(a, b)) return true
    if (!isObject(a) || !isObject(b)) return false

    const keysA = Object.keys(a)
    const keysB = Object.keys(b)
    if (keysA.length !== keysB.length) return false

    return keysA.every(
        (k) => Object.hasOwn(b, k) && isEqual(a[k], b[k])
    )
}

export function getDifferences<T extends object>(options: T, defaults: T): Partial<T> {
    const defaultValues = defaults as Record<string, unknown>
    const optionValues = options as Record<string, unknown>

    return Object.fromEntries(
        [...new Set([...Object.keys(options), ...Object.keys(defaults)])]
            .filter(key => !isEqual(optionValues[key], defaultValues[key]))
            .map(key => [key, optionValues[key]])
    ) as Partial<T>
}

export type SaveResult<T> = boolean | (keyof T)[]

export default function Save<T extends Record<string, unknown>>({ options, setOptions, save }: { options: T, setOptions: (value: T) => void, save: (diff: Partial<T>, defaults: Partial<T>) => SaveResult<T> | Promise<SaveResult<T>> }) {
    const [defaults, setDefaults] = React.useState<T>(options)

    const differences = getDifferences(options, defaults)
    const count = Object.keys(differences).length

    async function handleSave() {
        const result = await save(differences, defaults)

        if (result === false) return

        if (result === true) {
            setDefaults({ ...options })
            return
        }

        const failed = new Set(result.map(String))
        const oldDefaults = defaults as Record<string, unknown>

        const next = Object.fromEntries(
            [...new Set([...Object.keys(options), ...Object.keys(defaults)])]
                .flatMap(key => {
                    if (failed.has(key)) {
                        return Object.hasOwn(oldDefaults, key) ? [[key, oldDefaults[key]]] : []
                    }
                    return Object.hasOwn(options, key) ? [[key, options[key]]] : []
                })
        ) as T

        setDefaults(next)
    }

    return (
        <>
            <div className={`save ${isEqual(options, defaults) ? "hid" : ""} fredoka`}>
                <p className="ml-2">{getKey("UNSAVED_CHANGES", String(count), count !== 1 ? getKey("MULTIPLE") : "")}</p>
                <div className="buttons rubik">
                    <Button onClick={() => setOptions({ ...defaults })} className="breset">
                        {getKey("RESET")}
                    </Button>
                    <Button onClick={handleSave} className="bsave">
                        {getKey("SAVE")}
                    </Button>
                </div>
            </div>
        </>
    )
}