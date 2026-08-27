import Overlay from "../overlay/overlay.tsx";

export default function Loading() {
    return (
        <>
            <Overlay />
            <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-151">
                <div className="h-10 w-10 animate-spin rounded-full border-4 border-gray-600 border-t-white" />
            </div>
        </>
    )
}