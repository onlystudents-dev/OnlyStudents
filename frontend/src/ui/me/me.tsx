import "./me.css";
import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../types/api.ts";
import Button from "./button.tsx";
import {
    faAddressCard, faCalendarDays,
    faEnvelope, faFloppyDisk, faLanguage,
    faLock, faPaintRoller, faPaperPlane,
    faRightFromBracket, faRotateLeft, faTable,
    faUser,
    faUserLock
} from "@fortawesome/free-solid-svg-icons";
import {useRef, useState} from "react";
import Loading from "../../util/loading.tsx";
import {fromResponse, getKey, getLanguage, languages} from "../../util/language.ts";
import Config, {DropdownConfig} from "./config.tsx";
import {toast} from "react-toastify";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import PasswordCheck from "../../auth/pwc.tsx";
import Save from "../../util/save/save.tsx";
import {applyTheme, getAutoTheme, type Theme, themes} from "../../util/theme.ts";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import type {Dispatch, StateUpdater} from "preact/hooks";

export type option = "lang" | "theme" | "timetable_display" | "timetable_next"

export default function Me({ me, fetchMe, setMe, reload }: {me: Me, fetchMe: () => Promise<void>, setMe: Dispatch<StateUpdater<Me | null>>, reload: () => void}) {
    const [waiting, setWaiting] = useState(false)
    const [loading, setLoading] = useState(false)

    const [active, setActive] = useState<"appearance" | "user" | "security">("appearance")

    const [cancelled, setCancelled] = useState(false)
    const [email, setEmail] = useState(me.email_address)
    const [settingEmail, setSettingEmail] = useState(false)
    const [settingEmailA, setSettingEmailA] = useState(false)
    const newEmail = useRef<HTMLInputElement>(null)
    const code = useRef<HTMLInputElement>(null)

    const [nickname, setNickname] = useState(me.preferences.nickname)
    const newNickname = useRef<HTMLInputElement>(null)

    const [currentPassword, setCurrentPassword] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [passed, setPassed] = useState(false)

    const [options, setOptions] = useState<Record<option, unknown>>({
        lang: me.preferences.lang,
        theme: me.preferences.theme,
        timetable_display: me.preferences.timetable_display,
        timetable_next: me.preferences.timetable_next,
    })

    const autolang = getLanguage(me, true)

    const methods: Record<option, (value: unknown) => void> = {
        lang: reload,
        theme: (value) => applyTheme(value as Theme),
        timetable_display: () => {},
        timetable_next: () => {},
    }

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
                                        <DropdownConfig text={getKey("LANGUAGE")} icon={faLanguage} value={options.lang as string} options={options} setOptions={setOptions} lkey={"lang"}>
                                            <option value="">{getKey("AUTOLANG", `${autolang?.emoji} ${autolang?.name}`)}</option>
                                            {languages.map((language) => (
                                                <option value={language.key}>{language.emoji} {language.name}</option>
                                            ))}
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("THEME")} icon={faPaintRoller} value={options.theme as string} options={options} setOptions={setOptions} lkey={"theme"}>
                                            <option value="">{getKey("AUTOTHEME", getKey(`THEMES.${getAutoTheme()}`))}</option>
                                            {themes.map(theme => (
                                                theme && <option value={theme}>{getKey(`THEMES.${theme}`)}</option>
                                            ))}
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("TIMETABLE_DISPLAY")} icon={faTable} value={options.timetable_display as string} options={options} setOptions={setOptions} lkey={"timetable_display"}>
                                            <option value="1">{getKey("TIMETABLE_DISPLAY.LESSONS")}</option>
                                            <option value="2">{getKey("TIMETABLE_DISPLAY.WEEKDAYS")}</option>
                                            <option value="3">{getKey("TIMETABLE_DISPLAY.ALL")}</option>
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("TIMETABLE_NEXT")} icon={faCalendarDays} value={options.timetable_next as string} options={options} setOptions={setOptions} lkey={"timetable_next"}>
                                            <option value="true">{getKey("TIMETABLE_NEXT.NEXT")}</option>
                                            <option value="false">{getKey("TIMETABLE_NEXT.THIS")}</option>
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
                                                <input type="password" placeholder={getKey("CURRENT_PASSWORD")} onChange={(e) => setCurrentPassword(e.currentTarget.value)} />
                                                <input type="password" placeholder={getKey("PASSWORD")} onChange={(e) => setPassword(e.currentTarget.value)} />
                                                <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder={getKey("CONFIRM_PASSWORD")} onChange={(e) => setConfirmPassword(e.currentTarget.value)} />
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
            <Sidebar me={me}>
                <Button onClick={() => setActive("appearance")} icon={faUser} text={getKey("APPEARANCE_SETTINGS")} />
                <Button onClick={() => setActive("user")} icon={faUser} text={getKey("USER_SETTINGS")} />
                <Button onClick={() => setActive("security")} icon={faLock} text={getKey("SECURITY")} />
                <Button onClick={logout} className="text-(--wrong-color)" icon={faRightFromBracket} text={getKey("LOGOUT")} />
            </Sidebar>
        </>
    )

    async function save(diff: Record<string, unknown>) {
        const new_preferences = {
            ...me.preferences,
            ...diff,

            ...(diff.timetable_display !== undefined && {
                timetable_display: Number(diff.timetable_display),
            }),

            ...(diff.timetable_next !== undefined && {
                timetable_next: diff.timetable_next === "true",
            }),
        }
        const response = await fetch("/api/v1/me/update_preferences", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ new_preferences }),
        })
        if (!response.ok) {
            toast.error(await fromResponse(response))
            return false
        }
        setMe({ ...me, preferences: new_preferences })
        Object.entries(diff).forEach(([key, value]) => {
            methods[key as option]?.(value)
        })
        return true
    }

    async function sendEmail() {
        if (!newEmail.current?.value) return
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
                new_email: newEmail.current.value,
            })
        })
        switch (response.status) {
            case 200:
                await resend(() => setTimeout(() => setSettingEmail(true), Math.max(500 - (performance.now() - start), 0)))
                await fetchMe()
                setEmail(newEmail.current.value)
                break
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
            default:
                toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }

    async function updateNickname() {
        if (!newNickname.current?.value && newNickname.current?.value !== "") return
        setWaiting(true)
        const response = await fetch("/api/v1/me/update_preferences", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                new_preferences: {
                    ...me.preferences,
                    nickname: newNickname.current.value,
                }
            })
        })
        switch (response.status) {
            case 200:
                toast.success(getKey("NICKNAME_CHANGED"))
                setNickname(newNickname.current.value)
                setMe({ ...me, preferences: { ...me.preferences, nickname: newNickname.current.value } })
                break
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
}
