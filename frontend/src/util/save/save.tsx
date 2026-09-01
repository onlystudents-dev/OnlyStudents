import "./save.css";
import React from "react";
import Button from "../button/button.tsx";

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

export default function Save({ options, setOptions, save }: {options: Record<string, unknown>, setOptions: React.Dispatch<React.SetStateAction<Record<string, unknown>>>, save: (options: Record<string, unknown>) => void}) {
    const [defaults, setDefaults] = React.useState<Record<string, unknown>>(options);

    const differences = getDifferences(options, defaults)
    const count = Object.keys(differences).length

    return (
        <>
            <div className={`save ${isEqual(options, defaults) && "hid"} fredoka`}>
                <p className="ml-2">You have {count} unsaved change{count !== 1 ? "s" : ""}.</p>
                <div className="buttons rubik">
                    <Button onClick={() => setOptions({ ...defaults })} className="breset">
                        Reset
                    </Button>
                    <Button onClick={() => {save(differences); setDefaults({ ...options })}} className="bsave">
                        Save
                    </Button>
                </div>
            </div>
        </>
    )
}
