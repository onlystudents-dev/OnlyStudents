import "./grades.css";
import type {Me, Grade, Subject} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {fetchWithSchool} from "../../util/api.ts";
import {useEffect} from "react";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {formatUnixDate} from "../../util/time.ts";
import {getKey, getLanguage} from "../../util/language.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";

const colors: Record<number, string> = {
    1: "wrong",
    2: "warning",
    3: "yellow",
    4: "right",
    5: "edit",
}

export default function StudentGrades(/*{ me }: { me: Me }*/) {
    const [subjects, setSubjects] = useState<Subject[]>([]);

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/student/subjects").then(async res => setSubjects(await res.json() as Subject[]))
    }

    useEffect(() => void Fetch(), [])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full flex justify-center items-center pt-8">
                <div className="subjects">
                    {subjects.map(subject => (
                        <div className="subject" key={subject.id}>
                            <div className="flex flex-col">
                                <h1 className="rubik truncate">{subject.name}</h1>
                                <p className="poppins truncate">{subject.teacher_first_name} {subject.teacher_last_name}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </>
    )
}

export function StudentGradesList({ me }: { me: Me }) {
    const [grades, setGrades] = useState<Grade[] | null>(null)
    const [current, setCurrent] = useState<Grade | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/student/grades").then(async res => setGrades(await res.json() as Grade[]))
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
                                    {current.note}
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
                                <p>{getKey("TERM")}: {current.term}</p>
                                <p>{getKey("TYPE")}: {current.type}</p>
                                <p>{getKey("TEACHER")}: {current.teacher_first_name} {current.teacher_last_name}</p>
                            </div>
                        </div>
                    </div>
                </>}
                {grades?.map((grade) => (
                    <>
                        <div className="grade" onClick={() => setCurrent(grade)}>
                            <div className="flex flex-col items-start justify-between">
                                <p className="max-w-100 truncate">{grade.note}</p>
                                <p className="max-w-100 truncate opacity-70">{grade.subject}</p>
                            </div>
                            <div className="flex flex-col items-start justify-between">
                                <div className={`rounded-full size-8 inline-flex items-center justify-center poppins text-2xl`} style={{ backgroundColor: `var(--${colors[grade.value]}-base-color)`, color: `var(--${colors[grade.value]}-color)`, fontWeight: "900" }}>{grade.value}</div>
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
