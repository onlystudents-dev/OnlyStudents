import "./me.css";
import Navbar from "../../navbar/navbar.tsx";
import {isTimeFormat, type Me, type Preferences} from "../../types/api.ts";
import Button, {DropdownConfig} from "../../util/sidebar/config.tsx";
import {
    faAddressCard, faCalendarDays, faClock,
    faEnvelope, faFloppyDisk, faLanguage,
    faLock, faPaintRoller, faPaperPlane,
    faRightFromBracket, faRotateLeft, faTable,
    faUser,
    faUserLock
} from "@fortawesome/free-solid-svg-icons";
import {useRef, useState} from "react";
import Loading from "../../util/loading.tsx";
import {fromResponse, getKey, getLanguage, languages} from "../../util/language.ts";
import Config from "./config.tsx";
import {toast} from "react-toastify";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import PasswordCheck from "../../auth/pwc.tsx";
import Save from "../../util/save/save.tsx";
import {applyTheme, getAutoTheme, isTheme, themes} from "../../util/theme.ts";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import type {Dispatch, StateUpdater} from "preact/hooks";
import {getTimeFormat} from "../../util/time.ts";
import {postJSON} from "../../util/api.ts";

type Options = Pick<Preferences, "lang" | "time_format" | "theme" | "timetable_display" | "timetable_next">

export default function Me({ me, fetchMe, setMe, reload }: {me: Me, fetchMe: () => Promise<Me | null>, setMe: Dispatch<StateUpdater<Me | null>>, reload: () => void}) {
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

    const [options, setOptions] = useState<Options>({
        lang: me.preferences.lang,
        time_format: me.preferences.time_format,
        theme: me.preferences.theme,
        timetable_display: me.preferences.timetable_display,
        timetable_next: me.preferences.timetable_next,
    })

    const autolang = getLanguage(me, true)

    function setOption<K extends keyof Options>(key: K, value: Options[K]) {
        setOptions(current => ({ ...current, [key]: value }))
    }

    return (
        <>
            <Navbar me={me} />
            <Save options={options} setOptions={setOptions} save={save} />
            {loading && <Loading />}
            <div className="main">
                <div className="border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                    {(() => {
                        switch(active) {
                            case "appearance":
                                return (
                                    <>
                                        <DropdownConfig text={getKey("LANGUAGE")} icon={faLanguage} value={options.lang} onChange={(raw) => setOption("lang", raw)}>
                                            <option value="">{getKey("AUTOLANG", `${autolang?.emoji} ${autolang?.name}`)}</option>
                                            {languages.map((language) => (
                                                <option value={language.key}>{language.emoji} {language.name}</option>
                                            ))}
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("TIME_FORMAT")} icon={faClock} value={options.time_format} onChange={(raw) => isTimeFormat(raw) && setOption("time_format", raw)}>
                                            <option value="">{getKey("AUTO_TIME_FORMAT", getKey(getTimeFormat(getLanguage(me)) ? "24_TIME_FORMAT" : "12_TIME_FORMAT"))}</option>
                                            <option value="h12">{getKey("12_TIME_FORMAT")}</option>
                                            <option value="h23">{getKey("24_TIME_FORMAT")}</option>
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("THEME")} icon={faPaintRoller} value={options.theme} onChange={(raw) => isTheme(raw) && setOption("theme", raw)}>
                                            <option value="">{getKey("AUTOTHEME", getKey(`THEMES.${getAutoTheme()}`))}</option>
                                            {themes.map(theme => (
                                                theme && <option value={theme}>{getKey(`THEMES.${theme}`)}</option>
                                            ))}
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("TIMETABLE_DISPLAY")} icon={faTable} value={String(options.timetable_display)} onChange={(raw) => setOption("timetable_display", Number(raw))}>
                                            <option value="1">{getKey("TIMETABLE_DISPLAY.LESSONS")}</option>
                                            <option value="2">{getKey("TIMETABLE_DISPLAY.WEEKDAYS")}</option>
                                            <option value="3">{getKey("TIMETABLE_DISPLAY.ALL")}</option>
                                        </DropdownConfig>
                                        <DropdownConfig text={getKey("TIMETABLE_NEXT")} icon={faCalendarDays} value={String(options.timetable_next)} onChange={(raw) => setOption("timetable_next", raw === "true")}>
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

                                            <button disabled className="absolute bottom-7 right-7" onClick={updatePassword}>
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

    async function busy(task: () => Promise<void>) {
        setWaiting(true)
        try {
            await task()
        } finally {
            setWaiting(false)
        }
    }

    async function post(url: string, body?: unknown, success?: string) {
        const response = await postJSON(url, body)
        if (response.status !== 200) {
            toast.error(await fromResponse(response))
            return false
        }
        if (success) toast.success(getKey(success))
        return true
    }

    async function save(diff: Partial<Options>) {
        const preferences = { ...me.preferences, ...diff }
        if (!await post("/api/v1/me/update_preferences", preferences)) return false

        setMe({ ...me, preferences })
        if (diff.lang !== undefined || diff.time_format !== undefined) reload()
        if (diff.theme !== undefined) applyTheme(diff.theme)
        return true
    }

    async function sendEmail() {
        const value = newEmail.current?.value
        if (!value || value === me.email_address) return

        setSettingEmailA(true)
        const start = performance.now()
        await busy(async () => {
            if (!await post("/api/v1/me/change_email", { new_email: value })) return
            await resend(() => setTimeout(() => setSettingEmail(true), Math.max(500 - (performance.now() - start), 0)))
            await fetchMe()
            setEmail(value)
        })
    }

    async function resend(done: () => void = () => {}) {
        if (await post("/api/v1/me/verify_email")) done()
    }

    async function updateEmail() {
        const value = code.current?.value
        if (!value) return

        await busy(async () => {
            if (await post("/api/v1/me/verify_email_confirm", { code: value }, "EMAIL_CHANGED")) await fetchMe()
        })
    }

    async function updateNickname() {
        const value = newNickname.current?.value
        if (value === undefined) return

        const preferences = { ...me.preferences, nickname: value }
        await busy(async () => {
            if (!await post("/api/v1/me/update_preferences", preferences, "NICKNAME_CHANGED")) return
            setNickname(value)
            setMe({ ...me, preferences })
        })
    }

    async function updatePassword() {
        if (!currentPassword || !password || !confirmPassword || !passed || currentPassword === password) return

        await busy(async () => {
            await post("/api/v1/me/change_password", {
                current_password: currentPassword,
                new_password: password,
                confirm_new_password: confirmPassword,
            }, "PASSWORD_CHANGED")
        })
    }

    async function logout() {
        setLoading(true)
        const response = await postJSON("/api/v1/me/logout")

        if (response.status === 200) {
            location.reload()
        } else {
            toast.error(await response.text() || response.statusText)
        }
        setLoading(false)
    }
}
