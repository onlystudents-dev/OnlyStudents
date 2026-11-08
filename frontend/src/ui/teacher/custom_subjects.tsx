import type {FeatureProps} from "../../types/props.ts";
import {useCallback, useEffect} from "preact/compat";
import {fetchWithSchool} from "../../util/api.tsx";
import {useState} from "preact/hooks";
import type {CustomSubject} from "../../types/api.ts";
import {useReloadSchoolStuff} from "../../util/school.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPlus, faTrash} from "@fortawesome/free-solid-svg-icons";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {fromResponse, getKey} from "../../util/language.ts";
import {toast} from "react-toastify";
import Save from "../../util/save/save.tsx";
import SubjectIcon from "../../util/subject_icon.tsx";

type SubjectEdit = { name: string, icon: string }

function withoutKey(record: Record<number, SubjectEdit>, key: number): Record<number, SubjectEdit> {
    return Object.fromEntries(
        Object.entries(record).filter(([k]) => Number(k) !== key)
    )
}

export default function CustomSubjects({ fetchMe }: FeatureProps) {
    const [icons, setIcons] = useState<string[]>([])
    const [subjects, setSubjects] = useState<CustomSubject[]>()
    const [data, setData] = useState<Record<number, SubjectEdit>>({})
    const [name, setName] = useState("")
    const [icon, setIcon] = useState("")

    const Fetch = useCallback(async () => Promise.all([
        fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject").then(async response => setSubjects(await response.json() as CustomSubject[])),
        fetch("/assets/fonts/material.json").then(async response => setIcons(await response.json() as string[]))
    ]), [fetchMe])

    useEffect(() => void Fetch(), [Fetch])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {subjects?.map((subject) => {
                    const current = getCurrent(subject)

                    return (
                        <div key={subject.id} className="flex flex-row items-center justify-between p-4 bg-(--hover-color) rounded-3xl h-16 w-160 border-b-2 border-(--border-color)">
                            <div className="flex flex-row">
                                <SubjectIcon icon={isIcon(current.icon)} />
                                <div className="flex flex-col">
                                    <input type="text" value={current.name} className="outline-none" onInput={e => edit(subject, { name: e.currentTarget.value })} placeholder={getKey("NAME")} />
                                    <input type="text" value={current.icon} className="outline-none opacity-70" onInput={e => edit(subject, { icon: e.currentTarget.value })} placeholder={getKey("ICON")} />
                                    {/* todo: little ? saying you can find these at https://fonts.google.com/icons */}
                                </div>
                            </div>
                            <FontAwesomeIcon icon={faTrash} className="action-button hover:text-(--wrong-color)" onClick={() => deleteSubject(subject.id)} />
                        </div>
                    )
                }) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-(--border-color) rounded-3xl!" />
                ))}
                <div className="flex flex-row items-center justify-between p-4 bg-(--hover-color) rounded-3xl h-16 w-160 border-b-2 border-(--border-color)">
                    <div className="flex flex-row">
                        <SubjectIcon icon={isIcon(icon)} />
                        <div className="flex flex-col">
                            <input type="text" value={name} className="outline-none" onInput={e => setName(e.currentTarget.value)} placeholder={getKey("NAME")} />
                            <input type="text" value={icon} className="outline-none opacity-70" onInput={e => setIcon(e.currentTarget.value)} placeholder={getKey("ICON")} />
                        </div>
                    </div>
                    <FontAwesomeIcon icon={faPlus} className="action-button hover:text-(--right-color)" onClick={createSubject} />
                </div>
            </div>
            <Save options={data} setOptions={setData} save={saveNames} />
        </>
    )

    function getCurrent(subject: CustomSubject): SubjectEdit {
        return data[subject.id] ?? { name: subject.subject_name, icon: subject.subject_icon }
    }

    function edit(subject: CustomSubject, patch: Partial<SubjectEdit>) {
        const next = { ...getCurrent(subject), ...patch }
        const unchanged = next.name === subject.subject_name && next.icon === subject.subject_icon

        setData(prev => unchanged
            ? withoutKey(prev, subject.id)
            : { ...prev, [subject.id]: next }
        )
    }

    function isIcon(name: string) {
        return icons.includes(name) ? name : ""
    }

    async function createSubject() {
        const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                subject_name: name,
                subject_icon: isIcon(icon),
            })
        })

        if (!response.ok) {
            toast.error(await fromResponse(response))
        } else {
            await Fetch()
            setName("")
            setIcon("")
        }
    }

    async function saveNames(diff: Partial<Record<number, SubjectEdit>>) {
        const failed: number[] = []

        for (const [ID, entry] of Object.entries(diff)) {
            const id = Number(ID)
            if (!entry || !subjects?.some(subject => subject.id === id)) continue

            const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject", {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    id,
                    subject_name: entry.name,
                    subject_icon: isIcon(entry.icon),
                })
            })

            if (!response.ok) {
                failed.push(id)
                toast.error(await fromResponse(response))
            } else {
                setSubjects(prev => prev?.map(subject => subject.id === id
                    ? { ...subject, subject_name: entry.name, subject_icon: isIcon(entry.icon) }
                    : subject))
            }
        }

        return failed
    }

    async function deleteSubject(id: number) {
        if (!confirm(getKey("SURE?"))) return

        const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject", {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                id
            })
        })

        if (!response.ok) {
            toast.error(await fromResponse(response))
        } else {
            setSubjects(prev => prev?.filter(subject => subject.id !== id))
            setData(prev => withoutKey(prev, id))
            toast.success(getKey("DELETED"))
        }
    }
}