import Overlay from "./overlay/overlay.tsx";

export default function Loading({ absolute }: { absolute?: boolean }) {
    return (
        <>
            <div className={`${absolute ? "absolute" : "fixed"} left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-161`}>
                <div className="h-10 w-10 animate-spin rounded-full border-4 border-(--border-color) border-t-(--txt-color)" />
            </div>
            <Overlay z={160} />
        </>
    )
}