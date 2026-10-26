import {lazy} from "preact/compat";

export const features = {
    timetable: {
        student: lazy(() => import("./ui/student/timetable.tsx")),
        teacher: lazy(() => import("./ui/teacher/timetable.tsx")),
        guardian: lazy(() => import("./ui/guardian/timetable.tsx"))
    },
    grades: {
        student: lazy(() => import("./ui/student/grades.tsx")),
        guardian: lazy(() => import("./ui/guardian/grades.tsx"))
    },
    exams: {
        student: lazy(() => import("./ui/student/exams.tsx")),
        guardian: lazy(() => import("./ui/guardian/exams.tsx"))
    },
    homework: {
        student: lazy(() => import("./ui/student/homework.tsx")),
        teacher: lazy(() => import("./ui/teacher/homework.tsx")),
        guardian: lazy(() => import("./ui/guardian/homework.tsx"))
    },
    absences: {
        student: lazy(() => import("./ui/student/absences.tsx")),
        guardian: lazy(() => import("./ui/guardian/absences.tsx"))
    },
}
