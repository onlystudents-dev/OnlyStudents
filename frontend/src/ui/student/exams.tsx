import "./exams.css";
import {getLanguage} from "../../util/language.ts";
import {formatUnixDate} from "../../util/time.ts";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import type {Exam, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {fetchWithSchool} from "../../util/api.ts";
import {useEffect} from "react";
import Overlay from "../../util/overlay/overlay.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";

export default function StudentExams({ me }: { me: Me }) {
    const [exams, setExams] = useState<Exam[] | null>(null)
    const [current, setCurrent] = useState<Exam | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/student/exams").then(async res => setExams(await res.json() as Exam[]))
    }

    useEffect(() => void Fetch(), [])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && <>
                    <Overlay onClick={() => setCurrent(null)} />
                    <div className="popup">
                        <div className="popup-content">
                            <div className="absolute top-6 left-6 right-16 flex flex-row gap-3">
                                <div className="text-3xl rubik flex-1 truncate">
                                    {current.title}
                                </div>
                                <div className="shrink min-w-0 max-w-[45%]">
                                    {formatUnixDate(current.date, getLanguage(me)?.key || "en-US")}
                                </div>
                            </div>
                            <div className="absolute top-15 left-6 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                                {current.subject}
                            </div>
                            <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={() => setCurrent(null)}>
                                <FontAwesomeIcon icon={faXmark} />
                            </div>

                            <div className="absolute top-28 left-6 right-6 h-fit fredoka">
                                {current.description}
                            </div>
                        </div>
                    </div>
                </>}
                {exams?.map((exam) => (
                    <>
                        <div className="exam" onClick={() => setCurrent(exam)}>
                            <div className="flex flex-col items-start justify-between">
                                <p className="max-w-100 truncate">{exam.title}</p>
                                <p className="max-w-100 truncate opacity-70">{exam.subject}</p>
                            </div>
                            <div className="flex flex-col items-start justify-between">
                                <p>{formatUnixDate(exam.date, getLanguage(me)?.key || "en-US")}</p>
                            </div>
                        </div>
                    </>
                )) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                ))}
            </div>
        </>
    )
}