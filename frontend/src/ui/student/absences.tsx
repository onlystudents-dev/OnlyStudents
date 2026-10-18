import "./absences.css";
import type {Me, Absence} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {fetchWithSchool} from "../../util/api.ts";
import {useEffect} from "react";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {formatUnixDate} from "../../util/time.ts";
import {getLanguage} from "../../util/language.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";

export default function StudentAbsences({ me }: { me: Me }) {
    const [absences, setAbsences] = useState<Absence[] | null>(null)
    const [current, setCurrent] = useState<Absence | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/student/absences").then(async res => setAbsences(await res.json() as Absence[]))
    }

    useEffect(() => void Fetch(), [])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && <>
                    <Overlay onClick={() => setCurrent(null)} />
                    <div className={`popup ${!current.justified && "wrong"}`}>
                        <div className="popup-content">
                            <div className="absolute top-6 left-6 right-16 flex flex-row gap-3">
                                <div className="text-3xl rubik flex-1 truncate">
                                    {current.has_note ? current.note : current.subject}
                                </div>
                                <div className="shrink min-w-0 max-w-[45%]">
                                    {formatUnixDate(current.date, getLanguage(me)?.key || "en-US")}
                                </div>
                            </div>
                            <div className="absolute top-15 left-6 text-2xl opacity-70 rubik max-w-[56vw] truncate"> {/* truncate just in case */}
                                {current.has_note ? current.subject : `${current.teacher_first_name} ${current.teacher_last_name}`}
                            </div>
                            <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={() => setCurrent(null)}>
                                <FontAwesomeIcon icon={faXmark} />
                            </div>
                        </div>
                    </div>
                </>}
                {absences?.map((absence) => (
                    <>
                        <div className={`absence ${!absence.justified && "wrong"}`} onClick={() => setCurrent(absence)}>
                            <div className="flex flex-col items-start justify-between">
                                <p className="max-w-100 truncate">{absence.has_note ? absence.note : absence.subject}</p>
                                <p className="max-w-100 truncate opacity-70">{absence.has_note ? absence.subject : `${absence.teacher_first_name} ${absence.teacher_last_name}`}</p>
                            </div>
                            <div className="flex flex-col items-start justify-between">

                            </div>
                        </div>
                    </>
                )) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                ))}
            </div>
        </>
    )
}