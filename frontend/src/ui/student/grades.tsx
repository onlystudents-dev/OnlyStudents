import "./grades.css";
import type {Me, Grade, Subject} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {fetchWithSchool} from "../../util/api.ts";
import {useCallback, useEffect} from "react";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {formatUnixDate} from "../../util/time.ts";
import {getKey, getLanguage} from "../../util/language.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";
import {useLocation, useSearch} from "wouter";
import SubjectAverageChart from "./charts.tsx";

export const colors: Record<number, string> = {
    1: "wrong",
    2: "warning",
    3: "yellow",
    4: "right",
    5: "edit",
}

export default function StudentGrades({ me }: { me: Me }) {
    const search = useSearch()
    const subject = new URLSearchParams(search).get("subject")

    useEffect(() => {
        window.scrollTo(0, 0)
    }, [subject])

    return subject ? <StudentGradesList me={me} subject={subject} /> : <StudentSubjects me={me} />
}

export function StudentSubjects({ me }: { me: Me }) {
    const [, navigate] = useLocation()
    const [subjects, setSubjects] = useState<Subject[] | null>(null)
    const [grades, setGrades] = useState<Grade[] | null>(null)

    async function Fetch() {
        return Promise.all([
            fetchWithSchool("/api/v1/me/student/subjects").then(async res => setSubjects(await res.json() as Subject[])),
            fetchWithSchool("/api/v1/me/student/grades").then(async res => setGrades(await res.json() as Grade[])),
        ])
    }

    useEffect(() => void Fetch(), [])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full flex justify-center items-center pt-8">
                <div className="subjects">
                    {grades ? <SubjectAverageChart grades={grades} me={me} /> : <Skeleton width={640} height={232} color="var(--hover-color)" className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />}
                    {subjects?.map(subject => {
                        if (!grades) return <Skeleton width={640} height={64} color="var(--hover-color)" key={subject.id} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                        const gradeS = grades.filter(grade => grade.subject === subject.name)
                        const average = gradeS.length
                            ? gradeS.reduce((sum, grade) => sum + grade.value * grade.weight, 0) / gradeS.reduce((sum, grade) => sum + grade.weight, 0)
                            : 0

                        const color = colors[Math.round(average)]
                        const rounded = average.toFixed(2)

                        return (
                            <div className="grade" key={subject.id} onClick={() => navigate(`/grades?subject=${subject.is_custom ? "t" : "f"}${subject.id}`)}>
                                <div className="flex flex-col">
                                    <h1 className="rubik truncate">{subject.name}</h1>
                                    <p className="poppins opacity-70 truncate">{subject.teacher_first_name} {subject.teacher_last_name}</p>
                                </div>
                                <div className="flex flex-col">
                                    <div className={`rounded-full h-8 w-fit pl-2 pr-2 inline-flex items-center justify-center poppins text-2xl`} style={{
                                        backgroundColor: `var(--${color}-base-color)`,
                                        color: `var(--${color}-color)`, fontWeight: "900"
                                    }}>{Number(rounded) ? rounded : "---"}</div>
                                </div>
                            </div>
                        )
                    }) || Array.from({ length: 9 }).map((_, i) => (
                        <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                    ))}
                </div>
            </div>
        </>
    )
}

export function StudentGradesList({ me, subject }: { me: Me; subject: string }) {
    const [grades, setGrades] = useState<Grade[] | null>(null)
    const [current, setCurrent] = useState<Grade | null>(null)

    const Fetch = useCallback(async () => {
        return fetchWithSchool(`/api/v1/me/student/grades?subject=${subject}`).then(
            async (res) => setGrades((await res.json()) as Grade[])
        )
    }, [subject])

    useEffect(() => void Fetch(), [Fetch])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && (
                    <>
                        <Overlay onClick={() => setCurrent(null)} />
                        <div className="popup">
                            <div className="popup-content">
                                <div className="absolute top-6 left-6 right-16 flex flex-row gap-3">
                                    <div className="text-3xl rubik flex-1 truncate">
                                        {current.note}
                                    </div>
                                    <div className="shrink min-w-0 max-w-[45%]">
                                        {formatUnixDate(current.date, getLanguage(me)?.key || 'en-US')}
                                    </div>
                                </div>
                                <div className="absolute top-15 left-6 text-2xl opacity-70 rubik max-w-[56vw] truncate">
                                    {current.subject}
                                </div>
                                <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={() => setCurrent(null)}>
                                    <FontAwesomeIcon icon={faXmark} />
                                </div>

                                <div className="absolute top-28 left-6 right-6 h-fit fredoka">
                                    <p>
                                        {getKey("TERM")}: {current.term}
                                    </p>
                                    <p>
                                        {getKey("TYPE")}: {current.type}
                                    </p>
                                    <p>
                                        {getKey("TEACHER")}: {current.teacher_first_name}{' '}
                                        {current.teacher_last_name}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </>
                )}

                {grades ? <SubjectAverageChart grades={grades} me={me} /> : <Skeleton width={640} height={232} color="var(--hover-color)" className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />}

                {grades?.map((grade, idx) => (
                        <div key={idx} className="grade" onClick={() => setCurrent(grade)}>
                            <div className="flex flex-col items-start justify-between">
                                <p className="max-w-100 truncate">{grade.note}</p>
                                <p className="max-w-100 truncate opacity-70">{grade.subject}</p>
                            </div>
                            <div className="flex flex-col items-start justify-between">
                                <div className="rounded-full size-8 inline-flex items-center justify-center poppins text-2xl" style={{ backgroundColor: `var(--${colors[grade.value]}-base-color)`, color: `var(--${colors[grade.value]}-color)`, fontWeight: "900" }}>
                                    {grade.value}
                                </div>
                            </div>
                        </div>
                    )) ||
                    Array.from({ length: 9 }).map((_, i) => (
                        <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                    ))}
            </div>
        </>
    )
}
