import "../homework/homework.css";
import type {HomeworkSummary, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {useEffect} from "react";
import {fetchWithSchool} from "../../util/api.ts";
import Homework, {type PopupProps} from "../homework/homework.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";

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
                    <div className="absolute top-6 left-6 text-3xl rubik max-w-[56vw] truncate">
                        {homework.title}
                    </div>
                    <div className="absolute top-15 left-6 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                        {homework.subject}
                    </div>
                    <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>

                    <div className="absolute top-28 left-6 right-6 h-fit fredoka">
                        {homework.description}
                    </div>
                </div>
            </div>
        </>
    )
}
