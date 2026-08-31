import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";

export default function Button({ onClick, onMouseEnter, className, icon, text }: {onClick?: (button: HTMLButtonElement) => void, onMouseEnter?: (button: HTMLButtonElement) => void, className?: string, icon: IconDefinition, text: string}) {
    return (
        <>
            <button onClick={(e) => onClick?.(e.currentTarget)}  onMouseEnter={(e) => onMouseEnter?.(e.currentTarget)} className={`z-1 ${className}`}>
                <FontAwesomeIcon icon={icon} /> {text}
            </button>
        </>
    )
}