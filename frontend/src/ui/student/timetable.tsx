import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";
import {fromResponse, getKey, getLanguage} from "../../util/language.ts";
import React, {useCallback, useEffect, useRef, useState} from "react";
import "./timetable.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faAngleLeft, faAngleRight, faHouseChimney, faPenToSquare} from "@fortawesome/free-solid-svg-icons";
import RateLimit from "../../util/ratelimit.tsx";

type Lesson = {
    room_id: number,
    lesson_num: number,
    day_of_week: number,
    effective_teacher_id: number,
    group_id: number,
    school_id: number,
    custom_subject: boolean,
    subject_id: number | null,
    custom_subject_id: number | null,
    subject_name: string | null,
    is_substitution: boolean,
    canceled: boolean,
    has_exam: boolean,
    has_homework: boolean,
    teacher_first_name: string,
    teacher_last_name: string,
}

type Room = {
    id: number,
    name: string,
    capacity: number,
}

type day = {
    date: string,
    lessons: Lesson[],
}

export default function StudentTimetable({ me }: {me: Me}) {
    const [rooms, setRooms] = React.useState<Room[]>([])

    const [monday, setMonday] = React.useState<day>({date: "", lessons: []})
    const [tuesday, setTuesday] = React.useState<day>({date: "", lessons: []})
    const [wednesday, setWednesday] = React.useState<day>({date: "", lessons: []})
    const [thursday, setThursday] = React.useState<day>({date: "", lessons: []})
    const [friday, setFriday] = React.useState<day>({date: "", lessons: []})
    const [saturday, setSaturday] = React.useState<day>({date: "", lessons: []})
    const [sunday, setSunday] = React.useState<day>({date: "", lessons: []})

    const [year, setYear] = React.useState<string>("")

    const [ratelimit, setRateLimit] = useState(-1)

    const fetchWeekLessons = useCallback(async (date = new Date()) => {
        const { start, end, monday, mondayTime } = getWeekRange(date)

        const response = await fetch(`/api/v1/me/timetable?start_date=${start}&end_date=${end}`)

        if (response.status === 429) {
            setRateLimit(response.headers.get("retry-after") as unknown as number)
            return
        } else if (!response.ok) {
            await fromResponse(response)
            return
        }

        if (monday.getMonth() >= 7) {
            setYear(`${monday.getFullYear()}-${monday.getFullYear() + 1}`)
        } else {
            setYear(`${monday.getFullYear() - 1}-${monday.getFullYear()}`)
        }

        const lessons: Lesson[] = await response.json()
        const day = 60 * 60 * 24
        const language = getLanguage(me)?.key || "en-US"

        setMonday({
            date: formatUnixDate(mondayTime, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 1)
        })

        setTuesday({
            date: formatUnixDate(mondayTime + day, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 2)
        })

        setWednesday({
            date: formatUnixDate(mondayTime + day * 2, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 3)
        })

        setThursday({
            date: formatUnixDate(mondayTime + day * 3, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 4)
        })

        setFriday({
            date: formatUnixDate(mondayTime + day * 4, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 5)
        })

        setSaturday({
            date: formatUnixDate(mondayTime + day * 5, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 6)
        })

        setSunday({
            date: formatUnixDate(mondayTime + day * 6, language),
            lessons: lessons.filter(lesson => lesson.day_of_week === 0)
        })
    }, [me])

    useEffect(() => {
        async function Fetch() {
            await fetchWeekLessons()
            const response = await fetch("/api/v1/me/timetable/room")
            if (!response) {
                await fromResponse(response)
                return
            }

            setRooms(await response.json())
        }

        Fetch()
    }, [fetchWeekLessons]);

    const days: Record<string, day> = {
        monday,
        tuesday,
        wednesday,
        thursday,
        friday,
        saturday,
        sunday,
    }

    const offset = useRef(0)

    return (
        <>
            <Navbar me={me} />
            {ratelimit !== -1 && <RateLimit retry={ratelimit} expire={() => setRateLimit(-1)} />}
            <div className="w-full flex flex-row justify-center pt-4">
                <h1 className="text-5xl rubik">{year} {getKey("TERM")}</h1>
            </div>
            <div className="w-full h-fit p-4 gap-4 flex flex-row justify-center items-start">
                <div className="flex flex-col justify-center h-14">
                    <span className="icon" onClick={() => fetchWeekLessons(addDays(-7))}>
                        <FontAwesomeIcon icon={faAngleLeft} />
                    </span>
                </div>
                {Object.entries(days).filter(([, day]) => day.lessons.length > 0).map(([name, day]) => (
                    <div className="day">
                        <div className="date">
                            <h1 className="rubik">{getKey(`DAYS.${name}`)}</h1>
                            <p className="poppins">{day.date}</p>
                        </div>
                        {day.lessons.map((lesson) => (
                            <div className="lesson">
                                <div className="flex flex-col justify-between items-start h-full">
                                    <div className="flex flex-col">
                                        <h1 className="rubik">{lesson.subject_name}</h1>
                                        <h2 className="poppins">{lesson.teacher_first_name} {lesson.teacher_last_name}</h2>
                                    </div>
                                    <p className="fredoka">{rooms.find(room => room.id === lesson.room_id)?.name}</p>
                                </div>
                                <div className="flex flex-col justify-between items-end h-full">
                                    <div className="flex flex-col gap-1">
                                        {lesson.has_exam && <span className="bg-(--wrong-base-color) text-(--wrong-color) rounded-full size-8 inline-flex items-center justify-center">
                                            <FontAwesomeIcon icon={faPenToSquare} />
                                        </span>}
                                        {lesson.has_homework && <span className="bg-(--warning-base-color) text-(--warning-color) rounded-full size-8 inline-flex items-center justify-center">
                                            <FontAwesomeIcon icon={faHouseChimney} />
                                        </span>}
                                    </div>
                                    <h2></h2>
                                </div>
                            </div>
                        ))}
                    </div>
                ))}
                <div className="flex flex-col justify-center h-14">
                    <span className="icon" onClick={() => fetchWeekLessons(addDays(+7))}>
                        <FontAwesomeIcon icon={faAngleRight} />
                    </span>
                </div>
            </div>
        </>
    )

    function formatUnixDate(unix: number, locale: string): string {
        return new Intl.DateTimeFormat(locale, {
            month: "short",
            day: "numeric",
        }).format(new Date(unix * 1000));
    }

    function getWeekRange(date: Date) {
        const day = date.getDay()
        const diffToMonday = day === 0 ? -6 : 1 - day

        const monday = new Date(date)
        monday.setDate(date.getDate() + diffToMonday)
        monday.setHours(0, 0, 0, 0)

        const sunday = new Date(monday)
        sunday.setDate(monday.getDate() + 6)

        const toISODate = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`

        return { start: toISODate(monday), end: toISODate(sunday), monday, mondayTime: Math.floor(monday.getTime() / 1000) }
    }

    function addDays(days: number, date = new Date()) {
        offset.current += days

        const newDate = new Date(date)
        newDate.setDate(newDate.getDate() + offset.current)

        return newDate
    }
}