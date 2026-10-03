import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import React, {useState} from "react";
import Overlay from "../../util/overlay/overlay.tsx";
import {cx} from "../../util/cx.ts";

export default function Config({ icon, text, value, children, className}: {icon: IconDefinition, text: string, value: string, children: React.ReactNode, className?: string}) {
    const [open, setOpen] = useState(false)

    return (
        <>
            <div className={cx("config", className)} onClick={() => setOpen(true)}>
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

export function DropdownConfig({ text, icon, value, children, className, onChange }: { text: string, icon: IconDefinition, value: string, children: React.ReactNode, className?: string, onChange: (raw: string) => void }) {
    return (
        <>
            <div className={cx("config", className)}>
                <p>
                    <FontAwesomeIcon icon={icon} /> {text}
                </p>
                <select className="poppins" value={value} onChange={(e) => onChange(e.currentTarget.value)}>
                    {children}
                </select>
            </div>
        </>
    )
}
