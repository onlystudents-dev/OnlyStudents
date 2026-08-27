import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faStar, faTable} from "@fortawesome/free-solid-svg-icons";
import Button from "../util/button/button.tsx";

export default function Navbar() {
    return (
        <>
            <nav>
                <div className="flex flex-row items-center gap-2">
                    <Button>
                        <FontAwesomeIcon icon={faTable} /> Timetable
                    </Button>
                    <Button>
                        <FontAwesomeIcon icon={faStar} /> Grades
                    </Button>
                </div>
                <div className="flex flex-row-reverse items-center gap-2">

                </div>
            </nav>
        </>
    )
}