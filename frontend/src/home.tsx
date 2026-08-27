import Login from "./auth/login/login.tsx";
import {useEffect, useState} from "react";
import Loading from "./util/loading/loading.tsx";
import Welcome from "./ui/welcome.tsx";

export type Me = {
    role: string,
    account_id: number,
}

export default function Home(){
    const [loading, setLoading] = useState(true);
    const [me, setMe] = useState<Me | null>(null);

    useEffect(() => {
        async function Fetch() {
            try {
                const meR = await fetch("/api/v1/me");
                const meJ = await meR.json();
                setMe(meJ);
            } catch {/* empty */}
        }

        Fetch().then(() => setLoading(false));
    }, []);

    return (
        <>
            {loading ? <Loading /> : !me ? <Login /> : <Welcome me={me} />}
        </>
    )
}