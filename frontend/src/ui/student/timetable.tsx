import type {Class, Exam, HomeworkSummary, Lesson, LessonTime, Room} from "../../types/api.ts";
import {fromResponse, getLanguage} from "../../util/language.ts";
import {useCallback, useEffect, useRef, useState} from "preact/compat";
import {toast} from "react-toastify";
import "../homework/homework.css";
import {fetchInto, fetchWithSchool, getRetryAfter} from "../../util/api.tsx";
import type { FeatureProps } from "../../types/props.ts";
import {type Day, type DayName, DAY_NAMES, Timetable, type PopupProps} from "../timetable/timetable.tsx"
import { formatUnixDate, getWeekRange } from "../../util/time.ts";
import Loading from "../../util/loading.tsx";
import RateLimit from "../../util/ratelimit.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHouseChimney, faPenToSquare, faXmark} from "@fortawesome/free-solid-svg-icons";
import Overlay from "../../util/overlay/overlay.tsx";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {Badge} from "../../util/button/button.tsx";

export default function StudentTimetable({ me, fetchMe }: FeatureProps) {
    const [rooms, setRooms] = useState<Room[]>([])
    const [lessonTime, setLessonTime] = useState<LessonTime[]>([])
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

        const response = await fetchWithSchool(fetchMe, `/api/v1/me/student/timetable?start_date=${start}&end_date=${end}`)

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
            fetchInto(fetchMe, "/api/v1/me/student/timetable/room", setRooms),
            fetchInto<Class>(fetchMe, "/api/v1/me/student/class").then(clazz => fetchInto(fetchMe, `/api/v1/me/student/timetable/lesson_time?type_id=${clazz?.bell_id}`, setLessonTime).then(time => {
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
            })),
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
            <Timetable me={me} year={year} days={days} rooms={rooms} lessonTimes={lessonTime} onPrevWeek={() => fetchWeekLessons(weekOffsetRef.current - 1)} onNextWeek={() => fetchWeekLessons(weekOffsetRef.current + 1)} reloadSchoolStuff={Fetch} Popup={Popup} fetchMe={fetchMe} />
        </>
    )
}

function Popup({ lesson, rooms, clear, fetchMe }: PopupProps) {
    const [homework, setHomework] = useState<HomeworkSummary | null>(null)
    const [exam, setExam] = useState<Exam | null>(null)

    useEffect(() => {
        if (lesson.homework_id != 0) {
            void fetchWithSchool(fetchMe, `/api/v1/me/student/homework/get?id=${lesson.homework_id}`).then(async res => setHomework(await res.json() as HomeworkSummary))
        }
        if (lesson.exam_id != 0) {
            void fetchWithSchool(fetchMe, `/api/v1/me/student/exams/get?id=${lesson.exam_id}`).then(async res => setExam(await res.json() as Exam))
        }
    }, [fetchMe, lesson.exam_id, lesson.homework_id]);

    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="popup-content">
                    <div className="absolute top-6 left-6 right-16 flex flex-row gap-3">
                        <div className="text-3xl rubik flex-1 truncate">
                            {lesson.subject_name}
                        </div>
                        <div className="shrink min-w-0 max-w-[45%]">
                            {rooms.find(r => r.id === lesson.room_id)?.name}
                        </div>
                    </div>
                    <div className={`absolute top-15 left-6 text-2xl opacity-70 rubik max-w-[56vw] truncate ${lesson.is_substitution && "txt-(--warning-color)"}`}> {/* truncate just in case */}
                        {lesson.teacher_first_name} {lesson.teacher_last_name}
                    </div>
                    <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>

                    <div className="absolute top-28 left-6 right-6 h-fit fredoka flex flex-row gap-2">
                        {lesson.homework_id != 0 && <div className="min-w-[50%] flex flex-col gap-1">
                            <div className="flex flex-row w-full items-center gap-2">
                                <Badge icon={faHouseChimney} className="bg-(--warning-base-color) text-(--warning-color)" />
                                <p className="text-3xl rubik truncate">{homework?.title || <Skeleton width={128} height={36} color="var(--hover-color)" />}</p>
                            </div>
                            <p className="fredoka wrap-break-word">{homework?.description || <Skeleton width={256} height={24} color="var(--hover-color)" />}</p>
                        </div>}
                        {lesson.exam_id != 0 && <div className="min-w-[50%] flex flex-col">
                            <div className="flex flex-row w-full items-center gap-2">
                                <Badge icon={faPenToSquare} className="bg-(--wrong-base-color) text-(--wrong-color)" />
                                <p className="text-3xl rubik truncate">{exam?.title || <Skeleton width={128} height={36} color="var(--hover-color)" />}</p>
                            </div>
                            <p className="fredoka wrap-break-word">{exam?.description || <Skeleton width={256} height={24} color="var(--hover-color)" />}</p>
                        </div>}
                    </div>
                </div>
            </div>
        </>
    )
}
