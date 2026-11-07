import type {Absence, GradeType, Group, Homework, Lesson, LessonTime, Room, Student} from "../../types/api.ts";
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
import {faFloppyDisk, faHouseChimney, faList, faStar, faTrash, faUser} from "@fortawesome/free-solid-svg-icons";
import Sidebar from "../../util/sidebar/sidebar.tsx";
import Button from "../../util/sidebar/config.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {useLayoutEffect} from "preact/hooks";
import {colors} from "../globals.tsx";
import Select from "react-select";
import {selectStyles} from "../../util/select.ts";
import {getDifferences} from "../../util/save/save.tsx";
import SubmissionAttachmentInput, {type attachmentFile} from "../submission.tsx";

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
            <Timetable me={me} year={year} days={days} rooms={rooms} lessonTimes={lessonTime} onPrevWeek={() => fetchWeekLessons(weekOffsetRef.current - 1)} onNextWeek={() => fetchWeekLessons(weekOffsetRef.current + 1)} reloadSchoolStuff={Fetch} Popup={Popup} groups={groups} fetchMe={fetchMe} refresh={() => fetchWeekLessons(weekOffsetRef.current)} />
        </>
    )
}

type absenceChangeType = {"type": "DELETE", "id": number, "student_id": number} | {"type": "ADD", "absence_type": "absent" | "tardy", "student_id": number, "tardy_minutes": number} | {"type": "UPDATE", "absence_type": "absent" | "tardy", "id": number, "student_id": number, "tardy_minutes": number}

