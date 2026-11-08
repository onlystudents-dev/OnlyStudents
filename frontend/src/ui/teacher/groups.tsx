import type {FeatureProps} from "../../types/props.ts";
import {useCallback, useEffect} from "preact/compat";
import {useReloadSchoolStuff} from "../../util/school.ts";
import {fetchWithSchool} from "../../util/api.tsx";
import {useState} from "preact/hooks";
import type {Group} from "../../types/api.ts";

export default function Groups({ me, fetchMe }: FeatureProps) {
    const [groups, setGroups] = useState<Group[]>([])

    const Fetch = useCallback(async () => {
        return fetchWithSchool(fetchMe, "/api/v1/me/teacher/timetable/group").then(async response => setGroups(await response.json() as Group[]))
    }, [fetchMe])

    useEffect(() => void Fetch(), [Fetch])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {groups.map((group) => (
                    <div className="flex flex-row items-center justify-between p-4 bg-(--hover-color) rounded-3xl h-16 w-160 border-b-2 border-(--border-color)">
                        <p className="text-2xl">{group.group_name}</p>
                    </div>
                ))}
            </div>
        </>
    )
}