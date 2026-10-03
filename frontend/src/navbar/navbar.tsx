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

export default function Navbar({ me }: {me: Me}) {
    const [pfpLoaded, setPfpLoaded] = useState(false)

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
                </div>
            </nav>
        </>
    )
}