import "./panel.css";
import type {AdminStatusData} from "../../types/admin.ts";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import {useState} from "preact/compat";
import Button from "../../util/sidebar/button.tsx";
import {faHeartPulse} from "@fortawesome/free-solid-svg-icons";
import {getKey} from "../../util/language.ts";

export default function AdminPanel({ status }: {status: AdminStatusData}) {
    const [active, setActive] = useState<"status">("status")

    return (
        <>
            <div className="main">
                <div className="border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                    {(() => {
                        switch(active) {
                            // i hate these new checks
                            // i when i add more i dont want to add this again so its staying
                            // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
                            case "status":
                                return (
                                    <>
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
                                    </>
                                )
                            default:
                                return null
                        }
                    })()}
                </div>
            </div>
            <Sidebar>
                <Button icon={faHeartPulse} text={getKey("ADMIN_STATUS")} onClick={() => setActive("status")} />
            </Sidebar>
        </>
    )
}