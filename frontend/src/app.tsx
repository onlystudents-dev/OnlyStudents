import {Route, Switch} from "wouter";
import Home from "./home.tsx";
import {ToastContainer} from "react-toastify";
import React, {lazy, Suspense, useEffect, useRef, useState} from "react";
import Loading from "./util/loading.tsx";
import type { Me, Role } from "./types/api.ts";
import { fetchLanguage } from "./util/language.ts";
import RateLimit from "./util/ratelimit.tsx";
import { readJSON } from "./util/api.ts";
import type {FeatureProps} from "./types/props.ts";

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
}

const Login = lazy(() => import("./auth/login/login.tsx"))
const MeSettings = lazy(() => import("./ui/me/me.tsx"))

export default function App({ reload }: {reload: () => void}) {
    const [popupLogin, setPopupLogin] = useState(false)
    const loginResolver = useRef<((me: Me | null) => void) | null>(null)

    const byFeature = (feature: keyof typeof features, me: Me) => {
        const group = features[feature] as Partial<Record<Role, React.ComponentType<FeatureProps>>>
        const Feature = group[me.role]
        if (!Feature) return <p className="poppins">Not available for your role.</p>
        return <Feature me={me} unauthorized={unauthorized} />
    }

    function unauthorized() {
        setPopupLogin(true)

        return new Promise<Me | null>(resolve => {
            loginResolver.current = resolve
        })
    }

    async function popupLoginF() {
        const me = await fetchMe()
        setPopupLogin(false)
        loginResolver.current?.(me)
        loginResolver.current = null
        return me
    }

    const [ratelimit, setRatelimit] = useState<number>(-1)
    const [loading, setLoading] = useState(true)
    const [me, setMe] = useState<Me | null>(null)

    async function fetchMe(): Promise<Me | null> {
        const initial = window.__INITIAL_STATUS__
        if (initial !== undefined) {
            window.__INITIAL_STATUS__ = undefined
            if (typeof initial === "number") setRatelimit(initial)
            if (typeof initial !== "object") return null

            setMe(initial)
            return initial
        }
        try {
            const response = await fetch("/api/v1/me/status")
            if (response.status === 429) {
                const seconds = response.headers.get("Retry-After")
                if (seconds != null) setRatelimit(Number(seconds))
                return null
            }

            const me = response.status === 401 ? null : await readJSON<Me>(response)
            setMe(me)
            return me
        } catch {
            return null
        }
    }

    useEffect(() => {
        const path = window.location.pathname

        function Fetch() {
            void fetchMe()
                .then(me => fetchLanguage(me))
                .then(() => setLoading(false))
        }

        switch (path) {
            case "/admin":
                void fetchLanguage(null)
                    .then(() => setLoading(false))
                break
            default:
                Fetch()
        }
    }, [])

    return (
        <>
            <Suspense fallback={<Loading />}>
                <Switch>
                    <Route path="/admin">
                        {loading ? <Loading /> : lazy(() => import("./admin/admin.tsx"))}
                    </Route>
                    {ratelimit !== -1 ? <RateLimit retry={ratelimit} expire={async () => {setRatelimit(-1); await fetchMe()}} /> : loading ? <Loading /> : !me ? <Login /> : (
                        <>
                            <Route path="/"><Home me={me} /></Route>
                            <Route path="/me"><MeSettings me={me} fetchMe={fetchMe} setMe={setMe} reload={reload} /></Route>
                            <Route path="/homeworks">{byFeature("homeworks", me)}</Route>
                            <Route path="/timetable">{byFeature("timetable", me)}</Route>
                            <Route path="/absences">{byFeature("absences", me)}</Route>
                            <Route path="/grades">{byFeature("grades", me)}</Route>
                        </>
                    )}
                </Switch>
            </Suspense>

            {me && popupLogin && (
                <Suspense fallback={null}>
                    <Login fetchMe={popupLoginF} />
                </Suspense>
            )}

            <ToastContainer theme={"dark"} position={"bottom-right"} />
        </>
    )
}
// this thing is held together with duct tape
