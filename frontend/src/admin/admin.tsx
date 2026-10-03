import AdminLogin from "./auth/login.tsx";
import {useEffect} from "react";
import {useState} from "preact/hooks";
import Loading from "../util/loading.tsx";
import {type AdminStatusData} from "../types/admin.ts";
import {fromResponse} from "../util/language.ts";
import {toast} from "react-toastify";
import AdminPanel from "./panel/panel.tsx";

export default function Admin() {
    const [loading, setLoading] = useState(true)

    const [status, setStatus] = useState<AdminStatusData | null>(null)

    useEffect(() => {
        void fetch("/api/v1/admin/status").then(async response => {
            switch (response.status) {
                case 200: {
                    const json: AdminStatusData = await response.json()
                    setStatus(json)
                    break
                }
                case 401: break
                default:
                    toast.error(await fromResponse(response))
            }
            setLoading(false)
        })
    }, [])

    return (
        <>
            {loading ? <Loading /> : status ? <AdminPanel status={status} /> : <AdminLogin />}
        </>
    )
}