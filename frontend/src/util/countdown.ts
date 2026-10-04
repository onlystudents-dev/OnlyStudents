import {useCallback, useEffect, useRef, useState} from "react";

export function useCountdown(onExpire?: () => void) {
    const [remaining, setRemaining] = useState("")
    const timer = useRef<number | null>(null)
    const expire = useRef(onExpire)

    useEffect(() => {
        expire.current = onExpire
    })

    const stop = useCallback(() => {
        if (timer.current != null) clearTimeout(timer.current)
        timer.current = null
    }, [])

    const start = useCallback(function step(seconds: number) {
        stop()

        if (seconds < 0) {
            setRemaining("")
            expire.current?.()
            return
        }

        setRemaining(String(seconds))
        timer.current = window.setTimeout(() => step(seconds - 1), 1000)
    }, [stop])

    useEffect(() => stop, [stop])

    return [remaining, start] as const
}
