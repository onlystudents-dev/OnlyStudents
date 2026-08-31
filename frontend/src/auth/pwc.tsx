import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCheck, faXmark} from "@fortawesome/free-solid-svg-icons";
import React, {useEffect} from "react";

export default function PasswordCheck({ password, confirmPassword, setPassed }: {password: string, confirmPassword?: string, setPassed: React.Dispatch<React.SetStateAction<boolean>>}) {
    const checks = [
        { label: "At least 12 characters",    ok: password.length >= 12 },
        { label: "Contains lowercase letters", ok: /[a-z]/.test(password) },
        { label: "Contains uppercase letters", ok: /[A-Z]/.test(password) },
        { label: "Contains numbers",           ok: /[0-9]/.test(password) },
        { label: "Contains special characters", ok: /[^a-zA-Z0-9]/.test(password) },
    ];
    if (confirmPassword !== undefined) {
        checks.push({ label: "The two passwords match", ok: password === confirmPassword })
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
