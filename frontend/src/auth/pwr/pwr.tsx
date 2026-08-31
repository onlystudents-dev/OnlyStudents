import "./pwr.css";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowLeft, faPaperPlane, faPlane} from "@fortawesome/free-solid-svg-icons";
import React, {useCallback, useEffect, useState} from "react";
import {toast} from "react-toastify";

export default function PasswordReset({ pwrA, unsPwr, role, setRole, id, setId, red, checkUserID, setWaiting }: { pwrA: boolean, unsPwr: () => void, role: string, setRole: React.Dispatch<React.SetStateAction<string>>, id: string, setId: React.Dispatch<React.SetStateAction<string>>, red: boolean, checkUserID: (id: string) => void, setWaiting: React.Dispatch<React.SetStateAction<boolean>> }) {
    const [reset, setReset] = useState(false);
    const [resetA, setResetA] = useState(false);

    const [code, setCode] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const email = useCallback(async () => {
        if (red) return
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
                setTimeout(() => setReset(true), 200)
                break
            default:
                toast.error(await response.text() || response.statusText)
        }
        setWaiting(false)
    }, [red, id, role, setWaiting])

    const change = useCallback(async () => {
        if (!code) return
        if (!password) return
        if (password !== confirmPassword) return
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
            default:
                toast.error(await response.text() || response.statusText)
        }
        setWaiting(false)
    }, [code, password, confirmPassword, setWaiting, unsPwr])

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
            <div className={`content in outback ${!pwrA && "hid"}`}>
                <h1 className="self-center text-5xl font-bold mb-8 rubik">Reset password</h1>
                {!reset && <div className={`loginput fredoka out ${resetA && "hid"}`}>
                    <select value={role} onChange={(e) => setRole(e.target.value)}>
                        <option value="guardian">Guardian</option>
                        <option value="student">Student</option>
                        <option value="teacher">Teacher</option>
                    </select>
                    <input className={`${red && "wrong"}`} type="text" placeholder="User ID" value={id} onChange={(e) => {setId(e.target.value); checkUserID(e.target.value)}} />
                </div>}
                {reset && <div className="loginput fredoka in">
                    <input type="text" placeholder="Code" onChange={(e) => setCode(e.target.value)} />
                    <input type="password" placeholder="Password" onChange={(e) => setPassword(e.target.value)} />
                    <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder="Confirm password" onChange={(e) => setConfirmPassword(e.target.value)} />
                </div>}

                <div className={`pwrbutton ${reset ? "mt-4" : "mt-20"}`}>
                    <Button onClick={unsPwr}>
                        <FontAwesomeIcon icon={faArrowLeft} /> Back
                    </Button>
                    {!reset && <Button onClick={email} className={`outbottom ${resetA && "hid"}`}>
                        Send email <FontAwesomeIcon icon={faPaperPlane} />
                    </Button>}
                    {reset && <Button onClick={change} className="in">
                        Reset <FontAwesomeIcon icon={faPlane} />
                    </Button>}
                </div>
            </div>
        </>
    )
}