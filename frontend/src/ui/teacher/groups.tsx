import type {FeatureProps} from "../../types/props.ts";
import {useCallback, useEffect} from "preact/compat";
import {useReloadSchoolStuff} from "../../util/school.ts";
import {fetchWithSchool} from "../../util/api.tsx";
import {useState} from "preact/hooks";
import type {Group} from "../../types/api.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faTrash} from "@fortawesome/free-solid-svg-icons";
import {toast} from "react-toastify";
import {fromResponse, getKey} from "../../util/language.ts";
import Skeleton from "../../util/skeleton/skeleton.tsx";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function Groups({ me, fetchMe }: FeatureProps) {
    const [groups, setGroups] = useState<Group[]>()

    const Fetch = useCallback(async () => {
        return fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/group").then(async response => setGroups(await response.json() as Group[]))
    }, [fetchMe])

    useEffect(() => void Fetch(), [Fetch])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {groups?.map((group) => (
                    <div className="flex flex-row items-center justify-between p-4 bg-(--hover-color) rounded-3xl h-16 w-160 border-b-2 border-(--border-color)">
                        <p className="text-2xl">{group.group_name}</p>
                        <FontAwesomeIcon icon={faTrash} className="action-button hover:text-(--wrong-color)" onClick={() => deleteGroup(group.id)} />
                    </div>
                )) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-(--border-color) rounded-3xl!" />
                ))}
            </div>
        </>
    )

    async function deleteGroup(id: number) {
        const response = await fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/group", {
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
            setGroups(groups?.filter(group => group.id !== id))
            toast.success(getKey("DELETED"))
        }
    }
}