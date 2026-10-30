import {Route, Switch} from "wouter";
import {ToastContainer} from "react-toastify";
import {lazy, Suspense, useEffect, useState} from "preact/compat";
import Loading from "./util/loading.tsx";
import type { Me } from "./types/api.ts";
import { fetchLanguage } from "./util/language.ts";
import {applyTheme} from "./util/theme.ts";
import Layout from "./layout.tsx";

declare global {
    interface Window { __INITIAL_STATUS__?: Me | string | number }
}

const Login = lazy(() => import("./auth/login/login.tsx"))
const RateLimit = lazy(() => import("./util/ratelimit.tsx"))

export default function App({ reload }: {reload: () => void}) {

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
            localStorage.setItem("user", `${initial.account_id}${initial.role}`)
            return initial
        }
        try {
            const response = await fetch("/api/v1/me/status")
            if (response.status === 429) {
                const seconds = response.headers.get("Retry-After")
                if (seconds != null) setRatelimit(Number(seconds))
                return null
            }

            const me = response.status === 401 ? null : await response.json() as Me
            setMe(me)
            if (me) localStorage.setItem("user", `${me.account_id}${me.role}`)
            return me
        } catch {
            return null
        }
    }

    useEffect(() => {
        const path = window.location.pathname

        switch (path) {
            case "/admin":
                void fetchLanguage(null)
                    .then(() => setLoading(false))
                break
            default:
                void fetchMe()
                    .then(me => fetchLanguage(me))
                    .then(() => setLoading(false))
        }
    }, [])

    return (
        <>
            <Suspense fallback={<Loading />}>
                <Switch>
                    <Route path="/admin">
                        {loading ? <Loading /> : lazy(() => import("./admin/admin.tsx"))}
                    </Route>
                    <Route path="/enroll">
                        {loading ? <Loading /> : lazy(() => import("./auth/enroll/enroll.tsx"))}
                    </Route>
                    {ratelimit !== -1 ? <RateLimit retry={ratelimit} expire={async () => {setRatelimit(-1); await fetchMe().then(me => {if (me !== null) applyTheme(me)})}} /> : loading ? <Loading /> : !me ? <Login /> : (
                        <Layout me={me} fetchMe={fetchMe} setMe={setMe} reload={reload} />
                    )}
                </Switch>
            </Suspense>

            <ToastContainer theme={"dark"} position={"bottom-right"} />
        </>
    )
}
// this thing is held together with duct tape
