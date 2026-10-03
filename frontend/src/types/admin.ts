export type AdminStatusData = {
    db_ping: number,
    redis_ping: string,
    db_query_success: boolean,
    school_count: number,
    account_count: number,
    student_count: number,
    teacher_count: number,
    guardian_count: number,
}

export type Log = {
    name: string
    content: string
}