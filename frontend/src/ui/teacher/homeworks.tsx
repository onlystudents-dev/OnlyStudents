import "../homeworks/homeworks.css";
import type {HomeworkSubmission, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {useEffect} from "react";
import {fetchWithSchool} from "../../util/api.ts";
import Homeworks, {type PopupProps} from "../homeworks/homeworks.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";

export default function TeacherHomeworks({ me }: { me: Me }) {
    const [homeworks, setHomeworks] = useState<HomeworkSubmission[] | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/teacher/homework").then(async res => setHomeworks(await res.json() as HomeworkSubmission[]))
    }

    useEffect(() => void Fetch(), [])

    return <Homeworks me={me} homeworks={homeworks} Popup={Popup} reloadSchoolStuff={Fetch} />
}

function Popup({ homework: hw, clear }: PopupProps) {
    const homework = hw as HomeworkSubmission

    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="popup-content">
                    <div className="absolute top-4 left-6 right-16 flex flex-row items-center gap-3">
                        <div className="min-w-0 flex-1 text-3xl rubik bg-transparent resize-none outline-none p-2">
                            {homework.title}
                        </div>
                    </div>
                    <div className="absolute top-18 left-8 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                        {homework.subject}
                    </div>
                    <div className="absolute top-7.5 right-7.5 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>

                    <div className="absolute top-30 left-6 right-6 h-fit min-h-24 fredoka bg-transparent resize-none outline-none p-2">
                        {homework.content}
                    </div>
                </div>
            </div>
        </>
    )
}