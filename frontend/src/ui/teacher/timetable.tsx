import type {Group, Lesson, LessonTime, Room, Student} from "../../types/api.ts";
import {fromResponse, getKey, getLanguage} from "../../util/language.ts";
import {useCallback, useEffect, useRef, useState} from "preact/compat";
import {toast} from "react-toastify";
import {fetchInto, fetchWithSchool, getRetryAfter} from "../../util/api.tsx";
import type { FeatureProps } from "../../types/props.ts";
import {type Day, type DayName, DAY_NAMES, Timetable, type PopupProps} from "../timetable/timetable.tsx"
import { formatUnixDate, getWeekRange } from "../../util/time.ts";
import Loading from "../../util/loading.tsx";
import RateLimit from "../../util/ratelimit.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {faList, faStar} from "@fortawesome/free-solid-svg-icons";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import Button from "../../util/sidebar/config.tsx";

export default function TeacherTimetable({ me, fetchMe }: FeatureProps) {
    const [rooms, setRooms] = useState<Room[]>([])
    const [lessonTime, setLessonTime] = useState<LessonTime[]>([])
    const [groups, setGroups] = useState<Group[]>([])
    const [days, setDays] = useState<Day[]>(() => DAY_NAMES.map(name => ({name, date: "", lessons: [], today: false})))
    const [year, setYear] = useState<string>("")
    const [ratelimit, setRateLimit] = useState(-1)
    const [loading, setLoading] = useState(true)

    const weekOffsetRef = useRef(0)

    const fetchWeekLessons = useCallback(async (target: number) => {
        setLoading(true)

        const date = new Date()
        date.setDate(date.getDate() + target * 7)
        const { start, end, monday, mondayTime } = getWeekRange(date)

        const response = await fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/mine?start_date=${start}&end_date=${end}`)

        if (response.status === 429) {
            setRateLimit(getRetryAfter(response))
            setLoading(false)
            return
        }

        if (!response.ok) {
            toast.error(await fromResponse(response))
            setLoading(false)
            return
        }

        weekOffsetRef.current = target

        if (monday.getMonth() >= 7) {
            setYear(`${monday.getFullYear()}-${monday.getFullYear() + 1}`)
        } else {
            setYear(`${monday.getFullYear() - 1}-${monday.getFullYear()}`)
        }

        const lessons = await response.json() as Lesson[]
        const daySeconds = 86400
        const language = getLanguage(me)?.key || "en-US"
        const now = new Date()

        const createDayData = (name: DayName, dayOffset: number, dayOfWeek: number): Day => {
            const targetDate = new Date(monday)
            targetDate.setDate(monday.getDate() + dayOffset)

            return {
                name,
                date: formatUnixDate(mondayTime + daySeconds * dayOffset, language),
                lessons: lessons.filter(lesson => lesson.day_of_week === dayOfWeek),
                today: targetDate.toDateString() === now.toDateString(),
            }
        }

        setDays(DAY_NAMES.map((name, i) => createDayData(name, i, (i + 1) % 7)))

        setLoading(false)
        return lessons
    }, [fetchMe, me])

    const Fetch = useCallback(async () => {
        const lessons = await fetchWeekLessons(0)

        await Promise.all([
            fetchInto(fetchMe, "/api/v1/me/teacher/timetable/my_rooms", setRooms),
            fetchInto(fetchMe, "/api/v1/me/teacher/timetable/group", setGroups),
            fetchInto(fetchMe, "/api/v1/me/teacher/timetable/my_lesson_times", setLessonTime).then(time => {
                if (!me.preferences.timetable_next) return

                if (!time) return

                const lastLesson = lessons?.at(-1)
                if (!lastLesson) return

                const t = time.find(lt => lt.lesson_number === lastLesson.lesson_num)
                if (!t) return

                const date = new Date()
                const current = date.getDay() || 7

                date.setDate(date.getDate() - (current - lastLesson.day_of_week))
                date.setHours(0, 0, 0, 0)

                const end = date.getTime() + t.at_end * 1000

                if (Date.now() > end) {
                    void fetchWeekLessons(1)
                }
            }),
        ])
    }, [fetchMe, fetchWeekLessons, me.preferences.timetable_next])

    useEffect(() => {
        void Fetch()
    }, [Fetch])

    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            {loading && <Loading />}
            {ratelimit !== -1 && <RateLimit retry={ratelimit} expire={() => setRateLimit(-1)} />}
            <Timetable me={me} year={year} days={days} rooms={rooms} lessonTimes={lessonTime} onPrevWeek={() => fetchWeekLessons(weekOffsetRef.current - 1)} onNextWeek={() => fetchWeekLessons(weekOffsetRef.current + 1)} reloadSchoolStuff={Fetch} Popup={Popup} groups={groups} fetchMe={fetchMe} />
        </>
    )
}

function Popup({ lesson, rooms, groups, clear, fetchMe }: PopupProps) {
    const [active, setActive] = useState<"overview" | "grade">("overview")
    const [students, setStudents] = useState<Student[]>([])

    const room = rooms.find(r => r.id === lesson.room_id)
    const group = groups?.find(g => g.id === lesson.group_id)

    useEffect(() => {
        void fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/group/student?group_id=${group?.id}`).then(async response => setStudents(await response.json() as Student[]))
    }, [fetchMe, group])

    if (!group || !room) return

    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="absolute bottom-0 right-0 left-90 h-full p-4 flex flex-col gap-4">
                    {(() => {
                        switch (active) {
                            case "overview":
                                return (
                                    <>
                                        <div className="flex flex-row justify-between">
                                            <div className="flex flex-col">
                                                <p className="text-3xl truncate">{lesson.subject_name}</p>
                                                <p className="text-2xl opacity-70 truncate fredoka">{group.group_name}</p>
                                            </div>
                                            <div className="flex flex-col">
                                                <p className="text-2xl truncate fredoka">{room.name}</p>
                                            </div>
                                        </div>
                                        <div className="border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                                            {students.map(student => (
                                                <div key={student.id} className="flex flex-row justify-between items-center w-full">
                                                    <p>{student.first_name} {student.last_name}</p>
                                                </div>
                                            ))}
                                        </div>
                                    </>
                                )
                            default:
                                return null
                        }
                    })()}
                </div>
                <Sidebar className="border-(--txt-color)!">
                    <Button icon={faList} text={getKey("TIMETABLE_OVERVIEW")} onClick={() => setActive("overview")} />
                    <Button icon={faStar} text={getKey("TIMETABLE_GRADE")} onClick={() => setActive("grade")} />
                </Sidebar>
            </div>
        </>
    )
}
