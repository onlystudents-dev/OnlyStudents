import "../homework/homework.css";
import type {HomeworkSubmission, Me} from "../../types/api.ts";
import {useEffect, useState} from "preact/hooks";
import {fetchWithSchool} from "../../util/api.ts";
import Homework, {type PopupProps} from "../homework/homework.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
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

export default function TeacherHomework({ me }: { me: Me }) {
    const [homework, setHomework] = useState<HomeworkSubmission[] | null>(null)
    async function Fetch() {
        const res = await fetchWithSchool("/api/v1/me/teacher/homework")
        setHomework(await res.json() as HomeworkSubmission[])
    }
    useEffect(() => {
        void Fetch()
    }, [])
    return <Homework me={me} homework={homework} Popup={Popup} reloadSchoolStuff={Fetch} />
}

function Popup({ homework: hw, clear }: PopupProps) {
    const homework_submission = hw as HomeworkSubmission
    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="popup-content">
                    <div className="absolute top-4 left-6 right-16 flex flex-row items-center gap-3">
                        <div className="min-w-0 flex-1 text-3xl rubik bg-transparent resize-none outline-none p-2">
                            {homework_submission.title}
                        </div>
                    </div>
                    <div className="absolute top-18 left-8 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                        {homework_submission.subject}
                    </div>
                    <div className="absolute top-7.5 right-7.5 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>
                    {homework_submission.type == "note" && (
                        <div className="absolute top-30 left-6 right-6 h-fit min-h-24 fredoka bg-transparent resize-none outline-none p-2">
                            {homework_submission.note}
                        </div>
                    )}
                    {homework_submission.type == "file" && (
                        <div className="absolute top-30 left-6 right-6 h-fit min-h-24">
                            <div className="flex flex-row items-center gap-3 bg-white/5 border border-white/10 rounded-lg p-3">
                                <div className={`flex items-center justify-center w-12 h-12 rounded-lg text-2xl ${
                                    homework_submission.file_type === "document" ? "bg-blue-500/20 text-blue-400" :
                                    homework_submission.file_type === "pptx" ? "bg-orange-500/20 text-orange-400" :
                                    homework_submission.file_type === "pdf" ? "bg-red-500/20 text-red-400" :
                                    homework_submission.file_type === "image" ? "bg-green-500/20 text-green-400" :
                                    homework_submission.file_type === "workbook" ? "bg-purple-500/20 text-purple-400" :
                                    "bg-gray-500/20 text-gray-400"
                                }`}>
                                    <FontAwesomeIcon
                                        icon={
                                            homework_submission.file_type === "document" ? faFileLines :
                                            homework_submission.file_type === "pptx" ? faFilePowerpoint :
                                            homework_submission.file_type === "pdf" ? faFilePdf :
                                            homework_submission.file_type === "image" ? faFileImage :
                                            homework_submission.file_type === "workbook" ? faBookOpen :
                                            homework_submission.file_type === "archive" || homework_submission.file_type === "zip" ? faFileZipper :
                                            faFolder
                                        }
                                    />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="fredoka text-lg truncate">{homework_submission.file_name}</div>
                                    <div className="text-sm opacity-60">{(homework_submission.file_size / 1024).toFixed(1)} KB</div>
                                </div>
                                <a href={`/uploads/${homework_submission.file_uuid}`} className="action-button w-10 h-10 flex items-center justify-center rounded-lg bg-white/10 hover:bg-white/20 transition" download>
                                    <FontAwesomeIcon icon={faDownload} />
                                </a>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </>
    )
}
