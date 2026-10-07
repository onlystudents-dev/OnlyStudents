import "./panel.css";
import {type AdminStatusData, enrollKeys, type Log, type School} from "../../types/admin.ts";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import {useEffect, useRef, useState} from "preact/compat";
import Button, {DropdownConfig, InputConfig} from "../../util/sidebar/config.tsx";
import {
    faArrowLeft,
    faFileLines, faFloppyDisk,
    faHeartPulse,
    faPaintRoller,
    faPenToSquare,
    faPlus,
    faTrash,
    faUserPlus
} from "@fortawesome/free-solid-svg-icons";
import {fromResponse, getKey} from "../../util/language.ts";
import Loading from "../../util/loading.tsx";
import {toast} from "react-toastify";
import Save from "../../util/save/save.tsx";
import type {Theme} from "../../types/api.ts";
import {postJSON} from "../../util/api.ts";

type enrollable = "guardian" | "student" | "teacher"

function rootColors(): Record<string, string> {
    const style = document.getElementById("root-theme")
    if (!style) return {}

    const root = /:root\s*\{([^}]+)\}/.exec(style.textContent)
    const body = root?.[1]
    if (body === undefined) return {}

    const colors: Record<string, string> = {}
    for (const [, key, value] of body.matchAll(/--([a-zA-Z0-9_-]+)\s*:\s*([^;]+);?/g)) {
        if (key === undefined || value === undefined) continue
        colors[`--${key}`] = value.trim()
    }

    if (Object.keys(colors).length !== 0) return colors

    return {}
}

function colorOf(theme: Theme, key: string): string {
    return theme.colors[key] ?? theme.colors[key.slice(2)] ?? ""
}

