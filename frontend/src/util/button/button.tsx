import "./button.css";
import type {ReactNode} from "react";
import type {Ref} from "preact";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import type {IconDefinition} from "@fortawesome/fontawesome-svg-core";

export default function Button({ href, onClick, onMouseEnter, className, children, refA, refB }: {href?: string, onClick?: (e: MouseEvent) => void, onMouseEnter?: (e: MouseEvent) => void, className?: string, children?: ReactNode, refA?: Ref<HTMLAnchorElement>, refB?: Ref<HTMLButtonElement>}) {
    if (href) return (
        <>
            <a href={href} onClick={onClick} onMouseEnter={onMouseEnter} className={`button ${className}`} ref={refA}>
                {children}
            </a>
        </>
    )

    return (
        <>
            <button onClick={onClick} onMouseEnter={onMouseEnter} className={`button ${className}`} ref={refB}>
                {children}
            </button>
        </>
    )
}

export function Badge({ icon, className }: {icon: IconDefinition, className: string}) {
    return (
        <span className={`rounded-full size-8 inline-flex items-center justify-center ${className}`}>
            <FontAwesomeIcon icon={icon} />
        </span>
    )
}