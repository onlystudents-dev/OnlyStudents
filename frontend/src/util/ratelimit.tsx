import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faTriangleExclamation} from "@fortawesome/free-solid-svg-icons";
import {useCallback, useEffect, useRef, useState} from "react";
import {getKey} from "./language.ts";
import Overlay from "./overlay/overlay.tsx";

export default function RateLimit({ retry, expire, standalone }: {retry: number, expire: () => void, standalone?: boolean}) {
    const [remaining, setRemaining] = useState("");

    const timerRef = useRef<number | null>(null);

    const updateRemaining = useCallback(function step(remaining: number) {
        if (timerRef.current != null) {
            clearTimeout(timerRef.current)
        }

        if (remaining < 0) {
            expire()
            return
        }

        setRemaining(String(remaining));

        timerRef.current = setTimeout(() => {
            step(remaining - 1)
        }, 1000)
    }, [expire])

    useEffect(() => {
        updateRemaining(retry)
    }, [updateRemaining, retry])

    return (
        <>
            <div className={`cantar w-108 h-72 flex flex-col gap-4 justify-center items-center z-151 bg-(--bg-color) ${!standalone && "border-4 border-(--border-color) border-solid rounded-2xl p-4"}`}>
                <FontAwesomeIcon icon={faTriangleExclamation} size="7x" color="var(--wrong-color)" />
                <h1 className="text-3xl text-(--wrong-color) rubik font-bold">{getKey("RATE_LIMITED")}</h1>
                <p className="text-(--wrong-color) rubik font-bold">{getKey("BACK_IN", remaining)}</p>
            </div>
            {!standalone && <Overlay />}
        </>
    )
}