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
    has_teacher_first_name: boolean,
    has_teacher_last_name: boolean,
    teacher_first_name: string,
    teacher_last_name: string,
    group_id: number,
    school_id: number,
    custom_subject: boolean,
    has_subject_id: boolean,
    subject_id: number,
    has_custom_subject_id: boolean,
    custom_subject_id: number,
    subject_name: string,
    has_exam_id: boolean,
    exam_id: number,
    has_homework_id: boolean,
    homework_id: number,
    is_substitution: boolean,
    canceled: boolean,
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
    co_teacher_id: number,
    bell_id: number,
}

export type LessonTime = {
    id: number,
    school_id: number,
    type_id: number,
    lesson_number: number,
    at_start: number,
    at_end: number,
}

export type HomeworkSummary = {
    id: number,
    subject: string,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    title: string,
    description: string,
    due_date: number,
    created_at: number,
    submitted_at: number,
    graded_value: number,
    submission_type: string,
    submission_note: string,
    submission_file_uuid: string
    submission_file_name: string,
    submission_file_type: string,
    submission_file_size: number,
}

export type HomeworkSubmission = {
    id: number,
    homework_id: number,
    title: string,
    subject: string,
    student_id: number,
    first_name: string,
    last_name: string,
    file_name: string,
    file_type: string,
    file_size: number,
    file_uuid: string,
    note: string,
    created_at: number,
    due_date: number,
}

export type Exam = {
    id: number,
    subject: string,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    title: string,
    description: string,
    date: number,
    start_time: number,
    end_time: number,
    room: string,
}

export type Grade = {
    id:	number,
    subject: string,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    term: string,
    type: string,
    weight: number,
    value: number,
    date: number,
    note: string,
}

export type Absence = {
    id: number,
    subject: string,
    subject_code: string,
    teacher_last_name: string,
    teacher_first_name: string,
    date: number,
    type: string,
    justified: boolean,
    file_name: string,
    file_size: number,
    file_type: string,
    file_uuid: string,
    reason_note: string,
    verified_by: string,
    class_subjects_id: number,
}

export type Subject = {
    id: number,
    name: string,
    is_custom: boolean,
    teacher_id: number,
    teacher_first_name: string,
    teacher_last_name: string,
    group_id: number,
    group_name: string,
}

export type ForgetPasswordCheckResponse = {
    wrapped_mk_recovery: string
    recovery_key_hash: string
    recovery_key_salt: string
    recovery_key_params: string
}

export type ForgetPasswordConfirmResponse = {
    enroll_token: string
    crypto_grant_handle: string
}

export type Group = {
    id: number,
    school_id: number,
    bell_id: number,
    group_name: string,
}

export type Student = {
    id: number,
    first_name: string,
    last_name: string,
}