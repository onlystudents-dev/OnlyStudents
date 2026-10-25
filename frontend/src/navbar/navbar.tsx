import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHouseChimney, faPenToSquare, faStar, faStopwatch, faTable} from "@fortawesome/free-solid-svg-icons";
import Button from "../util/button/button.tsx";
import type {Me} from "../types/api.ts";
import {getKey} from "../util/language.ts";
import {useCallback, useRef, useState} from "react";
import Skeleton from "../util/skeleton/skeleton.tsx";
import {features} from "../features.ts";
import {reloadSchoolStuff} from "../util/school.ts";
import {useLayoutEffect} from "preact/hooks";
import Select from "react-select";
import {selectStyles} from "../util/select.tsx";

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
        if (path === href && !window.location.search) return
        navigate(href)
    }

    function getLeft(button: Element): number {
        const container = button.closest(".buttons")

        if (!container) return 0

        return button.getBoundingClientRect().left - container.getBoundingClientRect().left
    }

    function handleEnter(e: MouseEvent) {
        const element = e.currentTarget as Element

        setHidden(false)
        setHoverLeft(getLeft(element))
        setHoverWidth(element.getBoundingClientRect().width)
    }

    function handleClick(e: MouseEvent) {
        e.preventDefault()
        redirect((e.currentTarget as HTMLAnchorElement).pathname)

        const element = e.currentTarget as Element

        setSelectLeft(getLeft(element))
        setSelectWidth(element.getBoundingClientRect().width)
    }

    const handleSelect = useCallback(() => {
        const button = buttons.current.get(path)

        if (!button) {
            setSelected(false)
            return
        }

        setSelected(true)
        setSelectLeft(getLeft(button))
        setSelectWidth(button.getBoundingClientRect().width)
    }, [path])

    useLayoutEffect(() => {
        function handleStuff() {
            handleSelect()

            const hover = [...buttons.current.values()].find(el => el.matches(":hover"))
            if (hover) {
                setHoverLeft(getLeft(hover))
                setHoverWidth(hover.getBoundingClientRect().width)
            }
        }

        handleStuff()

        window.addEventListener("resize", handleStuff)
        return () => window.removeEventListener("resize", handleStuff)
    }, [handleSelect])

    const membershipOptions = memberships?.map(membership => ({ label: membership.school_name, value: String(membership.school_id) }))
    const childrenOptions = me.children?.map(child => ({ label: `${child.first_name} ${child.last_name}`, value: String(child.account_id) }))

    return (
        <>
            <nav>
                <div className="buttons" onMouseLeave={() => setHidden(true)}>
                    {links.map(({href, icon, key}) => (
                        <Button key={href} href={href} className="hover:bg-inherit! z-1" onMouseEnter={handleEnter} onClick={handleClick} refA={element => {if (element) buttons.current.set(href, element); else buttons.current.delete(href)}}>
                            <FontAwesomeIcon icon={icon} /> {getKey(key)}
                        </Button>
                    ))}
                    <div className={`navbbg ${hidden && "hid"}`} style={{ transform: `translateX(${hoverLeft}px)`, width: `${hoverWidth}px` }} />
                    <div className={`navbbg navbbgs ${!selected && "hid"}`} style={{ transform: `translateX(${selectLeft}px)`, width: `${selectWidth}px` }} />
                </div>
                <div className="flex flex-row-reverse items-center gap-2">
                    <a className="rounded-[50%] h-10 cursor-pointer" href="/me" onClick={e => {e.preventDefault(); redirect("/me")}}>
                        <img className={`rounded-[inherit] max-h-full ${!pfpLoaded && "hidden"}`} src={me.preferences.pfp_url} alt="" onLoad={() => setPfpLoaded(true)} />
                        {!pfpLoaded && <Skeleton width={40} height={40} color={"var(--bg-color)"} className="rounded-[inherit]!" />}
                    </a>
                    <Select options={membershipOptions} value={membershipOptions?.find(membership => membership.value === String(school))} onChange={val => val && setSchool(val.value)} styles={selectStyles("fit-content", "100%", "var(--bg-color)", "var(--card-color)", "var(--hover-color)")} />
                    {me.role === "guardian" && <Select options={childrenOptions} value={childrenOptions?.find(chilD => chilD.value === String(child))} onChange={val => val && setChild(val.value)} styles={selectStyles("fit-content", "100%", "var(--bg-color)", "var(--card-color)", "var(--hover-color)")} />}
                </div>
            </nav>
        </>
    )
}