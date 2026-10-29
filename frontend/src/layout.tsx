import {Fragment, Suspense, lazy, useEffect} from "preact/compat";
import type {ComponentType} from "preact/compat";
import type {Dispatch, StateUpdater} from "preact/hooks";
import {useLocation} from "wouter";
import type {Me, Role} from "./types/api.ts";
import type {FeatureProps, FetchMe} from "./types/props.ts";
import Navbar from "./navbar/navbar.tsx";
import Home from "./home.tsx";
import Loading from "./util/loading.tsx";
import {features} from "./features.ts";

const MeSettings = lazy(() => import("./ui/me/me.tsx"))

export type LayoutProps = {
    me: Me,
    fetchMe: FetchMe,
    setMe: Dispatch<StateUpdater<Me | null>>,
    reload: () => void,
}

export default function Layout({ me, fetchMe, setMe, reload }: LayoutProps) {
    const [path, navigate] = useLocation()

    useEffect(() => {
        window.scrollTo(0, 0)
    }, [path])

    const byFeature = (feature: keyof typeof features) => {
        const group = features[feature] as Partial<Record<Role, ComponentType<FeatureProps>>>
        const Feature = group[me.role]
        if (!Feature) return <p className="poppins">Not available for your role.</p>
        return <Feature me={me} fetchMe={fetchMe} />
    }

    function page() {
        switch (path) {
            case "/": return <Home me={me} />
            case "/me": return <MeSettings me={me} fetchMe={fetchMe} setMe={setMe} reload={reload} />
            case "/homework": return byFeature("homework")
            case "/timetable": return byFeature("timetable")
            case "/absences": return byFeature("absences")
            case "/grades": return byFeature("grades")
            case "/exams": return byFeature("exams")
            default: return null
        }
    }

    return (
        <>
            <Navbar me={me} path={path} navigate={navigate} />
            <Suspense fallback={<Loading />}>
                <Fragment key={path}>{page()}</Fragment>
            </Suspense>
        </>
    )
}
