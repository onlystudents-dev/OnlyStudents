import "./overlay.css";

export default function Overlay({ onClick, time = 0, z = 150 }: { onClick?: () => void, time?: number, z?: number }) {
    return (
        <>
            <div className={`overlay`} style={{ zIndex: z, animation: `fadeIn ${time} ease-in-out` }} onClick={onClick}></div>
        </>
    )
}