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
                <Overlay onClick={() => {setOpen(false)}} z={1} time={10000} />
            </>}
        </>
    )
}

export function DropdownConfig({ text, icon, value, children, className, options, setOptions, lkey }: { text: string, icon: IconDefinition, value: string, children: React.ReactNode, className?: string, options: Record<option, unknown>, setOptions: React.Dispatch<React.SetStateAction<Record<option, unknown>>>, lkey: option}) {
    const ref = useRef<HTMLSelectElement>(null)

    return (
        <>
            <div className={`config ${className}`}>
                <p>
                    <FontAwesomeIcon icon={icon} /> {text}
                </p>
                <select className="poppins" ref={ref} value={value} onChange={(e) => {
                    const raw = e.currentTarget.value
                    const original = options[lkey]

                    const newValue = (() => {
                        switch (typeof original) {
                            case "number":
                                return Number(raw)
                            case "boolean":
                                return raw === "true"
                            default:
                                return raw
                        }
                    })()

                    setOptions({ ...options, [lkey]: newValue })
                }}>
                    {children}
                </select>
            </div>
        </>
    )
}