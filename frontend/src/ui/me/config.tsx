import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import React, {useState} from "react";
import Overlay from "../../util/overlay/overlay.tsx";

export default function Config({ icon, text, value, children, className}: {icon: IconDefinition, text: string, value: string, children: React.ReactNode | ((close: () => void) => React.ReactNode), className?: string}) {
    const [open, setOpen] = useState(false)
    const close = () => setOpen(false)

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
                    {typeof children === "function" ? children(close) : children}
                </div>
                <Overlay onClick={close} z={1} time={10000} />
            </>}
        </>
    )
}