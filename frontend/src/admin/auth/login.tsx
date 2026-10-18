import "./login.css";
import {fromResponse, getKey} from "../../util/language.ts";
import {useCallback, useEffect, useRef, useState} from "preact/compat";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight} from "@fortawesome/free-solid-svg-icons";
import Loading from "../../util/loading.tsx";
import {toast} from "react-toastify";

export default function AdminLogin() {
    const [invalidToken, setInvalidToken] = useState(false)

    const [adminToken, setAdminToken] = useState("")

    const [waiting, setWaiting] = useState(false)
    const [secondsLeft, setSecondsLeft] = useState("")

    const timerRef = useRef<number | null>(null)

    const startCountdown = useCallback(function step(seconds: number) {
        if (timerRef.current) {
            clearTimeout(timerRef.current)
        }

        if (seconds < 0) {
            setSecondsLeft("")
            return
        }

        setSecondsLeft(String(seconds))

        timerRef.current = setTimeout(() => {
            step(seconds - 1)
        }, 1000)
    }, [])

    const login = useCallback(async () => {
        if (invalidToken || secondsLeft || !adminToken) return
        setWaiting(true)
        try {
            const response = await fetch("/api/v1/admin_token_login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    admin_token: adminToken
                })
            })

            switch (response.status) {
                case 200:
                    location.reload()
                    break
                case 429: {
                    const after = Number(response.headers.get("retry-after"))
                    startCountdown(after)
                    await fromResponse(response)
                    break
                }
                default:
                    toast.error(await fromResponse(response))
            }
        } finally { setWaiting(false) }
    }, [invalidToken, secondsLeft, adminToken, startCountdown])

    useEffect(() => {
        const handleKeyDown = async (e: KeyboardEvent) => {
            if (e.key === "Enter") await login()
        }

        window.addEventListener("keydown", e => void handleKeyDown(e))

        return () => {
            window.removeEventListener("keydown", e => void handleKeyDown(e))
        }
    }, [login])

    return (
        <>
            <div className="cantar box">
                <div className="content">
                    <h1 className="self-center text-5xl font-bold mb-4 rubik">{getKey("LOGIN_TITLE")}</h1>

                    <div className="loginput fredoka">
                        <input className={`${invalidToken && "wrong"}`} type="password" placeholder={getKey("ADMIN_TOKEN")} value={adminToken} onChange={(e) => {setAdminToken(e.currentTarget.value); setInvalidToken(false)}} />
                    </div>
                    <div className="logbutton">
                        <p className="text-(--wrong-color) text-3xl fredoka">{secondsLeft}</p>
                        <Button onClick={login}>
                            {getKey("LOGIN")} <FontAwesomeIcon icon={faArrowRight} />
                        </Button>
                    </div>
                </div>
            </div>
            {waiting && <Loading />}
        </>
    )
}
