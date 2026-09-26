import {Route, Switch} from "wouter";
import Home from "./home.tsx";
import {ToastContainer} from "react-toastify";
import {lazy, Suspense, useEffect, useState} from "react";
import Loading from "./util/loading.tsx";
import type { Me, Role } from "./types/api.ts";
import { fetchLanguage } from "./util/language.ts";
import { applyTheme, getAutoTheme } from "./util/theme.ts";
import RateLimit from "./util/ratelimit.tsx";

declare global {
    interface Window { __INITIAL_STATUS__?: Me | string | number }
}

const features = {
  homeworks: {
    student: lazy(() => import("./ui/student/homeworks.tsx")),
    teacher: lazy(() => import("./ui/teacher/homeworks.tsx")),
    guardian: lazy(() => import("./ui/guardian/homeworks.tsx"))
  },
  absences: {
    student: lazy(() => import("./ui/student/absences.tsx")),
    teacher: lazy(() => import("./ui/teacher/absences.tsx")),
    guardian: lazy(() => import("./ui/guardian/absences.tsx"))
  },
  timetable: {
    student: lazy(() => import("./ui/student/timetable.tsx")),
    teacher: lazy(() => import("./ui/teacher/timetable.tsx")),
    guardian: lazy(() => import("./ui/guardian/timetable.tsx"))
  },
  grades: {
    student: lazy(() => import("./ui/student/grades.tsx")),
    teacher: lazy(() => import("./ui/teacher/grades.tsx")),
    guardian: lazy(() => import("./ui/guardian/grades.tsx"))
  },
};

const Login = lazy(() => import("./auth/login/login.tsx"))
const MeSettings = lazy(() => import("./ui/me/me.tsx"))

export default function App() {
    const byFeature = (feature: keyof typeof features, me: Me) => {
        const group = features[feature] as Partial<Record<Role, typeof features.homeworks.student>>
        const Feature = group[me.role]
        if (!Feature) return <p className="poppins">Not available for your role.</p>
        return <Feature me={me} />
    }

    const [ratelimit, setRatelimit] = useState<number>(-1)
    const [loading, setLoading] = useState(true)
    const [me, setMe] = useState<Me | null>(null)

    async function fetchMe() {
        const me = window.__INITIAL_STATUS__
        if (typeof me !== "undefined") {
            if (typeof me === "string") return true
            if (typeof me === "number") {
                await fetchLanguage(null)
                setRatelimit(me)
                return false
            }

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

    return (
        <>
            <Suspense fallback={<Loading />}>
                {ratelimit !== -1 ? <RateLimit retry={ratelimit} expire={() => location.reload()} /> : loading ? <Loading /> : !me ? <Login /> : (
                    <Switch>
                        <Route path="/"><Home me={me} /></Route>
                        <Route path="/me"><MeSettings me={me} fetchMe={fetchMe} /></Route>
                        <Route path="/homeworks">{() => byFeature("homeworks", me)}</Route>
                        <Route path="/timetable">{() => byFeature("timetable", me)}</Route>
                        <Route path="/absences">{() => byFeature("absences", me)}</Route>
                        <Route path="/grades">{() => byFeature("grades", me)}</Route>
                    </Switch>
                )}
            </Suspense>
            <ToastContainer theme={"dark"} position={"bottom-right"} />
        </>
    )
}
// this thing is held together with duct tape