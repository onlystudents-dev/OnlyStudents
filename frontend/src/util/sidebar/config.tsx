import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import type {HTMLInputTypeAttribute, MouseEventHandler} from "preact";
import React from "preact/compat";
import "./config.css";

export default function Button({ onClick, onMouseEnter, className, icon, text }: {onClick?: MouseEventHandler<HTMLButtonElement>, onMouseEnter?: MouseEventHandler<HTMLButtonElement>, className?: string, icon: IconDefinition, text: string}) {
    return (
        <>
            <button onClick={onClick} onMouseEnter={onMouseEnter} className={`z-1 ${className}`}>
                <FontAwesomeIcon icon={icon} /> {text}
            </button>
        </>
    )
}

export function DropdownConfig({ text, icon, value, children, className, onChange }: { text?: string, icon?: IconDefinition, value: string, children: React.ReactNode, className?: string, onChange: (raw: string) => void }) {
    return (
        <>
            <div className={`config ${className}`}>
                <p>
                    {icon && <FontAwesomeIcon icon={icon} />} {text}
                </p>
                <select className="poppins" value={value} onChange={(e) => onChange(e.currentTarget.value)}>
                    {children}
                </select>
            </div>
        </>
    )
}

export function InputConfig({ text, icon, value, className, onChange, onFocusIn, onFocusOut, type }: {text?: string, icon?: IconDefinition, value: string | number | boolean, className?: string, onChange: (raw: string) => void, onFocusIn?: () => void, onFocusOut?: () => void, type?: HTMLInputTypeAttribute}) {
    let val: string | number | undefined
    let checked: boolean | undefined
    switch (typeof value) {
        case "boolean":
            checked = value
            break
        case "string":
            val = value
            break
        case "number":
            val = value
            break
    }

    return (
        <>
            <div className={`config ${className}`}>
                <p>
                    {icon && <FontAwesomeIcon icon={icon} />} {text}
                </p>
                <input className="poppins" value={val} checked={checked} type={type} onChange={(e) => onChange(e.currentTarget.value)} onFocusIn={onFocusIn} onFocusOut={onFocusOut} />
            </div>
        </>
    )
}