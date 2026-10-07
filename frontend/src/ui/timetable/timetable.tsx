import type { IconDefinition } from "@fortawesome/fontawesome-svg-core"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import type { Lesson, LessonTime, Me, Room } from "../../types/api"
import { getKey } from "../../util/language"
import {formatSecondsToHourAndMinute} from "../../util/time.ts";
import {faAngleRight, faHouseChimney, faPenToSquare} from "@fortawesome/free-solid-svg-icons";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import "./timetable.css";
import { faAngleLeft } from "@fortawesome/free-solid-svg-icons/faAngleLeft";

export const DAY_NAMES = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"] as const

export type DayName = typeof DAY_NAMES[number]

export type Day = {
    name: DayName,
    date: string,
    lessons: Lesson[],
    today: boolean
}

type TimetableProps = {
    me: Me
    year: string
    days: Day[]
    rooms: Room[]
    lessonTimes: LessonTime[]
    onPrevWeek: () => void
    onNextWeek: () => void
}

export function Timetable(props: TimetableProps) {
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

  return (<>
    <div className="w-full flex flex-row justify-center pt-4">
      <h1 className="text-5xl rubik">
        {props.year} {getKey("TERM")}
      </h1>
    </div>
    <div className="w-full h-fit p-4 gap-4 flex flex-row justify-center items-stretch">
      <div className="flex flex-col justify-center h-14">
        <span className="icon" onClick={() => props.onPrevWeek()}>
          <FontAwesomeIcon icon={faAngleLeft} />
        </span>
      </div>
      {d.map((day) => (
          <div className="day" key={day.name}>
              <div className={`date ${day.today && "rounded-2xl bg-(--border-color)"}`}>
                  <h1 className="rubik">{getKey(`DAYS.${day.name}`)}</h1>
                  <p className="poppins">{day.date}</p>
              </div>
              {day.lessons.map((lesson) => (
                  <div className="lesson" key={lesson.actual_date + lesson.lesson_num}>
                      <div className="flex flex-row justify-between items-start w-full gap-2">
                          <div className="flex flex-col min-w-0">
                              <h1 className="rubik truncate">{lesson.subject_name}</h1>
                              <h2 className="poppins truncate">{[
                                  lesson.has_teacher_first_name && lesson.teacher_first_name,
                                  lesson.has_teacher_last_name && lesson.teacher_last_name,
                              ].filter(Boolean).join(" ")}</h2>
                          </div>
                          <div className="flex flex-col gap-1 items-end shrink-0">
                              {lesson.has_exam && <Badge icon={faPenToSquare} className="bg-(--wrong-base-color) text-(--wrong-color)" />}
                              {lesson.has_homework && <Badge icon={faHouseChimney} className="bg-(--warning-base-color) text-(--warning-color)" />}
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
    const time = lessonTime.find(lt => lt.has_lesson_number && lt.lesson_number === lesson.lesson_num)
    if (!time) return ""

    return `${formatSecondsToHourAndMinute(time.at_start, me)}-${formatSecondsToHourAndMinute(time.at_end, me)}`
}

function Badge({ icon, className }: {icon: IconDefinition, className: string}) {
    return (
        <span className={`rounded-full size-8 inline-flex items-center justify-center ${className}`}>
            <FontAwesomeIcon icon={icon} />
        </span>
    )
}