function Popup({ lesson, rooms, groups, clear, fetchMe, me, lessonTimes }: PopupProps) {
    const [active, setActive] = useState<"overview" | "grade" | "homework">("overview")
    const [students, setStudents] = useState<Student[]>([])
    const [gradeTypes, setGradeTypes] = useState<GradeType[]>([])

    const [absencesChanges, setAbsencesChanges] = useState<absenceChangeType[]>([])
    const [absences, setAbsences] = useState<Absence[]>([])
    const presentRef = useRef<HTMLButtonElement>(null)
    const absentRef = useRef<HTMLButtonElement>(null)
    const tardyRef = useRef<HTMLButtonElement>(null)
    const [presentWidth, setPresentWidth] = useState<number>(0)
    const [absentWidth, setAbsentWidth] = useState<number>(0)
    const [tardyWidth, setTardyWidth] = useState<number>(0)

    type EditingGrade = { id:number, title: string, description: string, type: number, class_subjects_id: number, actual_date: number, grades: Record<number, number>, defGrades: Record<number, number>, gradeIDs: Record<number, number>, new: boolean }
    const [editingGrade, setEditingGrade] = useState<number | null>(null)
    const [grades, setGrades] = useState<Record<number, EditingGrade>>({})
    const [defaultGrades, setDefaultGrades] = useState<Record<number, EditingGrade>>({})
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

    const [homework, setHomework] = useState<Homework>()
    const [file, setFile] = useState<attachmentFile>()

    const room = rooms.find(r => r.id === lesson.room_id)
    const group = groups?.find(g => g.id === lesson.group_id)

    const AbsenceChangeFor = (student_id: number, absence_type: "present" | "absent" | "tardy", tardy_minutes: number): absenceChangeType | null => {
        const absence = absences.find(a => a.student_id === student_id)

        if (absence_type === "present") {
            return absence
                ? { type: "DELETE", id: absence.id, student_id: student_id }
                : null
        }

        if (!absence) {
            return { type: "ADD", absence_type: absence_type, student_id: student_id, tardy_minutes: tardy_minutes }
        }

        return absence.type === absence_type
            ? null
            : { type: "UPDATE", id: absence.id, absence_type: absence_type, student_id: student_id, tardy_minutes: tardy_minutes }
    }

    const ChangeAbsence = (student_id: number, absence_type: "present" | "absent" | "tardy", tardy_minutes: number) => {
        setAbsencesChanges(prev => {
            const withoutStudent = prev.filter(absence_change => absence_change.student_id != student_id)
            const change = AbsenceChangeFor(student_id, absence_type, tardy_minutes)
            return change ? [...withoutStudent, change] : withoutStudent
        })
    };

    const getTardyValue = (student_id: number): number => {
        const change = absencesChanges.find(absence_change => absence_change.student_id === student_id)

        if (change) {
            switch (change.type) {
                case "ADD":
                    return change.tardy_minutes
                case "DELETE":
                    return 1
                case "UPDATE":
                    return change.tardy_minutes
            }
        }

        const absence = absences.find(absence => absence.student_id === student_id)

        if (absence) {
            return absence.tardy_minutes
        }

        return 1
    }

    const AbsenceTypeFor = (student_id: number): "present" | "absent" | "tardy" => {
        const change = absencesChanges.find(absence_change => absence_change.student_id === student_id)

        if (change) {
            return change.type === "DELETE" ? "present" : change.absence_type
        }

        const absence = absences.find(absence => absence.student_id === student_id)

        if (absence?.type === "absent") return "absent"
        if (absence?.type === "tardy") return "tardy"
        return "present"
    }

    const AbsencePillStyle = (student_id: number) => {
        const absence_type = AbsenceTypeFor(student_id)
        const left = absence_type === "present"
            ? 0
            : absence_type === "absent"
                ? presentWidth
                : presentWidth + absentWidth
        const width = absence_type === "present"
            ? presentWidth
            : absence_type === "absent"
                ? absentWidth
                : tardyWidth

        return {
            left: `${left}px`,
            width: `${width}px`,
            borderTopLeftRadius: absence_type === "present" ? '6px' : '0px',
            borderBottomLeftRadius: absence_type === "present" ? '6px' : '0px',
            borderTopRightRadius: absence_type === "tardy" ? '6px' : '0px',
            borderBottomRightRadius: absence_type === "tardy" ? '6px' : '0px',
        }
    }

    useEffect(() => {
        void fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/students_in_my_group?group_id=${group?.id}`).then(async students_response => { const json = await students_response.json() as Student[]; setStudents(json); return json }).then(students => {
            void fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/grades?actual_date=${lesson.actual_date}&subject_id=${lesson.subject_id}`).then(async response => {
                const grades = await response.json() as EditingGrade[]
                let gradeS: Record<number, EditingGrade> = []

                let i = 0
                for (let grade of grades) {
                    const res = await fetchWithSchool(fetchMe, `/api/v1/me/teacher/grades?lesson_id=${grade.id}`)
                    const json = await res.json() as { student_id: number, value: number, id: number }[]
                    const gradeIDs = Object.fromEntries(json.map(obj => [obj.student_id, obj.id]))
                    const gradES = Object.fromEntries(json.map(obj => [obj.student_id, obj.value]))
                    const graDES = Object.fromEntries(json.map(obj => [obj.student_id, obj.value]))
                    students.forEach(student => {
                        if (!(student.id in gradES)) {
                            gradES[student.id] = 0
                        }
                        if (!(student.id in graDES)) {
                            graDES[student.id] = -1
                        }
                    })
                    grade = { ...grade, grades: gradES, defGrades: graDES, gradeIDs, new: false }

                    gradeS = { ...gradeS, [i]: grade }
                    i++
                }

                setGrades(gradeS)
                setDefaultGrades(gradeS)
            })
        })
        void fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/absences/lesson?class_subjects_id=${lesson.subject_id}?date=${lesson.actual_date}`).then(async absences_response => setAbsences(await absences_response.json() as Absence[]))
        void fetchWithSchool(fetchMe, "/api/v1/me/teacher/grades/types").then(async response => setGradeTypes(await response.json() as GradeType[]))
        if (lesson.has_homework_id) void fetchWithSchool(fetchMe, `/api/v1/me/teacher/homework/get?homework_id=${lesson.homework_id}`).then(async response => setHomework(await response.json() as Homework))
    }, [fetchMe, group, lesson.actual_date, lesson.has_homework_id, lesson.homework_id, lesson.subject_id])

    useLayoutEffect(() => {
        if (!students.length) return

        setPresentWidth(presentRef.current?.offsetWidth || 0)
        setAbsentWidth(absentRef.current?.offsetWidth || 0)
        setTardyWidth(tardyRef.current?.offsetWidth || 0)
    }, [students])

    useLayoutEffect(() => {
        if (active !== "grade" || typeof editingGrade != "number") return

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
                                                    <div className="flex flex-row gap-2">
                                                        {AbsenceTypeFor(student.id) === "tardy" &&
                                                            <input className="text-center p-1 cursor-text z-1 w-16 rounded-md border border-(--border-color) bg-(--card-color) outline-none [appearance:textfield]"
                                                                   min={1} max={60} type="number" placeholder={String(getTardyValue(student.id))} onChange={(e) => ChangeAbsence(student.id, "tardy", Number(e.currentTarget.value))} />
                                                        }
                                                        <div className="flex flex-row items-center h-full justify-end rounded-lg border-2 border-(--border-color) relative">
                                                            <button className="text-center p-1 cursor-pointer z-1" onClick={() => ChangeAbsence(student.id, "present", 0)} ref={presentRef}>
                                                                {getKey("PRESENT")}
                                                            </button>
                                                            <button className="text-center p-1 cursor-pointer z-1" onClick={() => ChangeAbsence(student.id, "absent", 0)} ref={absentRef}>
                                                                {getKey("ABSENT")}
                                                            </button>
                                                            <button className="text-center p-1 cursor-pointer z-1" onClick={() => ChangeAbsence(student.id, "tardy", getTardyValue(student.id))} ref={tardyRef}>
                                                                {getKey("TARDY")}
                                                            </button>
                                                            <div className="bg-(--card-color) absolute top-0 bottom-0 transition-[left_200ms_ease-in-out]" style={AbsencePillStyle(student.id)} />
                                                        </div>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </>
                                )
                            case "grade":
                                if (typeof editingGrade === "number") {
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
                                                    <Select options={gradeTypesOptions} value={gradeTypesOptions[eGrade.type]} onChange={val => val && setGrades({ ...grades, [editingGrade]: { ...eGrade, type: val.value } })}
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
                                            {Object.entries(grades).map(([i, grade]) => {
                                                const gradeS = Object.values(grade.grades)
                                                const weight = gradeTypes[grade.type]?.weight || 1
                                                const average = gradeS.length
                                                    ? gradeS.reduce((sum, grade) => sum + grade * weight, 0) / gradeS.reduce(sum => sum + weight, 0)
                                                    : 0

                                                const color = colors[Math.round(average)]
                                                const rounded = average.toFixed(2)

                                                return (
                                                    <div className="border-(--txt-color) w-80 h-40 flex flex-row justify-between p-4 border-2 rounded-2xl cursor-pointer" onClick={() => setEditingGrade(Number(i))}>
                                                        <div className="flex flex-col">
                                                            <p className="text-3xl rubik">{grade.title}</p>
                                                            <p className="opacity-70">{Object.values(grade.grades).filter(grade => grade !== 0).length}</p>
                                                        </div>
                                                        <div className="flex flex-col justify-between items-end">
                                                            <FontAwesomeIcon icon={faTrash} className="action-button hover:text-(--wrong-color)" onClick={e => { e.stopPropagation(); setGrades(Object.fromEntries(Object.entries(grades).filter(([key]) => key !== i))) }} /> {/* ofc js doesnt have remove why would it */}
                                                            <div className={`rounded-full h-8 w-fit pl-2 pr-2 inline-flex items-center justify-center poppins text-2xl`} style={{
                                                                backgroundColor: `var(--${color}-base-color)`,
                                                                color: `var(--${color}-color)`, fontWeight: "900"
                                                            }}>{Number(rounded) ? rounded : "---"}</div>
                                                        </div>
                                                    </div>
                                                )
                                            })}
                                            <div className="border-(--txt-color) w-80 h-40 flex justify-center items-center border-2 rounded-2xl cursor-pointer text-5xl" onClick={() => {
                                                const id = Math.max(-1, ...Object.keys(grades).map(grade => Number(grade))) + 1 // stupid js forces it to be string
                                                const grade: EditingGrade = { id: -1, title: "", description: "", type: 0, class_subjects_id: 0, actual_date: 0, grades: Object.fromEntries(students.map(student => [student.id, 0])), defGrades: Object.fromEntries(students.map(student => [student.id, -1])), gradeIDs: {}, new: true }
                                                setGrades({ ...grades, [id]: grade })
                                                setEditingGrade(id)
                                            }}>
                                                +
                                            </div>
                                        </div>
                                    </>
                                )
                            case "homework":
                                if (!homework) return (
                                    <>
                                        <div className="border-(--txt-color) w-80 h-40 flex justify-center items-center border-2 rounded-2xl cursor-pointer text-5xl"
                                             onClick={() => setHomework({ id: 0, due_date: new Date().getTime(), has_description: false, description: "", class_subjects_id: lesson.subject_id, subject: lesson.subject_name, title: "", created_at: new Date().getTime() })}>
                                            +
                                        </div>
                                    </>
                                )

                                return (
                                    <>
                                        <div className="flex flex-col justify-between gap-4">
                                            <div className="flex flex-row justify-between">
                                                <div className="flex flex-col">
                                                    <input type="text" className="text-3xl truncate outline-none rubik" defaultValue={homework.title} onChange={e => setHomework({ ...homework, title: e.currentTarget.value })} placeholder={getKey("TITLE")} />
                                                    <textarea rows={2} className="text-2xl opacity-70 truncate fredoka outline-none rubik" defaultValue={homework.description} onChange={e => setHomework({ ...homework, description: e.currentTarget.value, has_description: !!e.currentTarget.value })} placeholder={getKey("DESCRIPTION")} />
                                                </div>
                                                <div className="flex flex-col items-end">
                                                    <input type="date" defaultValue={new Date().toISOString().split("T")[0]} className="text-2xl outline-none rubik" onChange={e => setHomework({ ...homework, due_date: Math.floor(new Date(e.currentTarget.value).getTime() / 1000) })} />
                                                </div>
                                            </div>
                                            <SubmissionAttachmentInput onFile={file => setFile(file)} noLine />
                                        </div>
                                        <button className="absolute bottom-4 left-4 bg-(--card-color) p-2 rounded-lg cursor-pointer action-button hover:text-(--wrong-color)" onClick={() => setHomework(undefined)}>
                                            <FontAwesomeIcon icon={faTrash} /> {getKey("DELETE")}
                                        </button>
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
                    <Button icon={faHouseChimney} text={getKey("TIMETABLE_HOMEWORK")} onClick={() => setActive("homework")} />
                </Sidebar>
            </div>
        </>
    )

    async function save() {
        if (absencesChanges.length > 0) {
            const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/absences", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    absence_changes: absencesChanges,
                    class_subjects_id: lesson.subject_id,
                    date: lesson.actual_date
                }),
            })

            if (!response.ok) {
                toast.error(await fromResponse(response))
                return
            }
        }

        const diff = getDifferences(grades, defaultGrades)
        let j = -1
        for (const grade of Object.values(diff)) {
            j++
            if (!grade) {
                const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/grades", {
                    method: "DELETE",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        lesson_grade_id: defaultGrades[j]?.id
                    })
                })

                if (!response.ok) {
                    toast.error(await fromResponse(response))
                    return
                }

                continue
            }

            const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/grades", {
                method: grade.id === -1 ? "POST" : "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    id: grade.id,
                    actual_date: lesson.actual_date,
                    group_id: lesson.group_id,
                    class_subjects_id: lesson.subject_id,
                    title: grade.title,
                    description: grade.description,
                    type: grade.type,
                })
            })

            if (!response.ok) {
                toast.error(await fromResponse(response))
                return
            }
        }

        const lesson_grade_response = await fetchWithSchool(fetchMe, `/api/v1/me/teacher/timetable/grades?actual_date=${lesson.actual_date}&subject_id=${lesson.subject_id}`)
        if (!lesson_grade_response.ok) {
            toast.error(await fromResponse(lesson_grade_response))
            return
        }

        const json = await lesson_grade_response.json() as { id: number }[]
        let i = 0
        for (const grade of Object.values(grades)) {
            const id = json[i]?.id
            if (id === undefined) continue

            for (const [student_id, student_grade] of Object.entries(getDifferences(grade.grades, grade.defGrades))) {
                const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/grades", {
                    method: student_grade === 0 ? "DELETE" : grade.defGrades[Number(student_id)] !== -1 ? "PATCH" : "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        grade_id: grade.gradeIDs[Number(student_id)],
                        lesson_grade_id: id,
                        term_id: lesson.term_id,
                        date: lesson.actual_date,
                        value: student_grade,
                        has_note: grade.description != "",
                        note: grade.description,
                        student_id: Number(student_id),
                    })
                })
                if (!response.ok) {
                    toast.error(await fromResponse(response))
                    return
                }
            }
            i++
        }

        const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/homework", {
            method: homework ? homework.id === 0 ? "POST" : "PATCH" : "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                homework_id: homework?.id,
                title: homework?.title,
                has_description: !!homework?.description,
                description: homework?.description,
                due_date: homework?.due_date,
                has_assigned_date: !!homework?.created_at,
                assigned_date: homework?.created_at,
                class_subjects_id: lesson.subject_id,
                file_uuid: file?.file_uuid,
            })
        })
        if (!response.ok) {
            toast.error(await fromResponse(response))
            return
        }

        toast.success("Successfully Saved!")
        clear()
    }
}
// so long