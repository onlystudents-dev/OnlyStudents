import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import React, {useRef, useState} from "react";
import Overlay from "../../util/overlay/overlay.tsx";
import type {option} from "./me.tsx";

export default function Config({ icon, text, value, children, className}: {icon: IconDefinition, text: string, value: string, children: React.ReactNode, className?: string}) {
    const [open, setOpen] = useState(false)

    return (
        <>
            <div className={`config ${className}`} onClick={() => setOpen(true)}>
                <p>
                    <FontAwesomeIcon icon={icon} /> {text}
                </p>
                <p>
                    {value}
                </p>
            </div>
            {open && <>
                <div className="cantar configure configs">
                    {children}
                </div>
                <Overlay onClick={() => {setOpen(false)}} />
            </>}
        </>
    )
}

export function DropdownConfig({ icon, text, value, children, className, options, setOptions, lkey }: {icon: IconDefinition, text: string, value: string, children: React.ReactNode, className?: string, options: Record<option, unknown>, setOptions: React.Dispatch<React.SetStateAction<Record<option, unknown>>>, lkey: option}) {
    const ref = useRef<HTMLSelectElement>(null)

    return (
        <>
            <div className={`config ${className}`}>
                <p>
                    <FontAwesomeIcon icon={icon} /> {text}
                </p>
                <select className="poppins" ref={ref} value={value} onChange={(e) => setOptions({...options, [lkey]: e.currentTarget.value})}>
                    {children}
                </select>
            </div>
        </>
    )
}