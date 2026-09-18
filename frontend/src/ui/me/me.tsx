import "./me.css";
import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";
import Button from "./button.tsx";
import {
    faAddressCard,
    faEnvelope, faFloppyDisk, faLanguage,
    faLock, faPaperPlane,
    faRightFromBracket, faRotateLeft,
    faUser,
    faUserLock
} from "@fortawesome/free-solid-svg-icons";
import {useRef, useState} from "react";
import Loading from "../../util/loading.tsx";
import {fromResponse, getKey, getLanguage, getLanguages} from "../../util/language.ts";
import Config, {DropdownConfig} from "./config.tsx";
import {toast} from "react-toastify";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import PasswordCheck from "../../auth/pwc.tsx";
import Save from "../../util/save/save.tsx";

export type option = "language"

export default function Me({ me, fetchMe }: {me: Me, fetchMe: () => Promise<void>}) {
    const ref = useRef<HTMLDivElement>(null)
    const sref = useRef<HTMLDivElement>(null)

    const [waiting, setWaiting] = useState(false)
    const [loading, setLoading] = useState(false)
    const [hidden, setHidden] = useState(true)

    const [active, setActive] = useState<"appearance" | "user" | "security">("appearance")

    const [cancelled, setCancelled] = useState(false)
    const [email, setEmail] = useState(me.email_address)
    const [settingEmail, setSettingEmail] = useState(false)
    const [settingEmailA, setSettingEmailA] = useState(false)
    const newEmail = useRef<HTMLInputElement>(null)
    const emailPassword = useRef<HTMLInputElement>(null)
    const code = useRef<HTMLInputElement>(null)

    const [nickname, setNickname] = useState(me.nickname)
    const newNickname = useRef<HTMLInputElement>(null)

    const [currentPassword, setCurrentPassword] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [passed, setPassed] = useState(false)

    const [options, setOptions] = useState<Record<option, unknown>>({
        language: localStorage.getItem("language") || ""
    })

    const autolang = getLanguage(true)

    return (
        <>
            <Navbar me={me} />
            <Save options={options} setOptions={setOptions} save={save} />
            {loading && <Loading />}
            <div className="main">
                <div className="configs border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                {(() => {
                   switch(active) {
                       case "appearance":
                           return (
                               <>
                                   <DropdownConfig text={getKey("LANGUAGE")} icon={faLanguage} value={options.language as string} options={options} setOptions={setOptions} lkey={"language"}>
                                       <option value="">{getKey("AUTOLANG", `${autolang?.emoji} ${autolang?.name}`)}</option>
                                       {getLanguages().map((language) => (
                                           <option value={language.key}>{language.emoji} {language.name}</option>
                                       ))}
                                   </DropdownConfig>
                               </>
                           )
                       case "user":
                           return (
                               <>
                                   <Config text={getKey("EMAIL")} icon={faEnvelope} value={email} className={`${!me.email_verified && "bg-(--wrong-color)"}`} key={"email"}>
                                       <h1 className="text-center text-5xl font-bold mb-8 rubik">{getKey("CHANGE_EMAIL")}</h1>
                                       <div className="moving">
                                           {settingEmail || (!cancelled && !me.email_verified) ? (<div className={`moving-content ${!(!cancelled && !me.email_verified) && "second"}`}>
                                               <input type="text" placeholder={getKey("CODE")} ref={code} key={"code"} />

                                               <button className="absolute bottom-0 left-0" onClick={() => {setSettingEmail(false); setSettingEmailA(false); setCancelled(true)}}>
                                                   <FontAwesomeIcon icon={faRotateLeft} /> {getKey("USE_DIFFERENT_EMAIL")}
                                               </button>
                                               <div className="absolute bottom-0 right-0 flex flex-row gap-2">
                                                   <button onClick={async () => await resend()}>
                                                       <FontAwesomeIcon icon={faPaperPlane} /> {getKey("RESEND")}
                                                   </button>
                                                   <button onClick={updateEmail}>
                                                       <FontAwesomeIcon icon={faFloppyDisk} /> {getKey("SAVE")}
                                                   </button>
                                               </div>
                                           </div>) : (<div className={`moving-content ${settingEmailA && "first"}`}>
                                               <input type="email" placeholder={getKey("NEW_EMAIL")} ref={newEmail} key={"email"} />
                                               <input type="password" placeholder={getKey("PASSWORD")} ref={emailPassword} key={"password"} />

                                               <button className="absolute bottom-0 right-0" onClick={sendEmail}>
                                                   <FontAwesomeIcon icon={faPaperPlane} /> {getKey("SEND_EMAIL")}
                                               </button>
                                           </div>)}
                                       </div>
                                   </Config>
                                   <Config text={getKey("NICKNAME")} icon={faAddressCard} value={nickname} key={"nickname"}>
                                       <h1 className="text-center text-5xl font-bold mb-8 rubik w-140">{getKey("CHANGE_NICKNAME")}</h1>

                                       <input type="text" placeholder={getKey("NICKNAME")} defaultValue={nickname} ref={newNickname} />

                                       <button className="absolute bottom-7 right-7" onClick={updateNickname}>
                                           <FontAwesomeIcon icon={faFloppyDisk} /> {getKey("SAVE")}
                                       </button>
                                   </Config>
                               </>
                           )
                       case "security":
                           return (
                               <>
                                   <Config text={getKey("PASSWORD")} icon={faUserLock} value={""} key={"password"}>
                                       <h1 className="text-center text-5xl font-bold mb-8 rubik w-140">{getKey("CHANGE_PASSWORD")}</h1>

                                       <div className="flex flex-col gap-4 w-full">
                                           <input type="password" placeholder={getKey("CURRENT_PASSWORD")} onChange={(e) => setCurrentPassword(e.target.value)} />
                                           <input type="password" placeholder={getKey("PASSWORD")} onChange={(e) => setPassword(e.target.value)} />
                                           <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder={getKey("CONFIRM_PASSWORD")} onChange={(e) => setConfirmPassword(e.target.value)} />
                                           <PasswordCheck password={password} confirmPassword={confirmPassword} setPassed={setPassed} />
                                       </div>

                                       <button className="absolute bottom-7 right-7" onClick={updatePassword}>
                                           <FontAwesomeIcon icon={faFloppyDisk} /> {getKey("SAVE")}
                                       </button>
                                   </Config>
                               </>
                           )
                       default:
                           return null
                   }
                })()}
                </div>
                {waiting && <Loading absolute />}
            </div>
            <div className="sidebar">
                <div className="flex flex-row gap-4">
                    <img className="rounded-[50%] h-20" src={me.pfp_url} alt={me.last_name} />
                    <div className="flex flex-col gap-1 justify-center min-w-0">
                        <h1 className="text-2xl truncate">
                            {nickname || `${me.first_name} ${me.last_name}`}
                        </h1>
                        <p>{getKey("SETTINGS")}</p>
                    </div>
                </div>
                <br />
                <div className="buttons" onMouseLeave={() => setHidden(true)}>
                    <div ref={ref} className={`sbutton ${hidden && "hid"}`}></div>
                    <div ref={sref} className={`sbutton ssbutton`}></div>
                    <Button onClick={button => {setActive("appearance"); sreposition(button)}} onMouseEnter={reposition} icon={faUser} text={getKey("APPEARANCE_SETTINGS")} />
                    <Button onClick={button => {setActive("user"); sreposition(button)}} onMouseEnter={reposition} icon={faUser} text={getKey("USER_SETTINGS")} />
                    <Button onClick={button => {setActive("security"); sreposition(button)}} onMouseEnter={reposition} icon={faLock} text={getKey("SECURITY")} />
                    <Button onClick={async (button) => {await logout(); sreposition(button)}} onMouseEnter={reposition} className="text-(--wrong-color)" icon={faRightFromBracket} text={getKey("LOGOUT")} />
                </div>
            </div>
        </>
    )

    function save() {
        localStorage.setItem("language", options.language as string)
        location.reload()
    }

    async function sendEmail() {
        if (!newEmail.current?.value) return
        if (!emailPassword.current?.value) return
        if (newEmail.current.value === me.email_address) return
        setSettingEmailA(true)
        const start = performance.now()
        setWaiting(true)
        const response = await fetch("/api/v1/me/change_email", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                password: emailPassword.current.value,
                new_email: newEmail.current.value,
            })
        })
        switch (response.status) {
            case 200:
                await resend(() => setTimeout(() => setSettingEmail(true), Math.max(500 - (performance.now() - start), 0)))
                await fetchMe()
                setEmail(newEmail.current.value)
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                toast.error(await fromResponse(response, seconds))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }

    async function resend(done: () => void = () => {}) {
        const response = await fetch("/api/v1/me/verify_email", {
            method: "POST",
        })
        switch (response.status) {
            case 200:
                done()
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                toast.error(await fromResponse(response, seconds))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
    }

    async function updateEmail() {
        if (!code.current?.value) return
        setWaiting(true)
        const response = await fetch("/api/v1/me/verify_email_confirm", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                code: code.current.value,
            })
        })
        switch (response.status) {
            case 200:
                toast.success(getKey("EMAIL_CHANGED"))
                await fetchMe()
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                toast.error(await fromResponse(response, seconds))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }

    async function updateNickname() {
        if (!newNickname.current?.value) return
        setWaiting(true)
        const response = await fetch("/api/v1/me/change_nickname", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                nickname: newNickname.current.value,
            })
        })
        switch (response.status) {
            case 200:
                toast.success(getKey("NICKNAME_CHANGED"))
                setNickname(newNickname.current?.value || nickname)
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                toast.error(await fromResponse(response, seconds))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }

    async function updatePassword() {
        if (!currentPassword) return
        if (!password) return
        if (currentPassword === password) return
        if (!confirmPassword) return
        if (!passed) return
        setWaiting(true)
        const response = await fetch("/api/v1/me/change_password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                current_password: currentPassword,
                new_password: password,
                confirm_new_password: confirmPassword,
            })
        })
        switch (response.status) {
            case 200:
                toast.success(getKey("PASSWORD_CHANGED"))
                break
            case 429: {
                let seconds = response.headers.get("Retry-After")
                if (!seconds) seconds = "-1"
                toast.error(await fromResponse(response, seconds))
                break
            }
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }

    async function logout() {
        setLoading(true)
        const response = await fetch("/api/v1/me/logout", {
            method: "POST",
        });

        switch (response.status) {
            case 200:
                location.reload();
                break;
            default:
                toast.error(await response.text() || response.statusText);
        }
        setLoading(false)
    }

    function reposition(button: HTMLButtonElement) {
        const current = ref.current;
        if (!current) return;
        setHidden(false)

        const buttonRect = button.getBoundingClientRect();
        const containerRect = current.parentElement?.getBoundingClientRect();

        if (!containerRect) return;

        current.style.top = `${buttonRect.top - containerRect.top}px`;
    }

    function sreposition(button: HTMLButtonElement) {
        const current = sref.current;
        if (!current) return;

        const buttonRect = button.getBoundingClientRect();
        const containerRect = current.parentElement?.getBoundingClientRect();

        if (!containerRect) return;

        current.style.top = `${buttonRect.top - containerRect.top}px`;
    }
}