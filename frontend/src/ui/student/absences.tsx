import "./absences.css";
import type {Me, Absence} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {fetchInto} from "../../util/api.ts";
import {useEffect} from "react";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {formatUnixDate} from "../../util/time.ts";
import {getLanguage} from "../../util/language.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faXmark,
    faFileLines,
    faFilePowerpoint,
    faFilePdf,
    faFileImage,
    faBookOpen,
    faFolder,
    faDownload,
    faFileZipper,
} from "@fortawesome/free-solid-svg-icons";

function AbsenceFileCard({ absence }: { absence: Absence }) {
    return (
        <div className="flex flex-row items-center gap-3 bg-white/5 border border-white/10 rounded-lg p-3">
            <div className={`flex items-center justify-center w-12 h-12 rounded-lg text-2xl ${
                absence.file_type === "document" ? "bg-blue-500/20 text-blue-400" :
                absence.file_type === "pptx" ? "bg-orange-500/20 text-orange-400" :
                absence.file_type === "pdf" ? "bg-red-500/20 text-red-400" :
                absence.file_type === "image" ? "bg-green-500/20 text-green-400" :
                absence.file_type === "workbook" ? "bg-purple-500/20 text-purple-400" :
                "bg-gray-500/20 text-gray-400"
            }`}>
                <FontAwesomeIcon
                    icon={
                        absence.file_type === "document" ? faFileLines :
                        absence.file_type === "pptx" ? faFilePowerpoint :
                        absence.file_type === "pdf" ? faFilePdf :
                        absence.file_type === "image" ? faFileImage :
                        absence.file_type === "workbook" ? faBookOpen :
                        absence.file_type === "archive" || absence.file_type === "zip" ? faFileZipper :
                        faFolder
                    }
                />
            </div>
            <div className="flex-1 min-w-0">
                <div className="fredoka text-lg truncate">{absence.file_name}</div>
                <div className="text-sm opacity-60">{(absence.file_size / 1024).toFixed(1)} KB</div>
            </div>
            <a
                href={`/uploads/${absence.file_uuid}`}
                className="action-button w-10 h-10 flex items-center justify-center rounded-lg bg-white/10 hover:bg-white/20 transition"
                download
            >
                <FontAwesomeIcon icon={faDownload} />
            </a>
        </div>
    )
}

export default function StudentAbsences({ me }: { me: Me }) {
    const [absences, setAbsences] = useState<Absence[] | null>(null)
    const [current, setCurrent] = useState<Absence | null>(null)

    async function Fetch() {
        const list = await fetchInto<Absence[]>("/api/v1/me/student/absences")
        if (list) setAbsences(list)
    }

    useEffect(() => void Fetch(), [])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && <>
                    <Overlay onClick={() => setCurrent(null)} />
                    <div className={`popup ${!current.justified && "wrong"}`}>
                        <div className="popup-content">
                            <div className="absolute top-6 left-6 right-16 flex flex-row gap-3">
                                <div className="text-3xl rubik flex-1 truncate">
                                    {current.subject}
                                </div>
                                <div className="shrink min-w-0 max-w-[45%]">
                                    {formatUnixDate(current.date, getLanguage(me)?.key || "en-US")}
                                </div>
                            </div>
                            <div className="absolute top-15 left-6 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                                {current.has_teacher_first_name ? `${current.teacher_first_name} ${current.teacher_last_name}` : current.teacher_last_name}
                            </div>
                            <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={() => setCurrent(null)}>
                                <FontAwesomeIcon icon={faXmark} />
                            </div>
                            {current.reason_type === "note" && (
                                <div className="absolute top-30 left-6 right-6 h-fit min-h-24 fredoka bg-transparent resize-none outline-none p-2">
                                    {current.reason_note}
                                </div>
                            )}
                            {current.reason_type === "file" && (
                                <div className="absolute top-30 left-6 right-6 h-fit min-h-24">
                                    <AbsenceFileCard absence={current} />
                                </div>
                            )}
                        </div>
                    </div>
                </>}
                {absences?.map((absence) => (
                    <div key={absence.id} className={`absence ${!absence.justified && "wrong"}`} onClick={() => setCurrent(absence)}>
                        <div className="flex flex-col items-start justify-between">
                            <p className="max-w-100 truncate">
                                {absence.reason_type === "note" ? absence.reason_note :
                                 absence.reason_type === "file" ? absence.file_name :
                                 absence.subject}
                            </p>
                            <p className="max-w-100 truncate opacity-70">
                                {absence.reason_type ? absence.subject :
                                 `${absence.teacher_first_name} ${absence.teacher_last_name}`}
                            </p>
                        </div>
                        <div className="flex flex-col items-start justify-between">

                        </div>
                    </div>
                )) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                ))}
            </div>
        </>
    )
}
