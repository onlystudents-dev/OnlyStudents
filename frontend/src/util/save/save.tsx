import "./save.css";
import React from "react";
import Button from "../button/button.tsx";
import {getKey} from "../language.ts";

function isObject(v: unknown): v is Record<string, unknown> {
    return v !== null && typeof v === "object" && !Array.isArray(v);
}

function isEqual(a: unknown, b: unknown): boolean {
    if (Object.is(a, b)) return true;
    if (!isObject(a) || !isObject(b)) return false;

    const keysA = Object.keys(a);
    const keysB = Object.keys(b);
    if (keysA.length !== keysB.length) return false;

    return keysA.every(
        (k) => Object.hasOwn(b, k) && isEqual(a[k], b[k])
    );
}

function getDifferences(
    options: Record<string, unknown> = {},
    defaults: Record<string, unknown> = {},
): Record<string, unknown> {
    return Object.fromEntries(
        Object.entries(options).filter(([key, value]) =>
            key in defaults ? !isEqual(value, defaults[key]) : true
        )
    );
}

export default function Save({ options, setOptions, save }: {options: Record<string, unknown>, setOptions: React.Dispatch<React.SetStateAction<Record<string, unknown>>>, save: (options: Record<string, unknown>) => (boolean | Promise<boolean>)}) {
    const [defaults, setDefaults] = React.useState<Record<string, unknown>>(options);

    const differences = getDifferences(options, defaults)
    const count = Object.keys(differences).length

    return (
        <>
            <div className={`save ${isEqual(options, defaults) && "hid"} fredoka`}>
                <p className="ml-2">{getKey("UNSAVED_CHANGES", count as unknown as string, count !== 1 ? getKey("MULTIPLE") : "")}</p>
                <div className="buttons rubik">
                    <Button onClick={() => setOptions({ ...defaults })} className="breset">
                        {getKey("RESET")}
                    </Button>
                    <Button onClick={() => {if (!save(differences)) return; setDefaults({ ...options })}} className="bsave">
                        {getKey("SAVE")}
                    </Button>
                </div>
            </div>
        </>
    )
}
