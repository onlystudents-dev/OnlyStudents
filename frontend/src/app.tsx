import {BrowserRouter, Route, Routes} from "react-router";
import Home from "./home.tsx";
import {ToastContainer} from "react-toastify";
import {useEffect, useState} from "react";
import Login from "./auth/login/login.tsx";
import Loading from "./util/loading.tsx";
import GuardianHomeworks from "./ui/guardian/homeworks.tsx";
import Me from "./ui/me/me.tsx";

export type Me = {
    role: string,
    account_id: number,
    first_name: string,
    last_name: string,
    pfp_url: string,
}

export default function App() {
    const [loading, setLoading] = useState(true);
    const [me, setMe] = useState<Me | null>(null);

    useEffect(() => {
        async function Fetch() {
            try {
                const meR = await fetch("/api/v1/me/status");
                const meJ = await meR.json();
                setMe(meJ);
            } catch {/* empty */}
        }

        Fetch().then(() => setLoading(false));
    }, []);

    return (
        <>
            {loading ? <Loading /> : !me ? <Login /> : (
                <BrowserRouter>
                    <Routes>
                        <Route path="/" element={<Home me={me} />} />
                        <Route path="/me" element={<Me me={me} />} />
                        <Route path="/homeworks" element={(() => {
                            switch (me.role) {
                                case "guardian":
                                    return <GuardianHomeworks me={me} />
                                case "student":
                                    return null
                                case "teacher":
                                    return null
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