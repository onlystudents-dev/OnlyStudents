import {Route, Switch} from "wouter";
import Home from "./home.tsx";
import {ToastContainer} from "react-toastify";
import {useEffect, useState, type ReactNode} from "react";
import Login from "./auth/login/login.tsx";
import Loading from "./util/loading.tsx";
import GuardianHomeworks from "./ui/guardian/homeworks.tsx";
import Me from "./ui/me/me.tsx";
import TeacherHomeworks from "./ui/teacher/homeworks.tsx";
import StudentHomeworks from "./ui/student/homeworks.tsx";
import GuardianTimetable from "./ui/guardian/timetable.tsx";
import StudentTimetable from "./ui/student/timetable.tsx";
import TeacherTimetable from "./ui/teacher/timetable.tsx";
import GuardianAbsences from "./ui/guardian/absences.tsx";
import StudentAbsences from "./ui/student/absences.tsx";
import TeacherAbsences from "./ui/teacher/absences.tsx";
import GuardianGrades from "./ui/guardian/grades.tsx";
import StudentGrades from "./ui/student/grades.tsx";
import TeacherGrades from "./ui/teacher/grades.tsx";
import RateLimit from "./util/ratelimit.tsx";
import {fetchLanguage} from "./util/language.ts";
import {applyTheme, getAutoTheme, type Theme} from "./util/theme.ts";

export type Me = {
    role: string,
    account_id: number,
    first_name: string,
    last_name: string,
    pfp_url: string,
    email_address: string,
    email_verified: boolean,
    nickname: string,
    lang: string,
    theme: Theme,
}

declare global {
    interface Window { __INITIAL_STATUS__?: Me | string | null }
}

export default function App() {
    const [ratelimit, setRatelimit] = useState<number>(-1);
    const [loading, setLoading] = useState(true);
    const [me, setMe] = useState<Me | null>(null);

    async function fetchMe() {
        const me = window.__INITIAL_STATUS__
        if (typeof me !== "undefined") {
            if (typeof me === "string") return true
            if (me === null) return false

            window.__INITIAL_STATUS__ = undefined
            setMe(me)
            return me
        }
        try {
            const meR = await fetch("/api/v1/me/status")
            if (meR.status === 429) {
                const seconds = meR.headers.get("Retry-After")
                if (seconds == null) return false
                await fetchLanguage(null)
                setRatelimit(Number(seconds))
                return false
            }
            if (meR.status === 401) return true
            const meJ = await meR.json()
            setMe(meJ)
            return meJ
        } catch {/* empty */}
    }

    useEffect(() => {
        function Fetch() {
            fetchMe()
                .then(me => {
                    if (me) {
                        applyTheme(me)
                        fetchLanguage(me)
                            .then(() => setLoading(false))
                    } else setLoading(false)
                })
        }

        applyTheme(getAutoTheme())
        Fetch()
    }, [])

    const byRole = (guardian: ReactNode, student: ReactNode, teacher: ReactNode) => {
        if (!me) return <Loading />
        switch (me.role) {
            case "guardian":
                return guardian
            case "student":
                return student
            case "teacher":
                return teacher
            default:
                return <Loading />
        }
    }

    return (
        <>
            {ratelimit !== -1 ? <RateLimit retry={ratelimit} expire={() => location.reload()} standalone /> : loading ? <Loading /> : !me ? <Login /> : (
                <Switch>
                    <Route path="/">
                        <Home me={me} />
                    </Route>
                    <Route path="/me">
                        <Me me={me} fetchMe={fetchMe} />
                    </Route>
                    <Route path="/homeworks">
                        {() => byRole(
                            <GuardianHomeworks me={me} />,
                            <StudentHomeworks me={me} />,
                            <TeacherHomeworks me={me} />
                        )}
                    </Route>
                    <Route path="/timetable">
                        {() => byRole(
                            <GuardianTimetable me={me} />,
                            <StudentTimetable me={me} />,
                            <TeacherTimetable me={me} />
                        )}
                    </Route>
                    <Route path="/absences">
                        {() => byRole(
                            <GuardianAbsences me={me} />,
                            <StudentAbsences me={me} />,
                            <TeacherAbsences me={me} />
                        )}
                    </Route>
                    <Route path="/grades">
                        {() => byRole(
                            <GuardianGrades me={me} />,
                            <StudentGrades me={me} />,
                            <TeacherGrades me={me} />
                        )}
                    </Route>
                </Switch>
            )}
            <ToastContainer theme={"dark"} position={"bottom-right"} />
        </>
    )
}
// this thing is held together with duct tape