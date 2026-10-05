import "./panel.css";
import {type AdminStatusData, enrollKeys, type Log, type School} from "../../types/admin.ts";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import {useEffect, useRef, useState} from "preact/compat";
import Button, {DropdownConfig, InputConfig} from "../../util/sidebar/config.tsx";
import {faFileLines, faHeartPulse, faUserPlus} from "@fortawesome/free-solid-svg-icons";
import {getKey} from "../../util/language.ts";
import Loading from "../../util/loading.tsx";

type enrollable = "guardian" | "student" | "teacher"

export default function AdminPanel({ status }: {status: AdminStatusData}) {
    const [active, setActive] = useState<"status" | "logs" | "enroll" | "themes">("status")

    const [logs, setLogs] = useState<Log[] | null>(null)
    const [activeLog, setActiveLog] = useState(logs?.[0]?.name || "")
    const logsRef = useRef<HTMLDivElement>(null)

    const [enroll, setEnroll] = useState<enrollable>("guardian")
    const [options, setOptions] = useState<Record<string, string | number>>({})
    const [schools, setSchools] = useState<School[]>([])

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
                                    <DropdownConfig value={enroll} onChange={e => setEnroll(e as enrollable)} className="sselect">
                                        <option value="guardian">{getKey("ROLE.GUARDIAN")}</option>
                                        <option value="student">{getKey("ROLE.STUDENT")}</option>
                                        <option value="teacher">{getKey("ROLE.TEACHER")}</option>
                                    </DropdownConfig>
                                    <div className="box pr-2!">
                                        <div className="w-full overflow-y-auto pr-2" style={{height: "calc(100vh - 124px)"}}>
                                            {Object.entries(enrollKeys[enroll]).filter((([, type]) => type !== "boolean")).map(([key, type]) => {
                                                switch(key) {
                                                    case "school_id":
                                                        return (
                                                            <DropdownConfig value={schools[0]?.name || ""} onChange={e => setOptions({ ...options, [key]: Number(e) })} text={getKey(`ADMIN_ENROLL.${enroll.toUpperCase()}.${key.toUpperCase()}`)} className="enroll rubik">
                                                                {schools.map(school => (
                                                                    <option key={school.id} value={school.id}>{school.name}</option>
                                                                ))}
                                                            </DropdownConfig>
                                                        )
                                                    default:
                                                        return <InputConfig value={options[key] || ""} onChange={value => setOptions({ ...options, [key]: value })} key={key} text={getKey(`ADMIN_ENROLL.${enroll.toUpperCase()}.${key.toUpperCase()}`)} type={type} className="enroll rubik" />
                                                }
                                            })}
                                        </div>
                                    </div>
                                </>
                            )
                        case "themes":
                            return (
                                <>
                                    <div className="box flex-wrap">

                                    </div>
                                </>
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
}