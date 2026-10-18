import type {HomeworkSubmission, HomeworkSummary, Me} from "../../types/api.ts";
import React from "react";
import {formatUnixDate} from "../../util/time.ts";
import {getKey, getLanguage} from "../../util/language.ts";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {useState} from "preact/hooks";
import {useReloadSchoolStuff} from "../../util/school.ts";

export type PopupProps = {
    homework: HomeworkSummary | HomeworkSubmission,
    clear: () => void,
}

export default function Homeworks({ me, homeworks, Popup, reloadSchoolStuff }: { me: Me, homeworks: HomeworkSummary[] | HomeworkSubmission[] | null, Popup: React.ComponentType<PopupProps>, reloadSchoolStuff: () => void }) {
    const [current, setCurrent] = React.useState<HomeworkSummary | HomeworkSubmission | null>(null)

    const [now] = useState(Date.now)

    useReloadSchoolStuff(reloadSchoolStuff)

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && <Popup homework={current} clear={() => {setCurrent(null); reloadSchoolStuff()}} />}
                {homeworks?.length === 0 ? <p className="text-4xl">{getKey(`NO_HOMEWORK.${me.role}`)}</p> : homeworks?.map((homework) => (
                    <>
                        <div className={`homework ${homework.due_date < now && "wrong"}`} onClick={() => setCurrent(homework)}>
                            <div className="flex flex-col items-start justify-between">
                                <p className="max-w-100 truncate">{homework.title}</p>
                                <p className="max-w-100 truncate opacity-70">{homework.subject}</p>
                            </div>
                            <div className="flex flex-col items-start justify-between">
                                <p>{formatUnixDate(me.role === "teacher" ? homework.created_at : homework.due_date, getLanguage(me)?.key || "en-US")}</p>
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