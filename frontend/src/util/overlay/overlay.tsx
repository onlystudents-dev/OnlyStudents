import { useEffect } from "preact/hooks";
import "./overlay.css";

export default function Overlay({ onClick, time = 0, z = 150 }: { onClick?: () => void, time?: number, z?: number }) {
    useEffect(() => {
        if (!onClick) return

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") onClick()
        }

        window.addEventListener("keydown", handleKeyDown)
        return () => window.removeEventListener("keydown", handleKeyDown)
    }, [onClick])

    return (
        <>
            <div className={`overlay`} style={{ zIndex: z, animation: `fadeIn ${time} ease-in-out` }} onClick={onClick}></div>
        </>
    )
}
