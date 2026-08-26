import "./navbar.css";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faHome, faUser} from "@fortawesome/free-solid-svg-icons";
import {useState} from "react";
import Login from "../auth/login/login.tsx";
import Button from "../button/button.tsx";

export default function Navbar() {
    const [login, setLogin] = useState(false);

    return (
        <>
            <nav>
                <Button href="/">
                    <FontAwesomeIcon icon={faHome} /> Home
                </Button>

                <Button onClick={() => setLogin(!login)}>
                    <FontAwesomeIcon icon={faUser} /> Login
                </Button>
            </nav>

            {login && <Login setLogin={setLogin} />}
        </>
    )
}