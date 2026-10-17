import "../homeworks/homeworks.css";
import type {HomeworkSubmission, Me} from "../../types/api.ts";
import {useLayoutEffect, useState} from "preact/hooks";
import {useEffect, useRef} from "react";
import {fetchWithSchool} from "../../util/api.ts";
import Homeworks, {type PopupProps} from "../homeworks/homeworks.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faFloppyDisk, faXmark} from "@fortawesome/free-solid-svg-icons";
import {toast} from "react-toastify";
import {fromResponse, getKey} from "../../util/language.ts";

export default function TeacherHomeworks({ me }: { me: Me }) {
    const [homeworks, setHomeworks] = useState<HomeworkSubmission[] | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/teacher/homework").then(async res => setHomeworks(await res.json() as HomeworkSubmission[]))
    }

    useEffect(() => {
        void Fetch()
    }, [])

    return <Homeworks me={me} homeworks={homeworks} Popup={Popup} reloadSchoolStuff={Fetch} />
}

function Popup({ homework: hw, clear }: PopupProps) {
    const homework = hw as HomeworkSubmission

    const [title, setTitle] = useState(homework.title)
    const [content, setContent] = useState(homework.content)
    const [due, setDue] = useState<string>(new Date(homework.due_date * 1000).toISOString().slice(0, 10))
    const contentRef = useRef<HTMLTextAreaElement>(null)

    useLayoutEffect(() => {
        const textarea = contentRef.current
        if (!textarea) return

        textarea.style.height = 'auto'
        textarea.style.height = `${textarea.scrollHeight}px`
    }, [content])

    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="popup-content">
                    <div className="absolute top-4 left-6 right-16 flex flex-row items-center gap-3">
                        <textarea className="min-w-0 flex-1 text-3xl rubik bg-transparent resize-none outline-none p-2" value={title} onChange={(e) => setTitle(e.currentTarget.value)} rows={1} />
                        <input type="date" value={due} onChange={(e) => setDue(e.currentTarget.value)} className="shrink min-w-0 max-w-[45%] p-2 outline-none" />
                    </div>
                    <div className="absolute top-18 left-8 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                        {homework.subject}
                    </div>
                    <div className="absolute top-7.5 right-7.5 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>

                    <textarea className="absolute top-30 left-6 right-6 h-fit min-h-24 fredoka bg-transparent resize-none outline-none p-2" value={content} onChange={(e) => setContent(e.currentTarget.value)} ref={contentRef} />

                    <div className="absolute bottom-6 right-6 bgaction-button hover:bg-(--bg-color) duration-200! p-2 rounded-lg" onClick={save}>
                        <FontAwesomeIcon icon={faFloppyDisk} /> Save
                    </div>
                </div>
            </div>
        </>
    )

    async function save() {
        const response = await fetchWithSchool('/api/v1/me/teacher/homework', {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                homework_id: homework.homework_id,
                title,
                description: content,
                has_description: !!content,
                class_subjects_id: homework.class_subjects_id,
                due_date: Math.floor(Date.parse(due) / 1000),
            }),
        })

        switch (response.status) {
            case 200:
                toast.success(getKey("TEACHER_UPDATED_HOMEWORK"))
                break
            default:
                toast.error(await fromResponse(response))
        }
    }
}