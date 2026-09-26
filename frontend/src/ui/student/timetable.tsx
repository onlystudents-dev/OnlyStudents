import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";
import {fromResponse, getKey, getLanguage} from "../../util/language.ts";
import React, {useCallback, useEffect, useRef, useState} from "react";
import "./timetable.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faAngleLeft, faAngleRight, faHouseChimney, faPenToSquare} from "@fortawesome/free-solid-svg-icons";
import RateLimit from "../../util/ratelimit.tsx";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import Loading from "../../util/loading.tsx";

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

type Class = {
    id: number,
    school_id: number,
    name: string,
    teacher_id: number,
    has_co_teacher_id: boolean,
    co_teacher_id: number,
    has_bell_id: boolean,
    bell_id: number,
}

type LessonTime = {
    id: number,
    school_id: number,
    type_id: number,
    lesson_number: number,
    at_start: LessonTimeDate,
    at_end: LessonTimeDate,
}

type LessonTimeDate = {
    Microseconds: number,
    Valid: boolean,
}

type day = {
    date: string,
    lessons: Lesson[],
    today: boolean
}

export default function StudentTimetable({ me }: {me: Me}) {
    const [rooms, setRooms] = useState<Room[]>([])
    const [lessonTime, setLessonTime] = useState<LessonTime[]>([])

    const [monday, setMonday] = useState<day>({date: "", lessons: [], today: false})
    const [tuesday, setTuesday] = useState<day>({date: "", lessons: [], today: false})
    const [wednesday, setWednesday] = useState<day>({date: "", lessons: [], today: false})
    const [thursday, setThursday] = useState<day>({date: "", lessons: [], today: false})
    const [friday, setFriday] = useState<day>({date: "", lessons: [], today: false})
    const [saturday, setSaturday] = useState<day>({date: "", lessons: [], today: false})
    const [sunday, setSunday] = useState<day>({date: "", lessons: [], today: false})

    const [year, setYear] = useState<string>("")

    const [ratelimit, setRateLimit] = useState(-1)
    const [loading, setLoading] = useState(true)

    const fetchWeekLessons = useCallback(async (date = new Date()) => {
        setLoading(true)
        const { start, end, monday, mondayTime } = getWeekRange(date)

        const response = await fetch(`/api/v1/me/timetable?start_date=${start}&end_date=${end}`)

        if (response.status === 429) {
            setRateLimit(response.headers.get("retry-after") as unknown as number)
            setLoading(false)
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
        const daySeconds = 86400
        const language = getLanguage(me)?.key || "en-US"
        const now = new Date()

        const createDayData = (dayOffset: number, dayOfWeek: number) => {
            const targetDate = new Date(monday)
            targetDate.setDate(monday.getDate() + dayOffset)

            return {
                date: formatUnixDate(mondayTime + daySeconds * dayOffset, language),
                lessons: lessons.filter(lesson => lesson.day_of_week === dayOfWeek),
                today: targetDate.toDateString() === now.toDateString(),
            }
        }

        setMonday(createDayData(0, 1))
        setTuesday(createDayData(1, 2))
        setWednesday(createDayData(2, 3))
        setThursday(createDayData(3, 4))
        setFriday(createDayData(4, 5))
        setSaturday(createDayData(5, 6))
        setSunday(createDayData(6, 0))

        setLoading(false)
    }, [me])

    useEffect(() => {
        async function Fetch() {
            await fetchWeekLessons()

            async function fetchInto<T>(
                api: string,
                method?: React.Dispatch<React.SetStateAction<T>>
            ) {
                const response = await fetch(api)

                if (!response.ok) {
                    await fromResponse(response)
                    return
                }

                const json = await response.json() as T

                if (method) method(json)
                return json
            }

            await Promise.all([
                fetchInto("/api/v1/me/timetable/room", setRooms),
                fetchInto<Class>("/api/v1/me/class").then(clazz => fetchInto(`/api/v1/me/timetable/lesson_time?type_id=${clazz?.bell_id}`, setLessonTime)),
            ])
        }

        Fetch()
    }, [fetchWeekLessons])

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
            {loading && <Loading />}
            {ratelimit !== -1 && <RateLimit retry={ratelimit} expire={() => setRateLimit(-1)} />}
            <div className="w-full flex flex-row justify-center pt-4">
                <h1 className="text-5xl rubik">{year} {getKey("TERM")}</h1>
            </div>
            <div className="w-full h-fit p-4 gap-4 flex flex-row justify-center items-stretch">
                <div className="flex flex-col justify-center h-14">
                    <span className="icon" onClick={() => fetchWeekLessons(addDays(-7))}>
                        <FontAwesomeIcon icon={faAngleLeft} />
                    </span>
                </div>
                {Object.entries(days).filter(([, day]) => day.lessons.length > 0).map(([name, day]) => (
                    <div className="day">
                        <div className={`date ${day.today && "rounded-2xl bg-(--border-color)"}`}>
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
                                    <p className="fredoka">{rooms.find(room => room.id === lesson.room_id)?.name || <Skeleton width={48} height={16} color={"var(--card-color)"} />}</p>
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
                                    <h2>{toHoursAndMinutes(lesson) || <Skeleton width={96} height={16} color={"var(--hover-color)"} />}</h2>
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

    function toHoursAndMinutes(lesson: Lesson) {
        const time = lessonTime.find(lt => lt.lesson_number === lesson.lesson_num)
        if (!time) return ""

        const language = getLanguage(me)?.key || "en-US"

        const format = (microseconds: number) =>
            new Intl.DateTimeFormat(language, {
                hour: "numeric",
                minute: "2-digit",
            }).format(new Date(microseconds / 1000))

        return `${format(time.at_start.Microseconds)}-${format(time.at_end.Microseconds)}`
    }

    function formatUnixDate(unix: number, locale: string): string {
        return new Intl.DateTimeFormat(locale, {
            month: "short",
            day: "numeric",
        }).format(new Date(unix * 1000))
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