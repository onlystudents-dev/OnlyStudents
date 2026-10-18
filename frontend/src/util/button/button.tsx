import "./button.css";
import type {ReactNode} from "react";
import type {Ref} from "preact";

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