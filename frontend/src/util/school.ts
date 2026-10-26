import {useEffect, useRef} from "preact/compat";

let page: (() => void) | undefined

export function reloadSchoolStuff() {
    page?.()
}

export function useReloadSchoolStuff(fetch: () => void) {
    const ref = useRef(fetch)

    useEffect(() => {
        ref.current = fetch
    })

    useEffect(() => {
        page = () => ref.current()

        return () => {
            page = undefined
        }
    }, [])
}
