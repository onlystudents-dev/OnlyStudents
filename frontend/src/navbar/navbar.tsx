import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHouseChimney, faStar, faStopwatch, faTable} from "@fortawesome/free-solid-svg-icons";
import Button from "../util/button/button.tsx";

export default function Navbar() {
    return (
        <>
            <nav>
                <div className="flex flex-row items-center gap-2">
                    <Button href="/timetable">
                        <FontAwesomeIcon icon={faTable} /> Timetable
                    </Button>
                    <Button href="/grades">
                        <FontAwesomeIcon icon={faStar} /> Grades
                    </Button>
                    <Button href="/homeworks">
                        <FontAwesomeIcon icon={faHouseChimney} /> Homeworks
                    </Button>
                    <Button href="/absences">
                        <FontAwesomeIcon icon={faStopwatch} /> Absences
                    </Button>
                </div>
                <div className="flex flex-row-reverse items-center gap-2">

                </div>
            </nav>
        </>
    )
}