import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import type {MouseEventHandler} from "preact";

export default function Button({ onClick, onMouseEnter, className, icon, text }: {onClick?: MouseEventHandler<HTMLButtonElement>, onMouseEnter?: MouseEventHandler<HTMLButtonElement>, className?: string, icon: IconDefinition, text: string}) {
    return (
        <>
            <button onClick={onClick} onMouseEnter={onMouseEnter} className={`z-1 ${className}`}>
                <FontAwesomeIcon icon={icon} /> {text}
            </button>
        </>
    )
}