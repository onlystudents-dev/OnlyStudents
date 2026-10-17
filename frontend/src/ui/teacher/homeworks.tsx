import "../homeworks/homeworks.css";
import type {HomeworkSubmission, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {useEffect} from "react";
import {fetchWithSchool} from "../../util/api.ts";
import Homeworks, {type PopupProps} from "../homeworks/homeworks.tsx";

export default function TeacherHomeworks({ me }: { me: Me }) {
    const [homeworks, setHomeworks] = useState<HomeworkSubmission[] | null>(null)

    function Fetch() {
        return fetchWithSchool("/api/v1/me/teacher/homework").then(async res => setHomeworks(await res.json() as HomeworkSubmission[]))
    }

    useEffect(() => {
        void Fetch()
    }, [])

    return <Homeworks me={me} homeworks={homeworks} Popup={Popup} reloadSchoolStuff={Fetch} />
}

function Popup({ homework: hw, clear }: PopupProps) {
    const homework = hw as HomeworkSubmission

    return (
        <>
            {/* i like cats */}
        </>
    )
}