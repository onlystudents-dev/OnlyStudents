import type {Theme} from "../util/theme.ts";

export type Role = "student" | "teacher" | "guardian" | "admin";

export type Me = {
    role: Role,
    account_id: number,
    first_name: string,
    last_name: string,
    email_address: string,
    email_verified: boolean,
    preferences: Preferences,
}

export type Preferences = {
    pfp_url: string,
    nickname: string,
    lang: string,
    theme: Theme,
    timetable_display: number,
    timetable_next: boolean,
}

export type Lesson = {
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

export type Room = {
    id: number,
    name: string,
    capacity: number,
}

export type Class = {
    id: number,
    school_id: number,
    name: string,
    teacher_id: number,
    has_co_teacher_id: boolean,
    co_teacher_id: number,
    has_bell_id: boolean,
    bell_id: number,
}

export type LessonTime = {
    id: number,
    school_id: number,
    type_id: number,
    lesson_number: number,
    at_start: LessonTimeDate,
    at_end: LessonTimeDate,
}

export type LessonTimeDate = {
    Microseconds: number,
    Valid: boolean,
}
