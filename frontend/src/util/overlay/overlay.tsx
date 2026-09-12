import "./overlay.css";

export default function Overlay({ onClick, time = 0, z = 150 }: { onClick?: () => void, time?: number, z?: number }) {
    return (
        <>
            <div className={`overlay animate-[fadeIn_${time}ms_ease-in-out] z-${z}`} onClick={onClick}></div>
        </>
    )
}