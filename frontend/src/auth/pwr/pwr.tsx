import "./pwr.css";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowLeft, faPaperPlane, faPlane} from "@fortawesome/free-solid-svg-icons";
import React, {useCallback, useEffect, useState} from "react";
import {toast} from "react-toastify";
import PasswordCheck from "../pwc.tsx";
import {fromResponse, getKey} from "../../util/language.ts";
import {getPasswordConfig, type PasswordConfig} from "../passwordConfig.ts";
import RoleSelect from "../roleSelect.tsx";
import {useCountdown} from "../../util/countdown.ts";
import {getRetryAfter, postJSON} from "../../util/api.ts";
import {cx} from "../../util/cx.ts";

export default function PasswordReset({ pwrA, unsPwr, role, setRole, id, setId, red, checkUserID, setWaiting, setResetL }: {pwrA: boolean, unsPwr: () => void, role: string, setRole: React.Dispatch<React.SetStateAction<string>>, id: string, setId: React.Dispatch<React.SetStateAction<string>>, red: boolean, checkUserID: (id: string) => void, setWaiting: React.Dispatch<React.SetStateAction<boolean>>, setResetL:  React.Dispatch<React.SetStateAction<boolean>>}) {
    const [reset, setReset] = useState(false);
    const [resetA, setResetA] = useState(false);

    const [code, setCode] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [passed, setPassed] = useState(false);

    const [config, setConfig] = useState<PasswordConfig | null>(null);

    const [remaining, updateRemaining] = useCountdown();

    useEffect(() => {
        void getPasswordConfig().then(setConfig);
    }, []);

    const submit = useCallback(async (url: string, body: unknown, onOk: () => void) => {
        setWaiting(true)
        const response = await postJSON(url, body)

        if (response.status === 200) {
            onOk()
        } else {
            if (response.status === 429) updateRemaining(getRetryAfter(response))
            toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }, [setWaiting, updateRemaining])

    const email = useCallback(async () => {
        if (red || remaining || !id) return
        await submit("/api/forget_password", { user: Number(id), role }, () => {
            setResetA(true)
            setResetL(true)
            setTimeout(() => setReset(true), 200)
        })
    }, [red, remaining, id, role, setResetL, submit])

    const change = useCallback(async () => {
        if (remaining || !code || !password || !passed || password !== confirmPassword) return
        await submit("/api/forget_password_confirm", {
            password,
            confirm_password: confirmPassword,
            pending_password: code,
        }, unsPwr)
    }, [remaining, code, password, passed, confirmPassword, unsPwr, submit])

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Enter") void (reset ? change() : email())
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [reset, email, change])

    return (
        <>
            <div className={cx("content in outback", !pwrA && "hid", resetA ? "h-144" : "h-86")}>
                <h1 className="self-center text-5xl font-bold mb-8 rubik">{getKey("FORGOT_PASSWORD_TITLE")}</h1>
                {!reset && <div className={cx("loginput fredoka out", resetA && "hid")}>
<RoleSelect value={role} onChange={setRole} />
                    <input className={cx(red && "wrong")} type="text" placeholder="User ID" value={id} onChange={(e) => {setId(e.currentTarget.value); checkUserID(e.currentTarget.value)}} />
                </div>}
                {reset && <div className="loginput fredoka in">
                    <input type="text" placeholder={getKey("CODE")} onChange={(e) => setCode(e.currentTarget.value)} />
                    <input type="password" placeholder={getKey("PASSWORD")} maxLength={config?.maxLength} onChange={(e) => setPassword(e.currentTarget.value)} />
                    <input className={cx(password !== confirmPassword && "wrong")} type="password" placeholder={getKey("CONFIRM_PASSWORD")} maxLength={config?.maxLength} onChange={(e) => setConfirmPassword(e.currentTarget.value)} />
                    <PasswordCheck password={password} confirmPassword={confirmPassword} setPassed={setPassed} />
                </div>}

                <div className="pwrbutton">
                    <Button onClick={unsPwr}>
                        <FontAwesomeIcon icon={faArrowLeft} /> {getKey("BACK")}
                    </Button>
                    <div className="flex flex-row items-center gap-4">
                        <p className="text-(--wrong-color) text-3xl fredoka">{remaining}</p>
                        {!reset && <Button onClick={email} className={cx("outbottom", resetA && "hid")}>
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