import type {Group, Lesson, LessonTime, Room, Student} from "../../types/api.ts";
import {fromResponse, getKey, getLanguage} from "../../util/language.ts";
import {useCallback, useEffect, useRef, useState} from "preact/compat";
import {toast} from "react-toastify";
import {fetchInto, fetchWithSchool} from "../../util/api.tsx";
import type { FeatureProps } from "../../types/props.ts";
import {
    type Day,
    type DayName,
    DAY_NAMES,
    Timetable,
    type PopupProps,
    toHoursAndMinutes
} from "../timetable/timetable.tsx"
import { formatUnixDate, getWeekRange } from "../../util/time.ts";
import Loading from "../../util/loading.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {faFloppyDisk, faList, faStar, faUser} from "@fortawesome/free-solid-svg-icons";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import Button from "../../util/sidebar/config.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {useLayoutEffect} from "preact/hooks";

export default function TeacherTimetable({ me, fetchMe }: FeatureProps) {
    const [rooms, setRooms] = useState<Room[]>([])
    const [lessonTime, setLessonTime] = useState<LessonTime[]>([])
    const [groups, setGroups] = useState<Group[]>([])
    const [days, setDays] = useState<Day[]>(() => DAY_NAMES.map(name => ({name, date: "", lessons: [], today: false})))
    const [year, setYear] = useState<string>("")
    const [loading, setLoading] = useState(true)

    const weekOffsetRef = useRef(0)

    const fetchWeekLessons = useCallback(async (target: number) => {
        setLoading(true)

        const date = new Date()
        date.setDate(date.getDate() + target * 7)
        const { start, end, monday, mondayTime } = getWeekRange(date)

        const response = await fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/mine?start_date=${start}&end_date=${end}`)

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
            fetchInto(fetchMe, "/api/v1/me/teacher/timetable/room", setRooms),
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
            <Timetable me={me} year={year} days={days} rooms={rooms} lessonTimes={lessonTime} onPrevWeek={() => fetchWeekLessons(weekOffsetRef.current - 1)} onNextWeek={() => fetchWeekLessons(weekOffsetRef.current + 1)} reloadSchoolStuff={Fetch} Popup={Popup} groups={groups} fetchMe={fetchMe} />
        </>
    )
}

function Popup({ lesson, rooms, groups, clear, fetchMe, me, lessonTimes }: PopupProps) {
    const [active, setActive] = useState<"overview" | "grade">("overview")
    const [students, setStudents] = useState<Student[]>([])

    const [wasHere, setWasHere] = useState<Record<number, boolean>>({})
    const presentRef = useRef<HTMLButtonElement>(null)
    const absentRef = useRef<HTMLButtonElement>(null)
    const [presentWidth, setPresentWidth] = useState<number>(0)
    const [absendWidth, setAbsendWidth] = useState<number>(0)

    const room = rooms.find(r => r.id === lesson.room_id)
    const group = groups?.find(g => g.id === lesson.group_id)

    useEffect(() => {
        void fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/group/student?group_id=${group?.id}`).then(async response => {
            const data = await response.json() as Student[]
            setWasHere(Object.fromEntries(data.map(student => [student.id, true])))
            setStudents(data)
        })
    }, [fetchMe, group])

    useLayoutEffect(() => {
        if (!students.length) return

        setPresentWidth(presentRef.current?.offsetWidth || 0)
        setAbsendWidth(absentRef.current?.offsetWidth || 0)
    }, [students])

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
                                            <div className="flex flex-col items-end">
                                                <p className="text-2xl truncate fredoka">{room.name}</p>
                                                <p className="truncate rubik opacity-70">{toHoursAndMinutes(me, lesson, lessonTimes)}</p>
                                            </div>
                                        </div>
                                        <div className="border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                                            {students.map(student => (
                                                <div key={student.id} className="flex flex-row justify-between items-center w-full hover:bg-(--hover-color) transition-[background-color_400ms_ease-in-out] hover:transition-[background-color_400ms_ease-in-out] p-2 rounded-lg">
                                                    <p><FontAwesomeIcon icon={faUser} /> {student.first_name} {student.last_name}</p>
                                                    <div className="flex flex-row items-center h-full justify-end rounded-lg border-2 border-(--border-color) relative">
                                                        <button className="text-center p-1 cursor-pointer z-1" onClick={() => setWasHere({ ...wasHere, [student.id]: true })} ref={presentRef}>
                                                            {getKey("PRESENT")}
                                                        </button>
                                                        <button className="text-center p-1 cursor-pointer z-1" onClick={() => setWasHere({ ...wasHere, [student.id]: false })} ref={absentRef}>
                                                            {getKey("ABSENT")}
                                                        </button>
                                                        <div className="bg-(--card-color) absolute top-0 bottom-0 transition-[left_200ms_ease-in-out]" style={{
                                                            left: wasHere[student.id] ? '0px' : `${presentWidth}px`,
                                                            width: `${wasHere[student.id] ? presentWidth : absendWidth}px`,
                                                            borderTopLeftRadius: wasHere[student.id] ? '6px' : '0px',
                                                            borderBottomLeftRadius: wasHere[student.id] ? '6px' : '0px',
                                                            borderTopRightRadius: !wasHere[student.id] ? '6px' : '0px',
                                                            borderBottomRightRadius: !wasHere[student.id] ? '6px' : '0px',
                                                        }} />
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </>
                                )
                            default:
                                return null
                        }
                    })()}
                    <button className="absolute bottom-4 right-4 bg-(--card-color) p-2 rounded-lg cursor-pointer" onClick={save}>
                        <FontAwesomeIcon icon={faFloppyDisk} /> {getKey("SAVE")}
                    </button>
                </div>
                <Sidebar className="border-(--txt-color)!">
                    <Button icon={faList} text={getKey("TIMETABLE_OVERVIEW")} onClick={() => setActive("overview")} />
                    <Button icon={faStar} text={getKey("TIMETABLE_GRADE")} onClick={() => setActive("grade")} />
                </Sidebar>
            </div>
        </>
    )

    async function save() {
        for (const [student, here] of Object.entries(wasHere)) {
            if (here) continue;

            await fetchWithSchool(fetchMe, "/api/v1/me/teacher/absences", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    date: lesson.actual_date,
                    type: "absent",
                    has_note: false,
                    note: "",
                    student_id: Number(student),
                    class_subjects_id: lesson.subject_id,
                })
            })
        }
    }
}
