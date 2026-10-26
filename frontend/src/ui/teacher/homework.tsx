import "../homework/homework.css";
import type {HomeworkSubmission, Me} from "../../types/api.ts";
import {useEffect, useState} from "preact/hooks";
import {fetchWithSchool} from "../../util/api.ts";
import Homework, {type PopupProps} from "../homework/homework.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";
import Attachment from "../attachment.tsx";
import { getKey } from "../../util/language.ts";

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
                    <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>
                    <div className="flex flex-1 flex-col gap-3 min-h-0 overflow-y-auto p-6 pr-14">
                        <div className="text-3xl rubik truncate">{homework_submission.title}</div>
                        <div className="text-2xl opacity-70 rubik truncate">{homework_submission.subject}</div>
                        <Attachment missing_label={getKey("NO_SUBMISSION")} label={getKey("SUBMISSION")} attachment={homework_submission}/>
                    </div>
                </div>
            </div>
        </>
    )
}
