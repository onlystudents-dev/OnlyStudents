import "./pwr.css";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowLeft, faPaperPlane, faPlane} from "@fortawesome/free-solid-svg-icons";
import React, {useCallback, useEffect, useRef, useState} from "react";
import {toast} from "react-toastify";
import PasswordCheck from "../pwc.tsx";
import {fromResponse, getKey} from "../../util/language.ts";
import {getPasswordConfig, type PasswordConfig} from "../passwordConfig.ts";

export default function PasswordReset({ pwrA, unsPwr, role, setRole, id, setId, red, checkUserID, setWaiting, setResetL }: {pwrA: boolean, unsPwr: () => void, role: string, setRole: React.Dispatch<React.SetStateAction<string>>, id: string, setId: React.Dispatch<React.SetStateAction<string>>, red: boolean, checkUserID: (id: string) => void, setWaiting: React.Dispatch<React.SetStateAction<boolean>>, setResetL:  React.Dispatch<React.SetStateAction<boolean>>}) {
    const [reset, setReset] = useState(false);
    const [resetA, setResetA] = useState(false);

    const [code, setCode] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [passed, setPassed] = useState(false);

    const [config, setConfig] = useState<PasswordConfig | null>(null);

    const [remaining, setRemaining] = useState("");

    const timerRef = useRef<number | null>(null);

    useEffect(() => {
        getPasswordConfig().then(setConfig);
    }, []);

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

    const email = useCallback(async () => {
        if (red) return
        if (remaining) return
        if (!id) return
        setWaiting(true)
        const response = await fetch("/api/forget_password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                user: Number(id),
                role: role,
            })
        })

        switch (response.status) {
            case 200:
                setResetA(true)
                setResetL(true)
                setTimeout(() => setReset(true), 200)
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                updateRemaining(Number(seconds))
                toast.error(await fromResponse(response))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }, [red, remaining, id, setWaiting, role, setResetL, updateRemaining])

    const change = useCallback(async () => {
        if (remaining) return
        if (!code) return
        if (!password) return
        if (!passed) return
        if (password !== confirmPassword) return
        setWaiting(true)
        const response = await fetch("/api/forget_password_confirm", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                password,
                confirm_password: confirmPassword,
                pending_password: code
            })
        })

        switch (response.status) {
            case 200:
                unsPwr()
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                updateRemaining(Number(seconds))
                toast.error(await fromResponse(response))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }, [remaining, code, password, passed, confirmPassword, setWaiting, unsPwr, updateRemaining])

    useEffect(() => {
        const handleKeyDown = async (e: KeyboardEvent) => {
            if (e.key === "Enter") {
                if (!reset) await email()
                else await change()
            }
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [reset, email, change])

    return (
        <>
            <div className={`content in outback ${!pwrA && "hid"} ${resetA ? "h-144" : "h-86"}`}>
                <h1 className="self-center text-5xl font-bold mb-8 rubik">{getKey("FORGOT_PASSWORD_TITLE")}</h1>
                {!reset && <div className={`loginput fredoka out ${resetA && "hid"}`}>
                    <select className="poppins" value={role} onChange={(e) => setRole(e.currentTarget.value)}>
                        <option value="guardian">{getKey("ROLE.GUARDIAN")}</option>
                        <option value="student">{getKey("ROLE.STUDENT")}</option>
                        <option value="teacher">{getKey("ROLE.TEACHER")}</option>
                    </select>
                    <input className={`${red && "wrong"}`} type="text" placeholder="User ID" value={id} onChange={(e) => {setId(e.currentTarget.value); checkUserID(e.currentTarget.value)}} />
                </div>}
                {reset && <div className="loginput fredoka in">
                    <input type="text" placeholder={getKey("CODE")} onChange={(e) => setCode(e.currentTarget.value)} />
                    <input type="password" placeholder={getKey("PASSWORD")} maxLength={config?.maxLength} onChange={(e) => setPassword(e.currentTarget.value)} />
                    <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder={getKey("CONFIRM_PASSWORD")} maxLength={config?.maxLength} onChange={(e) => setConfirmPassword(e.currentTarget.value)} />
                    <PasswordCheck password={password} confirmPassword={confirmPassword} setPassed={setPassed} />
                </div>}

                <div className="pwrbutton">
                    <Button onClick={unsPwr}>
                        <FontAwesomeIcon icon={faArrowLeft} /> {getKey("BACK")}
                    </Button>
                    <div className="flex flex-row items-center gap-4">
                        <p className="text-(--wrong-color) text-3xl fredoka">{remaining}</p>
                        {!reset && <Button onClick={email} className={`outbottom ${resetA && "hid"}`}>
                            {getKey("SEND_EMAIL")} <FontAwesomeIcon icon={faPaperPlane} />
                        </Button>}
                        {reset && <Button onClick={change} className="in">
                            {getKey("RESET")} <FontAwesomeIcon icon={faPlane} />
                        </Button>}
                    </div>
                </div>
            </div>
        </>
    )
}