import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import type {Lesson, LessonTime, Me, Room} from "../../types/api"
import { getKey } from "../../util/language"
import {formatSecondsToHourAndMinute} from "../../util/time.ts";
import {faAngleRight, faHouseChimney, faPenToSquare} from "@fortawesome/free-solid-svg-icons";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import "./timetable.css";
import { faAngleLeft } from "@fortawesome/free-solid-svg-icons/faAngleLeft";
import React from "preact/compat";
import {Badge} from "../../util/button/button.tsx";

export const DAY_NAMES = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"] as const

export type DayName = typeof DAY_NAMES[number]

export type Day = {
    name: DayName,
    date: string,
    lessons: Lesson[],
    today: boolean
}

export type PopupProps = {
    lesson: Lesson,
    clear: () => void,
    rooms: Room[],
}

type TimetableProps = {
    me: Me,
    year: string,
    days: Day[],
    rooms: Room[],
    lessonTimes: LessonTime[],
    onPrevWeek: () => void,
    onNextWeek: () => void,
    reloadSchoolStuff: () => void,
    Popup:  React.ComponentType<PopupProps>,
}

export function Timetable(props: TimetableProps) {
    const [current, setCurrent] = React.useState<Lesson | null>(null)

    const d = (() => {
        switch (props.me.preferences.timetable_display) {
            case 1:
                return props.days.filter(day => day.lessons.length > 0)
            case 2:
                return props.days.slice(0, 5)
            default:
                return props.days
        }
    })()

    return (
        <>
            <div className="w-full flex flex-row justify-center pt-4">
                <h1 className="text-5xl rubik">
                    {getKey("SCHOOL_YEAR", props.year)}
                </h1>
            </div>
            <div className="w-full h-fit p-4 gap-4 flex flex-row justify-center items-stretch">
                <div className="flex flex-col justify-center h-14">
                    <span className="icon" onClick={() => props.onPrevWeek()}>
                        <FontAwesomeIcon icon={faAngleLeft} />
                    </span>
                </div>
                {current && <props.Popup lesson={current} clear={() => {setCurrent(null); if (props.me.role === "teacher") props.reloadSchoolStuff()}} rooms={props.rooms} />}
                {d.map((day) => (
                    <div className="day" key={day.name}>
                        <div className={`date ${day.today && "rounded-2xl bg-(--border-color)"}`}>
                            <h1 className="rubik">{getKey(`DAYS.${day.name}`)}</h1>
                            <p className="poppins">{day.date}</p>
                        </div>
                        {day.lessons.map((lesson) => (
                            <div className="lesson" key={lesson.actual_date + lesson.lesson_num} onClick={() => setCurrent(lesson)}>
                                <div className="flex flex-row justify-between items-start w-full gap-2">
                                    <div className="flex flex-col min-w-0">
                                        <h1 className="rubik truncate">{lesson.subject_name}</h1>
                                        <h2 className="poppins truncate">{[
                                            lesson.teacher_first_name != "" && lesson.teacher_first_name,
                                            lesson.teacher_last_name != "" && lesson.teacher_last_name,
                                        ].filter(Boolean).join(" ")}</h2>
                                    </div>
                                    <div className="flex flex-col gap-1 items-end shrink-0">
                                        {lesson.exam_id != 0 && <Badge icon={faPenToSquare} className="bg-(--wrong-base-color) text-(--wrong-color)" />}
                                        {lesson.homework_id != 0 && <Badge icon={faHouseChimney} className="bg-(--warning-base-color) text-(--warning-color)" />}
                                    </div>
                                </div>

                                <div className="flex flex-row justify-between items-end w-full gap-2">
                                    <div className="min-w-0">
                                        <p className="fredoka truncate">{props.rooms.find(room => room.id === lesson.room_id)?.name || <Skeleton width={48} height={16} color={"var(--card-color)"} />}</p>
                                    </div>
                                    <div className="whitespace-nowrap shrink-0 text-right">
                                        <h2>{toHoursAndMinutes(props.me, lesson, props.lessonTimes) || <Skeleton width={96} height={16} color={"var(--hover-color)"} />}</h2>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                ))}
                <div className="flex flex-col justify-center h-14">
                    <span className="icon" onClick={() => props.onNextWeek()}>
                        <FontAwesomeIcon icon={faAngleRight} />
                    </span>
                </div>
            </div>
        </>
    )
}

function toHoursAndMinutes(me: Me, lesson: Lesson, lessonTime: LessonTime[]) {
    const time = lessonTime.find(lt => lt.lesson_number != 0 && lt.lesson_number === lesson.lesson_num)
    if (!time) return ""

    return `${formatSecondsToHourAndMinute(time.at_start, me)}-${formatSecondsToHourAndMinute(time.at_end, me)}`
}
