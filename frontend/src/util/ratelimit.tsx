import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faTriangleExclamation} from "@fortawesome/free-solid-svg-icons";
import {useCallback, useEffect, useRef, useState} from "react";
import Overlay from "./overlay/overlay.tsx";

export default function RateLimit({ retry }: {retry: number}) {
    const [remaining, setRemaining] = useState("");

    const timerRef = useRef<number | null>(null);

    const updateRemaining = useCallback(function step(remaining: number) {
        if (timerRef.current) {
            clearTimeout(timerRef.current)
        }

        if (remaining < 0) {
            location.reload()
            return
        }

        setRemaining(String(remaining));

        timerRef.current = setTimeout(() => {
            step(remaining - 1)
        }, 1000)
    }, [])

    useEffect(() => {
        updateRemaining(retry)
    }, [updateRemaining, retry])

    return (
        <>
            <div className="fixed inset-0 flex flex-col gap-4 justify-center items-center z-151 border-4 border-(--border-color) bg-(--bg-color)">
                <FontAwesomeIcon icon={faTriangleExclamation} size="7x" color="var(--wrong-color)" />
                <h1 className="text-3xl text-(--wrong-color) rubik font-bold">You've been rate limited!</h1>
                <p className="text-(--wrong-color) rubik font-bold">Will auto-refresh in {remaining} seconds</p>
            </div>
            <Overlay />
        </>
    )
}