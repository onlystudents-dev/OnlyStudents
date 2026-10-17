import Navbar from "../../navbar/navbar.tsx";
import type {HomeworkSubmission, HomeworkSummary, Me} from "../../types/api.ts";
import React from "react";
import {formatUnixDate} from "../../util/time.ts";
import {getLanguage} from "../../util/language.ts";

export type PopupProps = {
    homework: HomeworkSummary | HomeworkSubmission,
    clear: () => void,
}

export default function Homeworks({ me, homeworks, Popup, reloadSchoolStuff }: { me: Me, homeworks: HomeworkSummary[] | HomeworkSubmission[] | null, Popup: React.ComponentType<PopupProps>, reloadSchoolStuff: () => void }) {
    const [current, setCurrent] = React.useState<HomeworkSummary | HomeworkSubmission | null>(null)

    return (
        <>
            <Navbar me={me} reloadSchoolStuff={reloadSchoolStuff} />
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && <Popup homework={current} clear={() => setCurrent(null)} />}
                {homeworks?.map((homework) => (
                    <>
                        <div className="homework" onClick={() => setCurrent(homework)}>
                            <div className="flex flex-col items-start justify-between">
                                <p>{homework.title}</p>
                                <p className="opacity-70">{homework.subject}</p>
                            </div>
                            <div className="flex flex-col items-start justify-between">
                                <p>{formatUnixDate(homework.created_at, getLanguage(me)?.key || "en-US")}</p>
                            </div>
                        </div>
                    </>
                ))}
            </div>
        </>
    )
}