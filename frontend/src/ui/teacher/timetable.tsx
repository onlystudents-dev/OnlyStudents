import type {GradeType, Group, Lesson, LessonTime, Room, Student} from "../../types/api.ts";
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
import {colors} from "../globals.tsx";
import Select from "react-select";
import {selectStyles} from "../../util/select.ts";

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
            fetchInto(fetchMe, "/api/v1/me/teacher/timetable/my_groups", setGroups),
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
    const [active, setActive] = useState<"overview" | "grade" | "homework">("overview")
    const [students, setStudents] = useState<Student[]>([])
    const [gradeTypes, setGradeTypes] = useState<GradeType[]>([])

    const [wasHere, setWasHere] = useState<Record<number, boolean>>({})
    const presentRef = useRef<HTMLButtonElement>(null)
    const absentRef = useRef<HTMLButtonElement>(null)
    const [presentWidth, setPresentWidth] = useState<number>(0)
    const [absendWidth, setAbsendWidth] = useState<number>(0)

    type EditingGrade = { title: string, description: string, type_id: number, grades: Record<number, number> }
    const [editingGrade, setEditingGrade] = useState<number | null>(null)
    const [grades, setGrades] = useState<Record<number, EditingGrade>>({})
    const ref1 = useRef<HTMLButtonElement>(null)
    const ref2 = useRef<HTMLButtonElement>(null)
    const ref3 = useRef<HTMLButtonElement>(null)
    const ref4 = useRef<HTMLButtonElement>(null)
    const ref5 = useRef<HTMLButtonElement>(null)
    const didntRef = useRef<HTMLButtonElement>(null)
    const [width1, setWidth1] = useState<number>(0)
    const [width2, setWidth2] = useState<number>(0)
    const [width3, setWidth3] = useState<number>(0)
    const [width4, setWidth4] = useState<number>(0)
    const [width5, setWidth5] = useState<number>(0)
    const [didntWidth, setDidntWidth] = useState<number>(0)
    const widthByVal: Record<number, number> = {
        0: didntWidth,
        1: width1,
        2: width2,
        3: width3,
        4: width4,
        5: width5,
    }

    const room = rooms.find(r => r.id === lesson.room_id)
    const group = groups?.find(g => g.id === lesson.group_id)

    useEffect(() => {
        void fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/students_in_my_group?group_id=${group?.id}`).then(async response => {
            const data = await response.json() as Student[]
            setWasHere(Object.fromEntries(data.map(student => [student.id, true])))
            setStudents(data)
        })
        void fetchWithSchool(fetchMe, "/api/v1/me/teacher/grades/types").then(async response => setGradeTypes(await response.json() as GradeType[]))
    }, [fetchMe, group])

    useLayoutEffect(() => {
        if (!students.length) return

        setPresentWidth(presentRef.current?.offsetWidth || 0)
        setAbsendWidth(absentRef.current?.offsetWidth || 0)
    }, [students])

    useLayoutEffect(() => {
        if (active !== "grade" || !editingGrade) return

        setWidth1(ref1.current?.offsetWidth || 0)
        setWidth2(ref2.current?.offsetWidth || 0)
        setWidth3(ref3.current?.offsetWidth || 0)
        setWidth4(ref4.current?.offsetWidth || 0)
        setWidth5(ref5.current?.offsetWidth || 0)
        setDidntWidth(didntRef.current?.offsetWidth || 0)
    }, [active, editingGrade])

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
                            case "grade":
                                if (editingGrade) {
                                    const eGrade = grades[editingGrade]
                                    if (!eGrade) return
                                    const gradeTypesOptions = gradeTypes.map((type, index) => ({ label: type.name, value: index }))

                                    return (
                                        <>
                                            <div className="flex flex-row justify-between">
                                                <div className="flex flex-col">
                                                    <input type="text" className="text-3xl truncate outline-none rubik" defaultValue={eGrade.title} onChange={e => setGrades({ ...grades, [editingGrade]: { ...eGrade, title: e.currentTarget.value } })} placeholder={getKey("TITLE")} />
                                                    <textarea rows={2} className="text-2xl opacity-70 truncate fredoka outline-none rubik" defaultValue={eGrade.description} onChange={e => setGrades({ ...grades, [editingGrade]: { ...eGrade, description: e.currentTarget.value } } )} placeholder={getKey("DESCRIPTION")} />
                                                </div>
                                                <div className="flex flex-col items-end">
                                                    <Select options={gradeTypesOptions} value={gradeTypesOptions[eGrade.type_id]} onChange={val => val && setGrades({ ...grades, [editingGrade]: { ...eGrade, type_id: val.value } })}
                                                            styles={selectStyles("fit-content", "100%", "var(--bg-color)", "var(--card-color)", "var(--hover-color)")} />
                                                </div>
                                            </div>
                                            <div className="border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                                                {Object.entries(eGrade.grades).map(([s, grade]) => {
                                                    const student = students.find(st => st.id === Number(s))
                                                    if (!student) return

                                                    return (
                                                        <div key={student.id} className="flex flex-row justify-between items-center w-full hover:bg-(--hover-color) transition-[background-color_400ms_ease-in-out] hover:transition-[background-color_400ms_ease-in-out] p-2 rounded-lg">
                                                            <p><FontAwesomeIcon icon={faUser} /> {student.first_name} {student.last_name}</p>
                                                            <div className="flex flex-row items-center h-full justify-end rounded-lg border-2 border-(--border-color) relative">
                                                                <button className="text-center p-1 cursor-pointer z-1 min-w-8" onClick={() => setGrades({ ...grades, [editingGrade]: { ...eGrade, grades: { ...eGrade.grades, [student.id]: 0 } } })} ref={didntRef}>
                                                                    -
                                                                </button>
                                                                <button className="text-center p-1 cursor-pointer z-1 min-w-8" onClick={() => setGrades({ ...grades, [editingGrade]: { ...eGrade, grades: { ...eGrade.grades, [student.id]: 1 } } })} ref={ref1} style={{ color: `var(--${colors[1]}-color)` }}>
                                                                    1
                                                                </button>
                                                                <button className="text-center p-1 cursor-pointer z-1 min-w-8" onClick={() => setGrades({ ...grades, [editingGrade]: { ...eGrade, grades: { ...eGrade.grades, [student.id]: 2 } } })} ref={ref2} style={{ color: `var(--${colors[2]}-color)` }}>
                                                                    2
                                                                </button>
                                                                <button className="text-center p-1 cursor-pointer z-1 min-w-8" onClick={() => setGrades({ ...grades, [editingGrade]: { ...eGrade, grades: { ...eGrade.grades, [student.id]: 3 } } })} ref={ref3} style={{ color: `var(--${colors[3]}-color)` }}>
                                                                    3
                                                                </button>
                                                                <button className="text-center p-1 cursor-pointer z-1 min-w-8" onClick={() => setGrades({ ...grades, [editingGrade]: { ...eGrade, grades: { ...eGrade.grades, [student.id]: 4 } } })} ref={ref4} style={{ color: `var(--${colors[4]}-color)` }}>
                                                                    4
                                                                </button>
                                                                <button className="text-center p-1 cursor-pointer z-1 min-w-8" onClick={() => setGrades({ ...grades, [editingGrade]: { ...eGrade, grades: { ...eGrade.grades, [student.id]: 5 } } })} ref={ref5} style={{ color: `var(--${colors[5]}-color)` }}>
                                                                    5
                                                                </button>
                                                                <div className="absolute top-0 bottom-0 transition-[left_200ms_ease-in-out]" style={{
                                                                    left: `${Object.entries(widthByVal).slice(0, grade).reduce((sum, [, width]) => sum + width, 0)}px`,
                                                                    width: `${widthByVal[grade]}px`,
                                                                    borderTopLeftRadius: grade === 0 ? '6px' : '0px',
                                                                    borderBottomLeftRadius: grade === 0 ? '6px' : '0px',
                                                                    borderTopRightRadius: grade === 5 ? '6px' : '0px',
                                                                    borderBottomRightRadius: grade === 5 ? '6px' : '0px',
                                                                    backgroundColor: grade === 0 ? "var(--card-color)" : `var(--${colors[grade]}-base-color)`
                                                                }} />
                                                            </div>
                                                        </div>
                                                    )
                                                })}
                                            </div>
                                        </>
                                    )
                                }

                                return (
                                    <>
                                        <div className="flex flex-wrap gap-4">
                                            {Object.entries(grades).map(([id, grade]) => {
                                                const gradeS = Object.values(grade.grades)
                                                const weight = gradeTypes[grade.type_id]?.weight || 1
                                                const average = gradeS.length
                                                    ? gradeS.reduce((sum, grade) => sum + grade * weight, 0) / gradeS.reduce(sum => sum + weight, 0)
                                                    : 0

                                                const color = colors[Math.round(average)]
                                                const rounded = average.toFixed(2)

                                                return (
                                                    <div className="border-(--txt-color) w-80 h-40 flex flex-row justify-between p-4 border-2 rounded-2xl cursor-pointer" onClick={() => setEditingGrade(Number(id))}>
                                                        <div className="flex flex-col">
                                                            <p className="text-3xl rubik">{grade.title}</p>
                                                            <p className="opacity-70">{Object.values(grade.grades).filter(grade => grade !== 0).length}</p>
                                                        </div>
                                                        <div className="flex flex-col justify-end">
                                                            <div className={`rounded-full h-8 w-fit pl-2 pr-2 inline-flex items-center justify-center poppins text-2xl`} style={{
                                                                backgroundColor: `var(--${color}-base-color)`,
                                                                color: `var(--${color}-color)`, fontWeight: "900"
                                                            }}>{Number(rounded) ? rounded : "---"}</div>
                                                        </div>
                                                    </div>
                                                )
                                            })}
                                            <div className="border-(--txt-color) w-80 h-40 flex justify-center items-center border-2 rounded-2xl cursor-pointer text-5xl" onClick={() => {
                                                const id = Math.max(0, ...Object.keys(grades).map(grade => Number(grade))) + 1 // stupid js forces it to be string
                                                const grade: EditingGrade = { title: "", description: "", type_id: 0, grades: Object.fromEntries(students.map(student => [student.id, 0])) }
                                                setGrades({ ...grades, [id]: grade })
                                                setEditingGrade(id)
                                            }}>
                                                +
                                            </div>
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
                    <Button icon={faStar} text={getKey("TIMETABLE_GRADE")} onClick={() => { setActive("grade"); setEditingGrade(null) }} />
                    <Button icon={faStar} text={getKey("TIMETABLE_HOMEWORK")} onClick={() => setActive("homework")} />
                </Sidebar>
            </div>
        </>
    )

    async function save() {
      switch (active) {
        case "overview":
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
          break
        case "grade":
          for (const grade of Object.values(grades)) {
            for (const [student_id, student_grade] of Object.entries(grade.grades)) {
              await fetchWithSchool(fetchMe, "/api/v1/me/teacher/grades", {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify({
                  term_id: lesson.term_id,
                  grade_type_id: grade.type_id,
                  date: lesson.actual_date,
                  value: student_grade,
                  has_note: grade.description != "",
                  note: grade.description,
                  student_id: Number(student_id),
                  class_subjects_id: lesson.subject_id,
                })
              })
            }
          }
          break
        default:
          break
        }
    }
}
