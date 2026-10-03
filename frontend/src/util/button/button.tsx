import "./button.css";
import type {ReactNode} from "react";
import {cx} from "../cx.ts";

export default function Button({ href, onClick, className, children }: {href?: string, onClick?: () => void, className?: string, children?: ReactNode}) {
    if (href) return (
        <>
            <a href={href} onClick={onClick} className={cx("button", className)}>
                {children}
            </a>
        </>
    )

    return (
        <>
            <button onClick={onClick} className={cx("button", className)}>
                {children}
            </button>
        </>
    )
}