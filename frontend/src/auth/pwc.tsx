import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCheck, faXmark} from "@fortawesome/free-solid-svg-icons";
import React, {useEffect, useState} from "react";
import {getKey} from "../util/language.ts";
import {getPasswordConfig, isPasswordExposed, type PasswordConfig} from "./passwordConfig.ts";

export default function PasswordCheck({ password, confirmPassword, setPassed }: {password: string, confirmPassword?: string, setPassed: React.Dispatch<React.SetStateAction<boolean>>}) {
    const [config, setConfig] = useState<PasswordConfig | null>(null);
    const [exposed, setExposed] = useState<boolean | null>(null);

    useEffect(() => {
        getPasswordConfig().then(setConfig);
    }, []);

    useEffect(() => {
        function Fetch() {
            setExposed(null)
            if (!config?.hibpCheckEnabled || !password) {
                return
            }

            let active = true
            const timer = setTimeout(() => {
                isPasswordExposed(password).then(result => {
                    if (active) setExposed(result)
                })
            }, 400)

            return () => {
                active = false
                clearTimeout(timer)
            }
        }

        Fetch()
    }, [password, config?.hibpCheckEnabled]);

    const minLength = config?.minLength ?? 12
    const maxLength = config?.maxLength ?? 256

    const checks = [
        { label: getKey("CONDITION.CHARS", String(minLength)), ok: password.length >= minLength && password.length <= maxLength },
        { label: getKey("CONDITION.LOWER"), ok: /[a-z]/.test(password) },
        { label: getKey("CONDITION.UPPER"), ok: /[A-Z]/.test(password) },
        { label: getKey("CONDITION.NUMBERS"), ok: /[0-9]/.test(password) },
        { label: getKey("CONDITION.SPECIAL"), ok: /[^a-zA-Z0-9]/.test(password) },
    ];
    if (confirmPassword !== undefined) {
        checks.push({ label: getKey("CONDITION.MATCH"), ok: password === confirmPassword })
    }
    if (config?.hibpCheckEnabled) {
        checks.push({ label: getKey("CONDITION.EXPOSED"), ok: exposed === false })
    }

    const passed = checks.filter(c => c.ok).length
    const didPassed = passed === checks.length

    useEffect(() => {
        setPassed(didPassed)
    }, [didPassed, setPassed])

    return (
        <div className="flex flex-col justify-center items-start p-4 gap-2 w-full">
            <div className="relative w-full h-1.5 mb-4">
                <div className="absolute rounded-lg bg-(--right-color) h-1.5 top-0 left-0 z-1 transition-all shadow-[0_0_10px_var(--right-color)]" style={{ width: `${(passed / checks.length) * 100}%` }} />
                <div className="absolute rounded-lg bg-(--wrong-color) h-1 top-px left-0 w-full shadow-[0_0_10px_var(--wrong-color)]" />
            </div>
            {checks.map(({label, ok}) => (
                <p key={label} className={`${ok ? "text-(--right-color)" : "text-(--wrong-color)"} transition-[color] duration-200`}>
                    <FontAwesomeIcon icon={ok ? faCheck : faXmark} /> {label}
                </p>
            ))}
        </div>
    );
}
