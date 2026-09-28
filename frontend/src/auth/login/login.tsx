import "./login.css";
import Button from "../../util/button/button.tsx";
import {useCallback, useEffect, useRef, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight, faQuestion} from "@fortawesome/free-solid-svg-icons";
import PasswordReset from "../pwr/pwr.tsx";
import {toast} from "react-toastify";
import Loading from "../../util/loading.tsx";
import {getKey} from "../../util/language.ts";
import { opaqueLogin } from "../opaqueLogin.ts";

export default function Login() {
    const [pwr, setPwr] = useState(false);
    const [pwrA, setPwrA] = useState(false);
    const [reset, setReset] = useState(false)

    const [red, setRed] = useState(false);
    const [wrong, setWrong] = useState(false);

    const [role, setRole] = useState("guardian");

    const [id, setId] = useState("");
    const [password, setPassword] = useState("");

    const [waiting, setWaiting] = useState(false);
    const [remaining, setRemaining] = useState("");

    const timerRef = useRef<number | null>(null);

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
      if (red || wrong || remaining || !id || !password) return;
      setWaiting(true);
      try {
        const res = await opaqueLogin({ userId: Number(id), role, password });
        switch (res.status) {
          case "ok": location.reload(); break;
          case "wrong": toast.error(getKey("WRONG_CREDENTIALS")); setWrong(true); break;
          case "ratelimited":
            toast.error(getKey("TOO_MANY_REQUESTS", String(res.retryAfter)));
            updateRemaining(res.retryAfter);
            break;
          case "error": toast.error(res.error ? getKey(res.error) : getKey("LOGIN_FAILED")); break;
        }
      } finally { setWaiting(false); }
    }, [red, wrong, remaining, id, password, role, updateRemaining]);

    useEffect(() => {
        const handleKeyDown = async (e: KeyboardEvent) => {
            if (e.key === "Enter") await login()
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [login])

    return (
        <>
            <div className={`box cantar ${reset ? "h-158" : "h-100"}`}>
                <div className={`content out ${pwrA && "hid"} h-86`}>
                    <h1 className="self-center text-5xl font-bold mb-8 rubik">{getKey("LOGIN_TITLE")}</h1>
                    <div className="loginput fredoka">
                        <select className="poppins" value={role} onChange={(e) => {setRole(e.currentTarget.value); setWrong(false)}}>
                            <option value="guardian">{getKey("ROLE.GUARDIAN")}</option>
                            <option value="student">{getKey("ROLE.STUDENT")}</option>
                            <option value="teacher">{getKey("ROLE.TEACHER")}</option>
                        </select>
                        <input className={`${(red || wrong) && "wrong"}`} type="text" placeholder={getKey("USER_ID")} value={id} onChange={(e) => {setId(e.currentTarget.value); checkUserID(e.currentTarget.value); setWrong(false)}} />
                        <input className={`${wrong && "wrong"}`} type="password" placeholder={getKey("PASSWORD")} value={password} onChange={(e) => {setPassword(e.currentTarget.value); setWrong(false)}} />
                    </div>
                    <div aria-disabled className="logbutton">
                        <Button onClick={sPwr}>
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
                {pwr && <PasswordReset pwrA={pwrA} unsPwr={unsPwr} role={role} setRole={setRole} id={id} setId={setId} red={red} checkUserID={checkUserID} setWaiting={setWaiting} setResetL={setReset} />}
            </div>
            {waiting && <Loading />}
        </>
    )

    function checkUserID(id: string) {
        if (id && Number.isNaN(Number(id))) {
            setRed(true);
        } else {
            setRed(false);
        }
    }

    function sPwr() {
        setPwr(true)
        setPwrA(true)
    }

    function unsPwr() {
        setReset(false)
        setPwrA(false)
        setTimeout(() => setPwr(false), 200)
    }
}
