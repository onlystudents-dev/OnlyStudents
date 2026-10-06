import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHouseChimney, faStar, faStopwatch, faTable} from "@fortawesome/free-solid-svg-icons";
import Button from "../util/button/button.tsx";
import type {Me} from "../types/api.ts";
import {getKey} from "../util/language.ts";
import {useState} from "react";
import Skeleton from "../util/skeleton/skeleton.tsx";

const links = [
    {href: "/", icon: faHouseChimney, key: "HOME"},
    {href: "/timetable", icon: faTable, key: "TIMETABLE"},
    {href: "/grades", icon: faStar, key: "GRADES"},
    {href: "/homeworks", icon: faHouseChimney, key: "HOMEWORKS"},
    {href: "/absences", icon: faStopwatch, key: "ABSENCES"},
]

export default function Navbar({ me, reloadSchoolStuff }: {me: Me, reloadSchoolStuff?: () => void}) {
    const [pfpLoaded, setPfpLoaded] = useState(false)

    const [child, setChildR] = useState<number | undefined>(() => {
        if (me.role !== "guardian") return undefined

        const stored = localStorage.getItem("child_id")
        const val = stored ? Number(stored) : me.children?.[0]?.account_id

        localStorage.setItem("child_id", String(val))
        return val
    })

    const [school, setSchoolR] = useState<number | undefined>(() => {
        const stored = localStorage.getItem("school_id")
        if (stored) return Number(stored)

        const memberships = me.role === "guardian"
            ? me.children?.find(c => c.account_id === child)?.school_memberships
            : me.school_memberships

        const val = memberships?.[0]?.school_id
        localStorage.setItem("school_id", String(val))

        return val
    })

    function setSchool(schoolId: string) {
        localStorage.setItem("school_id", schoolId)
        setSchoolR(Number(schoolId))
        reloadSchoolStuff?.()
    }

    function setChild(childId: string) {
        const id = Number(childId)
        const memberships = me.children?.find(c => c.account_id === id)?.school_memberships
        const schoolId = memberships?.[0]?.school_id

        localStorage.setItem("child_id", childId)
        setChildR(id)

        if (schoolId !== undefined) {
            localStorage.setItem("school_id", String(schoolId))
            setSchoolR(schoolId)
        }

        reloadSchoolStuff?.()
    }

    const memberships = me.role === "guardian"
        ? me.children?.find(c => c.account_id === child)?.school_memberships
        : me.school_memberships

    return (
        <>
            <nav>
                <div className="flex flex-row items-center gap-2">
                    {links.map(({href, icon, key}) => (
                        <Button href={href} key={href}>
                            <FontAwesomeIcon icon={icon} /> {getKey(key)}
                        </Button>
                    ))}
                </div>
                <div className="flex flex-row-reverse items-center gap-2">
                    <a className="rounded-[50%] h-10 cursor-pointer" href="/me">
                        <img className={`rounded-[inherit] max-h-full ${!pfpLoaded && "hidden"}`} src={me.preferences.pfp_url} alt="" onLoad={() => setPfpLoaded(true)} />
                        {!pfpLoaded && <Skeleton width={40} height={40} color={"var(--bg-color)"} className="rounded-[inherit]!" />}
                    </a>
                    {me.role === "guardian" && <select className="button text-start!" value={child} onChange={e => setChild(e.currentTarget.value)}>
                        {me.children?.map(child => (
                            <option value={child.account_id} key={child.account_id}>{child.first_name} {child.last_name}</option>
                        ))}
                    </select>}
                    <select className="button text-start!" value={school} onChange={e => setSchool(e.currentTarget.value)}>
                        {memberships?.map(membership => (
                            <option value={membership.school_id} key={membership.school_id}>{membership.school_name}</option>
                        ))}
                    </select>
                </div>
            </nav>
        </>
    )
}