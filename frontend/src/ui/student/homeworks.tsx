import "../homeworks/homeworks.css";
import type {HomeworkSummary, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {useEffect} from "react";
import {fetchWithSchool} from "../../util/api.ts";
import Homeworks, {type PopupProps} from "../homeworks/homeworks.tsx";

export default function StudentHomeworks({ me }: { me: Me }) {
    const [homeworks, setHomeworks] = useState<HomeworkSummary[] | null>(null)

    function Fetch() {
        return fetchWithSchool("/api/v1/me/student/homework").then(async res => setHomeworks(await res.json() as HomeworkSummary[]))
    }

    useEffect(() => {
        void Fetch()
    }, [])

    return <Homeworks me={me} homeworks={homeworks} Popup={Popup} reloadSchoolStuff={Fetch} />
}

function Popup({ homework: hw, clear }: PopupProps) {
    const homework = hw as HomeworkSummary

    return (
        <>
            {/* i like cats */}
        </>
    )
}