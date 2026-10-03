import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faTriangleExclamation} from "@fortawesome/free-solid-svg-icons";
import {useEffect} from "react";
import {getKey} from "./language.ts";
import Overlay from "./overlay/overlay.tsx";
import {useCountdown} from "./countdown.ts";
import {cx} from "./cx.ts";

export default function RateLimit({ retry, expire, standalone }: {retry: number, expire: () => void, standalone?: boolean}) {
    const [remaining, start] = useCountdown(expire)

    useEffect(() => {
        start(retry)
    }, [start, retry])

    return (
        <>
            <div className={cx("cantar w-108 h-72 flex flex-col gap-4 justify-center items-center z-151 bg-(--bg-color)", !standalone && "border-4 border-(--border-color) border-solid rounded-2xl p-4")}>
                <FontAwesomeIcon icon={faTriangleExclamation} size="7x" color="var(--wrong-color)" />
                <h1 className="text-3xl text-(--wrong-color) rubik font-bold">{getKey("RATE_LIMITED")}</h1>
                <p className="text-(--wrong-color) rubik font-bold">{getKey("BACK_IN", remaining)}</p>
            </div>
            {!standalone && <Overlay />}
        </>
    )
}