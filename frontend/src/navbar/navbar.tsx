import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHouseChimney, faStar, faStopwatch, faTable} from "@fortawesome/free-solid-svg-icons";
import Button from "../util/button/button.tsx";
import type {Me} from "../app.tsx";
import {getKey} from "../util/language.ts";
import {useState} from "react";
import Skeleton from "../util/skeleton/skeleton.tsx";

export default function Navbar({ me }: {me: Me}) {
    const [pfpLoaded, setPfpLoaded] = useState(false)

    return (
        <>
            <nav>
                <div className="flex flex-row items-center gap-2">
                    <Button href="/">
                        <FontAwesomeIcon icon={faHouseChimney} /> {getKey("HOME")}
                    </Button>
                    <Button href="/timetable">
                        <FontAwesomeIcon icon={faTable} /> {getKey("TIMETABLE")}
                    </Button>
                    <Button href="/grades">
                        <FontAwesomeIcon icon={faStar} /> {getKey("GRADES")}
                    </Button>
                    <Button href="/homeworks">
                        <FontAwesomeIcon icon={faHouseChimney} /> {getKey("HOMEWORKS")}
                    </Button>
                    <Button href="/absences">
                        <FontAwesomeIcon icon={faStopwatch} /> {getKey("ABSENCES")}
                    </Button>
                </div>
                <div className="flex flex-row-reverse items-center gap-2">
                    <a className="rounded-[50%] h-10 cursor-pointer" href="/me">
                        <img className={`rounded-[inherit] max-h-full ${!pfpLoaded && "hidden"}`} src={me.pfp_url} alt="" onLoad={() => setPfpLoaded(true)} />
                        {!pfpLoaded && <Skeleton width={40} height={40} color={"var(--bg-color)"} className="rounded-[inherit]!" />}
                    </a>
                </div>
            </nav>
        </>
    )
}