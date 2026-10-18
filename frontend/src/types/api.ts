export type Role = "student" | "teacher" | "guardian" | "admin"

export type Me = {
    role: Role,
    account_id: number,
    first_name: string,
    last_name: string,
    email_address: string,
    email_verified: boolean,
    preferences: Preferences,
    school_memberships: SchoolMembership[],
    children?: GuardianChildren[],
}

export const timeFormats = ["", "h12", "h23"] as const
export type TimeFormat = typeof timeFormats[number]

export function isTimeFormat(value: string): value is TimeFormat {
    return timeFormats.some(format => format === value)
}

export type Preferences = {
    pfp_url: string,
    nickname: string,
    lang: string,
    theme: string,
    timetable_display: number,
    timetable_next: boolean,
    time_format: TimeFormat,
}

export type SchoolMembership = {
    school_id: number,
    school_name: string,
    class_id: number,
}

export type GuardianChildren = {
    account_id: number,
    school_memberships: SchoolMembership[],
    first_name: string,
    last_name: string,
}

export type Theme = {
    name: string,
    colors: Record<string, string>,
}

export type Lesson = {
    actual_date: number,
    room_id: number,
    lesson_num: number,
    day_of_week: number,
    effective_teacher_id: number,
    group_id: number,
    school_id: number,
    custom_subject: boolean,
    has_subject_id: boolean,
    subject_id: number,
    has_custom_subject_id: boolean,
    custom_subject_id: number,
    subject_name: string,
    is_substitution: boolean,
    canceled: boolean,
    has_exam_id: boolean,
    exam_id: number,
    has_homework_id: boolean,
    homework_id: number,
    has_teacher_first_name: boolean,
    teacher_first_name: string,
    has_teacher_last_name: boolean,
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
    has_lesson_number: boolean,
    lesson_number: number,
    at_start: number,
    at_end: number,
}

export type HomeworkSummary = {
    id: number,
    subject: string,
    has_subject_code: boolean,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    title: string,
    has_description: number,
    description: string,
    due_date: number,
    created_at: number,
    submitted_at: number,
    has_graded_value: boolean,
    graded_value: number,
}

export type HomeworkSubmission = {
    id: number,
    homework_id: number,
    title: string,
    subject: string,
    student_id: number,
    first_name: string,
    last_name: string,
    has_content: boolean,
    content: string,
    created_at: number,
    due_date: number,
}

export type Exam = {
    id: number,
    subject: string,
    has_subject_code: boolean,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    title: string,
    has_description: boolean,
    description: string,
    date: number,
    start_time: number,
    end_time: number,
    has_room: boolean,
    room: string,
}

export type StudentGrades = {
    id:	number,
    subject: string,
    has_subject_code: boolean,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    term: string,
    type: string,
    value: number,
    date: number,
    has_note: boolean,
    note: string,
}