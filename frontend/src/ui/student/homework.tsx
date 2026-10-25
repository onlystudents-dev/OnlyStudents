import "../homework/homework.css";
import type {HomeworkSummary, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {useEffect} from "react";
import {fetchWithSchool} from "../../util/api.ts";
import Homework, {type PopupProps} from "../homework/homework.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";
import Attachment from "../attachment.tsx";
import { getKey } from "../../util/language.ts";

export default function StudentHomework({ me }: { me: Me }) {
    const [homework, setHomework] = useState<HomeworkSummary[] | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/student/homework").then(async res => setHomework(await res.json() as HomeworkSummary[]))
    }

    useEffect(() => {
        void Fetch()
    }, [])

    return <Homework me={me} homework={homework} Popup={Popup} reloadSchoolStuff={Fetch} />
}

function Popup({ homework: hw, clear }: PopupProps) {
    const homework = hw as HomeworkSummary

    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="popup-content">
                    <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>
                    <div className="flex flex-1 flex-col gap-3 min-h-0 overflow-y-auto p-6 pr-14">
                        <div className="text-3xl rubik truncate">{homework.title}</div>
                        <div className="text-2xl opacity-70 rubik truncate">{homework.subject}</div>
                        <div className="fredoka whitespace-pre-wrap wrap-break-word">{homework.description}</div>
                        <Attachment label={getKey("SUBMISSION")} attachment={
                      {
                        note: homework.submission_note,
                        file_type: homework.submission_file_type,
                        file_size: homework.submission_file_size,
                        file_name: homework.submission_file_name,
                        file_uuid: homework.submission_file_uuid
                      }
                    } />
                    </div>
                </div>
            </div>
        </>
    )
}