export default function AdminPanel({ status }: {status: AdminStatusData}) {
    const [active, setActive] = useState<"status" | "logs" | "enroll" | "themes">("status")

    const [logs, setLogs] = useState<Log[] | null>(null)
    const [activeLog, setActiveLog] = useState(logs?.[0]?.name || "")
    const logsRef = useRef<HTMLDivElement>(null)

    const [enroll, setEnroll] = useState<enrollable>("guardian")
    const [options, setOptions] = useState<Record<string, string | number>>({})
    const [schools, setSchools] = useState<School[]>([])
    const [focused, setFocused] = useState("")

    const [themes, setThemes] = useState<Theme[] | null>(null)
    const [editing, setEditing] = useState<Theme | null>(null)

    useEffect(() => {
        if (active === "logs" && logsRef.current) {
            logsRef.current.scrollTop = logsRef.current.scrollHeight
        }
    }, [active])

    return (
        <>
            <div className="main">
                {(() => {
                    switch(active) {
                        case "status":
                            return (
                                <div className="box">
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_DB_PING")}</p>
                                        <p>{status.db_ping}ms</p>
                                    </div>
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_REDIS_PING")}</p>
                                        <p>{status.redis_ping.replace(/\D/g, "")}ms</p>
                                    </div>
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_SCHOOL_COUNT")}</p>
                                        <p>{status.school_count}</p>
                                    </div>
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_ACCOUNT_COUNT")}</p>
                                        <p>{status.account_count}</p>
                                    </div>
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_GUARDIAN_COUNT")}</p>
                                        <p>{status.guardian_count}</p>
                                    </div>
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_STUDENT_COUNT")}</p>
                                        <p>{status.student_count}</p>
                                    </div>
                                    <div className="config status">
                                        <p>{getKey("ADMIN_STATUS_TEACHER_COUNT")}</p>
                                        <p>{status.teacher_count}</p>
                                    </div>
                                </div>
                            )

                        case "logs":
                            return (
                                <>
                                    {!logs ? <Loading /> : <>
                                        <DropdownConfig value={activeLog} onChange={e => setActiveLog(e)} className="sselect">
                                            {logs.map((log: Log) => (
                                                <option key={log.name} value={log.name}>{log.name}</option>
                                            ))}
                                        </DropdownConfig>
                                        <div className="box">
                                            <div ref={logsRef} className="w-full overflow-y-auto font-mono text-sm whitespace-pre" style={{height: "calc(100vh - 124px)"}}>
                                                <div key={logs.find(log => log.name === activeLog)?.name}>
                                                    {logs.find(log => log.name === activeLog)?.content}
                                                </div>
                                            </div>
                                        </div>
                                    </>}
                                </>
                            )
                        case "enroll":
                            return (
                                <>
                                    <DropdownConfig value={enroll} onChange={e => { setEnroll(e as enrollable); setOptions({}) }} className="sselect">
                                        <option value="guardian">{getKey("ROLE.GUARDIAN")}</option>
                                        <option value="student">{getKey("ROLE.STUDENT")}</option>
                                        <option value="teacher">{getKey("ROLE.TEACHER")}</option>
                                    </DropdownConfig>
                                    <div className="box pr-2!">
                                        <div className="w-full overflow-y-auto pr-2 flex flex-col" style={{height: "calc(100vh - 164px)"}}>
                                            {Object.entries(enrollKeys[enroll]).filter((([, type]) => type !== "boolean")).map(([key, type]) => {
                                                switch(key) {
                                                    case "school_id":
                                                        return (
                                                            <DropdownConfig value={String(options[key]) || schools[0]?.name || ""} onChange={e => setOptions({ ...options, [key]: Number(e) })} text={getKey(`ADMIN_ENROLL.${enroll.toUpperCase()}.${key.toUpperCase()}`)} className="enroll rubik">
                                                                {schools.map(school => (
                                                                    <option key={school.id} value={school.id}>{school.name}</option>
                                                                ))}
                                                            </DropdownConfig>
                                                        )
                                                    default:
                                                        // i have a 4k monitor so it fits for me
                                                        // well not anymore
                                                        return <InputConfig value={options[key] || ""} onChange={value => setOptions({ ...options, [key]: value })} onFocusIn={() => setFocused(key)} onFocusOut={() => focused === key && setFocused("")} key={key} text={getKey(`ADMIN_ENROLL.${enroll.toUpperCase()}.${key.toUpperCase()}`)}
                                                                            type={type} className={`enroll rubik ${focused === key && "bg-(--hover-color)"}`} />
                                                }
                                            })}
                                        </div>
                                        <Button className="enroll-save" icon={faFloppyDisk} text={getKey("SAVE")} onClick={() => void enrollUser()} />
                                    </div>
                                </>
                            )
                        case "themes":
                            if (editing) {
                                const colors = rootColors()
                                return (
                                    <>
                                        <div className="box">
                                            <Button className="self-start theme-action-button hover:text-(--warning-color) ml-1 mb-1" icon={faArrowLeft} text={getKey("ADMIN_THEMES_BACK")} onClick={() => setEditing(null)} />
                                            <div className="w-full overflow-y-auto pr-2 flex flex-col gap-2" style={{height: "calc(100vh - 164px)"}}>
                                                <InputConfig text={getKey("ADMIN_THEMES_NAME")} value={editing.name} onChange={value => setEditing({...editing, name: value})} className="enroll rubik" />
                                                {Object.entries(colors).map(([key]) => (
                                                    <div className="config theme-color rubik" key={key}>
                                                        <p>{key}</p>
                                                        <div className="flex flex-row items-center gap-2">
                                                            <input type="color" value={editing.colors[key]} onChange={e => setEditing({...editing, colors: {...editing.colors, [key]: e.currentTarget.value}})} />
                                                            <input value={editing.colors[key]} onChange={e => setEditing({...editing, colors: {...editing.colors, [key]: e.currentTarget.value}})} />
                                                        </div>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                        <Save options={editing} setOptions={setEditing} save={saveTheme} />
                                    </>
                                )
                            }

                            return (
                                <div className="box flex-wrap">
                                    <Button className="self-start theme-action-button hover:text-(--right-color) ml-1 mb-1" icon={faPlus} text={getKey("ADMIN_THEMES_NEW")} onClick={() => openEditor()} />
                                    {!themes ? <Loading /> : themes.length === 0 ? <p className="poppins">{getKey("ADMIN_THEMES_EMPTY")}</p> : themes.map(theme => (
                                        <div className="theme-card theme-action-shadow" key={theme.name}>
                                            <h1 className="rubik">{getKey(`THEMES.${theme.name}`) || theme.name}</h1>
                                            <div className="swatches">
                                                {Object.keys(rootColors()).map(key => (
                                                    <span key={key} title={key} style={{backgroundColor: colorOf(theme, key) || "transparent"}} />
                                                ))}
                                            </div>
                                            <div className="flex flex-row gap-2">
                                                <Button className="theme-action-button hover:text-(--edit-color)" icon={faPenToSquare} text={getKey("ADMIN_THEMES_EDIT")} onClick={() => openEditor(theme)} />
                                                <Button className="theme-action-button hover:text-(--wrong-color)" icon={faTrash} text={getKey("ADMIN_THEMES_DELETE")} onClick={() => void deleteTheme(theme.name, true)} />
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )
                        default:
                            return null
                    }
                })()}
            </div>

            <Sidebar>
                <Button icon={faHeartPulse} text={getKey("ADMIN_STATUS")} onClick={() => setActive("status")} />
                <Button icon={faFileLines} text={getKey("ADMIN_LOGS")} onClick={() => { setActive("logs"); void fetchLogs() }} />
                <Button icon={faUserPlus} text={getKey("ADMIN_ENROLL")} onClick={() => { setActive("enroll"); void fetchSchools() }} />
                <Button icon={faPaintRoller} text={getKey("ADMIN_THEMES")} onClick={() => { setActive("themes"); void fetchThemes() }} />
        </Sidebar>
        </>
    )

    async function fetchLogs() {
        const response = await fetch("/api/v1/admin/logs")
        const logs = await response.json() as Log[]

        setLogs(logs)
        if (!activeLog) setActiveLog(logs[0]?.name || "")
    }

    async function fetchSchools() {
        if (schools.length !== 0) return
        const response = await fetch("/api/v1/admin/schools")
        const schoolsJ = await response.json() as School[]

        setSchools(schoolsJ)
    }

    async function fetchThemes() {
        setEditing(null)

        const response = await fetch("/api/v1/themes")

        if (!response.ok) {
            void toast.error(await fromResponse(response))
            return
        }

        setThemes(await response.json() as Theme[])

        const element = document.getElementById("themes")
        if (element) {
            const themes = await fetch("/dynamic/themes.css")
            element.textContent = await themes.text()
        }
    }

    async function enrollUser() {
        const ops: Record<string, string | number | boolean> = {}

        for (const [key, type] of Object.entries(enrollKeys[enroll])) {
            if (type === "boolean") {
                const valueKey = key.startsWith("has_")
                    ? key.slice(4)
                    : key

                ops[key] = !!options[valueKey]
                continue
            }

            const value = options[key]

            if (value === undefined || value === "") {
                continue
            }

            switch (type) {
                case "number":
                    ops[key] = Number(value)
                    break

                case "date":
                    ops[key] = Math.floor(
                        new Date(String(value)).getTime() / 1000
                    )
                    break

                default:
                    ops[key] = String(value)
                    break
            }
        }

        const response = await postJSON(`/api/v1/admin/enroll/${enroll}`, ops)

        switch (response.status) {
            case 200:
                toast.success(getKey("ADMIN_SUCCESSFULLY_ENROLLED"))
                break
            default:
                toast.error(await fromResponse(response))
        }
    }

    function openEditor(theme?: Theme) {
        const defaults = rootColors()

        console.log("defaults", defaults)
        console.log("theme", theme)

        setEditing(theme
            ? {name: theme.name, colors: Object.fromEntries(Object.entries(defaults).map(([key, value]) => [key, colorOf(theme, key) || value]))}
            : {name: "", colors: {...defaults}})
    }

    async function saveTheme(diff: Partial<Theme>, defaults: Partial<Theme>): Promise<boolean> {
        if (!editing) return false

        if (!editing.name.trim()) {
            toast.error(getKey("ADMIN_THEMES_BAD_NAME"))
            return false
        }

        const missing = Object.keys(rootColors()).filter(key => !(editing.colors[key] ?? "").trim())
        if (missing.length !== 0) {
            toast.error(getKey("ADMIN_THEMES_MISSING", missing.join(", ")))
            return false
        }

        if (diff.name !== undefined && defaults.name !== undefined) {
            await deleteTheme(defaults.name)
        }

        const response = await fetch("/api/v1/admin/theme", {
            method: "PUT",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(editing),
        })

        if (!response.ok) {
            void toast.error(await fromResponse(response))
            return false
        }

        toast.success(getKey("ADMIN_THEMES_SAVED"))
        setEditing(null)
        void fetchThemes()
        return true
    }

    async function deleteTheme(name: string, stuff?: boolean) {
        if (stuff && !window.confirm(getKey("ADMIN_THEMES_DELETE_CONFIRM", name))) return

        const response = await fetch("/api/v1/admin/theme", {method: "DELETE", body: name})

        if (!response.ok) {
            void toast.error(await fromResponse(response))
            return
        }

        if (stuff) toast.success(getKey("ADMIN_THEMES_DELETED"))
        void fetchThemes()
    }
}
