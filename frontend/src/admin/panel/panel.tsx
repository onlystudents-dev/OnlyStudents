import "./panel.css";
import type {AdminStatusData, Log} from "../../types/admin.ts";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import {useEffect, useRef, useState} from "preact/compat";
import Button from "../../util/sidebar/button.tsx";
import {faFileLines, faHeartPulse} from "@fortawesome/free-solid-svg-icons";
import {getKey} from "../../util/language.ts";
import Loading from "../../util/loading.tsx";

export default function AdminPanel({ status }: {status: AdminStatusData}) {
    const [active, setActive] = useState<"status" | "logs">("status")

    const [logs, setLogs] = useState<Log[] | null>(null)
    const [activeLog, setActiveLog] = useState(logs?.[0]?.name || "")
    const logsRef = useRef<HTMLDivElement>(null)

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
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_DB_PING")}</p>
                                        <p>{status.db_ping}ms</p>
                                    </div>
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_REDIS_PING")}</p>
                                        <p>{status.redis_ping.replace(/\D/g, "")}ms</p>
                                    </div>
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_SCHOOL_COUNT")}</p>
                                        <p>{status.school_count}</p>
                                    </div>
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_ACCOUNT_COUNT")}</p>
                                        <p>{status.account_count}</p>
                                    </div>
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_GUARDIAN_COUNT")}</p>
                                        <p>{status.guardian_count}</p>
                                    </div>
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_STUDENT_COUNT")}</p>
                                        <p>{status.student_count}</p>
                                    </div>
                                    <div className="status">
                                        <p>{getKey("ADMIN_STATUS_TEACHER_COUNT")}</p>
                                        <p>{status.teacher_count}</p>
                                    </div>
                                </div>
                            )

                        case "logs":
                            return (
                                <>
                                    {!logs ? <Loading absolute /> : <>
                                        <div className="logselect">
                                            <select className="poppins" value={activeLog} onChange={(e) => setActiveLog(e.currentTarget.value)}>
                                                {logs.map((log: Log) => (
                                                    <option key={log.name} value={log.name}>{log.name}</option>
                                                ))}
                                            </select>
                                        </div>
                                        <div className="box">
                                            <div ref={logsRef} className="w-full h-84 overflow-y-auto font-mono text-sm whitespace-pre" style={{height: "calc(100vh - 124px)"}}>
                                                <div key={logs.find(log => log.name === activeLog)?.name}>
                                                    {logs.find(log => log.name === activeLog)?.content}
                                                </div>
                                            </div>
                                        </div>
                                    </>}
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
            </Sidebar>
        </>
    )

    async function fetchLogs() {
        const response = await fetch("/api/v1/admin/logs");
        const logs = await response.json() as Log[]

        setLogs(logs)
        if (!activeLog) setActiveLog(logs[0]?.name || "")
    }
}