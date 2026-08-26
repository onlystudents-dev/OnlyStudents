import "./overlay.css";

export default function Overlay({ onClick, time }: { onClick: () => void, time: number }) {
    return (
        <>
            <div className={`overlay animate-[fadeIn_${time}ms_ease-in-out]`} onClick={onClick}></div>
        </>
    )
}