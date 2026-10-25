import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {type HTMLInputTypeAttribute, type MouseEventHandler} from "preact";
import "./config.css";
import {type Option, selectStyles2} from "../select.tsx";
import Select, {type GroupBase, type StylesConfig} from "react-select";

export default function Button({ onClick, onMouseEnter, className, icon, text }: {onClick?: MouseEventHandler<HTMLButtonElement>, onMouseEnter?: MouseEventHandler<HTMLButtonElement>, className?: string, icon: IconDefinition, text: string}) {
    return (
        <>
            <button onClick={onClick} onMouseEnter={onMouseEnter} className={`z-1 ${className}`}>
                <FontAwesomeIcon icon={icon} /> {text}
            </button>
        </>
    )
}

export function DropdownConfig({ text, icon, value, children, className, onChange, styles = selectStyles2() }: { text?: string, icon?: IconDefinition, value: string, children: Option[], className?: string, onChange: (raw: string) => void, styles: StylesConfig<Option, false, GroupBase<Option>> | undefined }) {
    return (
        <>
            <div className={`config ${className}`}>
                <p>
                    {icon && <FontAwesomeIcon icon={icon} />} {text}
                </p>
                <Select
                    className="poppins"
                    value={children.find(val => val.value === value)}
                    options={children}
                    onChange={(val) => onChange(val?.value || "")}
                    isSearchable={false}
                    components={{ IndicatorSeparator: null }}
                    styles={styles}
                />
            </div>
        </>
    )
}

export function InputConfig({ text, icon, value, className, onChange, onFocusIn, onFocusOut, type, maxLength }: {text?: string, icon?: IconDefinition, value: string | number | boolean, className?: string, onChange: (raw: string) => void, onFocusIn?: () => void, onFocusOut?: () => void, type?: HTMLInputTypeAttribute, maxLength?: number}) {
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
                <input className="poppins" value={val} checked={checked} type={type} maxLength={maxLength} onChange={(e) => onChange(e.currentTarget.value)} onFocusIn={onFocusIn} onFocusOut={onFocusOut} />
            </div>
        </>
    )
}