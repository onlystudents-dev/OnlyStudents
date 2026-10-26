import "./login.css";
import Button from "../../util/button/button.tsx";
import {useCallback, useEffect, useState} from "preact/compat";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight, faQuestion} from "@fortawesome/free-solid-svg-icons";
import PasswordReset from "../pwr/pwr.tsx";
import {toast} from "react-toastify";
import Loading from "../../util/loading.tsx";
import {getKey} from "../../util/language.ts";
import { opaqueLogin } from "../opaque.ts";
import RoleSelect from "../roleSelect.tsx";
import { useCountdown } from "../../util/countdown.ts";
import type {Me} from "../../types/api.ts";
import Overlay from "../../util/overlay/overlay.tsx";

export default function Login({ fetchMe }: {fetchMe?: () => Promise<Me | null>}) {
    const [resetOpen, setResetOpen] = useState(false)
    const [resetActive, setResetActive] = useState(false)
    const [resetSent, setResetSent] = useState(false)

    const [red, setRed] = useState(false)
    const [wrong, setWrong] = useState(false)

    const [role, setRole] = useState("guardian")

    const [id, setId] = useState("")
    const [password, setPassword] = useState("")

    const [waiting, setWaiting] = useState(false)
    const [remaining, updateRemaining] = useCountdown()

    const login = useCallback(async () => {
      if (red || wrong || remaining || !id || !password) return
      setWaiting(true)
      try {
          const res = await opaqueLogin({ userId: Number(id), role, password });
          switch (res.status) {
              case "ok":
                  localStorage.removeItem("child_id")
                  localStorage.removeItem("school_id")
                  if (fetchMe !== undefined) {
                      await fetchMe()
                  } else {
                      location.reload()
                  }
                  break
              case "wrong":
                  toast.error(getKey("WRONG_CREDENTIALS"))
                  setWrong(true)
                  break
              case "ratelimited":
                  toast.error(getKey("TOO_MANY_REQUESTS", String(res.retryAfter)))
                  updateRemaining(res.retryAfter)
                  break
              case "error":
                  toast.error(res.error ? getKey(res.error) : getKey("LOGIN_FAILED"))
                  break
          }
      } finally {
          setWaiting(false)
      }
    }, [red, wrong, remaining, id, password, role, fetchMe, updateRemaining])

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Enter") void login()
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [login])

    return (
        <>
            <div className={`box cantar ${resetSent ? "h-184" : "h-100"}`}>
                <div className={`content out ${resetActive && "hid h-86"}`}>
                    <h1 className="self-center text-5xl font-bold mb-8 rubik">{getKey("LOGIN_TITLE")}</h1>
                    <div className="loginput fredoka">
                        <RoleSelect value={role} onChange={(r) => { setRole(r); setWrong(false) }} />
                        <input className={`${(red || wrong) && "wrong"}`} type="text" placeholder={getKey("USER_ID")} value={id} onChange={(e) => {setId(e.currentTarget.value); checkUserID(e.currentTarget.value); setWrong(false)}} />
                        <input className={`${wrong && "wrong"}`} type="password" placeholder={getKey("PASSWORD")} value={password} onChange={(e) => {setPassword(e.currentTarget.value); setWrong(false)}} />
                    </div>
                    <div className="logbutton">
                        <Button onClick={openReset}>
                            <FontAwesomeIcon icon={faQuestion} /> {getKey("FORGOT_PASSWORD")}
                        </Button>
                        <div className="flex flex-row items-center gap-4">
                            <p className="text-(--wrong-color) text-3xl fredoka">{remaining}</p>
                            <Button onClick={login}>
                                {getKey("LOGIN")} <FontAwesomeIcon icon={faArrowRight} />
                            </Button>
                        </div>
                    </div>
                </div>
                {resetOpen && <PasswordReset resetActive={resetActive} onClose={closeReset} role={role} setRole={setRole} id={id} setId={setId} red={red} checkUserID={checkUserID} setWaiting={setWaiting} setResetSent={setResetSent} />}
            </div>
            {waiting && <Loading />}
            {fetchMe !== undefined && <Overlay z={148} />}
        </>
    )

    function checkUserID(id: string) {
        if (id && Number.isNaN(Number(id))) {
            setRed(true)
        } else {
            setRed(false)
        }
    }

    function openReset() {
        setResetOpen(true)
        setResetActive(true)
    }

    function closeReset() {
        setResetSent(false)
        setResetActive(false)
        setTimeout(() => setResetOpen(false), 200)
    }
}
