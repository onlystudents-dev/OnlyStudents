import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHouseChimney, faPenToSquare, faStar, faStopwatch, faTable} from "@fortawesome/free-solid-svg-icons";
import Button from "../util/button/button.tsx";
import type {Me} from "../types/api.ts";
import {getKey} from "../util/language.ts";
import {useRef, useState} from "react";
import Skeleton from "../util/skeleton/skeleton.tsx";
import {features} from "../features.ts";
import {reloadSchoolStuff} from "../util/school.ts";
import {useLayoutEffect} from "preact/hooks";

export default function Navbar({ me, path, navigate }: {me: Me, path: string, navigate: (href: string) => void}) {
    const icons = {
        homeworks: faHouseChimney,
        absences: faStopwatch,
        timetable: faTable,
        grades: faStar,
        exams: faPenToSquare,
    }

    const links = [
        { href: "/", icon: faHouseChimney, key: "HOME" },
        ...Object.entries(features)
            .filter(([, roles]) => me.role in roles)
            .map(([feature]) => ({
                href: `/${feature}`,
                icon: icons[feature as keyof typeof icons],
                key: feature.toUpperCase(),
            })),
    ]

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
        reloadSchoolStuff()
    }

    function setChild(childId: string) {
        const id = Number(childId)
        const selectedChild = me.children?.find(c => c.account_id === id)
        const memberships = selectedChild?.school_memberships ?? []
        const schoolId = memberships[0]?.school_id

        localStorage.setItem("child_id", String(id))
        setChildR(id)

        if (schoolId !== undefined) {
            localStorage.setItem("school_id", String(schoolId))
            setSchoolR(schoolId)
        } else {
            localStorage.removeItem("school_id")
            setSchoolR(undefined)
        }

        reloadSchoolStuff()
    }

    const memberships = me.role === "guardian"
        ? me.children?.find(c => c.account_id === child)?.school_memberships
        : me.school_memberships

    const [hoverLeft, setHoverLeft] = useState(0)
    const [hoverWidth, setHoverWidth] = useState(0)
    const [selectLeft, setSelectLeft] = useState(0)
    const [selectWidth, setSelectWidth] = useState(0)

    const buttons = useRef(new Map<string, HTMLAnchorElement>())
    const [selected, setSelected] = useState(false)
    const [hidden, setHidden] = useState(true)

    function redirect(href: string) {
        if (path === href) return
        navigate(href)
    }

    function getLeft(e: MouseEvent): number {
        const button = e.currentTarget as Element
        const container = button.closest(".buttons")

        if (!container) return 0

        return button.getBoundingClientRect().left - container.getBoundingClientRect().left
    }

    function handleEnter(e: MouseEvent) {
        setHidden(false)
        setHoverLeft(getLeft(e))
        setHoverWidth((e.currentTarget as Element).getBoundingClientRect().width)
    }

    function handleClick(e: MouseEvent) {
        e.preventDefault()
        redirect((e.currentTarget as HTMLAnchorElement).pathname)

        setSelectLeft(getLeft(e))
        setSelectWidth((e.currentTarget as Element).getBoundingClientRect().width)
    }

    useLayoutEffect(() => {
        const button = buttons.current.get(path)

        if (!button) {
            setSelected(false)
            return
        }

        const container = button.closest(".buttons")
        if (!container) return

        setSelected(true)
        setSelectLeft(button.getBoundingClientRect().left - container.getBoundingClientRect().left)
        setSelectWidth(button.getBoundingClientRect().width)
    }, [path])

    return (
        <>
            <nav>
                <div className="buttons" onMouseLeave={() => setHidden(true)}>
                    {links.map(({href, icon, key}) => (
                        <Button key={href} href={href} className="hover:bg-inherit! z-1" onMouseEnter={handleEnter} onClick={handleClick} refA={element => {if (element) buttons.current.set(href, element); else buttons.current.delete(href)}}>
                            <FontAwesomeIcon icon={icon} /> {getKey(key)}
                        </Button>
                    ))}
                    <div className={`navbbg ${hidden && "hid"}`} style={{ left: `${hoverLeft}px`, width: `${hoverWidth}px` }} />
                    <div className={`navbbg navbbgs ${!selected && "hid"}`} style={{ left: `${selectLeft}px`, width: `${selectWidth}px` }} />
                </div>
                <div className="flex flex-row-reverse items-center gap-2">
                    <a className="rounded-[50%] h-10 cursor-pointer" href="/me" onClick={e => {e.preventDefault(); redirect("/me")}}>
                        <img className={`rounded-[inherit] max-h-full ${!pfpLoaded && "hidden"}`} src={me.preferences.pfp_url} alt="" onLoad={() => setPfpLoaded(true)} />
                        {!pfpLoaded && <Skeleton width={40} height={40} color={"var(--bg-color)"} className="rounded-[inherit]!" />}
                    </a>
                    <select className="button text-start!" value={school} onChange={e => setSchool(e.currentTarget.value)}>
                        {memberships?.map(membership => (
                            <option value={membership.school_id} key={membership.school_id}>{membership.school_name}</option>
                        ))}
                    </select>
                    {me.role === "guardian" && <select className="button text-start!" value={child} onChange={e => setChild(e.currentTarget.value)}>
                        {me.children?.map(child => (
                            <option value={child.account_id} key={child.account_id}>{child.first_name} {child.last_name}</option>
                        ))}
                    </select>}
                </div>
            </nav>
        </>
    )
}