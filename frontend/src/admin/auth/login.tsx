import "./login.css";
import {fromResponse, getKey} from "../../util/language.ts";
import {useCallback, useEffect, useRef, useState} from "preact/compat";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight} from "@fortawesome/free-solid-svg-icons";
import Loading from "../../util/loading.tsx";

export default function AdminLogin() {
    const [wrong, setWrong] = useState(false)

    const [token, setToken] = useState("")

    const [waiting, setWaiting] = useState(false)
    const [remaining, setRemaining] = useState("")

    const timerRef = useRef<number | null>(null)

    const updateRemaining = useCallback(function step(remaining: number) {
        if (timerRef.current) {
            clearTimeout(timerRef.current)
        }

        if (remaining < 0) {
            setRemaining("")
            return
        }

        setRemaining(String(remaining))

        timerRef.current = setTimeout(() => {
            step(remaining - 1)
        }, 1000)
    }, [])

    const login = useCallback(async () => {
        if (wrong || remaining || !token) return
        setWaiting(true)
        try {
            const response = await fetch("/api/admin_token_login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    admin_token: token
                })
            })

            switch (response.status) {
                case 200:
                    location.reload()
                    break
                case 429: {
                    const after = Number(response.headers.get("retry-after"))
                    updateRemaining(after)
                    await fromResponse(response)
                    break
                }
                default:
                    await fromResponse(response)
            }
        } finally { setWaiting(false) }
    }, [wrong, remaining, token, updateRemaining])

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
                        <input className={`${wrong && "wrong"}`} type="password" placeholder={getKey("ADMIN_TOKEN")} value={token} onChange={(e) => {setToken(e.currentTarget.value); setWrong(false)}} />
                    </div>
                    <div className="logbutton">
                        <p className="text-(--wrong-color) text-3xl fredoka">{remaining}</p>
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