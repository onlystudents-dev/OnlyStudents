import "./enroll.css";
import {getKey} from "../../util/language.ts";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight} from "@fortawesome/free-solid-svg-icons";
import {useCallback, useEffect, useState} from "preact/compat";
import PasswordCheck from "../pwc.tsx";
import {getPasswordConfig, type PasswordConfig} from "../passwordConfig.ts";
import {opaqueRegister} from "../opaque.ts";
import {toast} from "react-toastify";
import Loading from "../../util/loading.tsx";

export default function Enroll() {
    const enrollToken = new URLSearchParams(location.search).get("enroll_token");

    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")

    const [passed, setPassed] = useState(false)
    const [waiting, setWaiting] = useState(false)
    const [remaining, setRemaining] = useState("")

    const [config, setConfig] = useState<PasswordConfig | null>(null)

    const register = useCallback(async () => {
        if (!passed || remaining || !enrollToken) return

        setWaiting(true)

        try {
            const res = await opaqueRegister({enrollToken, password})

            switch (res.status) {
                case "ok":
                    location.href = "/me"
                    break
                case "ratelimited":
                    toast.error(getKey("TOO_MANY_REQUESTS", String(res.retryAfter)))
                    setRemaining(String(res.retryAfter))
                    break
                case "error":
                    toast.error(
                        res.error
                            ? getKey(res.error)
                            : getKey("ENROLL_FAILED")
                    )
                    break
            }
        } finally {
            setWaiting(false)
        }
    }, [passed, remaining, enrollToken, password])

    useEffect(() => {
        void getPasswordConfig().then(setConfig)

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Enter") void register()
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [register])

    return (
        <>
            <div className="cantar ebox">
                <div className="econtent">
                    <h1 className="self-center text-5xl font-bold mb-4 rubik">
                        {getKey("ENROLL")}
                    </h1>

                    <div className="loginput fredoka">
                        <input type="password" placeholder={getKey("PASSWORD")} maxLength={config?.maxLength} value={password} onChange={(e) => setPassword(e.currentTarget.value)} />
                        <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder={getKey("CONFIRM_PASSWORD")} maxLength={config?.maxLength} value={confirmPassword} onChange={(e) => setConfirmPassword(e.currentTarget.value)} />
                        <PasswordCheck password={password} confirmPassword={confirmPassword} setPassed={setPassed} />
                    </div>

                    <div className="elogbutton">
                        <p className="text-(--wrong-color) text-3xl fredoka">
                            {remaining}
                        </p>

                        <Button onClick={register}>
                            {getKey("ENROLL")} <FontAwesomeIcon icon={faArrowRight} />
                        </Button>
                    </div>
                </div>
            </div>

            {waiting && <Loading />}
        </>
    )
}
