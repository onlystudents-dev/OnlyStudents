import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import React, {useState} from "react";
import Overlay from "../../util/overlay/overlay.tsx";

export default function Config({ icon, text, value, children }: {icon: IconDefinition, text: string, value: string, children: React.ReactNode}) {
    const [open, setOpen] = useState(false)

    return (
        <>
            <div className="config">
                <button onClick={() => setOpen(true)}>
                    <FontAwesomeIcon icon={icon} /> {text}
                </button>
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