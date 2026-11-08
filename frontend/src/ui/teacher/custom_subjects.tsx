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

export default function CustomSubjects({ me, fetchMe }: FeatureProps) {
    const [subjects, setSubjects] = useState<CustomSubject[]>()
    const [names, setNames] = useState<Record<number, string>>({})
    const [name, setName] = useState("")

    const Fetch = useCallback(async () => {
        return fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject").then(async response => setSubjects(await response.json() as CustomSubject[]))
    }, [fetchMe])

    useEffect(() => void Fetch(), [Fetch])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {subjects?.map((subject) => (
                    <div className="flex flex-row items-center justify-between p-4 bg-(--hover-color) rounded-3xl h-16 w-160 border-b-2 border-(--border-color)">
                        <input type="text" defaultValue={subject.subject_name} className="text-2xl outline-none" onChange={e => setNames({ ...names, [subject.id]: e.currentTarget.value })} />
                        <FontAwesomeIcon icon={faTrash} className="action-button hover:text-(--wrong-color)" onClick={() => deleteSubject(subject.id)} />
                    </div>
                )) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-(--border-color) rounded-3xl!" />
                ))}
                <div className="flex flex-row items-center justify-between p-4 bg-(--hover-color) rounded-3xl h-16 w-160 border-b-2 border-(--border-color)">
                    <input type="text" value={name} className="text-2xl outline-none" onChange={e => setName(e.currentTarget.value)} />
                    <FontAwesomeIcon icon={faPlus} className="action-button hover:text-(--right-color)" onClick={createSubject} />
                </div>
            </div>
            <Save options={names} setOptions={setNames} save={saveNames} />
        </>
    )

    async function createSubject() {
        const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                subject_name: name,
            })
        })

        if (!response.ok) {
            toast.error(await fromResponse(response))
        } else {
            await Fetch()
            setName("")
        }
    }

    async function saveNames(diff: Partial<Record<number, string>>) {
        const failed: number[] = []
        for (const [id, name] of Object.entries(diff)) {
            const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/custom_subject", {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    id: Number(id),
                    subject_name: name,
                })
            })

            if (!response.ok) {
                failed.push(Number(id))
                toast.error(await fromResponse(response))
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
            setSubjects(subjects?.filter(subject => subject.id !== id))
            toast.success(getKey("DELETED"))
        }
    }
}