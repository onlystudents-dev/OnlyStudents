import {BrowserRouter, Route, Routes} from "react-router";
import Home from "./home.tsx";
import {ToastContainer} from "react-toastify";
import {useEffect, useState} from "react";
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

export type Me = {
    role: string,
    account_id: number,
    first_name: string,
    last_name: string,
    pfp_url: string,
}

export default function App() {
    const [ratelimit, setRatelimit] = useState<number>(-1);
    const [loading, setLoading] = useState(true);
    const [me, setMe] = useState<Me | null>(null);

    useEffect(() => {
        async function Fetch() {
            try {
                const meR = await fetch("/api/v1/me/status");
                if (meR.status === 429) {
                    const seconds = meR.headers.get("Retry-After")
                    if (seconds == null) return
                    setRatelimit(Number(seconds))
                    return
                }
                const meJ = await meR.json();
                setMe(meJ);
            } catch {/* empty */}
        }

        fetchLanguage().then(() => Fetch().then(() => setLoading(false)))
    }, [])

    return (
        <>
            {ratelimit !== -1 ? <RateLimit retry={ratelimit} /> : loading ? <Loading /> : !me ? <Login /> : (
                <BrowserRouter>
                    <Routes>
                        <Route path="/" element={<Home me={me} />} />
                        <Route path="/me" element={<Me me={me} />} />
                        <Route path="/homeworks" element={(() => {
                            switch (me.role) {
                                case "guardian":
                                    return <GuardianHomeworks me={me} />
                                case "student":
                                    return <StudentHomeworks me={me} />
                                case "teacher":
                                    return <TeacherHomeworks me={me} />
                                default:
                                    return <Loading />
                            }
                        })()} />
                        <Route path="/timetable" element={(() => {
                            switch (me.role) {
                                case "guardian":
                                    return <GuardianTimetable me={me} />
                                case "student":
                                    return <StudentTimetable me={me} />
                                case "teacher":
                                    return <TeacherTimetable me={me} />
                                default:
                                    return <Loading />
                            }
                        })()} />
                        <Route path="/absences" element={(() => {
                            switch (me.role) {
                                case "guardian":
                                    return <GuardianAbsences me={me} />
                                case "student":
                                    return <StudentAbsences me={me} />
                                case "teacher":
                                    return <TeacherAbsences me={me} />
                                default:
                                    return <Loading />
                            }
                        })()} />
                        <Route path="/grades" element={(() => {
                            switch (me.role) {
                                case "guardian":
                                    return <GuardianGrades me={me} />
                                case "student":
                                    return <StudentGrades me={me} />
                                case "teacher":
                                    return <TeacherGrades me={me} />
                                default:
                                    return <Loading />
                            }
                        })()} />
                    </Routes>
                </BrowserRouter>
            )}
            <ToastContainer theme={"dark"} position={"bottom-right"} />
        </>
    )
}
