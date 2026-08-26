import "./button.css";
import type {ReactNode} from "react";

export default function Button({ href, onClick, className, children }: {href?: string, onClick?: () => void, className?: string, children?: ReactNode}) {
    if (href) return (
        <>
            <a href={href} onClick={onClick} className={`button ${className}`}>
                {children}
            </a>
        </>
    )

    return (
        <>
            <button onClick={onClick} className={`button ${className}`}>
                {children}
            </button>
        </>
    )
}